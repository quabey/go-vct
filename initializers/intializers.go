package initializers

import (
	"bey/go-vct/common"
	"bey/go-vct/database"
	"log"

	"github.com/joho/godotenv"
)

func GetEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetSentMessages() {
	common.Messages, _ = database.GetSentMessages()
}
