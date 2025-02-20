package config

import (
	"os"
	"log"
	"github.com/joho/godotenv"
)

var (
	ClientID			string
	ClientSecret		string
	ClientCallbackURL 	string

	OpenAIKey 			string
	OpenAIOrgID			string
	OpenAIProjID 		string

	SessionKey			string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ClientID = os.Getenv("CLIENT_ID")
	ClientSecret = os.Getenv("CLIENT_SECRET")
	ClientCallbackURL = os.Getenv("CLIENT_CALLBACK_URL")

	OpenAIKey = os.Getenv("OPENAI_KEY")
	OpenAIOrgID = os.Getenv("OPENAI_ORG_ID")
	OpenAIProjID = os.Getenv("OPENAI_PROJ_ID")

	SessionKey = os.Getenv("SESSION_KEY")

	// Fail fast if any required environment variable is missing
	if OpenAIKey == "" || OpenAIOrgID == "" || OpenAIProjID == "" || ClientID == "" || ClientSecret == "" || ClientCallbackURL == "" || SessionKey == "" {
		log.Fatal("Missing one or more required environment variables")
	}
}
