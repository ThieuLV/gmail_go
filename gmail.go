package gmail

import (
	"fmt"
	"net/mail"

	"github.com/zond/gmail/imap"
	"github.com/zond/gmail/xmpp"
)

type Client struct {
	xmppClient  *xmpp.Client
	imapClient  *imap.Client
	mailHandler imap.MailHandler
}

func New(account, password string) (result *Client) {
	result = &Client{
		xmppClient: xmpp.New(account, password),
		imapClient: imap.New(account, password),
		mailHandler: func(msg *mail.Message) error {
			fmt.Println("Got", msg)
			return nil
		},
	}
	result.xmppClient.MailHandler(func() {
		result.imapClient.HandleNew(result.mailHandler)
	})
	return
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
	err1 := self.xmppClient.Close()
	err2 := self.imapClient.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
