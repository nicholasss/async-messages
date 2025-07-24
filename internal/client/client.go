// Package client implements a client to be created and perform common operations
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/nicholasss/async-messages/internal/msg"
)

// *** Types ***

// HealthCheck is what should be recieved from the server hitting 'GET /' endpoint
type HealthCheck struct {
	Health string `json:"health"`
}

type Config struct {
	SecretKey []byte
	Client    http.Client
	Outbox    *msg.PackagedQueue
	Inbox     *msg.PackagedQueue
	Name      string
	Vessel    string
	Server    string
	Online    *safeBool
}

type NewMessage struct {
	ToName   string
	ToVessel string
	Subject  string
	Body     string
}

// *** Internal Types ***

type safeBool struct {
	mux  sync.RWMutex
	bool bool
}

// *** Errors ***

// ErrServerOffline signifies that the server is offline
var ErrServerOffline = errors.New("server is offline")

// *** New Config ***

func NewClientConfig(name, vessel string) (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	secretKey := []byte(os.Getenv("HMAC_SECRET"))

	// inbox/outbox setup
	outbox := msg.NewQueue()
	inbox := msg.NewQueue()

	// online setup
	safeOnline := &safeBool{
		mux:  sync.RWMutex{},
		bool: false,
	}

	return &Config{
		SecretKey: secretKey,
		Client:    *http.DefaultClient,
		Outbox:    outbox,
		Inbox:     inbox,
		Name:      name,
		Vessel:    vessel,
		Server:    "http://localhost:8080",
		Online:    safeOnline,
	}, nil
}

// *** Functions ***

// StartClient will begin with checking if the server is online or not.
func (c *Config) StartClient() error {
	err := c.checkServerIsOnline()
	if err != nil {
		return err
	}

	return nil
}

// safely get the value of bool
func (bo *safeBool) getValue() bool {
	bo.mux.RLock()
	val := bo.bool
	bo.mux.RUnlock()
	return val
}

// safely set the value of bool
func (bo *safeBool) setValue(val bool) {
	bo.mux.Lock()
	bo.bool = val
	bo.mux.Unlock()
}

func (c *Config) checkServerIsOnline() error {
	res, err := c.Client.Get(c.Server + "/health")
	if err != nil {
		c.Online.setValue(false)
		return ErrServerOffline
	}
	defer res.Body.Close()

	// explicit response check
	var healthRes HealthCheck
	err = json.NewDecoder(res.Body).Decode(&healthRes)
	if err != nil {
		c.Online.setValue(false)
		return ErrServerOffline
	}
	if healthRes.Health != "OK" {
		c.Online.setValue(false)
		return ErrServerOffline
	}

	// health of server ok past this point
	c.Online.setValue(true)
	return nil
}

func (c *Config) getMessagesFromServer() error {
	// get response for specific user
	res, err := c.Client.Get(c.Server + "/get-messages")
	if err != nil {
		return err
	}

	// check return status
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("attempted to send message. response status code of '%s %d'", res.Status, res.StatusCode)
	}

	return nil
}

// WriteMessageIntoQueue crafts a message and inserts it into the clients queue.
func (c *Config) WriteMessageIntoQueue(toName, toVessel, subject, body string) error {
	newMessage := &msg.RawMessage{
		ToName:     toName,
		ToVessel:   toVessel,
		FromName:   c.Name,
		FromVessel: c.Vessel,
		Subject:    subject,
		Body:       body,
	}

	err := c.addToQueue(newMessage)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) addToQueue(rawMsg *msg.RawMessage) error {
	pkgMsg, err := rawMsg.ToPackagedMessage(c.SecretKey)
	if err != nil {
		return err
	}

	c.Outbox.Enqueue(*pkgMsg)
	return nil
}

// internal method for sending messages
func (c *Config) sendMessage(pkgMsg *msg.PackagedMessage) error {
	// verify message before sending
	// TODO: do something with messages that have invalid signatures?
	err := pkgMsg.VerifyMessage(c.SecretKey)
	if err != nil {
		return err
	}

	// marshal message into buffer
	msgData, err := json.Marshal(pkgMsg)
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

func (c *Config) SendOneFromQueue() error {
	online := c.Online.getValue()

	// perform additional check
	if !online {
		err := c.checkServerIsOnline()
		if err != nil {
			return err
		}
	}
	// continue if online

	msgToSend, ok := c.Outbox.Dequeue()
	if !ok {
		return errors.New("unable to dequeue message for sending")
	}

	c.sendMessage(&msgToSend)

	return nil
}

// SendAllFromQueue will go through the entire queue and
// send messages until its empty.
func (c *Config) SendAllFromQueue() error {
	online := c.Online.getValue()

	// perform additional check
	if !online {
		err := c.checkServerIsOnline()
		if err != nil {
			return err
		}
	}

	// send until queue is empty
	for !c.Outbox.IsEmpty() {
		msgToSend, ok := c.Outbox.Dequeue()
		if !ok {
			return errors.New("unable to dequeue message for sending")
		}

		err := c.sendMessage(&msgToSend)
		if err != nil {
			return err
		}
	}

	return nil
}
