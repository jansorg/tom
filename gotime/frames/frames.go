package frames

import (
	"fmt"
	"sort"
	"time"

	"github.com/jansorg/gotime/gotime/store"
)

type Bucket struct {
	From   time.Time
	To     time.Time
	Frames []*store.Frame
}

func (b *Bucket) Title() string {
	return fmt.Sprintf("%s - %s", b.From.String(), b.To.String())
}

func NewBucket(from, to time.Time) *Bucket {
	return &Bucket{From: from, To: to}
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
			return time.Date(y, m, 0, 0, 0, 0, 0, date.Location())
		},
		func(date time.Time) time.Time {
			y, m, _ := date.Date()
			return time.Date(y, m, 0, 0, 0, 0, 0, date.Location()).AddDate(0, 1, 0)
		},
	)
}

func SplitByYear(frames []*store.Frame) []*Bucket {
	return Split(frames,
		func(date time.Time) time.Time {
			y, _, _ := date.Date()
			return time.Date(y, time.January, 0, 0, 0, 0, 0, date.Location())
		},
		func(date time.Time) time.Time {
			y, _, _ := date.Date()
			return time.Date(y, time.January, 0, 0, 0, 0, 0, date.Location()).AddDate(1, 0, 0)
		},
	)
}

func Split(frames []*store.Frame, lowerBucketBound func(time.Time) time.Time, upperBucketBound func(time.Time) time.Time) []*Bucket {
	sort.SliceStable(frames, func(i, j int) bool {
		return frames[i].IsBefore(frames[j])
	})

	first := frames[0]
	dayStart := lowerBucketBound(*first.Start)
	dayEnd := upperBucketBound(*first.Start)

	bucket := NewBucket(dayStart, dayEnd)
	buckets := []*Bucket{bucket}

	for _, f := range frames {
		if f.Start != nil && f.Start.After(dayEnd) {
			dayStart = dayEnd
			dayEnd = upperBucketBound(dayStart)
			bucket = NewBucket(dayStart, dayEnd)
			buckets = append(buckets, bucket)
		}

		bucket.Frames = append(bucket.Frames, f)
	}
	return buckets
}
