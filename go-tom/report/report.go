package report

import (
	"fmt"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateTime"
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
	if !b.config.IncludeArchived {
		b.source.ExcludeArchived()
	}
	if !b.config.DateFilterRange.Empty() {
		b.source.FilterByDateRange(b.config.DateFilterRange, false)
	}

	projectIDs := b.config.ProjectIDs
	if len(projectIDs) == 0 {
		// we started with all top-level project if no project was selected
		// we add all projects if IncludeSubprojects is true
		for _, p := range b.ctx.Store.Projects() {
			if b.config.IncludeSubprojects || p.ParentID == "" {
				projectIDs = append(projectIDs, p.ID)
			}
		}
	} else if b.config.IncludeSubprojects {
		projectIDs = []string{}
		for _, parentID := range b.config.ProjectIDs {
			projectIDs = append(projectIDs, parentID)
			projectIDs = append(projectIDs, b.ctx.Query.CollectSubprojectIDs(parentID)...)
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
		config.DateFilterRange = b.source.DateRange(b.ctx.Locale).In(config.Timezone)
	} else if config.Timezone != nil {
		config.DateFilterRange = config.DateFilterRange.In(config.Timezone)
	}

	b.result = &ResultBucket{
		ctx:            b.ctx,
		config:         config,
		Frames:         b.source,
		Duration:       dateTime.NewDurationSumAll(b.config.EntryRounding, nil, nil),
		DailyTracked:   dateTime.NewTrackedDaily(nil),
		DailyUnTracked: dateTime.NewUntrackedDaily(nil),
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

func IsMatrix(bucket *ResultBucket, ignoreEmpty bool) bool {
	if bucket.Depth() != 2 {
		return false
	}

	// make sure that all child buckets are of the same depth
	for _, c := range bucket.ChildBuckets {
		depth := c.Depth()
		if depth == 0 && ignoreEmpty {
			continue
		} else if depth != 1 {
			return false
		}
	}

	// all buckets must have the same number of children
	refCol := bucket.FirstNonEmptyChild().ChildBuckets

	for _, b := range bucket.ChildBuckets {
		if b.Empty() && ignoreEmpty {
			continue
		}
		if len(b.ChildBuckets) != len(refCol) {
			return false
		}

		for i, col := range b.ChildBuckets {
			other := refCol[i]

			// if both are the same month names, then accepot
			if col.SplitByType == SplitByMonth && other.SplitByType == SplitByMonth && col.DateRange().IsMonthRange() && other.DateRange().IsMonthRange() && col.DateRange().Start.Month() == other.DateRange().Start.Month() {
				continue
			}

			if col.Title() != other.Title() {
				return false
			}
		}
	}

	return true
}
