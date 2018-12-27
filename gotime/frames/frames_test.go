package frames

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jansorg/gotime/gotime/store"
)

func TestSplitByDay(t *testing.T) {
	buckets := SplitByDay([]store.Frame{
		{
			Start: newTime(10, 12, 0),
			End:   newTime(10, 13, 30),
		},
		{
			Start: newTime(10, 14, 0),
			End:   newTime(10, 16, 0),
		},
		{
			Start: newTime(11, 12, 0),
			End:   newTime(11, 13, 30),
		},
		{
			Start: newTime(13, 1, 0),
			End:   newTime(13, 3, 0),
		},
	})

	assert.EqualValues(t, 3, len(buckets))
	assert.EqualValues(t, 2, len(buckets[0].Frames))
	assert.EqualValues(t, 1, len(buckets[1].Frames))
	assert.EqualValues(t, 1, len(buckets[1].Frames))
}

func newTime(day, hour, minute int) *time.Time {
	date := time.Date(2018, time.December, day, hour, minute, 0, 0, time.Local)
	return &date
}
