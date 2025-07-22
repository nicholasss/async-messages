// Package msg implements common message functions and types
package msg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// *** Types ***

// RawMessage is a raw message that needs to be processed further
// all fields are used to calculate the signature
type RawMessage struct {
	ToName     string
	ToVessel   string
	FromName   string
	FromVessel string
	Subject    string
	Body       string
}

// *** Functions ***

// ToPackageMessage takes a raw message and performs operations needed to package it into a packaged message
func (rawMsg RawMessage) ToPackagedMessage(secretKey []byte) (*PackagedMessage, error) {
	// TODO: perform validation of fields and not just 'not empty'
	if rawMsg.ToName == "" || rawMsg.ToVessel == "" {
		return nil, errors.New("one of the 'to' fields are empty")
	}
	toInfo := UserVessel{Name: rawMsg.ToName, Vessel: rawMsg.ToVessel}

	if rawMsg.FromName == "" || rawMsg.FromVessel == "" {
		return nil, errors.New("one of the 'from' fields are empty")
	}
	fromInfo := UserVessel{Name: rawMsg.FromName, Vessel: rawMsg.FromVessel}

	if rawMsg.Subject == "" {
		return nil, errors.New("subject field cannot be empty")
	}

	if rawMsg.Body == "" {
		return nil, errors.New("body field cannot be empty")
	}

	signature, err := rawMsg.createSignature(secretKey)
	if err != nil {
		return nil, err
	}

	packagedMsg := PackagedMessage{
		To:        toInfo,
		From:      fromInfo,
		Subject:   rawMsg.Subject,
		Body:      rawMsg.Body,
		Signature: signature,
		Packaged:  time.Now().UTC(),
	}

	return &packagedMsg, nil
}

// messageDataForSinging returns a stringified representation of the raw message
// this is needed for creating the signature
func (rawMsg *RawMessage) messageDataForSigning() []byte {
	messageData := fmt.Sprintf("%s@%s|%s@%s|%s|%s",
		rawMsg.ToName, rawMsg.ToVessel, rawMsg.FromName, rawMsg.FromVessel, rawMsg.Subject, rawMsg.Body)
	return []byte(messageData)
}

// createSignature returns the signature that validates the message itself
func (rawMsg *RawMessage) createSignature(secretKey []byte) (string, error) {
	messageData := rawMsg.messageDataForSigning()

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
