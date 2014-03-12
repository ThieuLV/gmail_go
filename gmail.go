package gmail

import (
	"fmt"

	"github.com/jhillyerd/go.enmime"
	"github.com/zond/gmail/imap"
	"github.com/zond/gmail/xmpp"
)

type Client struct {
	xmppClient   *xmpp.Client
	imapClient   *imap.Client
	mailHandler  imap.MailHandler
	errorHandler func(e error)
}

func New(account, password string) (result *Client) {
	result = &Client{
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

func (self *Client) Start() {
	self.xmppClient.Start()
	self.imapClient.HandleNew(self.mailHandler)
	return
}

func (self *Client) Close() error {
	return self.xmppClient.Close()
}
