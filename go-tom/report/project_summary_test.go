package report

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/test_setup"
)

func Test_ProjectReportTest(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p, err := ctx.Store.AddProject(model.Project{Name: "Project1"});
	require.NoError(t, err)

	// a wednesday
	refDate := time.Date(2018, time.July, 18, 12, 0, 0, 0, time.Local)

	startToday := refDate.Add(-10 * time.Minute)
	endToday := refDate

	startYesterday := refDate.Add(-24 * time.Hour)
	endYesterday := startYesterday.Add(10 * time.Minute)

	startTwoDaysAgo := refDate.Add(-2 * 24 * time.Hour)
	endTwoDaysAgo := startTwoDaysAgo.Add(10 * time.Minute)

	startTwoWeeksAgo := refDate.Add(-2 * 7 * 24 * time.Hour)
	endTwoWeeksAgo := startTwoWeeksAgo.Add(10 * time.Minute)

	startLastMonth := refDate.Add(-31 * 24 * time.Hour)
	endLastMonth := startLastMonth.Add(10 * time.Minute)

	// last month
	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p.ID, Start: &startLastMonth, End: &endLastMonth})
	require.NoError(t, err)

	// two days ago
	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p.ID, Start: &startTwoDaysAgo, End: &endTwoDaysAgo})
	require.NoError(t, err)

	// two weeks ago
	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p.ID, Start: &startTwoWeeksAgo, End: &endTwoWeeksAgo})
	require.NoError(t, err)

	// yesterday
	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p.ID, Start: &startYesterday, End: &endYesterday})
	require.NoError(t, err)

	// today
	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p.ID, Start: &startToday, End: &endToday})
	require.NoError(t, err)

	reports := CreateProjectReports(refDate, false, true, nil, "", ctx)
	require.EqualValues(t, 1, len(reports))
	require.EqualValues(t, 50*time.Minute, reports[p.ID].TrackedYear.Get())
	require.EqualValues(t, 40*time.Minute, reports[p.ID].TrackedMonth.Get())
	require.EqualValues(t, 30*time.Minute, reports[p.ID].TrackedWeek.Get())
	require.EqualValues(t, 10*time.Minute, reports[p.ID].TrackedYesterday.Get())
	require.EqualValues(t, 10*time.Minute, reports[p.ID].TrackedDay.Get())

	require.EqualValues(t, 50*time.Minute, reports[p.ID].TrackedTotalYear.Get())
	require.EqualValues(t, 40*time.Minute, reports[p.ID].TrackedTotalMonth.Get())
	require.EqualValues(t, 30*time.Minute, reports[p.ID].TrackedTotalWeek.Get())
	require.EqualValues(t, 10*time.Minute, reports[p.ID].TrackedTotalYesterday.Get())
	require.EqualValues(t, 10*time.Minute, reports[p.ID].TrackedTotalDay.Get())

	// now include active frames, no end defined here
	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p.ID, Start: &startToday})
	require.NoError(t, err)

	// adds 15 minutes for our active frame
	activeEnd := startToday.Add(15 * time.Minute)

	reports = CreateProjectReports(refDate, false, true, &activeEnd, "", ctx)
	require.EqualValues(t, 1, len(reports))
	require.EqualValues(t, 15*time.Minute+50*time.Minute, reports[p.ID].TrackedYear.Get())
	require.EqualValues(t, 15*time.Minute+40*time.Minute, reports[p.ID].TrackedMonth.Get())
	require.EqualValues(t, 15*time.Minute+30*time.Minute, reports[p.ID].TrackedWeek.Get())
	require.EqualValues(t, 10*time.Minute, reports[p.ID].TrackedYesterday.Get())
	require.EqualValues(t, 15*time.Minute+10*time.Minute, reports[p.ID].TrackedDay.Get())

	require.EqualValues(t, 15*time.Minute+50*time.Minute, reports[p.ID].TrackedTotalYear.Get())
	require.EqualValues(t, 15*time.Minute+40*time.Minute, reports[p.ID].TrackedTotalMonth.Get())
	require.EqualValues(t, 15*time.Minute+30*time.Minute, reports[p.ID].TrackedTotalWeek.Get())
	require.EqualValues(t, 10*time.Minute, reports[p.ID].TrackedTotalYesterday.Get())
	require.EqualValues(t, 15*time.Minute+10*time.Minute, reports[p.ID].TrackedTotalDay.Get())
}

func Test_MultipleDays(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p, err := ctx.Store.AddProject(model.Project{Name: "Project1"});
	require.NoError(t, err)

	// a wednesday, 12 am
	refDate := time.Date(2018, time.July, 18, 12, 0, 0, 0, time.Local)

	// 12 hours in the day, 11 in the next
	startToday := refDate
	endToday := startToday.Add(23 * time.Hour)

	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p.ID, Start: &startToday, End: &endToday})
	require.NoError(t, err)

	reports := CreateProjectReports(refDate, false, true, nil, "", ctx)
	require.EqualValues(t, 1, len(reports))
	require.EqualValues(t, 23*time.Hour, reports[p.ID].TrackedYear.Get())
	require.EqualValues(t, 23*time.Hour, reports[p.ID].TrackedMonth.Get())
	require.EqualValues(t, 23*time.Hour, reports[p.ID].TrackedWeek.Get())
	require.EqualValues(t, 12*time.Hour, reports[p.ID].TrackedDay.Get())
	require.EqualValues(t, 0*time.Hour, reports[p.ID].TrackedYesterday.Get())

	require.EqualValues(t, 23*time.Hour, reports[p.ID].TrackedTotalYear.Get())
	require.EqualValues(t, 23*time.Hour, reports[p.ID].TrackedTotalMonth.Get())
	require.EqualValues(t, 23*time.Hour, reports[p.ID].TrackedTotalWeek.Get())
	require.EqualValues(t, 12*time.Hour, reports[p.ID].TrackedTotalDay.Get())
	require.EqualValues(t, 0*time.Hour, reports[p.ID].TrackedTotalYesterday.Get())
}

func Test_MultipleDaysNoArchived(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p, err := ctx.Store.AddProject(model.Project{Name: "Project1"});
	require.NoError(t, err)

	// a wednesday, 12 am
	refDate := time.Date(2018, time.July, 18, 12, 0, 0, 0, time.Local)

	// 12 hours in the day, 11 in the next
	startToday := refDate
	endToday := startToday.Add(23 * time.Hour)

	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p.ID, Start: &startToday, End: &endToday, Archived: true})
	require.NoError(t, err)

	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p.ID, Start: &startToday, End: &endToday, Archived: false})
	require.NoError(t, err)

	// exclude archived
	reports := CreateProjectReports(refDate, false, false, nil, "", ctx)
	require.EqualValues(t, 1, len(reports))
	require.EqualValues(t, 23*time.Hour, reports[p.ID].TrackedYear.Get())
	require.EqualValues(t, 23*time.Hour, reports[p.ID].TrackedMonth.Get())
	require.EqualValues(t, 23*time.Hour, reports[p.ID].TrackedWeek.Get())
	require.EqualValues(t, 12*time.Hour, reports[p.ID].TrackedDay.Get())
	require.EqualValues(t, 0*time.Hour, reports[p.ID].TrackedYesterday.Get())

	require.EqualValues(t, 23*time.Hour, reports[p.ID].TrackedTotalYear.Get())
	require.EqualValues(t, 23*time.Hour, reports[p.ID].TrackedTotalMonth.Get())
	require.EqualValues(t, 23*time.Hour, reports[p.ID].TrackedTotalWeek.Get())
	require.EqualValues(t, 12*time.Hour, reports[p.ID].TrackedTotalDay.Get())
	require.EqualValues(t, 0*time.Hour, reports[p.ID].TrackedTotalYesterday.Get())
}

func Test_ActiveFrame(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p, err := ctx.Store.AddProject(model.Project{Name: "Project1"});
	require.NoError(t, err)

	// a wednesday, 12 am
	refDate := time.Date(2018, time.July, 18, 12, 0, 0, 0, time.Local)

	// 12 hours in the day, 11 in the next
	start := refDate.Add(-10 * time.Minute)
	end := refDate

	startActive := refDate.Add(-10 * time.Minute)

	// one closed frame, duration of 10 minutes
	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p.ID, Start: &start, End: &end})
	require.NoError(t, err)

	// one active frame, duration of 10 minutes so far
	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p.ID, Start: &startActive, End: nil})
	require.NoError(t, err)

	// exclude archived
	reports := CreateProjectReports(refDate, false, false, &refDate, "", ctx)
	require.EqualValues(t, 1, len(reports))
	require.EqualValues(t, 20*time.Minute, reports[p.ID].TrackedYear.Get())
	require.EqualValues(t, 20*time.Minute, reports[p.ID].TrackedMonth.Get())
	require.EqualValues(t, 20*time.Minute, reports[p.ID].TrackedWeek.Get())
	require.EqualValues(t, 20*time.Minute, reports[p.ID].TrackedDay.Get())
	require.EqualValues(t, 0*time.Minute, reports[p.ID].TrackedYesterday.Get())

	require.EqualValues(t, 20*time.Minute, reports[p.ID].TrackedTotalYear.Get())
	require.EqualValues(t, 20*time.Minute, reports[p.ID].TrackedTotalMonth.Get())
	require.EqualValues(t, 20*time.Minute, reports[p.ID].TrackedTotalWeek.Get())
	require.EqualValues(t, 20*time.Minute, reports[p.ID].TrackedTotalDay.Get())
	require.EqualValues(t, 0*time.Minute, reports[p.ID].TrackedTotalYesterday.Get())
}
