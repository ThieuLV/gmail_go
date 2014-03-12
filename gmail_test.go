package gmail

import (
	"os"
	"testing"
)

func TestNotifications(t *testing.T) {
	c := New(os.Getenv("GMAIL_ACCOUNT"), os.Getenv("GMAIL_PASSWORD"))
	c.Start()
	<-(make(chan bool))
	c.Close()
}
