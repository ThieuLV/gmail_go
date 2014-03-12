package gmail

import (
	"github.com/zond/gmail/imap"
	"github.com/zond/gmail/xmpp"
)

type Client struct {
	xmppClient *xmpp.Client
	imapClient *imap.Client
}

func New(account, password string) (result *Client) {
	result = &Client{
		xmppClient: xmpp.New(account, password),
		imapClient: imap.New(account, password),
	}
	return
}

func (self *Client) Start() {
	self.xmppClient.Start()
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
