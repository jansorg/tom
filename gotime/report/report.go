package report

import (
	"fmt"
	"log"
	"time"

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
	store  *store.Store
	source *frames.FrameList

	Result *ResultBucket `json:"result"`

	FilterRange dateUtil.DateRange `json:"dateRange,omitempty"`

	SplitOperations []SplitOperation `json:"splitOperations"`

	RoundingModeFrames dateUtil.RoundingMode `json:"roundingModeFrames"`
	RoundFramesTo      time.Duration         `json:"roundFramesTo"`

	RoundingModeTotals dateUtil.RoundingMode `json:"roundingModeTotals"`
	RoundTotalsTo      time.Duration         `json:"roundTotalsTo"`
}

func NewBucketReport(frameList *frames.FrameList) *BucketReport {
	report := &BucketReport{
		source: frameList,
	}
	return report
}

func (b *BucketReport) Update() {
	b.source.FilterByDatePtr(b.FilterRange.Start, b.FilterRange.End, false)
	b.Result = &ResultBucket{
		Source: b.source,
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
				leaf.DateRange = dateUtil.NewYearRange(*leaf.Source.First().Start)
				leaf.SplitBy = leaf.DateRange
			})
		case SplitByMonth:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *frames.FrameList) []*frames.FrameList {
					return list.SplitByMonth()
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.DateRange = dateUtil.NewMonthRange(*leaf.Source.First().Start)
				leaf.SplitBy = leaf.DateRange
			})
		case SplitByDay:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *frames.FrameList) []*frames.FrameList {
					return list.SplitByDay()
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.DateRange = dateUtil.NewDayRange(*leaf.Source.First().Start)
				leaf.SplitBy = leaf.DateRange
			})
		case SplitByProject:
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.Split(func(list *frames.FrameList) []*frames.FrameList {
					return list.SplitByProject()
				})
			})
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.SplitBy = leaf.Source.First().ProjectId
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

		if bucket.UsedDateRange.Start == nil {
			bucket.UsedDateRange.Start = first.UsedDateRange.Start
		}
		if bucket.UsedDateRange.End == nil {
			bucket.UsedDateRange.End = last.UsedDateRange.End
		}

		if bucket.DateRange.Empty() {
			bucket.DateRange.Start = first.DateRange.Start
			bucket.DateRange.End = last.DateRange.End
		}
	} else if !bucket.Source.Empty() {
		if bucket.UsedDateRange.Start == nil {
			bucket.UsedDateRange.Start = bucket.Source.First().Start
		}
		if bucket.UsedDateRange.End == nil {
			bucket.UsedDateRange.End = bucket.Source.Last().End
		}
	}
}
