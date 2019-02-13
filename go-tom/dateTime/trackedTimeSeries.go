package dateTime

import (
	"encoding/json"
	"time"
)

func NewTrackedDaily(filter *DateRange) TimeEntrySeries {
	location := time.Local
	return &trackedTimeSeries{
		loc:    location,
		filter: filter,
		rangeStart: func(t time.Time) time.Time {
			y, m, d := t.In(location).Date()
			return time.Date(y, m, d, 0, 0, 0, 0, location)
		},
	}
}

type trackedTimeSeries struct {
	// the key for a given frame is used to find out which frames belong to the same span of time in the series of entries
	rangeStart func(t time.Time) time.Time
	loc        *time.Location
	filter     *DateRange

	min        time.Duration
	max        time.Duration
	total      time.Duration
	frameCount int
	rangeCount int

	lastRangeStart time.Time
}

func (t *trackedTimeSeries) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Min time.Duration `json:"min"`
		Max time.Duration `json:"max"`
		Avg time.Duration `json:"average"`
	}{
		Min: t.Min(),
		Max: t.Max(),
		Avg: t.Avg(),
	})
}

func (t *trackedTimeSeries) Add(start time.Time, end time.Time) {
	if start.IsZero() || end.IsZero() || (t.filter != nil && !t.filter.Contains(start)) {
		return
	}

	t.frameCount++

	duration := end.Sub(start)
	t.total += duration
	if duration > t.max {
		t.max = duration
	}
	if duration < t.min || t.min == 0 {
		t.min = duration
	}

	frameRange := t.rangeStart(start)
	if t.rangeCount == 0 || frameRange != t.lastRangeStart {
		t.rangeCount++
	}

	t.lastRangeStart = frameRange
}

func (t *trackedTimeSeries) Min() time.Duration {
	return t.min
}

func (t *trackedTimeSeries) Max() time.Duration {
	return t.max
}

func (t *trackedTimeSeries) Avg() time.Duration {
	if t.rangeCount == 0 || t.total == 0 {
		return 0
	}
	return time.Duration(int64(float64(t.total.Nanoseconds()) / float64(t.rangeCount)))
}

func (t *trackedTimeSeries) Total() time.Duration {
	return t.total
}

func (t *trackedTimeSeries) DistinctRanges() int {
	return t.rangeCount
}
