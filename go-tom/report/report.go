package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/jansorg/tom/go-tom/cmd/util"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateUtil"
	"github.com/jansorg/tom/go-tom/model"
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

	return json.Marshal(name)
}

type BucketReport struct {
	ctx    *context.TomContext
	source *model.FrameList

	Result *ResultBucket `json:"result"`

	TargetLocation      *time.Location     `json:"timezone"`
	IncludeActiveFrames bool               `json:"includeActiveFrames"`
	ProjectIDs          []string           `json:"projectIDs,omitempty"`
	IncludeSubprojects  bool               `json:"includeSubprojects,omitempty"`
	FilterRange         dateUtil.DateRange `json:"dateRange,omitempty"`
	SplitOperations     []SplitOperation   `json:"splitOperations"`
	ShowEmptyBuckets    bool

	RoundingModeFrames dateUtil.RoundingMode `json:"roundingModeFrames"`
	RoundFramesTo      time.Duration         `json:"roundFramesTo"`

	RoundingModeTotals dateUtil.RoundingMode `json:"roundingModeTotals"`
	RoundTotalsTo      time.Duration         `json:"roundTotalsTo"`
}

func NewBucketReport(frameList *model.FrameList, context *context.TomContext) *BucketReport {
	report := &BucketReport{
		ctx:            context,
		source:         frameList,
		TargetLocation: time.Local,
	}
	return report
}

func (b *BucketReport) IsRounding() bool {
	return b.RoundFramesTo != 0 && b.RoundingModeFrames != dateUtil.RoundNone || b.RoundTotalsTo != 0 && b.RoundingModeTotals != dateUtil.RoundNone
}

func (b *BucketReport) Update() {
	b.source.FilterByDatePtr(b.FilterRange.Start, b.FilterRange.End, false)

	projectIDs := b.ProjectIDs
	if b.IncludeSubprojects {
		projectIDs = []string{}
		for _, p := range b.ctx.Store.Projects() {
			for _, parentID := range b.ProjectIDs {
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

	var dateRange dateUtil.DateRange
	if b.source.Empty() {
		dateRange = b.FilterRange
	} else {
		dateRange = dateUtil.NewDateRange(b.source.First().Start, b.source.Last().End, b.ctx.Locale)
	}
	b.Result = &ResultBucket{
		ctx:       b.ctx,
		Frames:    b.source,
		Duration:  dateUtil.NewDurationSumAll(b.RoundingModeFrames, b.RoundFramesTo, nil, nil),
		dateRange: dateRange,
	}

	for _, op := range b.SplitOperations {
		if op <= SplitByDay {
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.SplitByDateRange(op, b.ShowEmptyBuckets)
			})
		} else if op == SplitByProject {
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.SplitByProjectID(op, b.ShowEmptyBuckets, projectOf, projectIDs)
			})
		} else if op == SplitByParentProject {
			b.Result.WithLeafBuckets(func(leaf *ResultBucket) {
				leaf.SplitByProjectID(op, b.ShowEmptyBuckets, parentProjectOf(b.ctx), projectIDs)
			})
		} else {
			util.Fatal(fmt.Errorf("unknown split operation %d", op))
		}
	}

	updateBucket(b, b.Result)
}

func projectOf(frame *model.Frame) interface{} {
	return frame.ProjectId
}

func parentProjectOf(ctx *context.TomContext) func(*model.Frame) interface{} {
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
