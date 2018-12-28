package frames

import (
	"fmt"
	"sort"
	"time"

	"github.com/jansorg/gotime/gotime/store"
)

type Bucket struct {
	From   *time.Time     `json:"from,omitempty"`
	To     *time.Time     `json:"to,omitempty"`
	Frames []*store.Frame `json:"frames,omitempty"`
}

func (b *Bucket) Title() string {
	return fmt.Sprintf("%s - %s", b.From.String(), b.To.String())
}

func NewBucket(from, to time.Time) *Bucket {
	return &Bucket{From: &from, To: &to}
}

func SplitByDay(frames []*store.Frame) []*Bucket {
	return Split(frames,
		func(date time.Time) time.Time {
			y, m, d := date.Date()
			return time.Date(y, m, d, 0, 0, 0, 0, date.Location())
		},
		func(date time.Time) time.Time {
			y, m, d := date.Date()
			return time.Date(y, m, d, 24, 0, 0, 0, date.Location())
		},
	)
}

func SplitByMonth(frames []*store.Frame) []*Bucket {
	return Split(frames,
		func(date time.Time) time.Time {
			y, m, _ := date.Date()
			return time.Date(y, m, 1, 0, 0, 0, 0, date.Location())
		},
		func(date time.Time) time.Time {
			y, m, _ := date.Date()
			return time.Date(y, m, 1, 0, 0, 0, 0, date.Location()).AddDate(0, 1, 0)
		},
	)
}

func SplitByYear(frames []*store.Frame) []*Bucket {
	return Split(frames,
		func(date time.Time) time.Time {
			y, _, _ := date.Date()
			return time.Date(y, time.January, 1, 0, 0, 0, 0, date.Location())
		},
		func(date time.Time) time.Time {
			y, _, _ := date.Date()
			return time.Date(y, time.January, 1, 0, 0, 0, 0, date.Location()).AddDate(1, 0, 0)
		},
	)
}

func Split(frames []*store.Frame, lowerBucketBound func(time.Time) time.Time, upperBucketBound func(time.Time) time.Time) []*Bucket {
	if len(frames) == 0 {
		return []*Bucket{}
	}

	sort.SliceStable(frames, func(i, j int) bool {
		return frames[i].IsBefore(frames[j])
	})

	rangeStart := lowerBucketBound(*frames[0].Start)
	rangeEnd := upperBucketBound(*frames[0].End)

	bucket := NewBucket(rangeStart, rangeEnd)
	buckets := []*Bucket{bucket}

	for _, f := range frames {
		if f.Start != nil && f.Start.After(rangeEnd) {
			// fixme return no gaps
			rangeStart = lowerBucketBound(*f.Start)
			rangeEnd = upperBucketBound(*f.Start)

			bucket = NewBucket(rangeStart, rangeEnd)
			buckets = append(buckets, bucket)
		}

		bucket.Frames = append(bucket.Frames, f)
	}
	return buckets
}
