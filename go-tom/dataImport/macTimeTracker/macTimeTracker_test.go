package macTimeTracker

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/report"
	"github.com/jansorg/tom/go-tom/test_setup"
)

func Test_ImportCSV(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	csvPath, err := filepath.Abs("test/Time Tracker Data.csv")
	require.NoError(t, err)

	created, err := NewImporter().Import(csvPath, ctx)
	require.NoError(t, err)
	assert.EqualValues(t, 24, created.CreatedProjects)
	assert.EqualValues(t, 6112, created.CreatedFrames)

	projectReports := report.CreateProjectReports(time.Date(2019, 1, 21, 0, 0, 0, 0, time.UTC), false, ctx)

	acme, err := ctx.Query.ProjectByFullName([]string{"ACME Corp."})
	require.NoError(t, err)

	contacts, err := ctx.Query.ProjectByFullName([]string{"Contacts"})
	require.NoError(t, err)

	goodBurgers, err := ctx.Query.ProjectByFullName([]string{"Good Burger"})
	require.NoError(t, err)

	importInc, err := ctx.Query.ProjectByFullName([]string{"Import"})
	require.NoError(t, err)

	newYork, err := ctx.Query.ProjectByFullName([]string{"The New York Inquirer"})
	require.NoError(t, err)

	wonka, err := ctx.Query.ProjectByFullName([]string{"Wonka Industries"})
	require.NoError(t, err)

	assertSummary(t, d(682, 6, 15), projectReports[acme.ID])
	assertSummary(t, d(153, 17, 9), projectReports[contacts.ID])
	assertSummary(t, d(528, 15, 16), projectReports[goodBurgers.ID])
	assertSummary(t, d(45, 50, 42), projectReports[importInc.ID])
	assertSummary(t, d(2230, 20, 23), projectReports[newYork.ID])
	assertSummary(t, d(823, 2, 24), projectReports[wonka.ID])

}

func assertSummary(t *testing.T, expected time.Duration, summary *report.ProjectSummary) bool {
	return assert.EqualValues(t, expected, summary.TotalTrackedAll, "unexpected duration "+summary.TotalTrackedAll.String())
}

func d(hours, minutes int64, seconds int64) time.Duration {
	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
}
