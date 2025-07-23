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

// Config holds
type Config struct {
	SecretKey []byte
	Client    http.Client
	Queue     *msg.PackagedQueue
	Name      string
	Vessel    string
	Server    string
	Online    *safeBool
}

type safeBool struct {
	mux  sync.RWMutex
	bool bool
}

type NewMessage struct {
	ToName   string
	ToVessel string
	Subject  string
	Body     string
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

	queue := msg.NewQueue()

	safeBool := &safeBool{
		mux:  sync.RWMutex{},
		bool: false,
	}

	return &Config{
		SecretKey: secretKey,
		Client:    *http.DefaultClient,
		Queue:     queue,
		Name:      name,
		Vessel:    vessel,
		Server:    "http://localhost:8080",
		Online:    safeBool,
	}, nil
}

// *** Functions ***

// StartClient will begin with checking if the server is online or not
// then depending on the state, will start sending messages or
// will begin queueing and checking for connectivity
func (c *Config) StartClient() error {
	err := c.checkServerIsOnline()
	if err != nil {
		return err
	}

	return nil
}

func (s *safeBool) getValue() bool {
	s.mux.RLock()
	val := s.bool
	s.mux.RLock()
	return val
}

func (s *safeBool) setValue(val bool) {
	s.mux.Lock()
	s.bool = val
	s.mux.Unlock()
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

func (c *Config) AddToQueue(rawMsg *msg.RawMessage) error {
	pkgMsg, err := rawMsg.ToPackagedMessage(c.SecretKey)
	if err != nil {
		return err
	}

	c.Queue.Enqueue(*pkgMsg)
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

	msgToSend, ok := c.Queue.Dequeue()
	if !ok {
		return errors.New("unable to dequeue message for sending")
	}

	c.sendMessage(&msgToSend)

	return nil
}

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
	for !c.Queue.IsEmpty() {
		msgToSend, ok := c.Queue.Dequeue()
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
