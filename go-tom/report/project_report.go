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

func CreateProjectReports(referenceDay time.Time, showEmpty bool, ctx *context.GoTimeContext) map[string]*ProjectSummary {
	frames := model.NewFrameList(ctx.Store.Frames())

	year := dateUtil.NewYearRange(referenceDay, ctx.Locale)
	week := dateUtil.NewWeekRange(referenceDay, ctx.Locale)
	month := dateUtil.NewMonthRange(referenceDay, ctx.Locale)
	day := dateUtil.NewDayRange(referenceDay, ctx.Locale)

	result := map[string]*ProjectSummary{}
	if showEmpty {
		for _, p := range ctx.Store.Projects() {
			result[p.ID] = &ProjectSummary{Project: p}
		}
	}

	for _, frame := range frames.Frames() {
		duration := frame.Duration()

		isYear := year.ContainsP(frame.Start) && year.ContainsP(frame.End)
		isMonth := month.ContainsP(frame.Start) && month.ContainsP(frame.End)
		isWeek := week.ContainsP(frame.Start) && week.ContainsP(frame.End)
		isDay := day.ContainsP(frame.Start) && day.ContainsP(frame.End)

		ctx.Query.WithProjectAndParents(frame.ProjectId, func(project *model.Project) bool {
			target, ok := result[project.ID]
			if !ok {
				target = &ProjectSummary{Project: project}
				result[project.ID] = target
			}

			if project.ID == frame.ProjectId {
				target.addAll(duration)

				if isYear {
					target.addYear(duration)
				}
				if isMonth {
					target.addMonth(duration)
				}
				if isWeek {
					target.addWeek(duration)
				}
				if isDay {
					target.addDay(duration)
				}
			} else {
				target.addTotalAll(duration)

				if isYear {
					target.addTotalYear(duration)
				}
				if isMonth {
					target.addTotalMonth(duration)
				}
				if isWeek {
					target.addTotalWeek(duration)
				}
				if isDay {
					target.addTotalDay(duration)
				}
			}
			return true
		})
	}

	return result
}
