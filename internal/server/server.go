// Package server ...
package server

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nicholasss/async-messages/internal/msg"
)

type HealthCheck struct {
	Health string `json:"health"`
}

// Config holds all the configuration data
type Config struct {
	SecretKey []byte
	Queue     *msg.Queue
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("unable to load './.env'. error: %w", err)
	}

	rawHMAC := os.Getenv("HMAC_SECRET")
	if rawHMAC == "" {
		return nil, fmt.Errorf("unable to find 'HMAC_SECRET' in './.env'")
	}

	queue := msg.NewQueue()

	cfg := Config{
		SecretKey: []byte(rawHMAC),
		Queue:     queue,
	}

	return &cfg, nil
}

func (cfg *Config) SetupGinEngine() (*gin.Engine, error) {
	r := gin.Default()

	// allow clients to check health of server
	r.GET("/health", cfg.health)

	// allow clients to send messages
	r.POST("/send-message", cfg.sendMessage)

	// allow clients to check for messages
	// r.GET("/check-messages", cfg.checkMessages)

	return r, nil
}

func (cfg *Config) health(c *gin.Context) {
	healthRes := HealthCheck{Health: "OK"}
	c.JSON(200, healthRes)
}

func (cfg *Config) sendMessage(c *gin.Context) {
	requestMsg := &msg.Message{}
	c.Bind(requestMsg)

	err := requestMsg.VerifyMessage(cfg.SecretKey)
	if err != nil {
		log.Printf("unable to verify message due to: %q", err)
		c.Status(400) // bad request
		return
	}

	// add to server queue
	cfg.Queue.Enqueue(*requestMsg)
	log.Printf("Server message queue >\n%s", cfg.Queue.DumpToString())

	c.Status(200) // ok
}
