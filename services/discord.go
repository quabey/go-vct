package services

import (
	"bey/go-vct/common"
	"bey/go-vct/database"
	"bey/go-vct/helpers"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func sendToServices(url string, payload []byte) int {
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

func SendUpcomingToServices(match common.MatchDetail, addFields bool, addContent bool) {
    region := helpers.GetRegion(match.Tournament)
    title := "Live Match"
    if match.In != "" {
        timestamp, err := helpers.ParseDurationFromNow(match.In)
        if err != nil {
            fmt.Println("Error parsing duration:", err)
            return
        }
        title = fmt.Sprintf("Upcoming Match at <t:%d:t>", timestamp)
    }
    title = fmt.Sprintf("%s: %s", title, fmt.Sprintf("**%s** vs **%s**", match.Teams[0].Name, match.Teams[1].Name))
    embed := map[string]interface{}{
        "type":        "rich",
        "title":       title,
        "description": fmt.Sprintf("%s at %s - %s", match.Tournament, match.Event, match.Status),
        "color":       0x00FFFF,
        "footer": map[string]interface{}{
            "text": "Made with ❤️ by bey & Nate",
        },
    }

    if addFields {
        embed["fields"] = []map[string]interface{}{
            {
                "name": "Riot Streams",
                "value": fmt.Sprintf("[Twitch](%s)\n[YouTube](%s)",
                    helpers.GetTwitchLink(region), helpers.GetYoutubeLink(region)),
                "inline": true,
            },
            {
                "name":   "Watch Parties",
                "value":  BuildWatchPartyLinks(helpers.GetWatchParties(region)),
                "inline": true,
            },
        }
    }

	content := ""

	if addContent {
		content = fmt.Sprintf("## Upcoming matches for %s", helpers.GetRegion(match.Tournament))
	} 
    message := map[string]interface{}{
        "content": content,
        "embeds":  []map[string]interface{}{embed},
    }

    messageBytes, err := json.Marshal(message)
    if err != nil {
        fmt.Println("Error marshalling message:", err)
        return
    }

    messageId := sendToServices(common.WebhookURL, messageBytes)
    intMatchId, _ := strconv.Atoi(match.ID)
    database.AddSentMessage(intMatchId, messageId)
}

func BuildWatchPartyLinks(parties map[string]string) string {
	if len(parties) == 0 {
		return "No watch parties available"
	}
	var links []string
	for name, url := range parties {
		links = append(links, fmt.Sprintf("[%s](%s)", name, url))
	}
	return strings.Join(links, "\n")
}

func SendMatchStartToservices(match common.MatchDetail) {
	region := helpers.GetRegion(match.Tournament)
	title := fmt.Sprintf("Match Start: **%s** vs **%s**", match.Teams[0].Name, match.Teams[1].Name)
	embed := map[string]interface{}{
		"type":        "rich",
		"title":       title,
		"description": fmt.Sprintf("%s at %s - %s", match.Tournament, match.Event, match.Status),
		"color":       0x00FFFF,
		"footer": map[string]interface{}{
			"text": "Made with ❤️ by bey & Nate",
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
				"value":  BuildWatchPartyLinks(helpers.GetWatchParties(region)),
				"inline": true,
			},
		},
	}

	message := map[string]interface{}{
		"content": fmt.Sprintf("<@&%s>", common.Roles[region]),
		"embeds":  []map[string]interface{}{embed},
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return
	}

	messageId := sendToServices(common.WebhookURL, messageBytes)

	database.UpdateSentMessage(messageId, "starting_sent")
}

func SendResultsToservices(results common.MatchData) {
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
				"text": "Made with ❤️ by bey & Nate",
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

	messageId := sendToServices(common.WebhookURL, messageBytes)
	for range results.Data {
		database.UpdateSentMessage(messageId, "result_sent")
	}
}
