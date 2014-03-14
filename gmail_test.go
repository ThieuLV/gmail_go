package gmail

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jhillyerd/go.enmime"
)

func TestNotifications(t *testing.T) {
	inc := make(chan *enmime.MIMEBody)
	c, err := New(os.Getenv("GMAIL_ACCOUNT"), os.Getenv("GMAIL_PASSWORD")).MailHandler(func(msg *enmime.MIMEBody) error {
		inc <- msg
		return nil
	}).Start()
	if err != nil {
		t.Fatalf("%v", err)
	}
	rand.Seed(time.Now().UnixNano())
	subj := fmt.Sprint(rand.Int63())
	body := fmt.Sprint(rand.Int63())
	if err := c.Send(os.Getenv("GMAIL_ACCOUNT"), subj, body, os.Getenv("GMAIL_ACCOUNT")); err != nil {
		t.Fatalf("%v", err)
	}
	msg := <-inc
	if strings.TrimSpace(msg.Text) != body {
		t.Errorf("Wrong body. Wanted %#v but got %#v", body, strings.TrimSpace(msg.Text))
	}
	if strings.TrimSpace(msg.GetHeader("Subject")) != subj {
		t.Errorf("Wrong subject. Wanted %#v but got %#v", subj, strings.TrimSpace(msg.GetHeader("Subject")))
	}
	c.Close()
}
