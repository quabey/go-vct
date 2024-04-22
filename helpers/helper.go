package helpers

import (
	"bey/go-vct/common"
	"fmt"
	"strings"
	"time"

	str2duration "github.com/xhit/go-str2duration/v2"
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
	formattedDurationStr := formatDuration(durationStr)
	duration, err := str2duration.ParseDuration(formattedDurationStr)
	if err != nil {
		return 0, err
	}
	resultTime := time.Now().Add(duration).Round(time.Hour)
	fmt.Printf("%s -> %s (%s) \n", durationStr, duration, resultTime)

	return resultTime.Unix(), nil
}

func GetOffsetInHours(match1 common.MatchDetail, match2 common.MatchDetail) int {
	t1, _ := str2duration.ParseDuration(formatDuration(match1.In))
	t2, _ := str2duration.ParseDuration(formatDuration(match2.In))
	offset := t2 - t1
	fmt.Printf("Dur1: %s, Dur2: %s, Offset: ", t1, t2)
	fmt.Println(int(offset.Hours()))
	return int(offset.Hours())
}

func GetHoursFromNow(durationStr string) int {
	durationStr = formatDuration(durationStr)
	duration, _ := str2duration.ParseDuration(durationStr)
	fmt.Println()
	return int(duration.Hours())
}

func formatDuration(duration string) string {
	return strings.ReplaceAll(duration, " ", "")
}
