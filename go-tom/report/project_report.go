package report

import (
	"time"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateUtil"
	"github.com/jansorg/tom/go-tom/model"
)

type ProjectSummary struct {
	Project *model.Project

	TrackedAll   time.Duration
	TrackedYear  time.Duration
	TrackedMonth time.Duration
	TrackedWeek  time.Duration
	TrackedDay   time.Duration

	TotalTrackedAll   time.Duration
	TotalTrackedYear  time.Duration
	TotalTrackedMonth time.Duration
	TotalTrackedWeek  time.Duration
	TotalTrackedDay   time.Duration
}

func (p *ProjectSummary) addAll(d time.Duration) {
	p.TrackedAll += d
	p.TotalTrackedAll += d
}

func (p *ProjectSummary) addTotalAll(d time.Duration) {
	p.TotalTrackedAll += d
}

func (p *ProjectSummary) addYear(d time.Duration) {
	p.TrackedYear += d
	p.TotalTrackedYear += d
}

func (p *ProjectSummary) addTotalYear(d time.Duration) {
	p.TotalTrackedYear += d
}

func (p *ProjectSummary) addMonth(d time.Duration) {
	p.TrackedMonth += d
	p.TotalTrackedMonth += d
}

func (p *ProjectSummary) addTotalMonth(d time.Duration) {
	p.TotalTrackedMonth += d
}

func (p *ProjectSummary) addWeek(d time.Duration) {
	p.TrackedWeek += d
	p.TotalTrackedWeek += d
}

func (p *ProjectSummary) addTotalWeek(d time.Duration) {
	p.TotalTrackedWeek += d
}

func (p *ProjectSummary) addDay(d time.Duration) {
	p.TrackedDay += d
	p.TotalTrackedDay += d
}

func (p *ProjectSummary) addTotalDay(d time.Duration) {
	p.TotalTrackedDay += d
}

func (p *ProjectSummary) Add(v *ProjectSummary) {
	p.TotalTrackedDay += v.TotalTrackedDay
	p.TotalTrackedWeek += v.TotalTrackedWeek
	p.TotalTrackedMonth += v.TotalTrackedMonth
	p.TotalTrackedYear += v.TotalTrackedYear
	p.TotalTrackedAll += v.TotalTrackedAll

	p.TrackedDay += v.TrackedDay
	p.TrackedWeek += v.TrackedWeek
	p.TrackedMonth += v.TrackedMonth
	p.TrackedYear += v.TrackedYear
	p.TrackedAll += v.TrackedAll
}

func CreateProjectReports(referenceDay time.Time, showEmpty bool, activeEndRef *time.Time, overallSummaryID string, ctx *context.TomContext) map[string]*ProjectSummary {
	frames := model.NewFrameList(ctx.Store.Frames())

	year := dateUtil.NewYearRange(referenceDay, ctx.Locale)
	month := dateUtil.NewMonthRange(referenceDay, ctx.Locale)
	week := dateUtil.NewWeekRange(referenceDay, ctx.Locale)
	day := dateUtil.NewDayRange(referenceDay, ctx.Locale)

	result := map[string]*ProjectSummary{}
	if showEmpty {
		for _, p := range ctx.Store.Projects() {
			result[p.ID] = &ProjectSummary{Project: p}
		}
	}

	for _, frame := range frames.Frames() {
		duration := frame.ActiveDuration(activeEndRef)
		if duration == 0 {
			continue
		}

		yearDuration := frame.Intersection(activeEndRef, &year)
		monthDuration := frame.Intersection(activeEndRef, &month)
		weekDuration := frame.Intersection(activeEndRef, &week)
		dayDuration := frame.Intersection(activeEndRef, &day)

		ctx.Query.WithProjectAndParents(frame.ProjectId, func(project *model.Project) bool {
			target, ok := result[project.ID]
			if !ok {
				target = &ProjectSummary{Project: project}
				result[project.ID] = target
			}

			if project.ID == frame.ProjectId {
				target.addAll(duration)

				target.addYear(yearDuration)
				target.addMonth(monthDuration)
				target.addWeek(weekDuration)
				target.addDay(dayDuration)
			} else {
				target.addTotalAll(duration)

				target.addTotalYear(yearDuration)
				target.addTotalMonth(monthDuration)
				target.addTotalWeek(weekDuration)
				target.addTotalDay(dayDuration)
			}
			return true
		})
	}

	if overallSummaryID != "" {
		overall := ProjectSummary{Project: &model.Project{ID: overallSummaryID, FullName: []string{overallSummaryID}}}

		// sum up all project summaries of top-level projects
		for _, v := range result {
			if v.Project.ParentID == "" {
				overall.Add(v)
			}
		}
		result[overallSummaryID] = &overall;
	}

	return result
}
