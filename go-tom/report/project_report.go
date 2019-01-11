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

	TotalTrackedYear  time.Duration
	TotalTrackedMonth time.Duration
	TotalTrackedWeek  time.Duration
	TotalTrackedDay   time.Duration
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

func CreateProjectReports(frames *model.FrameList, referenceDay time.Time, ctx *context.GoTimeContext) map[string]*ProjectSummary {
	year := dateUtil.NewYearRange(referenceDay, ctx.Locale)
	week := dateUtil.NewWeekRange(referenceDay, ctx.Locale)
	month := dateUtil.NewMonthRange(referenceDay, ctx.Locale)
	day := dateUtil.NewDayRange(referenceDay, ctx.Locale)

	result := map[string]*ProjectSummary{}

	frames.FilterByDateRange(year, false)

	for _, frame := range frames.Frames() {
		isYear := year.ContainsP(frame.Start) && year.ContainsP(frame.End)
		isMonth := month.ContainsP(frame.Start) && month.ContainsP(frame.End)
		isWeek := week.ContainsP(frame.Start) && week.ContainsP(frame.End)
		isDay := day.ContainsP(frame.Start) && day.ContainsP(frame.End)

		ctx.Query.WithProjectAndParents(frame.ProjectId, func(project *model.Project) bool {
			target, ok := result[project.ID]
			if !ok {
				target = &ProjectSummary{ProjectID: project.ID}
				result[project.ID] = target
			}

			if project.ID == frame.ProjectId {
				if isYear {
					target.addYear(frame.Duration())
				}
				if isMonth {
					target.addMonth(frame.Duration())
				}
				if isWeek {
					target.addWeek(frame.Duration())
				}
				if isDay {
					target.addDay(frame.Duration())
				}
			} else {
				if isYear {
					target.addTotalYear(frame.Duration())
				}
				if isMonth {
					target.addTotalMonth(frame.Duration())
				}
				if isWeek {
					target.addTotalWeek(frame.Duration())
				}
				if isDay {
					target.addTotalDay(frame.Duration())
				}
			}
			return true
		})
	}

	return result
}
