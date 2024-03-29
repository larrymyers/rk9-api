package rk9

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRound(t *testing.T) {
	event := Event{
		PairingsURL: "/pairings/VAN01mctlk7ZeDoADoIh",
	}

	_, err := GetRound(&event, Masters, 1)
	assert.NoError(t, err)
}
