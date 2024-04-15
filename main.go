package main

import (
	"bey/go-vct/helpers"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type MatchData struct {
	Status string        `json:"status"`
	Size   int           `json:"size"`
	Data   []MatchDetail `json:"data"`
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

type Message struct {
	Id               int
	MessageId        string
	MatchId          string
	AnnouncementSent bool
	StartingSent     bool
	ResultSent       bool
}

var (
	webhookURL          string
	runningMatchId      string
	lastResultId        string
	lastUpcomingMatchId string
	dbPath              string

	roles = map[string]string{
	"Americas": "1227214059498115072",
	"China":    "1227214116834382009",
	"EMEA":     "1227214030616264734",
	"Pacific":  "1227213846268346440",
	}

)

func main() {
	fmt.Println("Starting...")
	godotenv.Load()
	webhookURL = os.Getenv("WEBHOOK_URL")
	dbPath = os.Getenv("SQLITE_DB")

	messages, err := getSentMessages()
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
	filter := []MatchDetail{}
	for _, match := range data.Data {
		if helpers.CheckVCT(match.Tournament) {
			filter = append(filter, match)
		}
	}
	data.Data = filter
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
	if results.Data[0].ID != lastResultId && helpers.CheckVCT(results.Data[0].Tournament) {
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

func sendToDiscord(url string, payload []byte) int {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return 5
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return 4
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return 2
	}
	var responseData map[string]interface{}
	json.Unmarshal(responseBytes, &responseData)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		fmt.Printf("Failed to send message, status code: %d, response: %s\n", resp.StatusCode, string(responseBytes))
	}

	if id, exists := responseData["id"].(string); exists {
		intId, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Error converting ID to int:", err)
			return 3
		}
		log.Printf("Message sent with ID: %d\n", intId)
		return intId
	}

	return 6
}

func sendUpcomingToDiscord(matches MatchData) {
	if len(matches.Data) > 3 {
		matches.Data = matches.Data[:3] // Only take the first 2 matches
	}

	embeds := make([]map[string]interface{}, len(matches.Data))
	for i, match := range matches.Data {
		region := helpers.GetRegion(match.Tournament)
		title := "Live Match"
		if match.In != "" {
			timestamp, err := helpers.ParseDurationFromNow(match.In)
			if err != nil {
				fmt.Println("Error parsing duration:", err)
				continue
			}
			title = fmt.Sprintf("Upcoming Match at <t:%d:t>", timestamp)
		}
		title = fmt.Sprintf("%s: %s", title, fmt.Sprintf("**%s** vs **%s**", match.Teams[0].Name, match.Teams[1].Name))
		embeds[i] = map[string]interface{}{
			"type":        "rich",
			"title":       title,
			"description": fmt.Sprintf("%s at %s - %s", match.Tournament, match.Event, match.Status),
			"color":       0x00FFFF,
			"footer": map[string]interface{}{
				"text": "Made with ❤️ by bey",
			},
			"fields": []map[string]interface{}{
				{
					"name": "Riot Streams",
					"value": fmt.Sprintf("[Twitch](%s)\n[YouTube](%s)",
						helpers.GetTwitchLink(region), helpers.GetYoutubeLink(region)),
					"inline": true,
				},
				{
					"name":   "Watch Parties",
					"value":  buildWatchPartyLinks(helpers.GetWatchParties(region)),
					"inline": true,
				},
			},
		}
	}

	message := map[string]interface{}{
		"content": "# Here are the upcoming matches:",
		"embeds":  embeds,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return
	}

	messageId := sendToDiscord(webhookURL, messageBytes)
	for _, match := range matches.Data {
		intMatchId, _ := strconv.Atoi(match.ID)
		addSentMessage(intMatchId, messageId)
	}
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
	region := helpers.GetRegion(match.Tournament)
	title := fmt.Sprintf("Match Start: **%s** vs **%s**", match.Teams[0].Name, match.Teams[1].Name)
	embed := map[string]interface{}{
		"type":        "rich",
		"title":       title,
		"description": fmt.Sprintf("%s at %s - %s", match.Tournament, match.Event, match.Status),
		"color":       0x00FFFF,
		"footer": map[string]interface{}{
			"text": "Made with ❤️ by bey",
		},
		"fields": []map[string]interface{}{
			{
				"name": "Riot Streams",
				"value": fmt.Sprintf("[Twitch](%s)\n[YouTube](%s)",
					helpers.GetTwitchLink(region), helpers.GetYoutubeLink(region)),
				"inline": true,
			},
			{
				"name":   "Watch Parties",
				"value":  buildWatchPartyLinks(helpers.GetWatchParties(region)),
				"inline": true,
			},
		},
	}

	message := map[string]interface{}{
		"content": fmt.Sprintf("<@&%s>", roles[region]),
		"embeds":  []map[string]interface{}{embed},
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return
	}

	messageId := sendToDiscord(webhookURL, messageBytes)

	updateSentMessage(messageId, "starting_sent")
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
			"type":        "rich",
			"title":       title,
			"description": fmt.Sprintf("||%s|| Wins: ||%s|| \n %s - %s", winner, score, result.Tournament, result.Event),
			"color":       0x00FFFF,
			"footer": map[string]interface{}{
				"text": "Made with ❤️ by bey",
			},
			"fields": []map[string]interface{}{},
		}
	}

	message := map[string]interface{}{
		"content": "",
		"embeds":  embeds,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return
	}

	messageId := sendToDiscord(webhookURL, messageBytes)
	for range results.Data {
		updateSentMessage(messageId, "result_sent")
	}
}

func getSentMessages() ([]Message, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM messages")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var message Message

		err = rows.Scan(&message.Id, &message.MessageId, &message.MatchId, &message.AnnouncementSent, &message.StartingSent, &message.ResultSent)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return messages, nil
}

func addSentMessage(matchId int, messageId int) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO messages (match_id, message_id, announcement_sent, starting_sent, result_sent) VALUES (?, ?, ?, ?, ?)",
		matchId, messageId, true, false, false)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func updateSentMessage(messageId int, field string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("UPDATE messages SET %s = ? WHERE message_id = ?", field), true, messageId)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
