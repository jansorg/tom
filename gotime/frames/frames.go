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

func newBucket(from, to time.Time) *Bucket {
	return &Bucket{From: from, To: to}
}

func SplitByDay(frames []store.Frame) []*Bucket {
	sort.SliceStable(frames, func(i, j int) bool {
		return frames[i].IsBefore(&frames[j])
	})

	first := frames[0]
	dayStart, dayEnd := DayRange(*first.Start)

	bucket := newBucket(dayStart, dayEnd)
	buckets := []*Bucket{bucket}

	for _, f := range frames {
		if f.Start != nil && f.Start.After(dayEnd) {
			dayStart = dayEnd
			dayEnd = dayStart.Add(time.Hour * 24)
			bucket = newBucket(dayStart, dayEnd)
			buckets = append(buckets, bucket)
		}

		bucket.Frames = append(bucket.Frames, &f)
	}
	return buckets
}

func DayRange(date time.Time) (time.Time, time.Time) {
	y, m, d := date.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, date.Location()), time.Date(y, m, d, 24, 0, 0, 0, date.Location())
}
