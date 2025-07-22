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

	// adding msg 1
	testMsg1, err := msg.NewMessage("Self", "Self", "Echoing message", "Hearing me ok?", c.SecretKey)
	if err != nil {
		log.Printf("could not create message due to: %q", err)
		return
	}
	c.AddToQueue(testMsg1)

	// adding msg 2
	testMsg2, err := msg.NewMessage("Self", "Self", "Echoing message", "Hearing me ok?", c.SecretKey)
	if err != nil {
		log.Printf("could not create message due to: %q", err)
		return
	}
	c.AddToQueue(testMsg2)

	// sending one of the messages
	err = c.SendFromQueue()
	if err != nil {
		log.Printf("unable to send from queue due to: %q", err)
	} else {
		log.Printf("sent queued message")
	}

	log.Printf("Queue dump >\n%s", c.Queue.DumpToString())
}
