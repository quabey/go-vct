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
			"FNS":                  "https://www.twitch.tv/gofns",
			"Sliggy":               "https://www.twitch.tv/sliggytv",
			"Sgares":               "https://www.twitch.tv/sgares",
			"ThinkingMansValorant": "https://www.twitch.tv/thinkingmansvalorant",
			"tarik":                "https://www.twitch.tv/tarik",
			"kyedae":               "https://www.twitch.tv/kyedae",
		},
		"EMEA": {
			"FNS":                  "https://www.twitch.tv/gofns",
			"Sliggy":               "https://www.twitch.tv/sliggytv",
			"Sgares":               "https://www.twitch.tv/sgares",
			"ThinkingMansValorant": "https://www.twitch.tv/thinkingmansvalorant",
			"tarik":                "https://www.twitch.tv/tarik",
			"kyedae":               "https://www.twitch.tv/kyedae",
		},
		"China": {
			"Ryancentral": "https://www.twitch.tv/ryancentral",
			"Yinsu":       "https://www.twitch.tv/yinsu",
		},
		"Americas": {
			"FNS":                  "https://www.twitch.tv/gofns",
			"Sliggy":               "https://www.twitch.tv/sliggytv",
			"Sgares":               "https://www.twitch.tv/sgares",
			"ThinkingMansValorant": "https://www.twitch.tv/thinkingmansvalorant",
			"tarik":                "https://www.twitch.tv/tarik",
			"kyedae":               "https://www.twitch.tv/kyedae",
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
