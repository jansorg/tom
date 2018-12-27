package report

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jansorg/gotime/gotime/frames"
	"github.com/jansorg/gotime/gotime/store"
)

func Test_Report(t *testing.T) {
	start := newDate(2018, time.March, 10, 10, 0)
	end := newDate(2018, time.March, 10, 12, 0)

	frameList := []*store.Frame{{Start: start, End: end}}
	report := NewBucketReport(frameList)
	report.Update()
	assert.EqualValues(t, report.Results, []*ResultBucket{
		{
			Duration:      2 * time.Hour,
			ExactDuration: 2 * time.Hour,
			FrameCount:    1,
			Source: &frames.Bucket{
				From:   *start,
				To:     *end,
				Frames: frameList,
			},
		},
	})
}

func newDate(year int, month time.Month, day, hour, minute int) *time.Time {
	date := time.Date(year, month, day, hour, minute, 0, 0, time.Local)
	return &date
}
