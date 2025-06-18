package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type config struct {
	YouTubeAPIKey string `env:"YOUTUBE_API_KEY,notEmpty"`
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	fmt.Printf("%+v\n", cfg)
}
