package report

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateUtil"
	"github.com/jansorg/tom/go-tom/model"
)

type SplitOperation int8

const (
	SplitByYear SplitOperation = iota + 1
	SplitByMonth
	SplitByWeek
	SplitByDay
	SplitByProject
	SplitByParentProject
)

func (s SplitOperation) MarshalJSON() ([]byte, error) {
	name := ""
	switch (s) {
	case SplitByYear:
		name = "year"
	case SplitByMonth:
		name = "month"
	case SplitByWeek:
		name = "week"
	case SplitByDay:
		name = "day"
	case SplitByProject:
		name = "project"
	case SplitByParentProject:
		name = "parentProject"
	}

	return json.Marshal(name)
}

type BucketReport struct {
	ctx    *context.TomContext
	source *model.FrameList

	Result *ResultBucket `json:"result"`

	TargetLocation      *time.Location     `json:"timezone"`
	IncludeActiveFrames bool               `json:"includeActiveFrames"`
	ProjectIDs          []string           `json:"projectIDs,omitempty"`
	IncludeSubprojects  bool               `json:"includeSubprojects,omitempty"`
	FilterRange         dateUtil.DateRange `json:"dateRange,omitempty"`
	SplitOperations     []SplitOperation   `json:"splitOperations"`
	ShowEmptyBuckets    bool

	RoundingModeFrames dateUtil.RoundingMode `json:"roundingModeFrames"`
	RoundFramesTo      time.Duration         `json:"roundFramesTo"`

	RoundingModeTotals dateUtil.RoundingMode `json:"roundingModeTotals"`
	RoundTotalsTo      time.Duration         `json:"roundTotalsTo"`
}

func NewBucketReport(frameList *model.FrameList, context *context.TomContext) *BucketReport {
	report := &BucketReport{
		ctx:            context,
		source:         frameList,
		TargetLocation: time.Local,
	}
	return report
}

func (b *BucketReport) IsRounding() bool {
	return b.RoundFramesTo != 0 && b.RoundingModeFrames != dateUtil.RoundNone || b.RoundTotalsTo != 0 && b.RoundingModeTotals != dateUtil.RoundNone
}

func (b *BucketReport) Update() {
	b.source.FilterByDatePtr(b.FilterRange.Start, b.FilterRange.End, false)
	if b.source.Empty() {
		return
	}

	projectIDs := b.ProjectIDs
	if b.IncludeSubprojects {
		projectIDs = []string{}
		for _, p := range b.ctx.Store.Projects() {
			for _, parentID := range b.ProjectIDs {
				if b.ctx.Store.ProjectIsSameOrChild(parentID, p.ID) {
					projectIDs = append(projectIDs, p.ID)
				}
			}
		}
	}

	if len(projectIDs) > 0 {
		// sort IDs to use binary search
		sort.Strings(projectIDs)
		b.source.Filter(func(frame *model.Frame) bool {
			i := sort.SearchStrings(projectIDs, frame.ProjectId)
			return i < len(projectIDs) && projectIDs[i] == frame.ProjectId
		})
	}

	b.Result = &ResultBucket{
		ctx:    b.ctx,
		Frames: b.source,
		// fixme filter?
		Duration: dateUtil.NewDurationSumAll(b.RoundingModeFrames, b.RoundFramesTo, nil, nil),
	}

	for _, op := range b.SplitOperations {
		switch op {
		case SplitByYear:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.SplitByDateRange(op, b.ShowEmptyBuckets, b.yearOf, func(r dateUtil.DateRange) dateUtil.DateRange {
					return r.Shift(1, 0, 0)
				})
			})
		case SplitByMonth:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.SplitByDateRange(op, b.ShowEmptyBuckets, b.monthOf, func(r dateUtil.DateRange) dateUtil.DateRange {
					return r.Shift(0, 1, 0)
				})
			})
		case SplitByWeek:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.SplitByDateRange(op, b.ShowEmptyBuckets, b.weekOf, func(r dateUtil.DateRange) dateUtil.DateRange {
					// fixme
					return r.Shift(1, 0, 0)
				})
			})
		case SplitByDay:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.SplitByDateRange(op, b.ShowEmptyBuckets, b.dayOf, func(r dateUtil.DateRange) dateUtil.DateRange {
					return r.Shift(0, 0, 1)
				})
			})
		case SplitByProject:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(op, projectOf, nil)
			})
		case SplitByParentProject:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(op, parentProjectOf(b.ctx), nil)
			})
		default:
			log.Fatal(fmt.Errorf("unknown split operation %d", op))
		}
	}

	updateBucket(b, b.Result)
}

func (b *BucketReport) yearOf(frame *model.Frame) dateUtil.DateRange {
	return dateUtil.NewYearRange(*frame.Start, b.ctx.Locale, frame.Start.Location())
}

func (b *BucketReport) monthOf(frame *model.Frame) dateUtil.DateRange {
	return dateUtil.NewMonthRange(*frame.Start, b.ctx.Locale, frame.Start.Location())
}

func (b *BucketReport) weekOf(frame *model.Frame) dateUtil.DateRange {
	return dateUtil.NewWeekRange(*frame.Start, b.ctx.Locale, frame.Start.Location())
}

func (b *BucketReport) dayOf(frame *model.Frame) dateUtil.DateRange {
	return dateUtil.NewDayRange(*frame.Start, b.ctx.Locale, frame.Start.Location())
}

func projectOf(frame *model.Frame) interface{} {
	return frame.ProjectId
}

func parentProjectOf(ctx *context.TomContext) func(*model.Frame) interface{} {
	return func(frame *model.Frame) interface{} {
		project, err := ctx.Store.ProjectByID(frame.ProjectId)
		if err != nil {
			return ""
		}
		return project.ParentID
	}
}

// depth first update of the buckets to aggregate stats from sub-buckets
func updateBucket(report *BucketReport, bucket *ResultBucket) {
	for _, sub := range bucket.ChildBuckets {
		updateBucket(report, sub)
	}

	for _, f := range bucket.Frames.Frames() {
		bucket.Duration.AddStartEndP(f.Start, f.End)
	}

	bucket.Update()
	bucket.SortChildBuckets()
}
