package tests

import (
	"bey/go-vct/common"
	"bey/go-vct/helpers"
	"testing"
)

func TestGetOffsetInHours(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		match1   common.MatchDetail
		match2   common.MatchDetail
		expected int
	}{
		{
			name: "Test Case 1: Positive offset",
			match1: common.MatchDetail{
				In: "2h 30m",
			},
			match2: common.MatchDetail{
				In: "1h   30m",
			},
			expected: -1,
		},
		{
			name: "Test Case 2: Negative offset",
			match1: common.MatchDetail{
				In: "1h 4m",
			},
			match2: common.MatchDetail{
				In: "2h 4m",
			},
			expected: 1,
		},
		{
			name: "Test Case 3: Zero offset",
			match1: common.MatchDetail{
				In: "1h 1m",
			},
			match2: common.MatchDetail{
				In: "1h 1m",
			},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := helpers.GetOffsetInHours(tc.match1, tc.match2)
			if result != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestParseDurationFromNow(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		durationStr string
		hasError    bool
	}{
		{
			name:        "Test Case 1: Valid duration",
			durationStr: "2h",
			hasError:    false,
		},
		{
			name:        "Test Case 2: Invalid duration",
			durationStr: "invalid",
			hasError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := helpers.ParseDurationFromNow(tc.durationStr)
			if (err != nil) != tc.hasError {
				t.Errorf("ParseDurationFromNow() error = %v, wantErr %v", err, tc.hasError)
				return
			}
		})
	}
}
