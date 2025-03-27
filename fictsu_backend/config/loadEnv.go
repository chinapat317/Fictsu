package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	GoogleClientID string
	DbHost         string
	DbPort         string
	DbUser         string
	DbPw           string
	DbName         string
	CoverPath      string
	BucketName     string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbUser = os.Getenv("DB_USER")
	DbPw = os.Getenv("DB_PW")
	DbName = os.Getenv("DB_NAME")
	CoverPath = os.Getenv("COVER_PATH")
	BucketName = os.Getenv("BUCKET_NAME")
}
