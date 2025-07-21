package main

import (
	"log"

	"github.com/nicholasss/async-messages/internal/client"
	"github.com/nicholasss/async-messages/internal/msg"
)

func main() {
	c, err := client.NewClient("Craig McMuffin")
	if err != nil {
		log.Printf("unable to create new client due to: %q", err)
		return
	}

	testMsg, err := msg.NewMessage("Me", "Me", "Echoing message", "Hearing me ok?", c.SecretKey)
	if err != nil {
		log.Printf("could not create message due to: %q", err)
		return
	}

	c.AddToQueue(testMsg)
	err = c.SendFromQueue()
	if err != nil {
		log.Printf("unable to send from queue due to: %q", err)
	} else {
		log.Printf("sent queued message")
	}
}
