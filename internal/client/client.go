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

// HealthCheck is what should be recieved from the server hitting 'GET /' endpoint
type HealthCheck struct {
	Health string `json:"health"`
}

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
		Server:    "http://localhost:8080",
	}, nil
}

func (c *Client) checkServerIsOnline() error {
	res, err := c.Client.Get(c.Server + "/health")
	if err != nil {
		return fmt.Errorf("unable to connect to server: '%d %s' due to: %w", res.StatusCode, res.Status, err)
	}
	defer res.Body.Close()

	// explicit response check
	var healthRes HealthCheck
	err = json.NewDecoder(res.Body).Decode(&healthRes)
	if err != nil {
		return fmt.Errorf("unable to decode server health response: %w", err)
	}
	if healthRes.Health != "OK" {
		return errors.New("unable to determine health of server")
	}

	// health of server ok past this point
	return nil
}

func (c *Client) AddToQueue(msg *msg.Message) {
	c.Queue.Enqueue(*msg)
}

func (c *Client) sendMessage(msgToSend *msg.Message) error {
	// verify message before sending
	// TODO: do something with messages that have invalid signatures?
	err := msgToSend.VerifyMessage(c.SecretKey)
	if err != nil {
		return err
	}

	// marshal message into buffer
	msgData, err := json.Marshal(msgToSend)
	if err != nil {
		return err
	}
	msgDataReader := bytes.NewBuffer(msgData)

	// post message
	res, err := c.Client.Post(c.Server+"/send-message", "application/json", msgDataReader)
	if err != nil {
		return err
	}

	// check return status
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("attempted to send message. response status code of '%s %d'", res.Status, res.StatusCode)
	}

	// successful send
	return nil
}

func (c *Client) SendOneFromQueue() error {
	err := c.checkServerIsOnline()
	if err != nil {
		return fmt.Errorf("server is not healthy. cannot send message due to: %w", err)
	}

	msgToSend, ok := c.Queue.Dequeue()
	if !ok {
		return errors.New("unable to dequeue message for sending")
	}

	c.sendMessage(&msgToSend)

	return nil
}

func (c *Client) SendAllFromQueue() error {
	err := c.checkServerIsOnline()
	if err != nil {
		return fmt.Errorf("server is not healthy. cannot send message due to: %w", err)
	}

	// send until queue is empty
	for !c.Queue.IsEmpty() {
		msgToSend, ok := c.Queue.Dequeue()
		if !ok {
			return errors.New("unable to dequeue message for sending")
		}

		err = c.sendMessage(&msgToSend)
		if err != nil {
			return err
		}
	}

	return nil
}
