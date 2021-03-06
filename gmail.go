package gmail

import (
	"fmt"
	"mime"
	"net/smtp"
	"regexp"
	"strings"

	"code.google.com/p/mahonia"

	"github.com/jhillyerd/go.enmime"
	"github.com/zond/gmail/imap"
	"github.com/zond/gmail/xmpp"
)

func DecodeText(body, mimeContent string) string {
	_, params, err := mime.ParseMediaType(mimeContent)
	if err != nil {
		return body
	}
	charset := params["charset"]
	if charset == "" {
		return body
	}
	decoder := mahonia.NewDecoder(charset)
	if decoder == nil {
		return body
	}
	return decoder.ConvertString(body)
}

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
		errorHandler: func(e error) {
			fmt.Println("Error", e)
		},
	}
	result.xmppClient.MailHandler(func() {
		if err := result.imapClient.HandleNew(result.mailHandler); err != nil {
			result.errorHandler(err)
		}
	}).ErrorHandler(func(e error) {
		result.errorHandler(e)
	})
	return
}

var AddrReg = regexp.MustCompile("(?i)[=A-Z0-9._%+-]+@[A-Z0-9.-]+\\.[A-Z]{2,4}")

func (self *Client) Send(from, subject, message string, recips ...string) (err error) {
	body := fmt.Sprintf("Content-Type: text/plain; charset=\"utf-8\"\r\nReply-To: %v\r\nFrom: %v\r\nTo: %v\r\nSubject: %v\r\n\r\n%v", from, from, strings.Join(recips, ", "), subject, message)
	auth := smtp.PlainAuth("", self.account, self.password, "smtp.gmail.com")
	actualRecips := []string{}
	for _, recip := range recips {
		if match := AddrReg.FindString(recip); match != "" {
			actualRecips = append(actualRecips, match)
		}
	}
	return smtp.SendMail("smtp.gmail.com:587", auth, self.account, actualRecips, []byte(body))
}

func (self *Client) Debug() *Client {
	self.xmppClient.Debug()
	return self
}

func (self *Client) ErrorHandler(f func(e error)) *Client {
	self.errorHandler = f
	return self
}

func (self *Client) MailHandler(f imap.MailHandler) *Client {
	self.mailHandler = f
	return self
}

func (self *Client) Start() (result *Client, err error) {
	if err = self.xmppClient.Start(); err != nil {
		return
	}
	if err = self.imapClient.HandleNew(self.mailHandler); err != nil {
		return
	}
	result = self
	return
}

func (self *Client) Close() error {
	return self.xmppClient.Close()
}
