package main

import (
	"fmt"

	"github.com/nicholasss/async-messages/internal/client"
)

func main() {
	c, err := client.NewClientConfig("Craig", "AlfredoExpress")
	if err != nil {
		fmt.Printf("unable to create new client config due to: %q\n", err)
		return
	}

	err = c.StartClient()
	if err != nil {
		fmt.Printf("unable to start new client due to: %q\n", err)
	}

	err = c.WriteMessageIntoQueue("Bob", "Snow", "Shovel", "We should get going on tuesday.")
	if err != nil {
		fmt.Printf("cannot write message due to: %q\n", err)
	}

	err = c.SendAllFromQueue()
	if err != nil {
		fmt.Printf("cannot send due to: %q\n", err)
	}

	fmt.Printf("Queue Summary: %s\n", c.Outbox.QueueSummary())
}
