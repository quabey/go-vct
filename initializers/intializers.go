package initializers

import (
	"bey/go-vct/common"
	"bey/go-vct/database"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func GetEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("this is only bad if you are not using docker")
	}
}

func GetSentMessages() {
	common.Messages, _ = database.GetSentMessages()

	for _, message := range common.Messages {
		fmt.Println(message)
	}
}
