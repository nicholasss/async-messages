package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nicholasss/async-messages/internal/msg"
)

// *** Server Config ***

type serverConfig struct {
	secretKey []byte
}

func loadConfig() (*serverConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("unable to load './.env'. error: %w", err)
	}

	rawHMAC := os.Getenv("HMAC_SECRET")
	if rawHMAC == "" {
		return nil, fmt.Errorf("unable to find 'HMAC_SECRET' in './.env'")
	}

	cfg := serverConfig{
		secretKey: []byte(rawHMAC),
	}

	return &cfg, nil
}

// *** Main ***

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Unable to load config due to: %q", err)
	}

	r := gin.Default()
	r.POST("/echo", cfg.echo)

	r.Run()
}

func (cfg *serverConfig) echo(c *gin.Context) {
	clientMsg := &msg.Message{}
	c.Bind(clientMsg)

	err := msg.VerifyMessage(clientMsg, cfg.secretKey)
	if err != nil {
		log.Printf("unable to verify message due to: %q", err)
		return
	}

	c.Status(200)
}
