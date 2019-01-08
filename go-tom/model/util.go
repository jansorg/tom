package model

import (
	"time"

	"github.com/satori/uuid"
)

func FilterFrames(frames []*Frame, start *time.Time, end *time.Time) []*Frame {
	if start == nil && end == nil {
		return frames
	}

	var result []*Frame

	for _, frame := range frames {
		if frame.Start != nil && start != nil && !start.IsZero() && frame.Start.Before(*start) && frame.End != nil && end != nil && !end.IsZero() && frame.End.After(*end) {
			continue
		}
		result = append(result, frame)
	}
	return result
}

func NextID() string {
	return uuid.NewV4().String()
}

