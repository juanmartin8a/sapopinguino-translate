//go:build !prod && !dev

package config

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error while loading .env file: %v", err)
	}
}
