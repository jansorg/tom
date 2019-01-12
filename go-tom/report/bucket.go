package report

import (
	"fmt"
	"time"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateUtil"
	"github.com/jansorg/tom/go-tom/model"
)

type ResultBucket struct {
	ctx *context.GoTimeContext

	DateRange        dateUtil.DateRange `json:"dateRange,omitempty"`
	TrackedDateRange dateUtil.DateRange `json:"trackedTime,omitempty"`

	FrameCount    int           `json:"frameCount"`
	Duration      time.Duration `json:"duration"`
	ExactDuration time.Duration `json:"duration_exact"`

	SplitBy interface{}      `json:"splitBy,omitempty"`
	Frames  *model.FrameList `json:"frames,omitempty"`
	Results []*ResultBucket  `json:"results,omitempty"`
}

func (b *ResultBucket) Empty() bool {
	return len(b.Results) == 0
}

func (b *ResultBucket) EmptySource() bool {
	return b.Frames.Empty()
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

func (b *ResultBucket) FindProjectBucket() (*model.Project, error) {
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

func (b *ResultBucket) Split(splitter func(list *model.FrameList) []*model.FrameList) {
	parts := splitter(b.Frames)
	// defer func() {
	// 	b.Frames = &model.FrameList{}
	// }()

	b.Results = []*ResultBucket{}
	for _, p := range parts {
		b.Results = append(b.Results, &ResultBucket{
			ctx:    b.ctx,
			Frames: p,
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
			if p, ok := value.(*model.Project); ok {
				return fmt.Sprintf("%s", p.GetFullName("/"))
			}

			if t, ok := value.(*model.Tag); ok {
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
