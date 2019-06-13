package report

import (
	"time"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateTime"
	"github.com/jansorg/tom/go-tom/model"
)

func NewProjectSummary(year *dateTime.DateRange, month *dateTime.DateRange, week *dateTime.DateRange, yesterday *dateTime.DateRange, day *dateTime.DateRange, refTime *time.Time, project *model.Project) *ProjectSummary {
	summary := &ProjectSummary{
		Project:               project,
		TrackedAll:            dateTime.NewDurationSum(),
		TrackedTotalAll:       dateTime.NewDurationSum(),
		TrackedYear:           dateTime.NewDurationSumFiltered(year, refTime),
		TrackedTotalYear:      dateTime.NewDurationSumFiltered(year, refTime),
		TrackedMonth:          dateTime.NewDurationSumFiltered(month, refTime),
		TrackedTotalMonth:     dateTime.NewDurationSumFiltered(month, refTime),
		TrackedWeek:           dateTime.NewDurationSumFiltered(week, refTime),
		TrackedTotalWeek:      dateTime.NewDurationSumFiltered(week, refTime),
		TrackedYesterday:      dateTime.NewDurationSumFiltered(yesterday, refTime),
		TrackedTotalYesterday: dateTime.NewDurationSumFiltered(yesterday, refTime),
		TrackedDay:            dateTime.NewDurationSumFiltered(day, refTime),
		TrackedTotalDay:       dateTime.NewDurationSumFiltered(day, refTime),
	}

	summary.trackedAll = []*dateTime.DurationSum{
		summary.TrackedAll,
		summary.TrackedYear,
		summary.TrackedMonth,
		summary.TrackedWeek,
		summary.TrackedYesterday,
		summary.TrackedDay,
	}

	summary.totalTrackedAll = []*dateTime.DurationSum{
		summary.TrackedTotalAll,
		summary.TrackedTotalYear,
		summary.TrackedTotalMonth,
		summary.TrackedTotalWeek,
		summary.TrackedTotalYesterday,
		summary.TrackedTotalDay,
	}

	return summary
}

type ProjectSummary struct {
	Project *model.Project

	trackedAll       []*dateTime.DurationSum
	TrackedAll       *dateTime.DurationSum
	TrackedYear      *dateTime.DurationSum
	TrackedMonth     *dateTime.DurationSum
	TrackedWeek      *dateTime.DurationSum
	TrackedYesterday *dateTime.DurationSum
	TrackedDay       *dateTime.DurationSum

	totalTrackedAll       []*dateTime.DurationSum
	TrackedTotalAll       *dateTime.DurationSum
	TrackedTotalYear      *dateTime.DurationSum
	TrackedTotalMonth     *dateTime.DurationSum
	TrackedTotalWeek      *dateTime.DurationSum
	TrackedTotalYesterday *dateTime.DurationSum
	TrackedTotalDay       *dateTime.DurationSum
}

func (p *ProjectSummary) add(frame *model.Frame) {
	for _, r := range p.trackedAll {
		r.AddStartEndP(frame.Start, frame.End)
	}
	p.addTotal(frame)
}

func (p *ProjectSummary) addTotal(frame *model.Frame) {
	for _, r := range p.totalTrackedAll {
		r.AddStartEndP(frame.Start, frame.End)
	}
}

func (p *ProjectSummary) Add(v *ProjectSummary) {
	p.TrackedAll.AddSum(v.TrackedAll)
	p.TrackedYear.AddSum(v.TrackedYear)
	p.TrackedMonth.AddSum(v.TrackedMonth)
	p.TrackedWeek.AddSum(v.TrackedWeek)
	p.TrackedYesterday.AddSum(v.TrackedYesterday)
	p.TrackedDay.AddSum(v.TrackedDay)

	p.TrackedTotalAll.AddSum(v.TrackedTotalAll)
	p.TrackedTotalYear.AddSum(v.TrackedTotalYear)
	p.TrackedTotalMonth.AddSum(v.TrackedTotalMonth)
	p.TrackedTotalWeek.AddSum(v.TrackedTotalWeek)
	p.TrackedTotalYesterday.AddSum(v.TrackedTotalYesterday)
	p.TrackedTotalDay.AddSum(v.TrackedTotalDay)
}

func CreateProjectReports(referenceDay time.Time, showEmpty bool, includeArchived bool, activeEndRef *time.Time, overallSummaryID string, ctx *context.TomContext) map[string]*ProjectSummary {
	frames := ctx.Store.Frames()
	if !includeArchived {
		frames.ExcludeArchived()
	}

	year := dateTime.NewYearRange(referenceDay, ctx.Locale, time.Local)
	month := dateTime.NewMonthRange(referenceDay, ctx.Locale, time.Local)
	week := dateTime.NewWeekRange(referenceDay, ctx.Locale, time.Local)
	day := dateTime.NewDayRange(referenceDay, ctx.Locale, time.Local)
	yesterday := day.Shift(0, 0, -1)

	result := map[string]*ProjectSummary{}
	if showEmpty {
		for _, p := range ctx.Store.Projects() {
			result[p.ID] = &ProjectSummary{Project: p}
		}
	}

	// fixme optimize, projects first, then frames
	mapping := frames.MapByProject()
	for projectID, frames := range mapping {
		ctx.Query.WithProjectAndParents(projectID, func(project *model.Project) bool {
			target, ok := result[project.ID]
			if !ok {
				target = NewProjectSummary(&year, &month, &week, &yesterday, &day, activeEndRef, project)
				result[project.ID] = target
			}

			for _, frame := range frames.Frames() {
				if project.ID == frame.ProjectId {
					target.add(frame)
				} else {
					target.addTotal(frame)
				}
			}
			return true
		})
	}

	if overallSummaryID != "" {
		overall := NewProjectSummary(nil, nil, nil, nil, nil, nil, &model.Project{ID: overallSummaryID, FullName: []string{overallSummaryID}})

		// sum up all project summaries of top-level projects
		for _, v := range result {
			if v.Project != nil && v.Project.ParentID == "" {
				overall.Add(v)
			}
		}
		result[overallSummaryID] = overall
	}

	return result
}
