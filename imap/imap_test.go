package imap

import (
	"os"
	"testing"
)

func TestIMAPGet(t *testing.T) {
	c := New(os.Getenv("GMAIL_ACCOUNT"), os.Getenv("GMAIL_PASSWORD"))
	_, err := c.GetNew()
	if err != nil {
		panic(err)
	}
}
