package report

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/test_setup"
)

func TestSplitEmptyReport(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top")
	require.NoError(t, err)
	p2, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top", "child")
	require.NoError(t, err)

	start := newDate(2018, time.March, 10, 10, 0)
	end := newDate(2018, time.March, 10, 12, 0)

	frameList := []*model.Frame{{
		Start: start, End: end, ProjectId: p2.ID,
	}}
	report := NewBucketReport(model.NewFrameList(frameList), ctx)
	report.ProjectIDs = []string{p1.ID}
	report.IncludeSubprojects = true
	report.SplitOperations = []SplitOperation{SplitByProject, SplitByMonth}
	report.Update()

	assert.NotNil(t, report.Result.DateRange().Start, "empty start value in top-level date range")
	assert.NotNil(t, report.Result.DateRange().End, "empty end value in top-level date range")

	assert.NotNil(t, report.Result.TrackedDateRange().Start, "empty start value in top-level date range")
	assert.NotNil(t, report.Result.TrackedDateRange().End, "empty end value in top-level date range")

	for i, b := range report.Result.ChildBuckets {
		assert.NotNil(t, b.DateRange().Start, "empty date range at index %d", i)
		assert.NotNil(t, b.DateRange().End, "empty date range at index %d", i)
	}
}

func TestReport(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	start := newDate(2018, time.March, 10, 10, 0)
	end := newDate(2018, time.March, 10, 12, 0)

	frameList := []*model.Frame{{Start: start, End: end}}
	report := NewBucketReport(model.NewFrameList(frameList), ctx)
	report.Update()
	assert.EqualValues(t, 1, report.Result.FrameCount)
	assert.EqualValues(t, 2*time.Hour, report.Result.Duration.Get())
	assert.EqualValues(t, 2*time.Hour, report.Result.Duration.GetExact())
	assert.EqualValues(t, start, report.Result.TrackedDateRange().Start)
	assert.EqualValues(t, end, report.Result.TrackedDateRange().End)
	assert.EqualValues(t, frameList, report.source.Frames())
}

func TestReportSplitYear(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	// two hours
	start := newDate(2018, time.March, 10, 10, 0)
	end := newDate(2018, time.March, 10, 12, 0)

	// one hour
	start2 := newDate(2019, time.March, 10, 9, 0)
	end2 := newDate(2019, time.March, 10, 10, 0)

	frameList := []*model.Frame{
		{Start: start2, End: end2},
		{Start: start, End: end},
	}

	report := NewBucketReport(model.NewSortedFrameList(frameList), ctx)
	report.SplitOperations = []SplitOperation{SplitByYear}
	report.Update()

	require.NotNil(t, report.Result, "expected one top-level group (containing two years)")
	require.EqualValues(t, 2, len(report.Result.ChildBuckets), "expected a sub-report for each year")

	assert.EqualValues(t, 2, report.Result.FrameCount)
	assert.EqualValues(t, 3*time.Hour, report.Result.Duration.Get())
	assert.EqualValues(t, 3*time.Hour, report.Result.Duration.GetExact())
	assert.EqualValues(t, newDate(2018, time.January, 1, 0, 0), report.Result.DateRange().Start)
	assert.EqualValues(t, newDate(2020, time.January, 1, 0, 0), report.Result.DateRange().End)
	assert.EqualValues(t, start, report.Result.TrackedDateRange().Start)
	assert.EqualValues(t, end2, report.Result.TrackedDateRange().End)
	assert.EqualValues(t, frameList, report.source.Frames())

	firstYear := report.Result.ChildBuckets[0]
	assert.EqualValues(t, newDate(2018, time.January, 1, 0, 0), firstYear.DateRange().Start)
	assert.EqualValues(t, newDate(2019, time.January, 1, 0, 0), firstYear.DateRange().End)
	assert.EqualValues(t, 1, firstYear.FrameCount)
	assert.EqualValues(t, 2*time.Hour, firstYear.Duration.Get())
	assert.EqualValues(t, 2*time.Hour, firstYear.Duration.GetExact())

	secondYear := report.Result.ChildBuckets[1]
	assert.EqualValues(t, newDate(2019, time.January, 1, 0, 0), secondYear.DateRange().Start)
	assert.EqualValues(t, newDate(2020, time.January, 1, 0, 0), secondYear.DateRange().End)
	assert.EqualValues(t, 1, secondYear.FrameCount)
	assert.EqualValues(t, 1*time.Hour, secondYear.Duration.Get())
	assert.EqualValues(t, 1*time.Hour, secondYear.Duration.GetExact())
}

func TestReportDateRanges(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	// two hours
	start := newDate(2018, time.March, 10, 10, 0)
	end := newDate(2018, time.March, 10, 12, 0)

	// one hour
	start2 := newDate(2019, time.March, 10, 9, 0)
	end2 := newDate(2019, time.March, 10, 10, 0)

	frameList := []*model.Frame{
		{Start: start2, End: end2},
		{Start: start, End: end},
	}

	var op SplitOperation = 0
	for ; op < SplitByParentProject; op += 1 {
		report := NewBucketReport(model.NewSortedFrameList(frameList), ctx)
		report.SplitOperations = []SplitOperation{op}
		report.Update()

		require.NotNil(t, report.Result, "expected one top-level group (containing two years)")

		assert.EqualValues(t, 2, report.Result.FrameCount)
		assert.EqualValues(t, 3*time.Hour, report.Result.Duration.Get())
		assert.EqualValues(t, 3*time.Hour, report.Result.Duration.GetExact())
		assert.EqualValues(t, start, report.Result.TrackedDateRange().Start, "unexpected tracked time for "+op.String())
		assert.EqualValues(t, end2, report.Result.TrackedDateRange().End, "unexpected tracked time for "+op.String())
		assert.EqualValues(t, frameList, report.source.Frames())
	}
}

// tests that frame dates in different time zones are not ending up in different split intervals
func TestReportSplitDifferentZones(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	// two hours, UTC+0
	start := newDate(2018, time.March, 10, 10, 0).UTC()
	end := newDate(2018, time.March, 10, 12, 0).UTC()

	// one hour, UTC+2
	utc2, err := time.LoadLocation("Europe/Berlin")
	require.NoError(t, err)
	start2 := newDate(2018, time.March, 10, 9, 0).In(utc2)
	end2 := newDate(2018, time.March, 10, 10, 0).In(utc2)

	frameList := []*model.Frame{
		{Start: &start, End: &end},
		{Start: &start2, End: &end2},
	}

	report := NewBucketReport(model.NewSortedFrameList(frameList), ctx)
	report.SplitOperations = []SplitOperation{SplitByDay}
	report.Update()

	require.NotNil(t, report.Result, "expected one top-level group (containing two days)")
	assert.EqualValues(t, 2, report.Result.FrameCount)

	assert.EqualValues(t, 1, len(report.Result.ChildBuckets), "expected a single day bucket, even if different time zones were used")
}

// tests that frame dates in different time zones are not ending up in different split intervals
func TestReportSplitDifferentZonesYear(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	// two hours, UTC+0
	start := newDate(2018, time.January, 10, 10, 0).UTC()
	end := newDate(2018, time.January, 10, 12, 0).UTC()

	// one hour, UTC+2
	utc2, err := time.LoadLocation("Europe/Berlin")
	require.NoError(t, err)
	start2 := newDate(2018, time.January, 10, 9, 0).In(utc2)
	end2 := newDate(2018, time.January, 10, 10, 0).In(utc2)

	frameList := []*model.Frame{
		{Start: &start, End: &end},
		{Start: &start2, End: &end2},
	}

	report := NewBucketReport(model.NewSortedFrameList(frameList), ctx)
	report.SplitOperations = []SplitOperation{SplitByYear}
	report.Update()

	require.NotNil(t, report.Result, "expected one top-level group (containing two days)")
	assert.EqualValues(t, 2, report.Result.FrameCount)

	assert.EqualValues(t, 1, len(report.Result.ChildBuckets), "expected a single day bucket, even if different time zones were used")
}

func TestReportEmptyRanges(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1")
	require.NoError(t, err)

	p2, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project2")
	require.NoError(t, err)

	start := newDate(2017, time.January, 10, 10, 0).UTC()
	end := newDate(2017, time.January, 10, 12, 0).UTC()

	start2 := newDate(2019, time.January, 10, 9, 0).UTC()
	end2 := newDate(2019, time.January, 10, 10, 0).UTC()

	frameList := []*model.Frame{
		{Start: &start, End: &end, ProjectId: p.ID},
		{Start: &start2, End: &end2, ProjectId: p2.ID},
	}

	report := NewBucketReport(model.NewSortedFrameList(frameList), ctx)
	report.SplitOperations = []SplitOperation{SplitByProject, SplitByYear, SplitByMonth}
	report.ShowEmptyBuckets = true
	report.Update()

	require.NotNil(t, report.Result, "expected one top-level group (containing two frames)")
	assert.EqualValues(t, 2, report.Result.FrameCount)

	assert.EqualValues(t, 2, len(report.Result.ChildBuckets), "expected two project buckets")

	for i, projectBucket := range report.Result.ChildBuckets {
		require.EqualValues(t, 3, len(projectBucket.ChildBuckets), "expected three year bucket, 2017 .. 2019")
		assert.True(t, projectBucket.IsProjectBucket())

		for j, yearBucket := range projectBucket.ChildBuckets {
			assert.True(t, yearBucket.IsDateBucket())
			require.EqualValues(t, 12, len(yearBucket.ChildBuckets), "bucket doesn't contain months, index %d | %d", i, j)
		}
	}
}

func newDate(year int, month time.Month, day, hour, minute int) *time.Time {
	date := time.Date(year, month, day, hour, minute, 0, 0, time.Local)
	return &date
}
