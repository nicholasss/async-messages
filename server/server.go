package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nicholasss/async-messages/internal/message"
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
	r.GET("/msg", cfg.msg)

	r.Run()
}

func (cfg *serverConfig) msg(c *gin.Context) {
	msg, err := message.NewMessage("You", "me", "hello", "testing, testing", cfg.secretKey)
	if err != nil {
		log.Printf("could not create message due to %q", err)
		return
	}

	c.JSON(200, msg)
}
