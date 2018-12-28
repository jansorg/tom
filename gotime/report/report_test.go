package report

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jansorg/gotime/gotime/store"
)

func Test_Report(t *testing.T) {
	start := newDate(2018, time.March, 10, 10, 0)
	end := newDate(2018, time.March, 10, 12, 0)

	frameList := []*store.Frame{{Start: start, End: end}}
	report := NewBucketReport(frameList)
	report.Update()
	assert.EqualValues(t, 1, report.Result.FrameCount)
	assert.EqualValues(t, 2*time.Hour, report.Result.Duration)
	assert.EqualValues(t, 2*time.Hour, report.Result.ExactDuration)
	assert.EqualValues(t, start, report.Result.From)
	assert.EqualValues(t, end, report.Result.To)
	assert.EqualValues(t, frameList, report.source)
}

func Test_ReportSplitYear(t *testing.T) {
	// two hours
	start := newDate(2018, time.March, 10, 10, 0)
	end := newDate(2018, time.March, 10, 12, 0)

	// one hour
	start2 := newDate(2019, time.March, 10, 9, 0)
	end2 := newDate(2019, time.March, 10, 10, 0)

	frameList := []*store.Frame{
		{Start: start2, End: end2},
		{Start: start, End: end},
	}

	report := NewBucketReport(frameList)
	report.GroupByYear = true
	report.Update()

	require.NotNil(t, report.Result, "expected one top-level group (containing two years)")
	assert.EqualValues(t, 4, report.Result.FrameCount)
	assert.EqualValues(t, 3*time.Hour, report.Result.Duration)
	assert.EqualValues(t, 3*time.Hour, report.Result.ExactDuration)
	assert.EqualValues(t, newDate(2018, time.January, 1, 0, 0), report.Result.From)
	assert.EqualValues(t, newDate(2020, time.January, 1, 0, 0), report.Result.To)
	assert.EqualValues(t, frameList, report.source)

	require.EqualValues(t, 2, len(report.Result.Results), "expected a sub-report for each year")

	firstYear := report.Result.Results[0]
	require.EqualValues(t, newDate(2018, time.January, 1, 0, 0), firstYear.From)
	require.EqualValues(t, newDate(2019, time.January, 1, 0, 0), firstYear.To)
	require.EqualValues(t, 1, firstYear.FrameCount)
	require.EqualValues(t, 2*time.Hour, firstYear.ExactDuration)
	require.EqualValues(t, 2*time.Hour, firstYear.Duration)

	secondYear := report.Result.Results[1]
	require.EqualValues(t, newDate(2019, time.January, 1, 0, 0), secondYear.From)
	require.EqualValues(t, newDate(2020, time.January, 1, 0, 0), secondYear.To)
	require.EqualValues(t, 1, secondYear.FrameCount)
	require.EqualValues(t, 1*time.Hour, secondYear.ExactDuration)
	require.EqualValues(t, 1*time.Hour, secondYear.Duration)
}

func newDate(year int, month time.Month, day, hour, minute int) *time.Time {
	date := time.Date(year, month, day, hour, minute, 0, 0, time.Local)
	return &date
}
