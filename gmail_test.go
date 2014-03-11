package gmail

import (
	"fmt"
	"net/mail"
	"os"
	"testing"
)

func TestNotifications(t *testing.T) {
	c := NewClient(os.Getenv("GMAIL_ACCOUNT"), os.Getenv("GMAIL_PASSWORD"), func(msg *mail.Message) error {
		fmt.Println(msg)
		return nil
	})
	err := c.Start()
	if err != nil {
		t.Fatalf("%v", err)
	}
	<-(make(chan bool))
	c.Close()
}
