package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	start := newDate(2019, time.January, 1, 10, 0)
	end := newDate(2019, time.January, 1, 12, 0)
	f := Frame{
		Start: start,
		End:   end,
	}

	ref := start.Add(-10 * time.Minute)
	assert.False(t, f.Contains(&ref))

	assert.True(t, f.Contains(start))

	ref = start.Add(10 * time.Minute)
	assert.True(t, f.Contains(&ref))

	assert.True(t, f.Contains(end))

	ref = end.Add(10 * time.Minute)
	assert.False(t, f.Contains(&ref))
}
