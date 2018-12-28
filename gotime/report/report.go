package report

import (
	"fmt"
	"log"
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

type ResultBucket struct {
	Start         *time.Time    `json:"start,omitempty"`
	End           *time.Time    `json:"end,omitempty"`
	Duration      time.Duration `json:"duration"`
	ExactDuration time.Duration `json:"duration_exact"`
	FrameCount    int64         `json:"frameCount"`

	Source  *frames.Bucket  `json:"source,omitempty"`
	Results []*ResultBucket `json:"results,omitempty"`
}

func (b *ResultBucket) Title(ctx *context.GoTimeContext) string {
	if b.Source != nil && b.Source.GroupedBy != nil {
		if id, ok := b.Source.GroupedBy.(string); ok {
			if value, err := ctx.Query.AnyByID(id); err == nil {
				if p, ok := value.(*store.Project); ok {
					return fmt.Sprintf("Project: %s", p.FullName)
				}

				if t, ok := value.(*store.Tag); ok {
					return fmt.Sprintf("Tag: %s", t.Name)
				}
			}
		}
	}

	return ""
}

type BucketReport struct {
	store  *store.Store
	source []*store.Frame

	Result *ResultBucket `json:"result"`

	FromDate *time.Time `json:"from,omitempty"`
	ToDate   *time.Time `json:"to,omitempty"`

	SplitOperations []SplitOperation `json:"splitOperations"`

	RoundingFrames dateUtil.RoundingMode `json:"roundingModeFrames"`
	RoundFramesTo  time.Duration         `json:"roundFramesTo"`

	RoundingTotals dateUtil.RoundingMode `json:"roundingModeTotals"`
	RoundTotalsTo  time.Duration         `json:"roundTotalsTo"`
}

func NewBucketReport(frameList []*store.Frame) *BucketReport {
	report := &BucketReport{
		source: frameList,
	}
	return report
}

func (b *BucketReport) Update() {
	b.source = frames.FilterFrames(b.source, b.FromDate, b.ToDate)
	b.Result = &ResultBucket{
		Source: &frames.Bucket{
			Frames: b.source,
		},
	}

	for _, op := range b.SplitOperations {
		switch op {
		case SplitByYear:
			splitLeafBuckets([]*ResultBucket{b.Result}, frames.SplitByYear)
		case SplitByMonth:
			splitLeafBuckets([]*ResultBucket{b.Result}, frames.SplitByMonth)
		case SplitByDay:
			splitLeafBuckets([]*ResultBucket{b.Result}, frames.SplitByDay)
		case SplitByProject:
			splitLeafBuckets([]*ResultBucket{b.Result}, frames.SplitByProject)
		default:
			log.Fatal(fmt.Errorf("unknown split operation %d", op))
		}
	}

	updateBucket(b, b.Result)
}

func splitLeafBuckets(buckets []*ResultBucket, splitter func([]*store.Frame) []*frames.Bucket) {
	for _, b := range buckets {
		if len(b.Results) != 0 {
			splitLeafBuckets(b.Results, splitter)
		} else {
			splitBuckets := splitter(b.Source.Frames)
			for _, s := range splitBuckets {
				b.Results = append(b.Results, &ResultBucket{
					Start:  s.Start,
					End:    s.End,
					Source: s,
				})
			}
		}
	}
}

// depth first update of the buckets to aggregate stats from sub-buckets
func updateBucket(report *BucketReport, bucket *ResultBucket) {
	bucket.FrameCount = int64(len(bucket.Source.Frames))
	for _, sub := range bucket.Results {
		updateBucket(report, sub)
		bucket.FrameCount += sub.FrameCount
	}

	for _, f := range bucket.Source.Frames {
		d := f.Duration()
		bucket.ExactDuration += d
		bucket.Duration += dateUtil.RoundDuration(d, report.RoundingFrames, report.RoundFramesTo)
	}

	if len(bucket.Results) > 0 {
		if bucket.Start == nil {
			bucket.Start = bucket.Results[0].Start
		}
		if bucket.End == nil {
			bucket.End = bucket.Results[len(bucket.Results)-1].End
		}
	} else {
		if bucket.Start == nil {
			bucket.Start = bucket.Source.Frames[0].Start
		}
		if bucket.End == nil {
			bucket.End = bucket.Source.Frames[len(bucket.Source.Frames)-1].End
		}
	}
}
