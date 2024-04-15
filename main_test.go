package main

import (
	"bey/go-vct/helpers"
	"testing"
)

func TestTimestamp(t *testing.T) {
	helpers.ParseDurationFromNow("1h 15m")
	helpers.ParseDurationFromNow("d")
}
