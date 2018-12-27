package frames

import (
	"time"

	"github.com/jansorg/gotime/gotime/store"
)

func FilterFrames(frames []store.Frame, start *time.Time, end *time.Time) []store.Frame {
	var result []store.Frame

	for _, frame := range frames {
		if frame.Start != nil && start != nil && !start.IsZero() && frame.Start.Before(*start) || frame.End != nil && end != nil && !end.IsZero() && frame.End.After(*end) {
			continue
		}
		result = append(result, frame)
	}
	return result
}
