package main

import (
	"bey/go-vct/common"
	"bey/go-vct/database"
	"bey/go-vct/discord"
	"bey/go-vct/helpers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	runningMatchId      string
	lastResultId        string
	lastUpcomingMatchId string
)

func main() {
	fmt.Println("Starting...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// database.InitDB("./sqlite/messages.db")
	// webhookURL = os.Getenv("WEBHOOK_URL")
	// dbPath = os.Getenv("SQLITE_DB")

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

	getUpcoming()
	checkAndSendResults()
	for range ticker.C {
		getUpcoming()
		checkAndSendResults()
	}
}

func getUpcoming() {
	fmt.Println("Fetching upcoming matches...")
	data := fetchData("https://vlr.orlandomm.net/api/v1/matches")
	filter := []common.MatchDetail{}
	for _, match := range data.Data {
		if helpers.CheckVCT(match.Tournament) {
			filter = append(filter, match)
		}
	}
	data.Data = filter
	updateUpcomingMatches(data)
	checkGameStart(data)
}

func updateUpcomingMatches(currentUpcoming common.MatchData) {
	if currentUpcoming.Data[0].ID != lastUpcomingMatchId {
		fmt.Println("New upcoming match found!")
		lastUpcomingMatchId = currentUpcoming.Data[0].ID
		discord.SendUpcomingToDiscord(currentUpcoming)
		return
	}
	fmt.Println("No new upcoming matches found.")
}

func checkGameStart(currentUpcoming common.MatchData) {
	if currentUpcoming.Data[0].ID != runningMatchId && currentUpcoming.Data[0].In == "" {
		fmt.Println("Match is starting!")
		runningMatchId = currentUpcoming.Data[0].ID
		discord.SendMatchStartToDiscord(currentUpcoming.Data[0])
		return
	}
	fmt.Println("No new match has started.")
}

func checkAndSendResults() {
	results := fetchData("https://vlr.orlandomm.net/api/v1/results?page=1")
	if results.Data[0].ID != lastResultId && helpers.CheckVCT(results.Data[0].Tournament) {
		fmt.Println("New result found!")
		lastResultId = results.Data[0].ID
		discord.SendResultsToDiscord(results)
		return
	}
	fmt.Println("No new results found.")
}

func fetchData(url string) common.MatchData {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return common.MatchData{}
	}
	defer resp.Body.Close()

	var data common.MatchData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return common.MatchData{}
	}

	return data
}
