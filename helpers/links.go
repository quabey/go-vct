package helpers

import "strings"

var (
	twitchLinks = map[string]string{
		"Americas": "https://www.twitch.tv/valorant_americas",
		"China":    "https://www.twitch.tv/valorantesports_cn",
		"EMEA":     "https://www.twitch.tv/valorant_emea",
		"Pacific":  "https://www.twitch.tv/valorant_pacific",
	}

	youtubeLinks = map[string]string{
		"Americas": "https://www.youtube.com/@valorant_americas/live",
		"China":    "https://www.youtube.com/@VALORANTEsportsCN/live",
		"EMEA":     "https://www.youtube.com/@VALORANTEsportsEMEA/live",
		"Pacific":  "https://www.youtube.com/@VCTPacific/live",
	}

	watchParties = map[string]map[string]string{
		"Pacific": {
			"Sliggy":                  "https://www.twitch.tv/sliggytv",
			"FNS":                     "https://www.twitch.tv/gofns",
			"Sean Gares":              "https://www.twitch.tv/sgares",
			"Thinking Man's Valorant": "https://www.twitch.tv/thinkingmansvalo",
			"Tarik":                   "https://www.twitch.tv/tarik",
			"Kyedae":                  "https://www.twitch.tv/kyedae",
		},
		"EMEA": {
			"Sliggy":                  "https://www.twitch.tv/sliggytv",
			"FNS":                     "https://www.twitch.tv/gofns",
			"Sean Gares":              "https://www.twitch.tv/sgares",
			"Thinking Man's Valorant": "https://www.twitch.tv/thinkingmansvalo",
			"Tarik":                   "https://www.twitch.tv/tarik",
			"Kyedae":                  "https://www.twitch.tv/kyedae",
		},
		"China": {
			"Yinsu":       "https://www.twitch.tv/yinsu",
			"Ryancentral": "https://www.twitch.tv/ryancentral",
		},
		"Americas": {
			"Sliggy":                  "https://www.twitch.tv/sliggytv",
			"FNS":                     "https://www.twitch.tv/gofns",
			"Sean Gares":              "https://www.twitch.tv/sgares",
			"Thinking Man's Valorant": "https://www.twitch.tv/thinkingmansvalo",
			"Tarik":                   "https://www.twitch.tv/tarik",
			"Kyedae":                  "https://www.twitch.tv/kyedae",
		},
	}
)

func CheckVCT(tournament string) bool {
	return strings.Contains(tournament, "Champions")
}

func GetTwitchLink(region string) string {
	return twitchLinks[region]
}

func GetYoutubeLink(region string) string {
	return youtubeLinks[region]
}

func GetWatchParties(region string) map[string]string {
	return watchParties[region]
}
