package report

import (
	"fmt"
	"time"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateUtil"
	"github.com/jansorg/tom/go-tom/frames"
	"github.com/jansorg/tom/go-tom/store"
)

type ResultBucket struct {
	ctx *context.GoTimeContext

	DateRange        dateUtil.DateRange `json:"dateRange,omitempty"`
	TrackedDateRange dateUtil.DateRange `json:"trackedTime,omitempty"`

	FrameCount    int           `json:"frameCount"`
	Duration      time.Duration `json:"duration"`
	ExactDuration time.Duration `json:"duration_exact"`

	SplitBy interface{}
	Source  *frames.FrameList `json:"source,omitempty"`
	Results []*ResultBucket   `json:"results,omitempty"`
}

func (b *ResultBucket) Empty() bool {
	return len(b.Results) == 0
}

func (b *ResultBucket) EmptySource() bool {
	return b.Source.Empty()
}

func (b *ResultBucket) IsRounded() bool {
	return b.Duration != b.ExactDuration
}

func (b *ResultBucket) IsDateBucket() bool {
	_, ok := b.SplitBy.(dateUtil.DateRange)
	return ok
}

func (b *ResultBucket) FilterResults(accepted func(bucket *ResultBucket) bool) {
	var result []*ResultBucket
	for _, r := range b.Results {
		if accepted(r) {
			result = append(result, r)
		}
	}
	b.Results = result
}

func (b *ResultBucket) IsProjectBucket() bool {
	_, err := b.FindProjectBucket()
	return err == nil
}

func (b *ResultBucket) FindProjectBucket() (*store.Project, error) {
	if id, ok := b.SplitBy.(string); ok {
		if p, err := b.ctx.Query.ProjectByID(id); err == nil {
			return p, nil
		}
	}
	return nil, fmt.Errorf("no project found for bucket")
}

func (b *ResultBucket) ProjectResults() []*ResultBucket {
	var result []*ResultBucket
	for _, r := range b.Results {
		if r.IsProjectBucket() {
			result = append(result, r)
		}
	}
	return result
}

func (b *ResultBucket) HasRoundedChildren() bool {
	for _, r := range b.Results {
		if r.IsRounded() {
			return true
		}
	}
	return false
}

func (b *ResultBucket) LeafChildren() []*ResultBucket {
	var result []*ResultBucket
	for _, r := range b.Results {
		if r.Empty() {
			result = append(result, r)
		}
	}
	return result
}

func (b *ResultBucket) Split(splitter func(list *frames.FrameList) []*frames.FrameList) {
	parts := splitter(b.Source)

	b.Results = []*ResultBucket{}
	for _, p := range parts {
		b.Results = append(b.Results, &ResultBucket{
			ctx:    b.ctx,
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

func (b *ResultBucket) Title() string {
	if id, ok := b.SplitBy.(string); ok {
		if value, err := b.ctx.Query.AnyByID(id); err == nil {
			if p, ok := value.(*store.Project); ok {
				return fmt.Sprintf("%s", p.FullName)
			}

			if t, ok := value.(*store.Tag); ok {
				return fmt.Sprintf("#%s", t.Name)
			}
		}
	}

	if dates, ok := b.SplitBy.(dateUtil.DateRange); ok {
		return dates.MinimalString()
	}

	if b.SplitBy != nil {
		return fmt.Sprintf("%v", b.SplitBy)
	}
	return ""
}
