package main

import (
	"testing"
)

func TestTimestamp(t *testing.T) {
	ParseDurationFromNow("1h 15m")
	ParseDurationFromNow("d")
}
