package report

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/dateUtil"
	"github.com/jansorg/gotime/gotime/frames"
	"github.com/jansorg/gotime/gotime/store"
)

type SplitOperation int8

const (
	SplitByYear SplitOperation = iota + 1
	SplitByMonth
	SplitByDay
	SplitByProject
)

type BucketReport struct {
	ctx    *context.GoTimeContext
	source *frames.FrameList

	Result *ResultBucket `json:"result"`

	ProjectID   string             `json:"projectID,omitempty"`
	FilterRange dateUtil.DateRange `json:"dateRange,omitempty"`

	SplitOperations []SplitOperation `json:"splitOperations"`

	RoundingModeFrames dateUtil.RoundingMode `json:"roundingModeFrames"`
	RoundFramesTo      time.Duration         `json:"roundFramesTo"`

	RoundingModeTotals dateUtil.RoundingMode `json:"roundingModeTotals"`
	RoundTotalsTo      time.Duration         `json:"roundTotalsTo"`
}

func NewBucketReport(frameList *frames.FrameList, context *context.GoTimeContext) *BucketReport {
	report := &BucketReport{
		ctx:    context,
		source: frameList,
	}
	return report
}

func (b *BucketReport) IsRounding() bool {
	return b.RoundFramesTo != 0 && b.RoundingModeFrames != dateUtil.RoundNone || b.RoundTotalsTo != 0 && b.RoundingModeTotals != dateUtil.RoundNone
}

func (b *BucketReport) Update() {
	b.source.FilterByDatePtr(b.FilterRange.Start, b.FilterRange.End, false)

	if b.ProjectID != "" {
		b.source.Filter(func(frame *store.Frame) bool {
			return b.ctx.Store.ProjectIsChild(b.ProjectID, frame.ProjectId)
		})
	}

	b.Result = &ResultBucket{
		ctx:     b.ctx,
		Source:  b.source,
		SplitBy: b.ProjectID,
	}

	for _, op := range b.SplitOperations {
		switch op {
		case SplitByYear:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *frames.FrameList) []*frames.FrameList {
					return list.SplitByYear()
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				if !leaf.Source.Empty() {
					leaf.DateRange = dateUtil.NewYearRange(*leaf.Source.First().Start)
				}
				leaf.SplitBy = leaf.DateRange
			})
		case SplitByMonth:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *frames.FrameList) []*frames.FrameList {
					return list.SplitByMonth()
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				if !leaf.Source.Empty() {
					leaf.DateRange = dateUtil.NewMonthRange(*leaf.Source.First().Start)
				}
				leaf.SplitBy = leaf.DateRange
			})
		case SplitByDay:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *frames.FrameList) []*frames.FrameList {
					return list.SplitByDay()
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				if !leaf.Source.Empty() {
					leaf.DateRange = dateUtil.NewDayRange(*leaf.Source.First().Start)
				}
				leaf.SplitBy = leaf.DateRange
			})
		case SplitByProject:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *frames.FrameList) []*frames.FrameList {
					return list.SplitByProject()
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				if !leaf.Source.Empty() {
					leaf.SplitBy = leaf.Source.First().ProjectId
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
	bucket.FrameCount = bucket.Source.Size()

	for _, sub := range bucket.Results {
		updateBucket(report, sub)
	}

	for _, f := range bucket.Source.Frames {
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
	} else if !bucket.Source.Empty() {
		if bucket.TrackedDateRange.Start == nil {
			bucket.TrackedDateRange.Start = bucket.Source.First().Start
		}
		if bucket.TrackedDateRange.End == nil {
			bucket.TrackedDateRange.End = bucket.Source.Last().End
		}
	}
}
