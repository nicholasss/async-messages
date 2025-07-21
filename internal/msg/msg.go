// Package msg implements common message functions and types
package msg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// *** Types ***

// Message defines a raw message that is signed.
// It should not be used directly.
type Message struct {
	To        string `json:"to"`
	From      string `json:"fr"`
	Subject   string `json:"subj"`
	Body      string `json:"body"`
	Signature string `json:"sign"`
}

// *** Functions ***

// Message Functions

// ToString returns a formatted string for the message
func (m *Message) ToString() string {
	format := `To: %s\n
	From: %s\n
	Subject: %s\n
	Body: %s\n
	Signature: %s\n
	`
	return fmt.Sprintf(format, m.To, m.From, m.Subject, m.Body, m.Signature)
}

// NewMessage create a new message object and hashes a valid signature
// The message object should not be altered after this function
func NewMessage(to, from, subject, body string, secretKey []byte) (*Message, error) {
	if to == "" || from == "" || subject == "" || body == "" {
		return nil, fmt.Errorf("empty message property of 'to', 'from', 'subject', or 'body'")
	}

	newMessage := Message{To: to, From: from, Subject: subject, Body: body}
	sig, err := createSignature(&newMessage, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	newMessage.Signature = sig
	return &newMessage, nil
}

// VerifyMessage verifies the signature on the received message
// The message object should not be altered before this function
func VerifyMessage(m *Message, sercretKey []byte) (bool, error) {
	ok, err := verifySignature(m, sercretKey)
	if err != nil {
		return false, fmt.Errorf("failed to verify message: %w", err)
	}

	return ok, nil
}

func prepMessageForSigning(m *Message) []byte {
	messageData := fmt.Sprintf("%s|%s|%s|%s", m.To, m.From, m.Subject, m.Body)
	return []byte(messageData)
}

// create a signature
func createSignature(m *Message, secretKey []byte) (string, error) {
	// prepare message data for signing
	messageData := prepMessageForSigning(m)

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
func verifySignature(m *Message, secretKey []byte) (bool, error) {
	messageData := prepMessageForSigning(m)

	// recalculate hash using message
	h := hmac.New(sha256.New, secretKey)
	_, err := h.Write([]byte(messageData))
	if err != nil {
		return false, fmt.Errorf("failed to write message to hmac: %w", err)
	}
	calulatedSignatureData := h.Sum(nil)

	// decode the hash from the message
	receivedSignatureData, err := hex.DecodeString(m.Signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode recieved signature: %w", err)
	}

	// securely perform comparison of signatures
	return hmac.Equal(calulatedSignatureData, receivedSignatureData), nil
}
