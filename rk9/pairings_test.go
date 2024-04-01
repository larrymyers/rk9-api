package rk9

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRound(t *testing.T) {
	event := Event{
		ID: "VAN01mctlk7ZeDoADoIh",
	}

	_, err := GetRound(&event, Masters, 1)
	assert.NoError(t, err)
}

func TestParseRecord(t *testing.T) {
	wins, losses, ties, points, err := parseRecord("(6-1-2) 20 pts")

	assert.NoError(t, err)
	assert.Equal(t, 6, wins)
	assert.Equal(t, 1, losses)
	assert.Equal(t, 2, ties)
	assert.Equal(t, 20, points)
}

func TestParseName(t *testing.T) {
	name, country := parseName("Larry Myers [US]")

	assert.Equal(t, "Larry Myers", name)
	assert.Equal(t, "US", country)
}
