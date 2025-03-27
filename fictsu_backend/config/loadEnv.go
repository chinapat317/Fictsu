package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	GoogleClientID string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")

	if GoogleClientID == "" {
		log.Fatal("Missing one or more required environment variables")
	}
}
