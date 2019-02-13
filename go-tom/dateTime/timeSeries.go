package dateTime

import "time"

// fixme support rounding?
type TimeEntrySeries interface {
	Min() time.Duration
	Max() time.Duration
	Avg() time.Duration

	Total() time.Duration
	DistinctRanges() int

	Add(start time.Time, end time.Time)
}
