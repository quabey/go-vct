package helpers

import (
	"bey/go-vct/common"
	"fmt"
	"strings"
	"time"
)

func GetRegion(tournament string) (region string) {
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

func ParseDurationFromNow(durationStr string) (int64, error) {
	durationStr = strings.ReplaceAll(durationStr, " ", "")
	duration, _ := time.ParseDuration(durationStr)
	fmt.Println("Duration:", durationStr, "->", duration)

	return time.Now().Add(duration).Round(time.Hour).Unix(), nil
}

func GetOffsetInHours(match1 common.MatchDetail, match2 common.MatchDetail) int {
	// Get the time of the first match

	return 0
}