package main

import (
	"log"

	"github.com/nicholasss/async-messages/internal/client"
)

func main() {
	c, err := client.NewClient("Craig", "AlfredoExpress")
	if err != nil {
		log.Printf("unable to create new client due to: %q", err)
		return
	}

	// sending one of the messages
	err = c.SendAllFromQueue()
	if err != nil {
		log.Printf("unable to send from queue due to: %q", err)
	} else {
		log.Printf("sent queued message")
	}

	log.Printf("Queue dump >\n%s", c.Queue.QueueSummary())
}
