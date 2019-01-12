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

	reports := CreateProjectReports(refDate, false, ctx)
	require.EqualValues(t, 1, len(reports))
	require.EqualValues(t, 50*time.Minute, reports[p.ID].TrackedYear)
	require.EqualValues(t, 40*time.Minute, reports[p.ID].TrackedMonth)
	require.EqualValues(t, 30*time.Minute, reports[p.ID].TrackedWeek)
	require.EqualValues(t, 10*time.Minute, reports[p.ID].TrackedDay)

	require.EqualValues(t, 50*time.Minute, reports[p.ID].TotalTrackedYear)
	require.EqualValues(t, 40*time.Minute, reports[p.ID].TotalTrackedMonth)
	require.EqualValues(t, 30*time.Minute, reports[p.ID].TotalTrackedWeek)
	require.EqualValues(t, 10*time.Minute, reports[p.ID].TotalTrackedDay)
}
