package report

import (
	"fmt"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/util"
)

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

func (b *BucketReport) Result() *ResultBucket {
	return b.result
}

func (b *BucketReport) Update() *ResultBucket {
	if !b.config.DateFilterRange.Empty() {
		b.source.FilterByDateRange(b.config.DateFilterRange, false)
	}

	projectIDs := b.config.ProjectIDs
	if b.config.IncludeSubprojects {
		projectIDs = []string{}
		if len(b.config.ProjectIDs) == 0 {
			for _, p := range b.ctx.Store.Projects().Projects() {
				projectIDs = append(projectIDs, p.ID)
			}
		} else {
			for _, parentID := range b.config.ProjectIDs {
				projectIDs = append(projectIDs, parentID)
				projectIDs = append(projectIDs, b.ctx.Query.CollectSubprojectIDs(parentID)...)
			}
		}
	}

	// we need to filter our source by project ID
	if len(projectIDs) > 0 {
		b.source.Filter(func(frame *model.Frame) bool {
			_, err := b.ctx.Query.FindSuitableProject(frame.ProjectId, projectIDs)
			return err == nil
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
				leaf.SplitByProjectID(op, func(frame *model.Frame) interface{} {
					if len(projectIDs) == 0 {
						return frame.ProjectId
					}

					if id, err := b.ctx.Query.FindSuitableProject(frame.ProjectId, projectIDs); err == nil {
						return id
					} else {
						// fixme shouldn't happen as source is already filtered
						util.Fatalf("unexpected error, unable to find project bucket for %s in %v", frame.ProjectId, projectIDs, err.Error())
						return ""
					}
				}, projectIDs)
			})
		} else {
			util.Fatal(fmt.Errorf("unknown split operation %d", op))
		}
	}

	updateBucket(b, b.result)

	return b.result
}

// depth first update of the buckets to aggregate stats from sub-buckets
func updateBucket(report *BucketReport, bucket *ResultBucket) {
	for _, sub := range bucket.ChildBuckets {
		updateBucket(report, sub)
	}

	bucket.Update()
	bucket.SortChildBuckets()
}
