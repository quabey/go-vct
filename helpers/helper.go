package helpers

import (
	"bey/go-vct/common"
	"fmt"
	"strings"
	"time"
)

var NowFunc = time.Now

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
	durationStr = formatDuration(durationStr);
	duration, _ := time.ParseDuration(durationStr)
	fmt.Println("Duration:", durationStr, "->", duration)

	return time.Now().Add(duration).Round(time.Hour).Unix(), nil
}

func GetOffsetInHours(match1 common.MatchDetail, match2 common.MatchDetail) int {
    t1, _ := time.ParseDuration(formatDuration(match1.In))
    t2, _ := time.ParseDuration(formatDuration(match2.In))
    duration := t2 - t1
    return int(duration.Hours())
}

func GetHoursFromNow(durationStr string) int {
	durationStr = formatDuration(durationStr);
	duration, _ := time.ParseDuration(durationStr)
	return int(duration.Hours())
}

func formatDuration(duration string) string {
	return strings.ReplaceAll(duration, " ", "");
}