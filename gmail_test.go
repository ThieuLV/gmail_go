package gmail

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNotifications(t *testing.T) {
	c := NewClient(os.Getenv("GMAIL_ACCOUNT"), os.Getenv("GMAIL_PASSWORD"), time.Now().Add(-21*time.Minute), func(i interface{}) {
		fmt.Println(i)
	})
	err := c.Start()
	if err != nil {
		t.Fatalf("%v", err)
	}
	<-(make(chan bool))
	c.Close()
}
