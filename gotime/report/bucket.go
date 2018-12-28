package report

import (
	"fmt"
	"time"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/dateUtil"
	"github.com/jansorg/gotime/gotime/frames"
	"github.com/jansorg/gotime/gotime/store"
)

type ResultBucket struct {
	DateRange        dateUtil.DateRange `json:"dateRange,omitempty"`
	TrackedDateRange dateUtil.DateRange `json:"trackedTime,omitempty"`

	FrameCount    int           `json:"frameCount"`
	Duration      time.Duration `json:"duration"`
	ExactDuration time.Duration `json:"duration_exact"`

	SplitBy interface{}
	Source  *frames.FrameList `json:"source,omitempty"`
	Results []*ResultBucket   `json:"results,omitempty"`
}

func (b *ResultBucket) Split(splitter func(list *frames.FrameList) []*frames.FrameList) {
	parts := splitter(b.Source)

	b.Results = []*ResultBucket{}
	for _, p := range parts {
		b.Results = append(b.Results, &ResultBucket{
			Source: p,
		})
	}
}

func (b *ResultBucket) WithLeafBuckets(handler func(leaf *ResultBucket)) {
	if len(b.Results) == 0 {
		handler(b)
		return
	}

	for _, sub := range b.Results {
		sub.WithLeafBuckets(handler)
	}
}

func (b *ResultBucket) Title(ctx *context.GoTimeContext) string {
	if id, ok := b.SplitBy.(string); ok {
		if value, err := ctx.Query.AnyByID(id); err == nil {
			if p, ok := value.(*store.Project); ok {
				return fmt.Sprintf("Project: %s", p.FullName)
			}

			if t, ok := value.(*store.Tag); ok {
				return fmt.Sprintf("Tag: %s", t.Name)
			}
		}
	}

	return ""
}
