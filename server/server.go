package main

import (
	"log"

	"github.com/nicholasss/async-messages/internal/server"
)

// *** Main ***

func main() {
	cfg, err := server.LoadConfig()
	if err != nil {
		log.Printf("could not load server config due to: %q", err)
		return
	}

	r, err := cfg.SetupGinEngine()
	if err != nil {
		log.Printf("could not setup gin engine(router) due to: %q", err)
		return
	}

	log.Fatalf("gin crashed due to: %q", r.Run())
}
