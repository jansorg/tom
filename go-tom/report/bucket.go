package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateUtil"
	"github.com/jansorg/tom/go-tom/model"
)

type ResultBucket struct {
	ctx              *context.TomContext
	parent           *ResultBucket
	dateRange        dateUtil.DateRange
	trackedDateRange dateUtil.DateRange

	Frames       *model.FrameList      `json:"frames,omitempty"`
	FrameCount   int                   `json:"frameCount"`
	Duration     *dateUtil.DurationSum `json:"duration"`
	SplitByType  SplitOperation        `json:"split_type,omitempty"`
	SplitBy      interface{}           `json:"split_by,omitempty"`
	ChildBuckets []*ResultBucket       `json:"results,omitempty"`
}

func (b *ResultBucket) Update() {
	b.FrameCount = b.Frames.Size()
	for _, f := range b.Frames.Frames() {
		b.Duration.AddStartEndP(f.Start, f.End)
	}

	if !b.Empty() && !b.IsDateBucket() {
		b.dateRange = dateUtil.NewDateRange(b.ChildBuckets[0].DateRange().Start, b.ChildBuckets[len(b.ChildBuckets)-1].DateRange().End, b.ctx.Locale)
	}

	// tracked range
	if !b.EmptySource() {
		first := b.Frames.First()
		last := b.Frames.Last()
		b.trackedDateRange = dateUtil.NewDateRange(first.Start, last.End, b.ctx.Locale)
	} else if !b.Empty() {
		first := b.ChildBuckets[0]
		last := b.ChildBuckets[len(b.ChildBuckets)-1]
		b.trackedDateRange = dateUtil.NewDateRange(first.TrackedDateRange().Start, last.TrackedDateRange().End, b.ctx.Locale)
	} else {
		b.trackedDateRange = dateUtil.DateRange{}
	}
}

func (b *ResultBucket) ParentDateRange() dateUtil.DateRange {
	if b.parent != nil && b.parent.IsDateBucket() {
		return b.parent.dateRange
	} else if b.dateRange.IsClosed() {
		return b.dateRange
	} else if b.parent != nil {
		return b.parent.ParentDateRange()
	}
	return b.dateRange
}

func (b *ResultBucket) TrackedDateRange() dateUtil.DateRange {
	return b.trackedDateRange
}

func (b *ResultBucket) DateRange() dateUtil.DateRange {
	return b.dateRange
}

func (b *ResultBucket) Depth() int {
	if b.Empty() {
		return 0
	}
	return 1 + b.ChildBuckets[0].Depth()
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

func (b *ResultBucket) Title() string {
	if id, ok := b.SplitBy.(string); ok {
		if value, err := b.ctx.Query.AnyByID(id); err == nil {
			if p, ok := value.(*model.Project); ok {
				return fmt.Sprintf("%s", p.GetFullName("/"))
			}

			if t, ok := value.(*model.Tag); ok {
				return fmt.Sprintf("#%s", t.Name)
			}
		}
	}

	if dates, ok := b.SplitBy.(dateUtil.DateRange); ok {
		return dates.MinimalString()
	}

	if b.SplitBy != nil {
		return fmt.Sprintf("%v", b.SplitBy)
	}
	return ""
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

		return strings.Compare(b1.Title(), b2.Title()) < 0
	});
}

func (b *ResultBucket) SplitByProjectID(splitType SplitOperation, showEmpty bool, splitValue func(frame *model.Frame) interface{}, minProjectIDs []string) {
	mapping := make(map[string]bool)

	b.ChildBuckets = []*ResultBucket{}
	for _, segment := range b.Frames.Split(splitValue) {
		value := splitValue(segment.First())
		mapping[value.(string)] = true

		b.ChildBuckets = append(b.ChildBuckets, &ResultBucket{
			ctx:         b.ctx,
			parent:      b,
			Frames:      segment,
			Duration:    dateUtil.NewEmptyCopy(b.Duration),
			SplitByType: splitType,
			SplitBy:     value,
		})
	}

	for _, id := range minProjectIDs {
		if !mapping[id] {
			b.ChildBuckets = append(b.ChildBuckets, &ResultBucket{
				ctx:         b.ctx,
				parent:      b,
				Frames:      model.NewFrameList([]*model.Frame{}),
				Duration:    dateUtil.NewDurationSum(),
				SplitByType: splitType,
				SplitBy:     id,
			})
		}
	}
}

func (b *ResultBucket) SplitByDateRange(splitType SplitOperation, showEmpty bool) {
	b.ChildBuckets = []*ResultBucket{}

	var start, end *time.Time
	if !b.EmptySource() {
		start = b.Frames.First().Start
		end = b.Frames.Last().End
	}

	parentRange := b.ParentDateRange()
	if parentRange.IsClosed() {
		start = parentRange.Start
		end = parentRange.End
	}

	if start == nil || end == nil {
		return
	}

	var value dateUtil.DateRange
	switch splitType {
	case SplitByYear:
		value = dateUtil.NewYearRange(*start, b.ctx.Locale, start.Location())
	case SplitByMonth:
		value = dateUtil.NewMonthRange(*start, b.ctx.Locale, start.Location())
	case SplitByWeek:
		value = dateUtil.NewWeekRange(*start, b.ctx.Locale, start.Location())
	case SplitByDay:
		value = dateUtil.NewDayRange(*start, b.ctx.Locale, start.Location())
	}

	for value.IsClosed() && value.Start.Before(*end) {
		matchingFrames := b.Frames.Copy()
		matchingFrames.FilterByDateRange(value, false)

		if showEmpty || !matchingFrames.Empty() {
			b.ChildBuckets = append(b.ChildBuckets, &ResultBucket{
				ctx:         b.ctx,
				parent:      b,
				dateRange:   value,
				Frames:      matchingFrames,
				Duration:    dateUtil.NewEmptyCopy(b.Duration),
				SplitByType: splitType,
				SplitBy:     value,
			})
		}

		switch splitType {
		case SplitByYear:
			value = value.Shift(1, 0, 0)
		case SplitByMonth:
			value = value.Shift(0, 1, 0)
		case SplitByWeek:
			value = value.Shift(0, 1, 0)
		case SplitByDay:
			value = value.Shift(0, 0, 1)
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
