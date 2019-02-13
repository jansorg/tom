package dateTime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTrackedSeries(t *testing.T) {
	s := NewTrackedDaily()

	// day 1: 1 hour
	start := time.Date(2018, time.February, 10, 12, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	s.Add(start, end)

	// day 2: 2 hours
	start = start.Add(24 * time.Hour)
	end = start.Add(2 * time.Hour)
	s.Add(start, end)

	assert.EqualValues(t, 2, s.DistinctRanges())
	assert.EqualValues(t, 1*time.Hour, s.Min())
	assert.EqualValues(t, 2*time.Hour, s.Max())
	assert.EqualValues(t, 1.5*60*time.Minute, s.Avg())

	// day 3: 4 hours + 0.5 hours + 8 hours
	start = start.Add(24 * time.Hour)
	end = start.Add(4 * time.Hour)
	s.Add(start, end)

	end = start.Add(30 * time.Minute)
	s.Add(start, end)

	end = start.Add(8 * time.Hour)
	s.Add(start, end)

	assert.EqualValues(t, 3, s.DistinctRanges())
	assert.EqualValues(t, 30*time.Minute, s.Min())
	assert.EqualValues(t, 8*time.Hour, s.Max())
	assert.EqualValues(t, 5*time.Hour+10*time.Minute, s.Avg(), "expected a daily average of 5:10h (15:30h on 3 days)")

}
