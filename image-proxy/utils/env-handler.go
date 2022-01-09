package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ENV struct {
	CloudName      string
	ApiKey         string
	ApiSecret      string
	ClarifaiApiKey string
	ClarifaiUserID string
	ClarifaiAppID  string
}

func GetENV() ENV {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := ENV{
		CloudName:      os.Getenv("CLOUD_NAME"),
		ApiKey:         os.Getenv("API_KEY"),
		ApiSecret:      os.Getenv("API_SECRET"),
		ClarifaiApiKey: os.Getenv("CLARIFAI_API_KEY"),
		ClarifaiUserID: os.Getenv("CLARIFAI_USER_ID"),
		ClarifaiAppID:  os.Getenv("CLARIFAI_APP_ID"),
	}

	return env
}
