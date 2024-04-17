package services

import (
	"bey/go-vct/common"
	"bey/go-vct/helpers"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	runningMatchId      string
	lastResultId        string
	lastUpcomingMatchId string
)

func GetUpcoming() {
	fmt.Println("Fetching upcoming matches...")
	data := FetchData("https://vlr.orlandomm.net/api/v1/matches")
	filter := []common.MatchDetail{}
	for _, match := range data.Data {
		if helpers.CheckVCT(match.Tournament) {
			filter = append(filter, match)
		}
	}
	data.Data = filter
	UpdateUpcomingMatches(data)
	CheckGameStarts(data)
}

func UpdateUpcomingMatches(currentUpcoming common.MatchData) {
	if currentUpcoming.Data[0].ID != lastUpcomingMatchId {
		fmt.Println("New upcoming match found!")
		lastUpcomingMatchId = currentUpcoming.Data[0].ID
		SendUpcomingToservices(currentUpcoming)
		return
	}
	fmt.Println("No new upcoming matches found.")
}

func CheckGameStarts(currentUpcoming common.MatchData) {
	if currentUpcoming.Data[0].ID != runningMatchId && currentUpcoming.Data[0].In == "" {
		fmt.Println("Match is starting!")
		runningMatchId = currentUpcoming.Data[0].ID
		SendMatchStartToservices(currentUpcoming.Data[0])
		return
	}
	fmt.Println("No new match has started.")
}

func CheckAndSendResults() {
	results := FetchData("https://vlr.orlandomm.net/api/v1/results?page=1")
	if results.Data[0].ID != lastResultId && helpers.CheckVCT(results.Data[0].Tournament) {
		fmt.Println("New result found!")
		lastResultId = results.Data[0].ID
		SendResultsToservices(results)
		return
	}
	fmt.Println("No new results found.")
}

func FetchData(url string) common.MatchData {
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
