// Package client implements a client to be created and perform common operations
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nicholasss/async-messages/internal/msg"
)

type Client struct {
	SecretKey []byte
	Client    http.Client
	Queue     *msg.Queue
	Name      string
	Server    string
}

func NewClient(name string) (*Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	secretKey := []byte(os.Getenv("HMAC_SECRET"))

	queue := msg.NewQueue()

	return &Client{
		SecretKey: secretKey,
		Client:    *http.DefaultClient,
		Queue:     queue,
		Name:      name,
		Server:    "http://localhost:8080/echo",
	}, nil
}

func (c *Client) checkServerIsOnline() error {
	res, err := c.Client.Get(c.Server)
	if err != nil {
		return fmt.Errorf("unable to connect to server: '%d %s' due to: %w", res.StatusCode, res.Status, err)
	}
	defer res.Body.Close()

	// explicit response check
	resBodyBuffer := make([]byte, 0)
	_, err = res.Body.Read(resBodyBuffer)
	if err != nil {
		return fmt.Errorf("unable to read body of health check response due to: %w", err)
	}
	if !bytes.Equal([]byte(`{"health":"200 OK"}`), resBodyBuffer) {
		return fmt.Errorf("server health unknown: '%s'", resBodyBuffer)
	}

	// health of server ok past this point
	return nil
}

func (c *Client) AddToQueue(msg *msg.Message) {
	c.Queue.Enqueue(*msg)
}

func (c *Client) SendFromQueue() error {
	msgToSend, ok := c.Queue.Dequeue()
	if !ok {
		return errors.New("unable to send message due to issue dequeuing message")
	}

	// verify message for sending
	err := msg.VerifyMessage(&msgToSend, c.SecretKey)
	if err != nil {
		return err
	}

	msgData, err := json.Marshal(msgToSend)
	if err != nil {
		return err
	}
	msgDataReader := bytes.NewBuffer(msgData)

	res, err := c.Client.Post(c.Server, "application/json", msgDataReader)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("attempted to send message. response status code of '%s %d'", res.Status, res.StatusCode)
	}

	return nil
}
