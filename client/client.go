package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nicholasss/async-messages/internal/message"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("could not load './.env' due to: %q", err)
		return
	}
	secretKey := []byte(os.Getenv("HMAC_SECRET"))

	c := http.DefaultClient
	res, err := c.Get("http://localhost:8080/msg")
	if err != nil {
		log.Printf("could not make request due to: %q", err)
		return
	}

	msg := &message.Message{}
	err = json.NewDecoder(res.Body).Decode(msg)
	if err != nil {
		log.Printf("could not decode the request due to: %q", err)
		return
	}
	res.Body.Close()

	msgOk, err := message.VerifyMessage(msg, secretKey)
	if err != nil {
		log.Printf("could not verify message due to: %q", err)
		return
	}

	if !msgOk {
		log.Printf("the message signature is not correct.")
		return
	}

	log.Printf("To: %q", msg.To)
	log.Printf("From: %q", msg.From)
	log.Printf("Subject: %q", msg.Subject)
	log.Printf("Body: %q", msg.Body)
}
