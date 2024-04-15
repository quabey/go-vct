package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MatchData struct {
	Status    string `json:"status"`
	Size      int    `json:"size"`
	Data      []MatchDetail `json:"data"`
}

type MatchDetail struct {
	ID         string `json:"id"`
	Teams      []Team `json:"teams"`
	Status     string `json:"status"`
	Event      string `json:"event"`
	Tournament string `json:"tournament"`
	Img        string `json:"img"`
	In         string `json:"in,omitempty"`
}

type Team struct {
	Name    string `json:"name"`
	Tag     string `json:"tag,omitempty"`
	Logo    string `json:"logo,omitempty"`
	Score   string `json:"score,omitempty"`
	Country string `json:"country,omitempty"`
	Won     bool   `json:"won,omitempty"`
}

var (
	upcomingMatchWebhookURL = ""
	runningMatchId string
	lastResultId string
	lastUpcomingMatchId string

	twitchLinks = map[string]string{
		"Americas": "https://www.twitch.tv/valorant_americas",
		"China": "https://www.twitch.tv/valorant_china",
		"EMEA": "https://www.twitch.tv/valorant_emea",
		"Pacific": "https://www.twitch.tv/valorant_pacific",
	}

	youtubeLinks = map[string]string{
		"Americas": "https://www.youtube.com/@valorant_americas/live",
		"China": "https://www.youtube.com/@VALORANTEsportsCN/live",
		"EMEA": "https://www.youtube.com/@valorant_emea/live",
		"Pacific": "https://www.youtube.com/@valorant_pacific/live",
	}

	roles = map[string]string{
		"Americas": "1227214059498115072",
		"China": "1227214116834382009",
		"EMEA": "1227214030616264734",
		"Pacific": "1227213846268346440",
	}

	watchParties = map[string]map[string]string{
	"Pacific": {
		"Sliggy": "https://www.twitch.tv/sliggytv",
		"FNS": "https://www.twitch.tv/gofns",
		"Sean Gares": "https://www.twitch.tv/sgares",
		"Thinking Mans Valorant": "https://www.twitch.tv/thinkingmansvalorant",
		"Tarik": "https://www.twitch.tv/tarik",
		"Kyedae": "https://www.twitch.tv/kyedae",
	},
	"EMEA": {
		"FNS": "https://www.twitch.tv/gofns",
		"Sliggy": "https://www.twitch.tv/sliggytv",
		"Sgares": "https://www.twitch.tv/sgares",
		"ThinkingMansValorant": "https://www.twitch.tv/thinkingmansvalorant",
		"tarik": "https://www.twitch.tv/tarik",
		"kyedae": "https://www.twitch.tv/kyedae",
	},
	"China": {
		"Ryancentral": "https://www.twitch.tv/ryancentral",
		"Yinsu": "https://www.twitch.tv/yinsu",
	},
	"Americas": {
		"FNS": "https://www.twitch.tv/gofns",
		"Sliggy": "https://www.twitch.tv/sliggytv",
		"Sgares": "https://www.twitch.tv/sgares",
		"ThinkingMansValorant": "https://www.twitch.tv/thinkingmansvalorant",
		"tarik": "https://www.twitch.tv/tarik",
		"kyedae": "https://www.twitch.tv/kyedae",
	},
}

)

func main() {
	fmt.Println("Starting...")
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
	updateUpcomingMatches(data)
	checkGameStart(data)
}

func updateUpcomingMatches(currentUpcoming MatchData) {
	if currentUpcoming.Data[0].ID != lastUpcomingMatchId {
		fmt.Println("New upcoming match found!")
		lastUpcomingMatchId = currentUpcoming.Data[0].ID
		sendUpcomingToDiscord(currentUpcoming)
		return
	}
	fmt.Println("No new upcoming matches found.")
}

func checkGameStart(currentUpcoming MatchData) {
	if currentUpcoming.Data[0].ID != runningMatchId && currentUpcoming.Data[0].In == "" {
		fmt.Println("Match is starting!")
		runningMatchId = currentUpcoming.Data[0].ID
		sendMatchStartToDiscord(currentUpcoming.Data[0])
		return
	}
	fmt.Println("No new match has started.")
}

func checkAndSendResults() {
	results := fetchData("https://vlr.orlandomm.net/api/v1/results?page=1")
	if results.Data[0].ID != lastResultId {
		fmt.Println("New result found!")
		lastResultId = results.Data[0].ID
		sendResultsToDiscord(results)
		return
	}
	fmt.Println("No new results found.")
}

func fetchData(url string) MatchData {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return MatchData{}
	}
	defer resp.Body.Close()

	var data MatchData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return MatchData{}
	}

	return data
}

func sendToDiscord(url string, payload []byte) {
    req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
    if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }
    req.Header.Set("Content-Type", "application/json")

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending message:", err)
        return
    }
    defer resp.Body.Close()

    // Check the response
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        responseBytes, _ := io.ReadAll(resp.Body)
        fmt.Printf("Failed to send message, status code: %d, response: %s\n", resp.StatusCode, string(responseBytes))
    }
}



func sendUpcomingToDiscord(matches MatchData) {
    if len(matches.Data) > 3 {
        matches.Data = matches.Data[:3] // Only take the first 2 matches
    }

    embeds := make([]map[string]interface{}, len(matches.Data))
    for i, match := range matches.Data {
		region := getRegion(match.Tournament)
		title := "Live Match"
		if match.In != "" {
		timestamp, err := parseDurationFromNow(match.In)
		if err != nil {
            fmt.Println("Error parsing duration:", err)
            continue
        }
		title = fmt.Sprintf("Upcoming Match at <t:%d:t>", timestamp)
		} 
		title = fmt.Sprintf("%s: %s", title, fmt.Sprintf("**%s** vs **%s**", match.Teams[0].Name, match.Teams[1].Name))
        embeds[i] = map[string]interface{}{
            "type": "rich",
            "title": title,
            "description": fmt.Sprintf("%s at %s - %s", match.Tournament, match.Event, match.Status),
            "color": 0x00FFFF,
			"footer": map[string]interface{}{
        		"text": "Made with ❤️ by bey",
    		},
            "fields": []map[string]interface{}{
                {
                    "name": "Riot Streams",
                    "value": fmt.Sprintf("[Twitch](%s)\n[YouTube](%s)", 
                        getTwitchLink(region), getYoutubeLink(region)),
					"inline": true,
                },
				{
					"name": "Watch Parties",
					"value": buildWatchPartyLinks(getWatchParties(region)),
					"inline": true,
				},
            },
        }
    }

    message := map[string]interface{}{
        "content": "# Here are the upcoming matches:",
        "embeds": embeds,
    }

    messageBytes, err := json.Marshal(message)
    if err != nil {
        fmt.Println("Error marshalling message:", err)
        return
    }

    sendToDiscord(upcomingMatchWebhookURL, messageBytes)
}

func buildWatchPartyLinks(parties map[string]string) string {
    if len(parties) == 0 {
        return "No watch parties available"
    }
    var links []string
    for name, url := range parties {
        links = append(links, fmt.Sprintf("[%s](%s)", name, url))
    }
    return strings.Join(links, "\n")
}

func sendMatchStartToDiscord(match MatchDetail) {
	region := getRegion(match.Tournament)
	title := fmt.Sprintf("Match Start: **%s** vs **%s**", match.Teams[0].Name, match.Teams[1].Name)
	embed := map[string]interface{}{
		"type": "rich",
		"title": title,
		"description": fmt.Sprintf("%s at %s - %s", match.Tournament, match.Event, match.Status),
		"color": 0x00FFFF,
		"footer": map[string]interface{}{
			"text": "Made with ❤️ by bey",
		},
		"fields": []map[string]interface{}{
			{
				"name": "Riot Streams",
				"value": fmt.Sprintf("[Twitch](%s)\n[YouTube](%s)", 
					getTwitchLink(region), getYoutubeLink(region)),
				"inline": true,
			},
			{
				"name": "Watch Parties",
				"value": buildWatchPartyLinks(getWatchParties(region)),
				"inline": true,
			},
		},
	}

	message := map[string]interface{}{
		"content": fmt.Sprintf("<@&%s>", roles[region]),
		"embeds": []map[string]interface{}{embed},
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return
	}

	sendToDiscord(upcomingMatchWebhookURL, messageBytes)
}

func sendResultsToDiscord(results MatchData) {
    if len(results.Data) > 0 {
    	results.Data = results.Data[:1] // Only take the latest result
    }
	embeds := make([]map[string]interface{}, len(results.Data))
	for i, result := range results.Data {
		score := fmt.Sprintf("**%s** - **%s**", result.Teams[0].Score, result.Teams[1].Score)
		winner := result.Teams[0].Name
		if result.Teams[1].Won {
			score = fmt.Sprintf("**%s** - **%s**", result.Teams[1].Score, result.Teams[0].Score)
			winner = result.Teams[1].Name
		}
		title := fmt.Sprintf("Match Result: **%s** vs **%s**", result.Teams[0].Name, result.Teams[1].Name)
		embeds[i] = map[string]interface{}{
			"type": "rich",
			"title": title,
			"description": fmt.Sprintf("||%s|| Wins: ||%s|| \n %s - %s", winner, score, result.Tournament, result.Event),
			"color": 0x00FFFF,
			"footer": map[string]interface{}{
        		"text": "Made with ❤️ by bey",
    		},
			"fields": []map[string]interface{}{},
		}
	}

	message := map[string]interface{}{
		"content": "",
		"embeds": embeds,
	}
	messageBytes, err := json.Marshal(message)
    if err != nil {
        fmt.Println("Error marshalling message:", err)
        return
    }
	
	sendToDiscord(upcomingMatchWebhookURL, messageBytes)
}



func getTwitchLink(region string) string {
	return twitchLinks[region]
}

func getYoutubeLink(region string) string {
	return youtubeLinks[region]
}

func getWatchParties(region string) map[string]string {
	return watchParties[region]
}

func getWatchPartyLink(region, streamer string) string {
	return watchParties[region][streamer]
}

func parseDurationFromNow(durationStr string) (int64, error) {
    // Split the string by spaces
    parts := strings.Split(durationStr, " ")
    if len(parts) != 2 {
        return 0, fmt.Errorf("invalid format, expected 'HHh MMm'")
    }

    // Parse hours
    hours, err := strconv.Atoi(strings.TrimSuffix(parts[0], "h"))
    if err != nil {
        return 0, fmt.Errorf("invalid hour format: %v", err)
    }

    // Parse minutes
    minutes, err := strconv.Atoi(strings.TrimSuffix(parts[1], "m"))
    if err != nil {
        return 0, fmt.Errorf("invalid minute format: %v", err)
    }

    // Get the current time
    now := time.Now()

    // Calculate the future time by adding the parsed duration and ensure that its on the hour
    futureTime := now.Add(time.Hour*time.Duration(hours) + time.Minute*time.Duration(minutes))
	futureTime = futureTime.Truncate(time.Hour)

    // Return the Unix timestamp of the future time
    return futureTime.Unix(), nil
}

func getRegion(tournament string) (region string) {
	if strings.Contains(tournament, "EMEA") {
		return "EMEA"
	} else if strings.Contains(tournament, "China") {
		return "China"
	} else if strings.Contains(tournament, "Americas") {
		return "Americas"
	} else if strings.Contains(tournament, "Pacific") {
		return "Pacific"
	}
	return ""
}