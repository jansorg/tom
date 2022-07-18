package report

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/dateTime"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/money"
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

	start := newLocalDate(2018, time.March, 10, 10, 0)
	end := newLocalDate(2018, time.March, 10, 12, 0)

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

	start := newLocalDate(2018, time.March, 10, 10, 0)
	end := newLocalDate(2018, time.March, 10, 12, 0)

	frameList := []*model.Frame{{Start: start, End: end}}
	report := NewBucketReport(model.NewFrameList(frameList), Config{}, ctx)
	report.Update()
	assert.EqualValues(t, 1, report.result.FrameCount)
	assert.EqualValues(t, 2*time.Hour, report.result.Duration.Get())
	assert.EqualValues(t, 2*time.Hour, report.result.Duration.GetExact())
	assert.EqualValues(t, start, report.result.TrackedDateRange().Start)
	assert.EqualValues(t, end, report.result.TrackedDateRange().End)
	assert.EqualValues(t, frameList, report.source.Frames())

	tracked := report.result.GetDailyTracked()
	require.NotNil(t, tracked)
	assert.EqualValues(t, 2*time.Hour, tracked.Min())
	assert.EqualValues(t, 2*time.Hour, tracked.Max())
	assert.EqualValues(t, 2*time.Hour, tracked.Avg())
	assert.EqualValues(t, 1, tracked.DistinctRanges())

	unTracked := report.result.GetDailyUnTracked()
	require.NotNil(t, unTracked)
	assert.EqualValues(t, 0, unTracked.Min())
	assert.EqualValues(t, 0, unTracked.Max())
	assert.EqualValues(t, 0, unTracked.Avg())
	assert.EqualValues(t, 1, unTracked.DistinctRanges())
}

func TestReportSplitYear(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	// two hours, 10th of March 2018
	start := newLocalDate(2018, time.March, 10, 10, 0)
	end := newLocalDate(2018, time.March, 10, 12, 0)

	// one hour. 9th of March 2019
	start2 := newLocalDate(2019, time.March, 10, 9, 0)
	end2 := newLocalDate(2019, time.March, 10, 10, 0)

	frameList := []*model.Frame{
		{Start: start2, End: end2},
		{Start: start, End: end},
	}

	report := NewBucketReport(model.NewSortedFrameList(frameList),
		Config{TimezoneName: NewTimezoneNameLocal(), Splitting: []SplitOperation{SplitByYear}},
		ctx)
	report.Update()

	require.NotNil(t, report.result, "expected one top-level group (containing two years)")
	require.EqualValues(t, 2, len(report.result.ChildBuckets), "expected a sub-report for each year")

	assert.EqualValues(t, 2, report.result.FrameCount)
	assert.EqualValues(t, 3*time.Hour, report.result.Duration.Get())
	assert.EqualValues(t, 3*time.Hour, report.result.Duration.GetExact())
	assert.EqualValues(t, newLocalDate(2018, time.January, 1, 0, 0), report.result.DateRange().Start)
	assert.EqualValues(t, newLocalDate(2020, time.January, 1, 0, 0), report.result.DateRange().End)
	assert.EqualValues(t, start, report.result.TrackedDateRange().Start)
	assert.EqualValues(t, end2, report.result.TrackedDateRange().End)
	assert.EqualValues(t, frameList, report.source.Frames())

	firstYear := report.result.ChildBuckets[0]
	assert.EqualValues(t, newLocalDate(2018, time.January, 1, 0, 0), firstYear.DateRange().Start)
	assert.EqualValues(t, newLocalDate(2019, time.January, 1, 0, 0), firstYear.DateRange().End)
	assert.EqualValues(t, 1, firstYear.FrameCount)
	assert.EqualValues(t, 2*time.Hour, firstYear.Duration.Get())
	assert.EqualValues(t, 2*time.Hour, firstYear.Duration.GetExact())

	secondYear := report.result.ChildBuckets[1]
	assert.EqualValues(t, newLocalDate(2019, time.January, 1, 0, 0), secondYear.DateRange().Start)
	assert.EqualValues(t, newLocalDate(2020, time.January, 1, 0, 0), secondYear.DateRange().End)
	assert.EqualValues(t, 1, secondYear.FrameCount)
	assert.EqualValues(t, 1*time.Hour, secondYear.Duration.Get())
	assert.EqualValues(t, 1*time.Hour, secondYear.Duration.GetExact())
}

func TestReportDateRanges(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	// two hours
	start := newLocalDate(2018, time.March, 10, 10, 0)
	end := newLocalDate(2018, time.March, 10, 12, 0)

	// one hour
	start2 := newLocalDate(2019, time.March, 10, 9, 0)
	end2 := newLocalDate(2019, time.March, 10, 10, 0)

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
	start := newLocalDate(2018, time.March, 10, 10, 0).UTC()
	end := newLocalDate(2018, time.March, 10, 12, 0).UTC()

	// one hour, UTC+2
	utc2, err := time.LoadLocation("Europe/Berlin")
	require.NoError(t, err)
	start2 := newLocalDate(2018, time.March, 10, 9, 0).In(utc2)
	end2 := newLocalDate(2018, time.March, 10, 10, 0).In(utc2)

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
	start := newLocalDate(2018, time.January, 10, 10, 0).UTC()
	end := newLocalDate(2018, time.January, 10, 12, 0).UTC()

	// one hour, UTC+2
	utc2, err := time.LoadLocation("Europe/Berlin")
	require.NoError(t, err)
	start2 := newLocalDate(2018, time.January, 10, 9, 0).In(utc2)
	end2 := newLocalDate(2018, time.January, 10, 10, 0).In(utc2)

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
	start := newLocalDate(2017, time.January, 10, 10, 0).UTC()
	end := newLocalDate(2017, time.January, 10, 12, 0).UTC()
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
	start := newLocalDate(2017, time.May, 10, 10, 0).UTC()
	end := newLocalDate(2017, time.May, 10, 12, 0).UTC()
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
	start := newLocalDate(2017, time.May, 10, 10, 0).UTC()
	end := newLocalDate(2017, time.May, 10, 12, 0).UTC()
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
	assert.EqualValues(t, 10, report.result.FrameCount, "expected 10 frames in total")
	assert.EqualValues(t, 20*time.Hour, report.result.Duration.SumExact)

	require.EqualValues(t, 2, report.result.Depth(), "expected bucket hierarchy project > month")
	assert.EqualValues(t, 2, len(report.result.ChildBuckets), "expected 2 project bucket (1 empty, 1 full)")
	for _, year := range report.result.ChildBuckets {
		if !year.Empty() {
			require.EqualValues(t, 1, len(year.ChildBuckets), "expected 1 month bucket for the project")
		}
	}
	assert.True(t, IsMatrix(report.result, true))
	assert.False(t, IsMatrix(report.result, false))
}

func TestReportWithoutArchived(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	pTop, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1")
	require.NoError(t, err)
	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1", "child1")
	require.NoError(t, err)

	start := newLocalDate(2017, time.May, 10, 10, 0).UTC()
	end := newLocalDate(2017, time.May, 10, 12, 0).UTC()

	// adding 5*2 hours = 10 hours on project p1, UNARCHIVED
	frames := model.NewEmptyFrameList()
	for i := 0; i < 5; i++ {
		newStart := start.AddDate(0, 0, 1)
		newEnd := end.AddDate(0, 0, 1)

		frames.Append(&model.Frame{Start: &newStart, End: &newEnd, ProjectId: p1.ID})
	}

	// adding 5*2 hours = 10 hours on project p1, ARCHIVED
	for i := 0; i < 5; i++ {
		newStart := start.AddDate(0, 0, 1)
		newEnd := end.AddDate(0, 0, 1)

		frames.Append(&model.Frame{Start: &newStart, End: &newEnd, ProjectId: p1.ID, Archived: true})
	}

	frames.Sort()

	report := NewBucketReport(frames.Copy(), Config{
		ProjectIDs:      []string{pTop.ID},
		Splitting:       []SplitOperation{SplitByProject, SplitByMonth},
		ShowEmpty:       false,
		IncludeArchived: false,
	}, ctx)
	report.Update()
	assert.EqualValues(t, 5, report.result.FrameCount, "expected 5 frames in total")
	assert.EqualValues(t, 10*time.Hour, report.result.Duration.SumExact)
}

func TestSalesStats(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	// pTop with 100 EUR / hour
	pTop, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1")
	require.NoError(t, err)
	pTop.SetHourlyRate(money.NewMoney(100*100, "EUR"))

	// p1 with 50 EUR / hour
	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1", "child1")
	require.NoError(t, err)
	p1.SetHourlyRate(money.NewMoney(50*100, "EUR"))

	// p2 with 75 USD / hour
	p2, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project1", "child2")
	require.NoError(t, err)
	p2.SetHourlyRate(money.NewMoney(75*100, "USD"))

	// 2 hours
	start := newLocalDate(2017, time.May, 10, 10, 0).UTC()
	end := newLocalDate(2017, time.May, 10, 12, 0).UTC()

	// adding 5*2 hours = 10 hours each on projects pTop, p1, p2
	frames := model.NewEmptyFrameList()
	for i := 0; i < 5; i++ {
		newStart := start.AddDate(0, 0, 1)
		newEnd := end.AddDate(0, 0, 1)
		frames.Append(&model.Frame{Start: &newStart, End: &newEnd, ProjectId: pTop.ID})
		frames.Append(&model.Frame{Start: &newStart, End: &newEnd, ProjectId: p1.ID})
		frames.Append(&model.Frame{Start: &newStart, End: &newEnd, ProjectId: p2.ID})
	}
	frames.Sort()

	report := NewBucketReport(frames.Copy(), Config{
		ProjectIDs: []string{pTop.ID},
		Splitting:  []SplitOperation{SplitByProject, SplitByMonth},
	}, ctx)
	report.Update()
	assert.EqualValues(t, 15, report.result.FrameCount, "expected 5 frames in total")
	assert.EqualValues(t, 30*time.Hour, report.result.Duration.SumExact)

	// EUR sales = p1 + p2 = 5*2hours * 100 EUR + 5*2hours * 50 EUR = 1500 EUR
	// USD sales = p3 = 5*2 hours * 75 USD = 750 USD
	assert.EqualValues(t, "€1,500.00", report.Result().Sales.values["EUR"].String())
	assert.EqualValues(t, "$750.00", report.Result().Sales.values["USD"].String())
}

func TestReportTimeFilter(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top")
	require.NoError(t, err)

	// different time zones
	start := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.FixedZone("UTC-2", -2*60*60))
	end := time.Date(2019, time.January, 2, 0, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))

	report := NewBucketReport(model.NewFrameList([]*model.Frame{}), Config{
		ProjectIDs:         []string{p1.ID},
		IncludeSubprojects: true,
		Splitting:          []SplitOperation{SplitByProject, SplitByMonth},
		DateFilterRange:    dateTime.NewDateRange(&start, &end, ctx.Locale),
		TimezoneName:       NewTimezoneNameUTC(),
	}, ctx)
	report.Update()
	assert.NotNil(t, report.config.DateFilterRange.Start, "start must be not non nil")
	assert.EqualValues(t, "2019-01-01 00:00:00 -0200 UTC-2 - 2019-01-02 00:00:00 +0800 UTC+8", report.config.DateFilterRange.String())
	assert.EqualValues(t, "2019-01-01 02:00:00 +0000 UTC - 2019-01-01 16:00:00 +0000 UTC", report.result.config.DateFilterRange.String())

	// different target zone
	report = NewBucketReport(model.NewFrameList([]*model.Frame{}), Config{
		ProjectIDs:         []string{p1.ID},
		IncludeSubprojects: true,
		Splitting:          []SplitOperation{SplitByProject, SplitByMonth},
		DateFilterRange:    dateTime.NewDateRange(&start, &end, ctx.Locale),
		TimezoneName:       NewTimezoneName(time.FixedZone("UTC+8", 8*60*60)),
	}, ctx)
	report.Update()
	assert.NotNil(t, report.config.DateFilterRange.Start, "start must be not non nil")
	assert.EqualValues(t, "2019-01-01 10:00:00 +0800 UTC+8 - 2019-01-02 00:00:00 +0800 UTC+8", report.result.config.DateFilterRange.String())
}

// test for https://github.com/jansorg/tom-ui/issues/91
// make sure that entries which overlap the filter range are properly handled
func TestReportTimeFilterOverlap(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top")
	require.NoError(t, err)

	start := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(1 * time.Hour)

	// the frame start 10 minutes before and end 10 minutes after the filter range
	frameStart := start.Add(-10 * time.Minute)
	frameEnd := end.Add(10 * time.Minute)

	frames := model.NewEmptyFrameList()
	frames.Append(&model.Frame{Start: &frameStart, End: &frameEnd, ProjectId: p1.ID})

	// with filter
	report := NewBucketReport(frames, Config{
		ProjectIDs:         []string{p1.ID},
		IncludeSubprojects: true,
		Splitting:          []SplitOperation{SplitByProject},
		DateFilterRange:    dateTime.NewDateRange(&start, &end, ctx.Locale),
		TimezoneName:       NewTimezoneNameUTC(),
	}, ctx)
	report.Update()
	assert.EqualValues(t, 1*time.Hour, report.Result().Duration.GetExact(), "1 hour max range expected for overlapping entries")
}

// test for https://github.com/jansorg/tom-ui/issues/91
// make sure that entries which overlap the filter range are properly handled
func TestReportTimeFilterNoOverlap(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top")
	require.NoError(t, err)

	start := time.Date(2019, time.January, 15, 12, 0, 0, 0, time.UTC)
	end := start.Add(1 * time.Hour)

	// the frame start 10 minutes before and end 10 minutes after the filter range
	frameStart := start.Add(-10 * time.Minute)
	frameEnd := end.Add(10 * time.Minute)

	frames := model.NewEmptyFrameList()
	frames.Append(&model.Frame{Start: &frameStart, End: &frameEnd, ProjectId: p1.ID})

	// without broader filter, i.e. without overlapping
	beforeStart := start.AddDate(0, 0, -1)
	afterEnd := end.AddDate(0, 0, 1)

	report := NewBucketReport(frames, Config{
		ProjectIDs:         []string{p1.ID},
		IncludeSubprojects: true,
		Splitting:          []SplitOperation{SplitByProject},
		DateFilterRange:    dateTime.NewDateRange(&beforeStart, &afterEnd, ctx.Locale),
		TimezoneName:       NewTimezoneNameUTC(),
	}, ctx)
	report.Update()
	assert.EqualValues(t, 1*time.Hour+20*time.Minute, report.Result().Duration.GetExact(), "expected full duration for non-overlapping frames")
}

func TestReportTimeFilterOverlapMultipleMonths(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top")
	require.NoError(t, err)

	// January 2019
	start := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)

	// the frame start 10 minutes before and end 10 minutes after the filter range
	frameStart := start.Add(-10 * time.Minute)
	frameEnd := end.Add(10 * time.Minute)

	frames := model.NewEmptyFrameList()
	frames.Append(&model.Frame{Start: &frameStart, End: &frameEnd, ProjectId: p1.ID})

	// no filter, but the month buckets must properly split the single frame
	report := NewBucketReport(frames, Config{
		ProjectIDs:         []string{p1.ID},
		IncludeSubprojects: true,
		Splitting:          []SplitOperation{SplitByDay},
		TimezoneName:       NewTimezoneNameUTC(),
	}, ctx)
	report.Update()

	assert.EqualValues(t, 24*time.Hour+20*time.Minute, report.Result().Duration.GetExact())

	subBuckets := report.Result().ChildBuckets

	require.EqualValues(t, 3, len(subBuckets))

	assert.EqualValues(t, 10*time.Minute, subBuckets[0].Duration.GetExact())
	assert.EqualValues(t, 1, subBuckets[0].FrameCount)
	assert.EqualValues(t, 10*time.Minute, (*subBuckets[0].Frames)[0].Duration())

	assert.EqualValues(t, 24*time.Hour, subBuckets[1].Duration.GetExact())
	assert.EqualValues(t, 1, subBuckets[1].FrameCount)

	assert.EqualValues(t, 10*time.Minute, subBuckets[2].Duration.GetExact())
	assert.EqualValues(t, 1, subBuckets[2].FrameCount)
}

func TestReportWithActiveFrame(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top")
	require.NoError(t, err)

	// January 2019
	start := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(1 * time.Hour)

	_, err = ctx.Store.AddFrame(model.Frame{Start: &start, End: &end, ProjectId: p1.ID})
	require.NoError(t, err)
	_, err = ctx.Store.AddFrame(model.Frame{Start: &start, End: nil, ProjectId: p1.ID})
	require.NoError(t, err)

	// no filter, but the month buckets must properly split the single frame
	frames := ctx.Store.Frames()
	report := NewBucketReport(&frames, Config{
		IncludeSubprojects: false,
		Splitting:          []SplitOperation{SplitByDay},
		TimezoneName:       NewTimezoneNameUTC(),
	}, ctx)
	report.Update()

	assert.EqualValues(t, 1*time.Hour, report.Result().Duration.GetExact())
	require.EqualValues(t, 1, len(report.Result().ChildBuckets))
}

func newLocalDate(year int, month time.Month, day, hour, minute int) *time.Time {
	date := time.Date(year, month, day, hour, minute, 0, 0, time.Local)
	return &date
}
