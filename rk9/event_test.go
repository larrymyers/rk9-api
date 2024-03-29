package rk9

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDateRange(t *testing.T) {
	start, end, err := parseDateRange("April 5-7, 2024")

	assert.NoError(t, err)
	assert.Equal(t, time.Date(2024, 4, 5, 0, 0, 0, 0, time.UTC), start)
	assert.Equal(t, time.Date(2024, 4, 7, 0, 0, 0, 0, time.UTC), end)

	start, end, err = parseDateRange("August 11–13, 2023")

	assert.NoError(t, err)
	assert.Equal(t, time.Date(2023, 8, 11, 0, 0, 0, 0, time.UTC), start)
	assert.Equal(t, time.Date(2023, 8, 13, 0, 0, 0, 0, time.UTC), end)

	start, end, err = parseDateRange("June 30–July 2, 2023")
	assert.NoError(t, err)
}
