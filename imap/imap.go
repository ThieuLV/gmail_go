package imap

import (
	"bytes"
	"io"
	"net/mail"

	"code.google.com/p/go-imap/go1/imap"
)

type MailHandler func(*mail.Message) error

var OldKeyword = "FETCHEDBYAPI"

type Client struct {
	user     string
	password string
}

func New(user, password string) *Client {
	return &Client{
		user:     user,
		password: password,
	}
}

func (self *Client) connect() (result *imap.Client, err error) {
	result, err = imap.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return
	}
	if _, err = result.Login(self.user, self.password); err != nil {
		return
	}
	if _, err = result.Select("INBOX", false); err != nil {
		return
	}
	return
}

func (self *Client) GetNew() (result []mail.Message, err error) {
	handler := func(msg *mail.Message) error {
		result = append(result, *msg)
		return nil
	}
	if err = self.HandleNew(handler); err != nil {
		return
	}
	return
}

func (self *Client) HandleNew(handler MailHandler) (err error) {
	client, err := self.connect()
	if err != nil {
		return
	}
	defer client.Close(false)
	cmd, err := imap.Wait(client.UIDSearch("UNKEYWORD " + OldKeyword))
	if err != nil {
		return
	}
	foundSeq := &imap.SeqSet{}
	for _, rsp := range cmd.Data {
		for _, res := range rsp.SearchResults() {
			foundSeq.AddNum(res)
		}
	}

	if !foundSeq.Empty() {
		var fetchCmd *imap.Command
		fetchCmd, err = imap.Wait(client.UIDFetch(foundSeq, "RFC822.TEXT", "RFC822.HEADER"))
		if err != nil {
			return
		}
		markSeq := &imap.SeqSet{}
		for _, rsp := range fetchCmd.Data {
			buf := &bytes.Buffer{}
			if _, err = rsp.MessageInfo().Attrs["RFC822.HEADER"].(io.WriterTo).WriteTo(buf); err != nil {
				return
			}
			if _, err = rsp.MessageInfo().Attrs["RFC822.TEXT"].(io.WriterTo).WriteTo(buf); err != nil {
				return
			}
			var msg *mail.Message
			if msg, err = mail.ReadMessage(buf); err != nil {
				return
			}
			if e := handler(msg); e == nil {
				markSeq.AddNum(rsp.MessageInfo().UID)
			}
		}
		if !markSeq.Empty() {
			if _, err = imap.Wait(client.Store(markSeq, "FLAGS", []imap.Field{OldKeyword})); err != nil {
				return
			}
		}
	}
	return
}
