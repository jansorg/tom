package report

import (
	"time"

	"github.com/jansorg/gotime/gotime/dateUtil"
	"github.com/jansorg/gotime/gotime/frames"
	"github.com/jansorg/gotime/gotime/store"
)

type ResultBucket struct {
	From          *time.Time `json:"from,omitempty"`
	To            *time.Time `json:"to,omitempty"`
	Duration      time.Duration
	ExactDuration time.Duration
	FrameCount    int64

	Results []*ResultBucket `json:"results"`
	Source  *frames.Bucket  `json:"source"`
}

type BucketReport struct {
	store  *store.Store
	source []*store.Frame

	Result *ResultBucket `json:"result"`

	FromDate *time.Time `json:"from,omitempty"`
	ToDate   *time.Time `json:"to,omitempty"`

	GroupByYear  bool `json:"groupByYear"`
	GroupByMonth bool `json:"groupByMonth"`
	GroupByDay   bool `json:"groupByDay"`

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

	if b.GroupByYear {
		splitLeafBuckets([]*ResultBucket{b.Result}, frames.SplitByYear)
	}

	if b.GroupByMonth {
		splitLeafBuckets([]*ResultBucket{b.Result}, frames.SplitByMonth)
	}

	if b.GroupByDay {
		splitLeafBuckets([]*ResultBucket{b.Result}, frames.SplitByDay)
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
					From:   s.From,
					To:     s.To,
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
		if bucket.From == nil {
			bucket.From = bucket.Results[0].From
		}
		if bucket.To == nil {
			bucket.To = bucket.Results[len(bucket.Results)-1].To
		}
	} else {
		if bucket.From == nil {
			bucket.From = bucket.Source.Frames[0].Start
		}
		if bucket.To == nil {
			bucket.To = bucket.Source.Frames[len(bucket.Source.Frames)-1].End
		}
	}
}
