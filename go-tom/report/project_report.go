package report

import (
	"time"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateUtil"
	"github.com/jansorg/tom/go-tom/model"
)

type ProjectSummary struct {
	ProjectID    string
	TrackedYear  time.Duration
	TrackedMonth time.Duration
	TrackedWeek  time.Duration
	TrackedDay   time.Duration
}

func CreateProjectReports(frames *model.FrameList, referenceDay time.Time, ctx *context.GoTimeContext) map[string]ProjectSummary {
	year := dateUtil.NewYearRange(referenceDay, ctx.Locale)
	week := dateUtil.NewWeekRange(referenceDay, ctx.Locale)
	month := dateUtil.NewMonthRange(referenceDay, ctx.Locale)
	day := dateUtil.NewDayRange(referenceDay, ctx.Locale)

	result := map[string]ProjectSummary{}

	// frames.FilterByDateRange(year, false)

	// update years
	report := NewBucketReport(frames, ctx)
	report.SplitOperations = []SplitOperation{SplitByProject}
	report.FilterRange = year
	report.Update()
	for _, r := range report.Result.Results {
		projectID := r.SplitBy.(string)
		report := result[projectID]
		report.ProjectID = projectID
		report.TrackedYear = r.ExactDuration
		result[projectID] = report
	}

	// update months
	report.FilterRange = month
	report.Update()
	for _, r := range report.Result.Results {
		projectID := r.SplitBy.(string)
		report := result[projectID]
		report.TrackedMonth = r.ExactDuration
		result[projectID] = report
	}

	// update week
	report.FilterRange = week
	report.Update()
	for _, r := range report.Result.Results {
		projectID := r.SplitBy.(string)
		report := result[projectID]
		report.TrackedWeek = r.ExactDuration
		result[projectID] = report
	}

	// update day
	report.FilterRange = day
	report.Update()
	for _, r := range report.Result.Results {
		projectID := r.SplitBy.(string)
		report := result[projectID]
		report.TrackedDay = r.ExactDuration
		result[projectID] = report
	}

	return result
}
