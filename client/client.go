package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nicholasss/async-messages/internal/msg"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("could not load './.env' due to: %q", err)
		return
	}
	secretKey := []byte(os.Getenv("HMAC_SECRET"))

	c := http.DefaultClient

	testMsg, err := msg.NewMessage("Me", "Me", "Echoing message", "Hearing me ok?", secretKey)
	if err != nil {
		log.Printf("could not create message due to: %q", err)
		return
	}

	bodyData, err := json.Marshal(testMsg)
	if err != nil {
		log.Printf("could not marshal data due to: %q", err)
		return
	}
	bodyReader := bytes.NewBuffer(bodyData)

	res, err := c.Post("http://localhost:8080/echo", "application/json", bodyReader)
	if err != nil {
		log.Printf("could not make request due to: %q", err)
		return
	}

	echoMsg := &msg.Message{}
	err = json.NewDecoder(res.Body).Decode(echoMsg)
	if err != nil {
		log.Printf("could not decode the request due to: %q", err)
		return
	}
	res.Body.Close()

	msgOk, err := msg.VerifyMessage(echoMsg, secretKey)
	if err != nil {
		log.Printf("could not verify message due to: %q", err)
		return
	}

	if !msgOk {
		log.Printf("the message signature is not correct.")
		return
	}

	log.Printf("To: %q", echoMsg.To)
	log.Printf("From: %q", echoMsg.From)
	log.Printf("Subject: %q", echoMsg.Subject)
	log.Printf("Body: %q", echoMsg.Body)
}
