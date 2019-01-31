package report

import (
	"fmt"
	"log"
	"sort"
	"strings"
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
		ctx:       b.ctx,
		Frames:    b.source,
		DateRange: dateUtil.NewDateRange(nil, nil, b.ctx.Locale),
	}

	for _, op := range b.SplitOperations {
		switch op {
		case SplitByYear:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *model.FrameList) []*model.FrameList {
					return list.SplitByYear(b.TargetLocation)
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				if !leaf.Frames.Empty() {
					leaf.DateRange = dateUtil.NewYearRange(*leaf.Frames.First().Start, b.ctx.Locale, b.TargetLocation).In(b.TargetLocation)
				}
				leaf.SplitBy = leaf.DateRange
			})
		case SplitByMonth:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *model.FrameList) []*model.FrameList {
					return list.SplitByMonth(b.TargetLocation)
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				if !leaf.Frames.Empty() {
					leaf.DateRange = dateUtil.NewMonthRange(*leaf.Frames.First().Start, b.ctx.Locale, b.TargetLocation).In(b.TargetLocation)
				}
				leaf.SplitBy = leaf.DateRange
			})
		case SplitByWeek:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *model.FrameList) []*model.FrameList {
					return list.SplitByWeek(b.TargetLocation)
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				if !leaf.Frames.Empty() {
					leaf.DateRange = dateUtil.NewWeekRange(*leaf.Frames.First().Start, b.ctx.Locale, b.TargetLocation).In(b.TargetLocation)
				}
				leaf.SplitBy = leaf.DateRange
			})
		case SplitByDay:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *model.FrameList) []*model.FrameList {
					return list.SplitByDay(b.TargetLocation)
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				if !leaf.Frames.Empty() {
					leaf.DateRange = dateUtil.NewDayRange(*leaf.Frames.First().Start, b.ctx.Locale, b.TargetLocation).In(b.TargetLocation)
				}
				leaf.SplitBy = leaf.DateRange
			})
		case SplitByProject:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *model.FrameList) []*model.FrameList {
					return list.SplitByProject()
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				if !leaf.Frames.Empty() {
					leaf.SplitBy = leaf.Frames.First().ProjectId
				}
			})
			sort.SliceStable(b.Result.Results, func(i, j int) bool {
				a := b.Result.Results[i]
				b := b.Result.Results[j]
				return strings.Compare(strings.ToLower(a.Title()), strings.ToLower(b.Title())) < 0
			})
		case SplitByParentProject:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *model.FrameList) []*model.FrameList {
					return list.SplitByParentProject(b.ctx.Store)
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				if !leaf.Frames.Empty() {
					project, err := b.ctx.Store.ProjectByID(leaf.Frames.First().ProjectId)
					if err == nil {
						leaf.SplitBy = project.ParentID
					}
				}
			})
			sort.SliceStable(b.Result.Results, func(i, j int) bool {
				a := b.Result.Results[i]
				b := b.Result.Results[j]
				return strings.Compare(strings.ToLower(a.Title()), strings.ToLower(b.Title())) < 0
			})
		default:
			log.Fatal(fmt.Errorf("unknown split operation %d", op))
		}
	}

	updateBucket(b, b.Result)

	if b.Result.DateRange.Empty() {
		b.Result.DateRange = b.FilterRange
	}
}

// depth first update of the buckets to aggregate stats from sub-buckets
func updateBucket(report *BucketReport, bucket *ResultBucket) {
	bucket.FrameCount = bucket.Frames.Size()

	for _, sub := range bucket.Results {
		updateBucket(report, sub)
	}

	for _, f := range bucket.Frames.Frames() {
		d := f.Duration()
		bucket.ExactDuration += d
		bucket.Duration += dateUtil.RoundDuration(d, report.RoundingModeFrames, report.RoundFramesTo)
	}

	if len(bucket.Results) > 0 {
		first := bucket.Results[0]
		last := bucket.Results[len(bucket.Results)-1]

		if bucket.TrackedDateRange.Start == nil {
			bucket.TrackedDateRange.Start = first.TrackedDateRange.Start
		}
		if bucket.TrackedDateRange.End == nil {
			bucket.TrackedDateRange.End = last.TrackedDateRange.End
		}

		if bucket.DateRange.Empty() {
			bucket.DateRange.Start = first.DateRange.Start
			bucket.DateRange.End = last.DateRange.End
		}
	} else if !bucket.Frames.Empty() {
		if bucket.TrackedDateRange.Start == nil {
			bucket.TrackedDateRange.Start = bucket.Frames.First().Start
		}
		if bucket.TrackedDateRange.End == nil {
			bucket.TrackedDateRange.End = bucket.Frames.Last().End
		}
	}

	bucket.Sort()
}
