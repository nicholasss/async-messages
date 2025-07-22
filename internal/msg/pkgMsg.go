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

// PackagedMessage is a signed and packaged message
// This kind of message is ready to be sent and recieved
type PackagedMessage struct {
	To        UserVessel `json:"to"`
	From      UserVessel `json:"from"`
	Subject   string     `json:"subject"`
	Body      string     `json:"body"`
	Signature string     `json:"signature"`
	Packaged  time.Time  `json:"packagedAt"`
	Recieved  time.Time  `json:"recievedAt"`
}

// UserVessel identifies a persons name and a vessel that they are on
type UserVessel struct {
	Name   string `json:"name"`
	Vessel string `json:"vessel"`
}

// *** Functions ***

// PackagedMessage.String() returns a stringified version of the struct
func (m *PackagedMessage) String() string {
	template := "To: %s\nFrom: %s\nSubject: %s\nBody: %s\nSignature: %s\n"
	return fmt.Sprintf(template, m.To.String(), m.From.String(), m.Subject, m.Body, m.Signature)
}

// UserVessel.String() returns a stringified version of the struct
func (uv *UserVessel) String() string {
	return fmt.Sprintf("%s@%s", uv.Name, uv.Vessel)
}

// VerifyMessage verifies the signature on the received message
// The message object should not be altered before this function
func (m *PackagedMessage) VerifyMessage(secretKey []byte) error {
	messageData := m.messageDataForSigning()

	// recalculate hash using message
	h := hmac.New(sha256.New, secretKey)
	_, err := h.Write([]byte(messageData))
	if err != nil {
		return fmt.Errorf("failed to create hash for signature check: %w", err)
	}
	calulatedSignatureData := h.Sum(nil)

	// decode the hash from the message
	receivedSignatureData, err := hex.DecodeString(m.Signature)
	if err != nil {
		return fmt.Errorf("failed to decode the messages signature: %w", err)
	}

	// securely perform comparison of signatures
	if hmac.Equal(calulatedSignatureData, receivedSignatureData) {
		return nil
	} else {
		return errors.New("signature of message is invalid")
	}
}

// messageDataForSigning is an internal function to prepare data for creating a signature
func (pkgMsg *PackagedMessage) messageDataForSigning() []byte {
	messageData := fmt.Sprintf("%s|%s|%s|%s",
		pkgMsg.To.String(), pkgMsg.From.String(), pkgMsg.Subject, pkgMsg.Body)
	return []byte(messageData)
}
