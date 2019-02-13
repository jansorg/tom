package dateTime

import (
	"encoding/json"
	"time"
)

// fixme support rounding?
type TimeEntrySeries interface {
	json.Marshaler

	Min() time.Duration
	Max() time.Duration
	Avg() time.Duration

	Total() time.Duration
	DistinctRanges() int

	Add(start time.Time, end time.Time)
}
