package imap

import (
	"bytes"
	"io"
	"net/mail"
	"strings"
	"sync"

	"code.google.com/p/go-imap/go1/imap"
)

var OldKeyword = "FETCHEDBYAPI"

type Client struct {
	user     string
	password string
	client   *imap.Client
	lock     *sync.Mutex
}

func New(user, password string) *Client {
	return &Client{
		user:     user,
		password: password,
		lock:     &sync.Mutex{},
	}
}

func (self *Client) connect() (err error) {
	self.client, err = imap.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return
	}
	if _, err = self.client.Login(self.user, self.password); err != nil {
		return
	}
	if _, err = self.client.Select("INBOX", false); err != nil {
		return
	}
	return
}

func (self *Client) Close() (err error) {
	_, err = self.client.Close(false)
	return
}

func (self *Client) GetNew() (result []mail.Message, err error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.client == nil {
		if err = self.connect(); err != nil {
			return
		}
	}
	handler := func(msg *mail.Message) error {
		result = append(result, *msg)
		return nil
	}
	if err = self.handleNew(handler); err != nil && strings.Contains(err.Error(), "closed") {
		if err = self.connect(); err != nil {
			return
		}
		result = nil
		err = self.handleNew(handler)
	}
	return
}

func (self *Client) handleNew(handler func(msg *mail.Message) error) (err error) {
	cmd, err := self.client.UIDSearch("UNKEYWORD " + OldKeyword)
	if err != nil {
		return
	}
	foundSeq := &imap.SeqSet{}
	for cmd.InProgress() {
		self.client.Recv(-1)
		for _, rsp := range cmd.Data {
			for _, res := range rsp.SearchResults() {
				foundSeq.AddNum(res)
			}
		}
		cmd.Data = nil
		self.client.Data = nil
	}

	if !foundSeq.Empty() {
		var fetchCmd *imap.Command
		fetchCmd, err = self.client.UIDFetch(foundSeq, "RFC822.TEXT", "RFC822.HEADER")
		if err != nil {
			return
		}
		markSeq := &imap.SeqSet{}
		for fetchCmd.InProgress() {
			self.client.Recv(-1)
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
				if e := handler(msg); e != nil {
					markSeq.AddNum(rsp.MessageInfo().UID)
				}
			}
			fetchCmd.Data = nil
			self.client.Data = nil
		}
		if !markSeq.Empty() {
			if _, err = imap.Wait(self.client.Store(markSeq, "FLAGS", []imap.Field{OldKeyword})); err != nil {
				return
			}
		}
	}
	return
}
