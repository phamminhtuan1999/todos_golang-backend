package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func LoadConfig() (*Config, error) {
	var err error = godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	config := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}

	if config.DatabaseURL == "" || config.Port == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}
	return config, nil
}
