package gmail

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/jhillyerd/go.enmime"
	"github.com/zond/gmail/imap"
	"github.com/zond/gmail/xmpp"
)

type Client struct {
	account      string
	password     string
	xmppClient   *xmpp.Client
	imapClient   *imap.Client
	mailHandler  imap.MailHandler
	errorHandler func(e error)
}

func New(account, password string) (result *Client) {
	result = &Client{
		account:    account,
		password:   password,
		xmppClient: xmpp.New(account, password),
		imapClient: imap.New(account, password),
		mailHandler: func(msg *enmime.MIMEBody) error {
			fmt.Println("Got", msg)
			return nil
		},
	}
	result.xmppClient.MailHandler(func() {
		result.imapClient.HandleNew(result.mailHandler)
	})
	return
}

func (self *Client) Send(subject, message string, recips ...string) (err error) {
	body := fmt.Sprintf("To: %v\r\nSubject: %v\r\n\r\n%v", strings.Join(recips, ", "), subject, message)
	auth := smtp.PlainAuth("", self.account, self.password, "smtp.gmail.com")
	return smtp.SendMail("smtp.gmail.com:587", auth, self.account, recips, []byte(body))
}

func (self *Client) Debug() *Client {
	self.xmppClient.Debug()
	return self
}

func (self *Client) ErrorHandler(f func(e error)) *Client {
	self.xmppClient.ErrorHandler(f)
	return self
}

func (self *Client) MailHandler(f imap.MailHandler) *Client {
	self.mailHandler = f
	return self
}

func (self *Client) Start() *Client {
	self.xmppClient.Start()
	self.imapClient.HandleNew(self.mailHandler)
	return self
}

func (self *Client) Close() error {
	return self.xmppClient.Close()
}
