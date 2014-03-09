package gmail

import (
	"fmt"
	"os"
	"testing"
)

func TestNotifications(t *testing.T) {
	c := NewClient(os.Getenv("GMAIL_ACCOUNT"), os.Getenv("GMAIL_PASSWORD"), func(i interface{}) {
		fmt.Println(i)
	})
	err := c.Start()
	if err != nil {
		t.Fatalf("%v", err)
	}
	<-(make(chan bool))
	c.Close()
}
