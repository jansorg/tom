package dateTime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUntrackedSeries(t *testing.T) {
	s := NewUntrackedDaily(nil)

	first := time.Date(2018, time.February, 10, 6, 0, 0, 0, time.UTC)

	// day 1: 30 min untracked between 1 and 2
	start := first
	end := start.Add(time.Hour)
	s.Add(start, end)

	start = end.Add(30 * time.Minute)
	end = start.Add(time.Hour)
	s.Add(start, end)

	// day 2: 1 hour untracked
	start = first.Add(24 * time.Hour)
	end = start.Add(2 * time.Hour)
	s.Add(start, end)

	start = end.Add(1 * time.Hour)
	end = start.Add(2 * time.Hour)
	s.Add(start, end)

	assert.EqualValues(t, 2, s.DistinctRanges())
	assert.EqualValues(t, 30*time.Minute, s.Min())
	assert.EqualValues(t, 1*time.Hour, s.Max())
	assert.EqualValues(t, 45*time.Minute, s.Avg())

	// day 3: 20 min + 30min + 1 hour + 5 hours + 40 min
	start = first.Add(2 * 24 * time.Hour)
	end = start.Add(1 * time.Hour)
	s.Add(start, end)

	start = end.Add(20 * time.Minute)
	end = start.Add(1 * time.Hour)
	s.Add(start, end)

	start = end.Add(30 * time.Minute)
	end = start.Add(1 * time.Hour)
	s.Add(start, end)

	start = end.Add(1 * time.Hour)
	end = start.Add(1 * time.Hour)
	s.Add(start, end)

	start = end.Add(5 * time.Hour)
	end = start.Add(1 * time.Hour)
	s.Add(start, end)

	start = end.Add(40 * time.Minute)
	end = start.Add(1 * time.Hour)
	s.Add(start, end)

	assert.EqualValues(t, 3, s.DistinctRanges())
	assert.EqualValues(t, 20*time.Minute, s.Min())
	assert.EqualValues(t, 5*time.Hour, s.Max())
	assert.EqualValues(t, 3*time.Hour, s.Avg(), "expected a daily average of 3h (9 hours untracked over 3 days)")
}
