package report

import (
	"sort"
	"strings"
	"time"

	"../context"
	"../store"
)

type RoundingMode int8

const (
	RoundNone RoundingMode = iota + 1
	RoundNearest
	RoundUp
)

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

func (t *TimeReport) Calc(frames []store.Frame, ctx *context.GoTimeContext) ([]Result, error) {
	projects := make(map[string]time.Duration)
	// tags := make(map[string]time.Duration)

	for _, frame := range frames {
		duration := frame.Duration()
		duration = round(duration, t.FrameRoundingMode, t.RoundFramesTo)

		projectData := projects[frame.ProjectId]
		projects[frame.ProjectId] = projectData + duration
	}

	var reports []Result

	for k, v := range projects {
		project, err := ctx.Store.FindProject(k)
		if err != nil {
			return nil, err
		}

		reports = append(reports, Result{
			Name:          project.ShortName,
			ExactDuration: v,
			Duration:      round(v, t.TotalRoundingMode, t.RoundTotalTo),
		})
	}

	sort.SliceStable(reports, func(i, j int) bool {
		return strings.Compare(reports[i].Name, reports[j].Name) < 0
	})
	return reports, nil
}
