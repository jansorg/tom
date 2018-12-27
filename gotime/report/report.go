package report

import (
	"time"

	"github.com/jansorg/gotime/gotime/dateUtil"
	"github.com/jansorg/gotime/gotime/frames"
	"github.com/jansorg/gotime/gotime/store"
)

type ResultBucket struct {
	From    *time.Time      `json:"from,omitempty"`
	To      *time.Time      `json:"to,omitempty"`
	Results []*ResultBucket `json:"results"`
	Source  *frames.Bucket  `json:"source"`
}

type BucketReport struct {
	store  *store.Store
	source *frames.Bucket

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
		source: &frames.Bucket{Frames: frameList},
	}
	return report
}

func (b *BucketReport) Update() {
	var buckets []*ResultBucket
	buckets = append(buckets, &ResultBucket{
		Source: b.source,
	})

	if b.GroupByYear {
		splitLeafBuckets(buckets, frames.SplitByYear)
	}

	if b.GroupByMonth {
		splitLeafBuckets(buckets, frames.SplitByMonth)
	}

	if b.GroupByDay {
		splitLeafBuckets(buckets, frames.SplitByDay)
	}

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
					From:   &s.From,
					To:     &s.To,
					Source: s,
				})
			}
		}
	}
}
