package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateTime"
	"github.com/jansorg/tom/go-tom/model"
)

type ResultBucket struct {
	ctx    *context.TomContext
	config Config
	parent *ResultBucket

	dateRange        dateTime.DateRange
	trackedDateRange dateTime.DateRange

	Frames       *model.FrameList      `json:"-"`
	FrameCount   int                   `json:"frameCount"`
	Duration     *dateTime.DurationSum `json:"duration"`
	Sales        *Sales                `json:"sales"`
	SplitByType  SplitOperation        `json:"split_type,omitempty"`
	SplitBy      interface{}           `json:"split_by,omitempty"`
	ChildBuckets []*ResultBucket       `json:"results,omitempty"`

	DailyTracked   dateTime.TimeEntrySeries `json:"daily_tracked"`
	DailyUnTracked dateTime.TimeEntrySeries `json:"daily_untracked"`
}

func (b *ResultBucket) Update() {
	b.Sales = NewSales(b.ctx, b.config.EntryRounding)

	b.FrameCount = b.Frames.Size()
	for _, f := range b.Frames.Frames() {
		b.Duration.AddStartEndP(f.Start, f.End)

		if !f.IsActive() {
			if b.DailyTracked != nil {
				b.DailyTracked.Add(*f.Start, *f.End)
			}
			if b.DailyUnTracked != nil {
				b.DailyUnTracked.Add(*f.Start, *f.End)
			}
		}
	}

	if b.dateRange.Empty() {
		if !b.Empty() && !b.IsDateBucket() {
			childBuckets := b.ChildBuckets
			b.dateRange = dateTime.NewDateRange(childBuckets[0].DateRange().Start, childBuckets[len(b.ChildBuckets)-1].DateRange().End, b.ctx.Locale)
		} else if b.Empty() && b.parent != nil {
			b.dateRange = b.parent.DateRange()
			if b.dateRange.Empty() {
				b.dateRange = b.config.DateFilterRange
			}
		}
	}

	// tracked range
	if !b.EmptySource() {
		first := b.Frames.First()
		last := b.Frames.Last()
		b.trackedDateRange = dateTime.NewDateRange(first.Start, last.End, b.ctx.Locale)
	} else if !b.Empty() {
		first := b.ChildBuckets[0]
		last := b.ChildBuckets[len(b.ChildBuckets)-1]
		b.trackedDateRange = dateTime.NewDateRange(first.TrackedDateRange().Start, last.TrackedDateRange().End, b.ctx.Locale)
	} else {
		b.trackedDateRange = dateTime.DateRange{}
	}

	if b.Empty() {
		for _, f := range b.Frames.Frames() {
			_ = b.Sales.Add(f)
		}
	} else {
		for _, c := range b.ChildBuckets {
			// fixme handle error
			if c.Sales != nil {
				_ = b.Sales.AddSales(c.Sales)
			}
		}
	}
}

func (b *ResultBucket) AppliedFilterRange() dateTime.DateRange {
	if !b.dateRange.Empty() {
		return b.dateRange
	}

	if b.parent != nil {
		return b.parent.AppliedFilterRange()
	}

	return b.config.DateFilterRange
}

func (b *ResultBucket) AddChild(child *ResultBucket) {
	child.ctx = b.ctx
	child.config = b.config
	child.parent = b
	b.ChildBuckets = append(b.ChildBuckets, child)
}

func (b *ResultBucket) TrackedDateRange() dateTime.DateRange {
	return b.trackedDateRange
}

func (b *ResultBucket) DateRange() dateTime.DateRange {
	return b.dateRange
}

func (b *ResultBucket) Depth() int {
	if b.Empty() {
		return 0
	}

	max := 0
	for _, c := range b.ChildBuckets {
		d := c.Depth()
		if d > max {
			max = d
		}
	}
	return 1 + max
}

func (b *ResultBucket) Empty() bool {
	return len(b.ChildBuckets) == 0
}

func (b *ResultBucket) EmptySource() bool {
	return b.Frames.Empty()
}

func (b *ResultBucket) IsRounded() bool {
	return b.Duration.IsRounded()
}

func (b *ResultBucket) IsDateBucket() bool {
	return b.SplitByType > 0 && b.SplitByType < SplitByProject && b.dateRange.IsClosed()
}

func (b *ResultBucket) IsProjectBucket() bool {
	_, err := b.FindProjectBucket()
	return err == nil
}

func (b *ResultBucket) HasDailyTracked() bool {
	return b.DailyTracked != nil
}

func (b *ResultBucket) HasDailyUnTracked() bool {
	return b.DailyUnTracked != nil
}

func (b *ResultBucket) GetDailyTracked() dateTime.TimeEntrySeries {
	return b.DailyTracked
}

func (b *ResultBucket) GetDailyUnTracked() dateTime.TimeEntrySeries {
	return b.DailyUnTracked
}

func (b *ResultBucket) FindProjectBucket() (*model.Project, error) {
	if id, ok := b.SplitBy.(string); ok {
		if p, err := b.ctx.Query.ProjectByID(id); err == nil {
			return p, nil
		}
	}
	return nil, fmt.Errorf("no project found for bucket")
}

func (b *ResultBucket) FilterResults(accepted func(bucket *ResultBucket) bool) {
	var result []*ResultBucket
	for _, r := range b.ChildBuckets {
		if accepted(r) {
			result = append(result, r)
		}
	}
	b.ChildBuckets = result
}

func (b *ResultBucket) ProjectResults() []*ResultBucket {
	var result []*ResultBucket
	for _, r := range b.ChildBuckets {
		if r.IsProjectBucket() {
			result = append(result, r)
		}
	}
	return result
}

func (b *ResultBucket) HasRoundedChildren() bool {
	for _, r := range b.ChildBuckets {
		if r.IsRounded() {
			return true
		}
	}
	return false
}

func (b *ResultBucket) EmptyChildren() []*ResultBucket {
	var result []*ResultBucket
	for _, r := range b.ChildBuckets {
		if r.Empty() {
			result = append(result, r)
		}
	}
	return result
}

func (b *ResultBucket) FirstNonEmptyChild() *ResultBucket {
	for _, r := range b.ChildBuckets {
		if !r.Empty() {
			return r
		}
	}
	return nil
}

func (b *ResultBucket) Title() string {
	if id, ok := b.SplitBy.(string); ok {
		if value, err := b.ctx.Query.AnyByID(id); err == nil {
			// Title of a project
			if p, ok := value.(*model.Project); ok {
				name := p.GetFullName(b.config.ProjectDelimiter)

				if b.config.ShortTitles {
					if parent := b.FindProjectParent(); parent != nil {
						parentName := parent.GetFullName(b.config.ProjectDelimiter)
						return strings.TrimPrefix(name, parentName+b.config.ProjectDelimiter)
					}
				}

				return name
			}

			// Title of a Tag
			if t, ok := value.(*model.Tag); ok {
				return fmt.Sprintf("#%s", t.Name)
			}
		}
	}

	// Title of a DateRange
	if dates, ok := b.SplitBy.(dateTime.DateRange); ok {
		return dates.MinimalString()
	}

	// Title of a Split container
	if b.SplitBy != nil {
		return fmt.Sprintf("%v", b.SplitBy)
	}

	// fallback
	return ""
}

func (b *ResultBucket) MatrixTitle() string {
	if b.SplitByType == SplitByMonth && b.IsDateBucket() && b.parent == nil || b.parent.dateRange.IsYearRange() {
		return b.ctx.Locale.MonthWide(b.dateRange.Start.Month())
	}
	return b.Title()
}

func (b *ResultBucket) SortChildBuckets() {
	if b.Empty() {
		return
	}

	sort.Slice(b.ChildBuckets, func(i, j int) bool {
		b1 := b.ChildBuckets[i]
		b2 := b.ChildBuckets[j]

		if b1.IsDateBucket() && b2.IsDateBucket() {
			return b1.DateRange().Start.Before(*b2.DateRange().Start)
		}

		return strings.Compare(strings.ToLower(b1.Title()), strings.ToLower(b2.Title())) < 0
	})
}

func (b *ResultBucket) SplitByProjectID(splitType SplitOperation, splitValue func(frame *model.Frame) interface{}, minProjectIDs []string) {
	mapping := make(map[string]bool)

	for _, frameSubset := range b.Frames.Split(splitValue) {
		value := splitValue(frameSubset.First())
		mapping[value.(string)] = true

		b.AddChild(&ResultBucket{
			Frames:         frameSubset,
			Duration:       dateTime.NewEmptyCopy(b.Duration),
			DailyTracked:   dateTime.NewTrackedDaily(nil),
			DailyUnTracked: dateTime.NewUntrackedDaily(nil),
			SplitByType:    splitType,
			SplitBy:        value,
		})
	}

	for _, id := range minProjectIDs {
		if !mapping[id] {
			b.AddChild(&ResultBucket{
				Frames:      model.NewFrameList([]*model.Frame{}),
				Duration:    dateTime.NewDurationSum(),
				SplitByType: splitType,
				SplitBy:     id,
			})
		}
	}
}

func (b *ResultBucket) SplitByDateRange(splitType SplitOperation) {
	b.ChildBuckets = []*ResultBucket{}

	filterRange := b.AppliedFilterRange()
	if filterRange.Empty() {
		filterRange = b.Frames.DateRange(b.ctx.Locale)
	}

	if !filterRange.IsClosed() {
		return
	}

	start := filterRange.Start
	end := filterRange.End

	var splitValue dateTime.DateRange
	switch splitType {
	case SplitByYear:
		splitValue = dateTime.NewYearRange(*start, b.ctx.Locale, start.Location())
	case SplitByMonth:
		splitValue = dateTime.NewMonthRange(*start, b.ctx.Locale, start.Location())
	case SplitByWeek:
		splitValue = dateTime.NewWeekRange(*start, b.ctx.Locale, start.Location())
	case SplitByDay:
		splitValue = dateTime.NewDayRange(*start, b.ctx.Locale, start.Location())
	}

	for splitValue.IsClosed() && splitValue.Start.Before(*end) {
		matchingFrames := b.Frames.Copy()
		matchingFrames.FilterByDateRange(splitValue, false, true)
		matchingFrames.CutEntriesTo(splitValue.Start, splitValue.End)

		rangeCopy := splitValue
		if b.config.ShowEmpty || !matchingFrames.Empty() {
			b.AddChild(&ResultBucket{
				dateRange:      splitValue,
				Frames:         matchingFrames,
				Duration:       dateTime.NewEmptyCopy(b.Duration),
				SplitByType:    splitType,
				SplitBy:        splitValue,
				DailyTracked:   dateTime.NewTrackedDaily(&rangeCopy),
				DailyUnTracked: dateTime.NewUntrackedDaily(&rangeCopy),
			})
		}

		switch splitType {
		case SplitByYear:
			splitValue = splitValue.Shift(1, 0, 0)
		case SplitByMonth:
			splitValue = splitValue.Shift(0, 1, 0)
		case SplitByWeek:
			splitValue = splitValue.Shift(0, 0, 7)
		case SplitByDay:
			splitValue = splitValue.Shift(0, 0, 1)
		}
	}

}

func (b *ResultBucket) WithLeafBuckets(handler func(leaf *ResultBucket)) {
	if len(b.ChildBuckets) == 0 {
		handler(b)
		return
	}

	for _, sub := range b.ChildBuckets {
		sub.WithLeafBuckets(handler)
	}
}

func (b *ResultBucket) FindFirstParent(accepted func(parent *ResultBucket) bool) *ResultBucket {
	parent := b.parent
	for parent != nil {
		if accepted(parent) {
			return parent
		}
		parent = parent.parent
	}
	return nil
}

func (b *ResultBucket) FindProjectParent() *model.Project {
	result := b.FindFirstParent(func(parent *ResultBucket) bool {
		return parent.IsProjectBucket()
	})

	if result != nil {
		if result, err := result.FindProjectBucket(); err == nil {
			return result
		}
	}

	// fallback in case only a single project is filtered for the report
	if len(b.config.ProjectIDs) == 1 {
		if p, err := b.ctx.Store.ProjectByID(b.config.ProjectIDs[0]); err == nil {
			return p
		}
	}

	return nil
}
