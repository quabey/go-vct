package main

import (
	"testing"
)

func TestTimestamp(t *testing.T) {
	parseDurationFromNow("1h 15m")
	parseDurationFromNow("d")
}
