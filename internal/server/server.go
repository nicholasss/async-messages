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
	secretKey []byte
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

	cfg := Config{
		secretKey: []byte(rawHMAC),
	}

	return &cfg, nil
}

func (cfg *Config) SetupGinEngine() (*gin.Engine, error) {
	r := gin.Default()

	r.GET("/health", cfg.health)
	r.POST("/send-message", cfg.sendMessage)

	return r, nil
}

func (cfg *Config) health(c *gin.Context) {
	healthRes := HealthCheck{Health: "OK"}
	c.JSON(200, healthRes)
}

func (cfg *Config) sendMessage(c *gin.Context) {
	requestMsg := &msg.Message{}
	c.Bind(requestMsg)

	err := requestMsg.VerifyMessage(cfg.secretKey)
	if err != nil {
		log.Printf("unable to verify message due to: %q", err)
		c.Status(400) // bad request
		return
	}

	c.Status(200) // ok
}
