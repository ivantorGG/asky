package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Domain      string
	DatabaseURL string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Domain:      os.Getenv("DOMAIN"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	return cfg
}
