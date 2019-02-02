package report

import (
	"fmt"
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
	report := NewBucketReport(model.NewFrameList(frameList), Config{
		ProjectIDs:         []string{p1.ID},
		IncludeSubprojects: true,
		Splitting:          []SplitOperation{SplitByProject, SplitByMonth},
	}, ctx)
	report.Update()

	assert.NotNil(t, report.result.DateRange().Start, "empty start value in top-level date range")
	assert.NotNil(t, report.result.DateRange().End, "empty end value in top-level date range")

	assert.NotNil(t, report.result.TrackedDateRange().Start, "empty start value in top-level date range")
	assert.NotNil(t, report.result.TrackedDateRange().End, "empty end value in top-level date range")

	for i, b := range report.result.ChildBuckets {
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
	report := NewBucketReport(model.NewFrameList(frameList), Config{}, ctx)
	report.Update()
	assert.EqualValues(t, 1, report.result.FrameCount)
	assert.EqualValues(t, 2*time.Hour, report.result.Duration.Get())
	assert.EqualValues(t, 2*time.Hour, report.result.Duration.GetExact())
	assert.EqualValues(t, start, report.result.TrackedDateRange().Start)
	assert.EqualValues(t, end, report.result.TrackedDateRange().End)
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

	report := NewBucketReport(model.NewSortedFrameList(frameList), Config{Splitting: []SplitOperation{SplitByYear}}, ctx)
	report.Update()

	require.NotNil(t, report.result, "expected one top-level group (containing two years)")
	require.EqualValues(t, 2, len(report.result.ChildBuckets), "expected a sub-report for each year")

	assert.EqualValues(t, 2, report.result.FrameCount)
	assert.EqualValues(t, 3*time.Hour, report.result.Duration.Get())
	assert.EqualValues(t, 3*time.Hour, report.result.Duration.GetExact())
	assert.EqualValues(t, newDate(2018, time.January, 1, 0, 0), report.result.DateRange().Start)
	assert.EqualValues(t, newDate(2020, time.January, 1, 0, 0), report.result.DateRange().End)
	assert.EqualValues(t, start, report.result.TrackedDateRange().Start)
	assert.EqualValues(t, end2, report.result.TrackedDateRange().End)
	assert.EqualValues(t, frameList, report.source.Frames())

	firstYear := report.result.ChildBuckets[0]
	assert.EqualValues(t, newDate(2018, time.January, 1, 0, 0), firstYear.DateRange().Start)
	assert.EqualValues(t, newDate(2019, time.January, 1, 0, 0), firstYear.DateRange().End)
	assert.EqualValues(t, 1, firstYear.FrameCount)
	assert.EqualValues(t, 2*time.Hour, firstYear.Duration.Get())
	assert.EqualValues(t, 2*time.Hour, firstYear.Duration.GetExact())

	secondYear := report.result.ChildBuckets[1]
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

	for op := SplitByYear; op < SplitByProject; op += 1 {
		report := NewBucketReport(model.NewSortedFrameList(frameList), Config{Splitting: []SplitOperation{op}}, ctx)
		report.Update()

		require.NotNil(t, report.result, "expected one top-level group (containing two years)")

		assert.EqualValues(t, 2, report.result.FrameCount)
		assert.EqualValues(t, 3*time.Hour, report.result.Duration.Get())
		assert.EqualValues(t, 3*time.Hour, report.result.Duration.GetExact())
		assert.EqualValues(t, start, report.result.TrackedDateRange().Start, "unexpected tracked time for "+op.String())
		assert.EqualValues(t, end2, report.result.TrackedDateRange().End, "unexpected tracked time for "+op.String())
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

	report := NewBucketReport(model.NewSortedFrameList(frameList), Config{Splitting: []SplitOperation{SplitByDay}}, ctx)
	report.Update()

	require.NotNil(t, report.result, "expected one top-level group (containing two days)")
	assert.EqualValues(t, 2, report.result.FrameCount)

	assert.EqualValues(t, 1, len(report.result.ChildBuckets), "expected a single day bucket, even if different time zones were used")
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

	report := NewBucketReport(model.NewSortedFrameList(frameList), Config{Splitting: []SplitOperation{SplitByYear}}, ctx)
	report.Update()

	require.NotNil(t, report.result, "expected one top-level group (containing two days)")
	assert.EqualValues(t, 2, report.result.FrameCount)

	assert.EqualValues(t, 1, len(report.result.ChildBuckets), "expected a single day bucket, even if different time zones were used")
}

func TestReportProjectHierarchySum(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	pTop, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1")
	require.NoError(t, err)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1", "child1")
	require.NoError(t, err)

	p2, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1", "child2")
	require.NoError(t, err)

	_, _, err = ctx.StoreHelper.GetOrCreateNestedProjectNames("unused project")
	require.NoError(t, err)

	// adding 3*10*2 hours = 60 hours, 20 for each project
	frames := model.NewEmptyFrameList()
	start := newDate(2017, time.January, 10, 10, 0).UTC()
	end := newDate(2017, time.January, 10, 12, 0).UTC()
	for i := 0; i < 10; i++ {
		frames.Append(&model.Frame{Start: &start, End: &end, ProjectId: pTop.ID})
		frames.Append(&model.Frame{Start: &start, End: &end, ProjectId: p1.ID})
		frames.Append(&model.Frame{Start: &start, End: &end, ProjectId: p2.ID})

		start = start.Add(10 * time.Minute)
		end = end.Add(10 * time.Minute)
	}
	frames.Sort()

	report := NewBucketReport(frames.Copy(), Config{
		ProjectIDs: []string{pTop.ID},
	}, ctx)
	report.Update()
	require.NotNil(t, report.result, "expected one top-level group with 30 frames")
	assert.EqualValues(t, 0, len(report.result.ChildBuckets))
	assert.EqualValues(t, 30, report.result.FrameCount)
	assert.EqualValues(t, 60*time.Hour, report.result.Duration.SumExact)

	// splitting without subprojects must result in only one project bucket
	report = NewBucketReport(frames.Copy(), Config{
		ProjectIDs: []string{pTop.ID},
		Splitting:  []SplitOperation{SplitByProject},
	}, ctx)
	report.Update()
	require.NotNil(t, report.result, "expected one top-level group with 30 frames")
	assert.EqualValues(t, 1, len(report.result.ChildBuckets))
	assert.EqualValues(t, 30, report.result.FrameCount)
	assert.EqualValues(t, 60*time.Hour, report.result.Duration.SumExact)
	assert.EqualValues(t, 60*time.Hour, report.result.ChildBuckets[0].Duration.SumExact)

	// splitting into p1 and p2 must result in only buckets for each
	report = NewBucketReport(frames.Copy(), Config{
		ProjectIDs: []string{p1.ID, p2.ID},
		Splitting:  []SplitOperation{SplitByProject},
	}, ctx)
	report.Update()
	require.NotNil(t, report.result, "expected one top-level group with 30 frames")
	assert.EqualValues(t, 2, len(report.result.ChildBuckets))
	assert.EqualValues(t, 20, report.result.FrameCount)
	assert.EqualValues(t, 40*time.Hour, report.result.Duration.SumExact)
	assert.EqualValues(t, 20*time.Hour, report.result.ChildBuckets[0].Duration.SumExact)
	assert.EqualValues(t, 20*time.Hour, report.result.ChildBuckets[1].Duration.SumExact)

	// splitting into pTop and p2 must result in only buckets for each, pTop's must contain frames of p2, too
	report = NewBucketReport(frames.Copy(), Config{
		ProjectIDs: []string{pTop.ID, p1.ID},
		Splitting:  []SplitOperation{SplitByProject},
	}, ctx)
	report.Update()
	require.NotNil(t, report.result, "expected one top-level group with 30 frames")
	assert.EqualValues(t, 2, len(report.result.ChildBuckets))
	assert.EqualValues(t, 30, report.result.FrameCount)
	assert.EqualValues(t, 60*time.Hour, report.result.Duration.SumExact)
	assert.EqualValues(t, 40*time.Hour, report.result.ChildBuckets[0].Duration.SumExact)
	assert.EqualValues(t, 20*time.Hour, report.result.ChildBuckets[1].Duration.SumExact)

	// with splitting into frame's project, i.e. 3 projects
	report = NewBucketReport(frames.Copy(), Config{
		ProjectIDs:         []string{pTop.ID},
		Splitting:          []SplitOperation{SplitByProject},
		IncludeSubprojects: true,
	}, ctx)
	report.Update()
	require.NotNil(t, report.result, "expected one top-level group with all 30 frames")
	assert.EqualValues(t, 30, report.result.FrameCount)
	require.EqualValues(t, 3, len(report.result.ChildBuckets), "expected 3 project buckets")

	bucket1 := report.result.ChildBuckets[0]
	bucket2 := report.result.ChildBuckets[1]
	bucket3 := report.result.ChildBuckets[2]

	assert.EqualValues(t, 10, bucket1.FrameCount)
	assert.EqualValues(t, 10, bucket2.FrameCount)
	assert.EqualValues(t, 10, bucket3.FrameCount)

	assert.EqualValues(t, 60*time.Hour, report.result.Duration.SumExact)
	assert.EqualValues(t, 20*time.Hour, bucket1.Duration.SumExact, "10h expected for each project bucket")
	assert.EqualValues(t, 20*time.Hour, bucket2.Duration.SumExact, "10h expected for each project bucket")
	assert.EqualValues(t, 20*time.Hour, bucket3.Duration.SumExact, "10h expected for each project bucket")

	// report without project filter must create buckets for all, even if it's empty
	report = NewBucketReport(frames.Copy(), Config{
		Splitting:          []SplitOperation{SplitByProject},
		IncludeSubprojects: true,
		ShowEmpty:          true,
	}, ctx)
	report.Update()
	require.NotNil(t, report.result, "expected one top-level group with all 30 frames")
	assert.EqualValues(t, 30, report.result.FrameCount)
	require.EqualValues(t, 4, len(report.result.ChildBuckets), "expected 3 project buckets")
}

// test fix for a bug where the month split wasn't visible
func TestReportSplitYearProjectMonth(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	pTop, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1")
	require.NoError(t, err)
	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1", "child1")
	require.NoError(t, err)
	p2, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1", "child2")
	require.NoError(t, err)

	// adding 3*16*2 hours = 96 hours, 32 for each project spread across 16 months
	frames := model.NewEmptyFrameList()
	start := newDate(2017, time.May, 10, 10, 0).UTC()
	end := newDate(2017, time.May, 10, 12, 0).UTC()
	for i := 0; i < 16; i++ {
		newStart := start.AddDate(0, i, 0)
		newEnd := end.AddDate(0, i, 0)

		frames.Append(&model.Frame{Start: &newStart, End: &newEnd, ProjectId: pTop.ID})
		frames.Append(&model.Frame{Start: &newStart, End: &newEnd, ProjectId: p1.ID})
		frames.Append(&model.Frame{Start: &newStart, End: &newEnd, ProjectId: p2.ID})
	}
	frames.Sort()

	report := NewBucketReport(frames.Copy(), Config{
		Splitting:          []SplitOperation{SplitByYear, SplitByProject, SplitByMonth},
		IncludeSubprojects: true,
	}, ctx)
	report.Update()
	require.NotNil(t, report.result, "expected one top-level group with 30 frames")
	assert.EqualValues(t, 48, report.result.FrameCount, "expected 3*16 frames in total")
	assert.EqualValues(t, 96*time.Hour, report.result.Duration.SumExact)

	require.EqualValues(t, 3, report.result.Depth(), "expected bucket hierarchy year > project > month")
	assert.EqualValues(t, 2, len(report.result.ChildBuckets), "expected 2 buckets mapping years 2017 and 2018: "+report.result.ChildBuckets[0].Title())
	for _, year := range report.result.ChildBuckets {
		fmt.Println(year.Title())
		require.EqualValues(t, 3, len(year.ChildBuckets), "expected 3 buckets, one for each project")

		for _, project := range year.ChildBuckets {
			require.EqualValues(t, 8, len(project.ChildBuckets), "expected 8 month buckets (16 spreach equally on 2017 and 2018)")

			for _, month := range project.ChildBuckets {
				require.Empty(t, month.ChildBuckets)
			}
		}
	}
}

// test fix for a where splitting into project and month wasn't rendering a matrix
func TestReportSplitProjectMonthMatrix(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	pTop, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1")
	require.NoError(t, err)
	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1", "child1")
	require.NoError(t, err)

	// adding 10*2 hours = 20 hours on project p1
	frames := model.NewEmptyFrameList()
	start := newDate(2017, time.May, 10, 10, 0).UTC()
	end := newDate(2017, time.May, 10, 12, 0).UTC()
	for i := 0; i < 10; i++ {
		newStart := start.AddDate(0, 0, 1)
		newEnd := end.AddDate(0, 0, 1)

		frames.Append(&model.Frame{Start: &newStart, End: &newEnd, ProjectId: p1.ID})
	}
	frames.Sort()

	report := NewBucketReport(frames.Copy(), Config{
		ProjectIDs:         []string{pTop.ID},
		Splitting:          []SplitOperation{SplitByProject, SplitByMonth},
		IncludeSubprojects: true,
		ShowEmpty:          false,
	}, ctx)
	report.Update()
	assert.EqualValues(t, 10, report.result.FrameCount, "expected 3*16 frames in total")
	assert.EqualValues(t, 20*time.Hour, report.result.Duration.SumExact)

	require.EqualValues(t, 2, report.result.Depth(), "expected bucket hierarchy project > month")
	assert.EqualValues(t, 2, len(report.result.ChildBuckets), "expected 2 project bucket (1 empty, 1 full)")
	for _, year := range report.result.ChildBuckets {
		if !year.Empty(){
			require.EqualValues(t, 1, len(year.ChildBuckets), "expected 1 month bucket for the project")
		}
	}
	assert.True(t, IsMatrix(report.result, true))
	assert.False(t, IsMatrix(report.result, false))
}

func newDate(year int, month time.Month, day, hour, minute int) *time.Time {
	date := time.Date(year, month, day, hour, minute, 0, 0, time.Local)
	return &date
}
