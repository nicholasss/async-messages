// Package msg implements common message functions and types
package msg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
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

// MissingFieldError is returned when there is a missing field
type MissingFieldError struct {
	Field string
}

func (err *MissingFieldError) Error() string {
	return fmt.Sprintf("field '%s' is missing from the raw message", err.Field)
}

// *** Functions ***

// ToPackagedMessage takes a raw message and performs operations needed to package it into a packaged message
func (rawMsg RawMessage) ToPackagedMessage(secretKey []byte) (*PackagedMessage, error) {
	// checking to fields
	if rawMsg.ToName == "" {
		return nil, &MissingFieldError{Field: "ToName"}
	}
	if rawMsg.ToVessel == "" {
		return nil, &MissingFieldError{Field: "ToVessel"}
	}
	toInfo := UserVessel{
		Name:   strings.ToLower(rawMsg.ToName),
		Vessel: strings.ToLower(rawMsg.ToVessel),
	}

	// checking from fields
	if rawMsg.FromName == "" {
		return nil, &MissingFieldError{Field: "FromName"}
	}
	if rawMsg.FromVessel == "" {
		return nil, &MissingFieldError{Field: "FromVessel"}
	}
	fromInfo := UserVessel{
		Name:   strings.ToLower(rawMsg.FromName),
		Vessel: strings.ToLower(rawMsg.FromVessel),
	}

	// checking subject
	if rawMsg.Subject == "" {
		return nil, &MissingFieldError{Field: "Subject"}
	}

	if rawMsg.Body == "" {
		return nil, &MissingFieldError{Field: "Body"}
	}
	// body does not get changed, as it could affect the message

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
