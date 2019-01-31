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

	if !b.Empty() {
		first := b.ChildBuckets[0]
		last := b.ChildBuckets[len(b.ChildBuckets)-1]
		start := first.DateRange().Start
		b.dateRange = dateUtil.NewDateRange(start, last.DateRange().End, b.ctx.Locale)
	} else if !b.EmptySource() {
		start := b.Frames.First().Start
		switch b.SplitByType {
		case SplitByYear:
			// fixme
			b.dateRange = dateUtil.NewYearRange(*start, b.ctx.Locale, time.Local)
		case SplitByMonth:
			// fixme
			b.dateRange = dateUtil.NewMonthRange(*start, b.ctx.Locale, time.Local)
		case SplitByWeek:
			// fixme
			b.dateRange = dateUtil.NewWeekRange(*start, b.ctx.Locale, time.Local)
		case SplitByDay:
			// fixme
			b.dateRange = dateUtil.NewDayRange(*start, b.ctx.Locale, time.Local)
		default:
			b.dateRange = dateUtil.NewDateRange(nil, nil, b.ctx.Locale)
		}
	} else {
		b.dateRange = dateUtil.NewDateRange(nil, nil, b.ctx.Locale)
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

func (b *ResultBucket) TrackedDateRange() dateUtil.DateRange {
	return b.trackedDateRange
}

func (b *ResultBucket) DateRange() dateUtil.DateRange {
	return b.dateRange
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
	return b.SplitByType < SplitByProject && b.dateRange.IsClosed()
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

func (b *ResultBucket) Split(splitType SplitOperation, splitValue func(frame *model.Frame) interface{}, minValues []interface{}) {
	parts := b.Frames.Split(splitValue)

	mapping := make(map[interface{}]bool)

	b.ChildBuckets = []*ResultBucket{}
	for _, segment := range parts {
		value := splitValue(segment.First())
		mapping[value] = true

		b.ChildBuckets = append(b.ChildBuckets, &ResultBucket{
			ctx:         b.ctx,
			Frames:      segment,
			Duration:    dateUtil.NewDurationCopy(b.Duration),
			SplitByType: splitType,
			SplitBy:     value,
		})
	}
}

func (b *ResultBucket) SplitByDateRange(splitType SplitOperation, addEmpty bool, splitValue func(frame *model.Frame) dateUtil.DateRange, nextValue func(dateUtil.DateRange) dateUtil.DateRange) {
	b.ChildBuckets = []*ResultBucket{}

	value := splitValue(b.Frames.First())
	lastFrameDate := b.Frames.Last().Start

	for value.IsClosed() && !value.Start.After(*lastFrameDate) {
		matchingFrames := b.Frames.Copy()
		matchingFrames.FilterByDateRange(value, false)
		if addEmpty || !matchingFrames.Empty() {
			b.ChildBuckets = append(b.ChildBuckets, &ResultBucket{
				ctx:         b.ctx,
				Frames:      matchingFrames,
				Duration:    dateUtil.NewDurationCopy(b.Duration),
				SplitByType: splitType,
				SplitBy:     value,
			})
		}

		value = nextValue(value)
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
