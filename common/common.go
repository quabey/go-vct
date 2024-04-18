package common

import (
	"fmt"
	"os"
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

var Roles = map[string]string{
	"Americas": "1227214059498115072",
	"China":    "1227214116834382009",
	"EMEA":     "1227214030616264734",
	"Pacific":  "1227213846268346440",
}

type Message struct {
	Id               int
	MessageId        string
	MatchId          string
	AnnouncementSent bool
	StartingSent     bool
	ResultSent       bool
	Timestamp        int64
}

var (
	WebhookURL string
	DbPath     string
)

func LoadEnvVariables() {
	WebhookURL = fmt.Sprintf("%s?wait=true", os.Getenv("WEBHOOK_URL"));
	DbPath = os.Getenv("SQLITE_DB")
}
