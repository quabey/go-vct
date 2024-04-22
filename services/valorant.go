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
	filter := make(map[string][]common.MatchDetail)
	for _, match := range data.Data {
		if helpers.CheckVCT(match.Tournament) {
			region := helpers.GetRegion(match.Tournament)
			if len(filter[region]) < 3 {
				filter[region] = append(filter[region], match)
			}
		}
	}

	for _, region := range filter {
		fmt.Printf("Message been sent for %s ?: %s", region[0].ID, helpers.CheckIfMessageBeenSent(region[0].ID, "upcoming"))
		if len(region) > 0 && region[0].In != "" && helpers.GetHoursFromNow(region[0].In) < 10 && !helpers.CheckIfMessageBeenSent(region[0].ID, "upcoming") {
			SendUpcomingToServices(region[0], false, true)
			if len(region) >= 2 && helpers.GetOffsetInHours(region[0], region[1]) <= 3 {
				SendUpcomingToServices(region[1], false, false)
				if len(region) >= 3 && helpers.GetOffsetInHours(region[1], region[2]) <= 3 {
					SendUpcomingToServices(region[2], false, false)
				}
			}
		}
		followingMatch = region[1]
		CheckGameStarts(region[0])
	}
}

var followingMatch common.MatchDetail

func CheckGameStarts(currentUpcoming common.MatchDetail) {
	if currentUpcoming.In != "" && !helpers.CheckIfMessageBeenSent(currentUpcoming.ID, "start") {
		fmt.Println("Match is starting!")
		runningMatchId = currentUpcoming.ID
		isFirst := helpers.GetHoursFromNow(followingMatch.In) <= 3
		SendMatchStartToServices(currentUpcoming, isFirst)
		return
	}
	fmt.Printf("No new match has started in %s", helpers.GetRegion(currentUpcoming.Tournament))
}

func CheckAndSendResults() {
	results := FetchData("https://vlr.orlandomm.net/api/v1/results?page=1")
	for index, match := range results.Data {
		if !helpers.CheckIfMessageBeenSent(match.ID, "result") && helpers.CheckVCT(match.Tournament) {
			fmt.Println("New result found!")
			lastResultId = results.Data[0].ID
			SendResultsToservices(match)
			return
		}
		if index > 3 {
			return
		}
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
