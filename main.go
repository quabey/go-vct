package main

import (
	"bey/go-vct/common"
	"bey/go-vct/database"
	"bey/go-vct/initializers"
	"bey/go-vct/services"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	initializers.GetEnvVariables()
	common.LoadEnvVariables()
	database.InitDB(common.DbPath)
	initializers.GetSentMessages()
}

func main() {
	services.GetUpcoming()

	return
	fmt.Println("Starting...")

	for _, message := range common.Messages {
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
