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

	Results []*ResultBucket `json:"results"`

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

	buckets := []*ResultBucket{{
		Source: &frames.Bucket{
			Frames: b.source,
		},
	}}

	if b.GroupByYear {
		splitLeafBuckets(buckets, frames.SplitByYear)
	}

	if b.GroupByMonth {
		splitLeafBuckets(buckets, frames.SplitByMonth)
	}

	if b.GroupByDay {
		splitLeafBuckets(buckets, frames.SplitByDay)
	}

	buckets = updateBuckets(b, buckets)
	b.Results = buckets
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
func updateBuckets(report *BucketReport, buckets []*ResultBucket) []*ResultBucket {
	for _, bucket := range buckets {
		for _, sub := range bucket.Results {
			updateBuckets(report, sub.Results)
			bucket.FrameCount += sub.FrameCount
		}

		bucket.FrameCount += int64(len(bucket.Source.Frames))
		if len(bucket.Source.Frames) > 0 {
			bucket.From = bucket.Source.Frames[0].Start
			bucket.To = bucket.Source.Frames[len(bucket.Source.Frames)-1].End
		}

		for _, f := range bucket.Source.Frames {
			d := f.Duration()
			bucket.ExactDuration += d
			bucket.Duration += dateUtil.RoundDuration(d, report.RoundingFrames, report.RoundFramesTo)
		}
	}

	return buckets
}
