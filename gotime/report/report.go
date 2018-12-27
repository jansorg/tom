package report

import (
	"sort"
	"strings"
	"time"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/store"
)

type RoundingMode int8

const (
	RoundNone RoundingMode = iota + 1
	RoundNearest
	RoundUp
)

type Results struct {
	From  *time.Time `json:"from_date,omitempty"`
	To    *time.Time `json:"to_date,omitempty"`
	Items []Result   `json:"items"`
}

type Result struct {
	Name          string        `json:"name"`
	Duration      time.Duration `json:"duration"`
	ExactDuration time.Duration `json:"exact_duration"`

	Projects []Result `json:"projects,omitempty"`
	Tags     []Result `json:"tags,omitempty"`
}

type TimeReport struct {
	RoundFramesTo     time.Duration
	FrameRoundingMode RoundingMode

	RoundTotalTo      time.Duration
	TotalRoundingMode RoundingMode
}

func round(value time.Duration, mode RoundingMode, roundTo time.Duration) time.Duration {
	if mode == RoundNone {
		return value
	}

	result := value.Round(roundTo)
	if mode == RoundNearest {
		return result
	}

	if result < value {
		// round up if Go rounded down
		result = result + roundTo
	}
	return result
}

func (t *TimeReport) Calc(start *time.Time, end *time.Time, ctx *context.GoTimeContext) (Results, error) {
	projects := make(map[string]time.Duration)
	projectsExact := make(map[string]time.Duration)
	// tags := make(map[string]time.Duration)

	frames := filterFrames(ctx.Store.Frames(), start, end)
	for _, frame := range frames {
		duration := frame.Duration()
		projectsExact[frame.ProjectId] = projectsExact[frame.ProjectId] + duration

		roundedDuration := round(duration, t.FrameRoundingMode, t.RoundFramesTo)
		projects[frame.ProjectId] = projects[frame.ProjectId] + roundedDuration
	}

	var reports []Result

	for k, v := range projects {
		project, err := ctx.Store.FindProject(k)
		if err != nil {
			return Results{}, err
		}

		reports = append(reports, Result{
			Name:          project.ShortName,
			ExactDuration: projectsExact[k],
			Duration:      round(v, t.TotalRoundingMode, t.RoundTotalTo),
		})
	}

	sort.SliceStable(reports, func(i, j int) bool {
		return strings.Compare(reports[i].Name, reports[j].Name) < 0
	})
	return Results{
		From:  start,
		To:    end,
		Items: reports,
	}, nil
}

func filterFrames(frames []store.Frame, start *time.Time, end *time.Time) []store.Frame {
	var result []store.Frame

	for _, frame := range frames {
		if frame.Start != nil && start != nil && !start.IsZero() && frame.Start.Before(*start) || frame.End != nil && end != nil && !end.IsZero() && frame.End.After(*end) {
			continue
		}
		result = append(result, frame)
	}
	return result
}
