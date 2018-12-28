package frames

import (
	"fmt"
	"sort"
	"time"

	"github.com/jansorg/gotime/gotime/store"
)

type Bucket struct {
	Start     *time.Time     `json:"from,omitempty"`
	End       *time.Time     `json:"to,omitempty"`
	Frames    []*store.Frame `json:"frames,omitempty"`
	GroupedBy interface{}    `json:"groupedBy,omitempty"`
}

func (b *Bucket) Title() string {
	if b.GroupedBy == nil {
		return ""
	}

	if p, ok := b.GroupedBy.(*store.Project); ok {
		return fmt.Sprintf("Project: %s", p.FullName)
	}

	return fmt.Sprintf("%v", b.GroupedBy)
}

func NewBucket(from, to time.Time) *Bucket {
	return &Bucket{Start: &from, End: &to}
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

func SplitByProject(frames []*store.Frame) []*Bucket {
	mapping := map[string][]*store.Frame{}

	for _, f := range frames {
		mapping[f.ProjectId] = append(mapping[f.ProjectId], f)
	}

	var result []*Bucket
	for k, v := range mapping {
		result = append(result, &Bucket{Frames: v, GroupedBy: k})
	}
	return result
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
			rangeStart = lowerBucketBound(*f.Start)
			rangeEnd = upperBucketBound(*f.End)

			bucket = NewBucket(rangeStart, rangeEnd)
			buckets = append(buckets, bucket)
		}

		bucket.Frames = append(bucket.Frames, f)
	}
	return buckets
}
