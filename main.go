package main

import (
	"bey/go-vct/common"
	"bey/go-vct/database"
	"bey/go-vct/initializers"
	"bey/go-vct/services"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	initializers.GetEnvVariables()
	common.LoadEnvVariables()
	database.InitDB(common.DbPath)
}

func main() {
	services.GetUpcoming()

	return
	fmt.Println("Starting...")

	messages, err := database.GetSentMessages()
	if err != nil {
		log.Fatal(err)
	}

	for _, message := range messages {
		fmt.Println(message)
	}

	interval := 1 * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	services.GetUpcoming()
	services.CheckAndSendResults()
	for range ticker.C {
		services.GetUpcoming()
		services.CheckAndSendResults()
	}
}
