package report

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jansorg/gotime/gotime/store"
)

func Test_Report(t *testing.T) {
	start := newDate(2018, time.March, 10, 10, 0)
	end := newDate(2018, time.March, 10, 12, 0)

	frameList := []*store.Frame{{Start: start, End: end}}
	report := NewBucketReport(frameList)
	report.Update()
	assert.EqualValues(t, 1, report.Results[0].FrameCount)
	assert.EqualValues(t, 2*time.Hour, report.Results[0].Duration)
	assert.EqualValues(t, 2*time.Hour, report.Results[0].ExactDuration)
	assert.EqualValues(t, start, report.Results[0].From)
	assert.EqualValues(t, end, report.Results[0].To)
	assert.EqualValues(t, frameList, report.source)
}

func newDate(year int, month time.Month, day, hour, minute int) *time.Time {
	date := time.Date(year, month, day, hour, minute, 0, 0, time.Local)
	return &date
}
