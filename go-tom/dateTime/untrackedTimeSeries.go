package dateTime

import (
	"fmt"
	"time"
)

func NewUntrackedDaily() TimeEntrySeries {
	location := time.Local
	return &untrackedTimeSeries{
		loc: location,
		rangeStart: func(t time.Time) time.Time {
			y, m, d := t.In(location).Date()
			return time.Date(y, m, d, 0, 0, 0, 0, location)
		},
	}
}

type untrackedTimeSeries struct {
	// the key for a given frame is used to find out which frames belong to the same span of time in the series of entries
	rangeStart func(t time.Time) time.Time
	loc        *time.Location

	min        time.Duration
	max        time.Duration
	total      time.Duration
	frameCount int
	rangeCount int

	lastRangeStart time.Time
	lastEnd        time.Time
}

func (t *untrackedTimeSeries) Add(start time.Time, end time.Time) {
	if start.IsZero() || end.IsZero() {
		return
	}

	// fixme handle location

	fmt.Printf("%s to %s\n", start.String(), end.String())
	t.frameCount++

	// first entry
	if t.lastEnd.IsZero() {
		t.lastEnd = end
		t.lastRangeStart = t.rangeStart(start)
		t.rangeCount++
		return
	}

	rangeStart := t.rangeStart(start)
	rangeEnd := t.rangeStart(end)
	if rangeStart != rangeEnd {
		// not yet supported
		return
	}

	if rangeStart != t.lastRangeStart {
		t.lastEnd = end
		t.rangeCount++
		t.lastRangeStart = rangeStart
		return
	}

	// same range as previous entry, add the time between lastEnd and current start
	untracked := start.Sub(t.lastEnd)
	t.lastEnd = end
	if untracked < 0 {
		// unsupported, probably invalid order of source
		return
	}

	if untracked > t.max {
		t.max = untracked
	}
	if untracked < t.min || t.min == 0 {
		t.min = untracked
	}
	t.total += untracked
}

func (t *untrackedTimeSeries) Min() time.Duration {
	return t.min
}

func (t *untrackedTimeSeries) Max() time.Duration {
	return t.max
}

func (t *untrackedTimeSeries) Avg() time.Duration {
	return time.Duration(int64(float64(t.total.Nanoseconds()) / float64(t.rangeCount)))
}

func (t *untrackedTimeSeries) Total() time.Duration {
	return t.total
}

func (t *untrackedTimeSeries) DistinctRanges() int {
	return t.rangeCount
}
