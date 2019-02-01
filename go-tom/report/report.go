package report

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/util"
)

type SplitOperation int8

const (
	SplitByYear SplitOperation = iota + 1
	SplitByMonth
	SplitByWeek
	SplitByDay
	SplitByProject
	SplitByParentProject
)

func (s SplitOperation) IsDateSplit() bool  {
	return s >= SplitByYear && s <= SplitByDay
}

func (s SplitOperation) String() string {
	name := ""
	switch (s) {
	case SplitByYear:
		name = "year"
	case SplitByMonth:
		name = "month"
	case SplitByWeek:
		name = "week"
	case SplitByDay:
		name = "day"
	case SplitByProject:
		name = "project"
	case SplitByParentProject:
		name = "parentProject"
	}
	return name
}
func (s SplitOperation) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

type BucketReport struct {
	ctx    *context.TomContext
	source *model.FrameList
	config Config
	result *ResultBucket
}

func NewBucketReport(frameList *model.FrameList, config Config, context *context.TomContext) *BucketReport {
	report := &BucketReport{
		ctx:    context,
		config: config,
		source: frameList,
	}
	return report
}

func (b *BucketReport) Update() *ResultBucket {
	if !b.config.DateFilterRange.Empty() {
		b.source.FilterByDateRange(b.config.DateFilterRange, false)
	}

	projectIDs := b.config.ProjectIDs
	if b.config.IncludeSubprojects {
		projectIDs = []string{}
		for _, p := range b.ctx.Store.Projects() {
			for _, parentID := range b.config.ProjectIDs {
				if b.ctx.Store.ProjectIsSameOrChild(parentID, p.ID) {
					projectIDs = append(projectIDs, p.ID)
				}
			}
		}
	}

	if len(projectIDs) > 0 {
		// sort IDs to use binary search
		sort.Strings(projectIDs)
		b.source.Filter(func(frame *model.Frame) bool {
			i := sort.SearchStrings(projectIDs, frame.ProjectId)
			return i < len(projectIDs) && projectIDs[i] == frame.ProjectId
		})
	}

	config := b.config
	if config.DateFilterRange.Empty() {
		config.DateFilterRange = b.source.DateRange(b.ctx.Locale)
	}

	b.result = &ResultBucket{
		ctx:      b.ctx,
		config:   config,
		Frames:   b.source,
		Duration: util.NewDurationSumAll(b.config.EntryRounding, nil, nil),
	}

	for _, op := range b.config.Splitting {
		if op.IsDateSplit() {
			b.result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.SplitByDateRange(op)
			})
		} else if op == SplitByProject {
			b.result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.SplitByProjectID(op, frameProject, projectIDs)
			})
		} else if op == SplitByParentProject {
			b.result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.SplitByProjectID(op, frameProjectParent(b.ctx), projectIDs)
			})
		} else {
			util.Fatal(fmt.Errorf("unknown split operation %d", op))
		}
	}

	updateBucket(b, b.result)

	return b.result
}

func frameProject(frame *model.Frame) interface{} {
	return frame.ProjectId
}

func frameProjectParent(ctx *context.TomContext) func(*model.Frame) interface{} {
	return func(frame *model.Frame) interface{} {
		project, err := ctx.Store.ProjectByID(frame.ProjectId)
		if err != nil {
			return ""
		}
		return project.ParentID
	}
}

// depth first update of the buckets to aggregate stats from sub-buckets
func updateBucket(report *BucketReport, bucket *ResultBucket) {
	for _, sub := range bucket.ChildBuckets {
		updateBucket(report, sub)
	}

	bucket.Update()
	bucket.SortChildBuckets()
}
