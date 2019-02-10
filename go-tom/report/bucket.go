package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/store"
	"github.com/jansorg/tom/go-tom/util"
)

type ResultBucket struct {
	ctx    *context.TomContext
	config Config
	parent *ResultBucket

	dateRange        util.DateRange
	trackedDateRange util.DateRange

	Frames       *model.FrameList  `json:"-"`
	FrameCount   int               `json:"frameCount"`
	Duration     *util.DurationSum `json:"duration"`
	Properties   []PropertyValue   `json:"properties"`
	SplitByType  SplitOperation    `json:"split_type,omitempty"`
	SplitBy      interface{}       `json:"split_by,omitempty"`
	ChildBuckets []*ResultBucket   `json:"results,omitempty"`
}

func (b *ResultBucket) Update() {
	b.FrameCount = b.Frames.Size()
	for _, f := range b.Frames.Frames() {
		b.Duration.AddStartEndP(f.Start, f.End)
	}

	if b.dateRange.Empty() {
		if !b.Empty() && !b.IsDateBucket() {
			childBuckets := b.ChildBuckets
			b.dateRange = util.NewDateRange(childBuckets[0].DateRange().Start, childBuckets[len(b.ChildBuckets)-1].DateRange().End, b.ctx.Locale)
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
		b.trackedDateRange = util.NewDateRange(first.Start, last.End, b.ctx.Locale)
	} else if !b.Empty() {
		first := b.ChildBuckets[0]
		last := b.ChildBuckets[len(b.ChildBuckets)-1]
		b.trackedDateRange = util.NewDateRange(first.TrackedDateRange().Start, last.TrackedDateRange().End, b.ctx.Locale)
	} else {
		b.trackedDateRange = util.DateRange{}
	}

	// property values
	// fixme optimize this
	for _, prop := range b.config.Properties {
		propValue := PropertyValue{
			Property:        prop,
			ValueForRounded: 0,
			ValueForExact:   0,
		}

		if b.Empty() {
			// we can't just multiply with the rounded duration as the frames might be from different projects with different property values
			// for _, frame := range b.Frames.Frames() {
			// fixme handle rounded
			// if value, err := b.ctx.Query.FindPropertyValue(prop.ID, frame.ProjectId); err == nil {
			// if floatValue, err := prop.ToFloat(value); err == nil {
			// 	propValue.ValueForExact += frame.Duration().Hours() * floatValue
			// 	propValue.ValueForRounded += frame.Duration().Hours() * floatValue
			// }
			// }
			// }
			// fixme make sure that we properly sum up the values of a project with subproject with hidden sub-projects
			// 		add a test for this scenario
		} else {
			// add child bucket values
			for _, child := range b.ChildBuckets {
				if childValue, err := child.PropertyValue(prop.ID); err == nil {
					propValue.ValueForExact += childValue.ValueForExact
					propValue.ValueForRounded += childValue.ValueForRounded
				}
			}
		}

		b.Properties = append(b.Properties, propValue)
	}
}

func (b *ResultBucket) PropertyValue(id string) (*PropertyValue, error) {
	for _, v := range b.Properties {
		if v.Property.ID == id {
			return &v, nil
		}
	}
	return nil, store.ErrPropertyNotFound
}

func (b *ResultBucket) AppliedFilterRange() util.DateRange {
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

func (b *ResultBucket) TrackedDateRange() util.DateRange {
	return b.trackedDateRange
}

func (b *ResultBucket) DateRange() util.DateRange {
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
			if p, ok := value.(*model.Project); ok {
				return fmt.Sprintf("%s", p.GetFullName("/"))
			}

			if t, ok := value.(*model.Tag); ok {
				return fmt.Sprintf("#%s", t.Name)
			}
		}
	}

	if dates, ok := b.SplitBy.(util.DateRange); ok {
		return dates.MinimalString()
	}

	if b.SplitBy != nil {
		return fmt.Sprintf("%v", b.SplitBy)
	}
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
	});
}

func (b *ResultBucket) SplitByProjectID(splitType SplitOperation, splitValue func(frame *model.Frame) interface{}, minProjectIDs []string) {
	mapping := make(map[string]bool)

	for _, frameSubset := range b.Frames.Split(splitValue) {
		value := splitValue(frameSubset.First())
		mapping[value.(string)] = true

		b.AddChild(&ResultBucket{
			Frames:      frameSubset,
			Duration:    util.NewEmptyCopy(b.Duration),
			SplitByType: splitType,
			SplitBy:     value,
		})
	}

	for _, id := range minProjectIDs {
		if !mapping[id] {
			b.AddChild(&ResultBucket{
				Frames:      model.NewFrameList([]*model.Frame{}),
				Duration:    util.NewDurationSum(),
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

	var value util.DateRange
	switch splitType {
	case SplitByYear:
		value = util.NewYearRange(*start, b.ctx.Locale, start.Location())
	case SplitByMonth:
		value = util.NewMonthRange(*start, b.ctx.Locale, start.Location())
	case SplitByWeek:
		value = util.NewWeekRange(*start, b.ctx.Locale, start.Location())
	case SplitByDay:
		value = util.NewDayRange(*start, b.ctx.Locale, start.Location())
	}

	for value.IsClosed() && value.Start.Before(*end) {
		matchingFrames := b.Frames.Copy()
		matchingFrames.FilterByDateRange(value, false)

		if b.config.ShowEmpty || !matchingFrames.Empty() {
			b.AddChild(&ResultBucket{
				dateRange:   value,
				Frames:      matchingFrames,
				Duration:    util.NewEmptyCopy(b.Duration),
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
			value = value.Shift(0, 0, 7)
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
