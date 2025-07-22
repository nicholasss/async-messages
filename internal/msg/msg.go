// Package msg implements common message functions and types
package msg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

// *** Types ***

// Message defines a raw message that is signed.
// It should not be used directly.
type Message struct {
	To        string `json:"to"`
	From      string `json:"from"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Signature string `json:"signature"`
}

// *** Functions ***

// Message Functions

// NewMessage create a new message object and hashes a valid signature
// The message object should not be altered after this function
func NewMessage(to, from, subject, body string, secretKey []byte) (*Message, error) {
	if to == "" || from == "" || subject == "" || body == "" {
		return nil, fmt.Errorf("empty message property of 'to', 'from', 'subject', or 'body'")
	}

	newMessage := Message{To: to, From: from, Subject: subject, Body: body}
	sig, err := newMessage.createSignature(secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	newMessage.Signature = sig
	return &newMessage, nil
}

// ToString returns a formatted string for the message
func (m *Message) ToString() string {
	template := "To: %s\nFrom: %s\nSubject: %s\nBody: %s\nSignature: %s\n"
	return fmt.Sprintf(template, m.To, m.From, m.Subject, m.Body, m.Signature)
}

// VerifyMessage verifies the signature on the received message
// The message object should not be altered before this function
func (m *Message) VerifyMessage(sercretKey []byte) error {
	err := m.verifySignature(sercretKey)
	if err != nil {
		return fmt.Errorf("failed to verify message: %w", err)
	}

	return nil
}

func (m *Message) messageDataForSigning() []byte {
	messageData := fmt.Sprintf("%s|%s|%s|%s", m.To, m.From, m.Subject, m.Body)
	return []byte(messageData)
}

// create a signature
func (m *Message) createSignature(secretKey []byte) (string, error) {
	// prepare message data for signing
	messageData := m.messageDataForSigning()

	// create hash with sha256 and secret key
	h := hmac.New(sha256.New, secretKey)
	_, err := h.Write([]byte(messageData))
	if err != nil {
		return "", fmt.Errorf("failed to write message to hmac: %w", err)
	}

	// get the final hash and encode to hex string
	calulatedSignature := hex.EncodeToString(h.Sum(nil))
	return calulatedSignature, nil
}

// verify the signature within a message
func (m *Message) verifySignature(secretKey []byte) error {
	messageData := m.messageDataForSigning()

	// recalculate hash using message
	h := hmac.New(sha256.New, secretKey)
	_, err := h.Write([]byte(messageData))
	if err != nil {
		return fmt.Errorf("failed to write message to hmac: %w", err)
	}
	calulatedSignatureData := h.Sum(nil)

	// decode the hash from the message
	receivedSignatureData, err := hex.DecodeString(m.Signature)
	if err != nil {
		return fmt.Errorf("failed to decode recieved signature: %w", err)
	}

	// securely perform comparison of signatures
	if hmac.Equal(calulatedSignatureData, receivedSignatureData) {
		return nil
	} else {
		return errors.New("signatures do not match")
	}
}
