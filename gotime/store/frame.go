package store

import "time"

type Frame struct {
	Id        string     `json:"id"`
	ProjectId string     `json:"project"`
	Start     *time.Time `json:"start,omitempty"`
	End       *time.Time `json:"end,omitempty"`
	Updated   *time.Time `json:"updated,omitempty"`
	Notes     string     `json:"notes,omitempty"`
}

func (f *Frame) IsStopped() bool {
	return f.End != nil && !f.End.IsZero()
}

func (f *Frame) IsActive() bool {
	return f.End == nil || f.End.IsZero()
}

func (f *Frame) Stop() {
	now := time.Now()
	f.End = &now
	f.Updated = &now
}

func (f *Frame) Duration() time.Duration {
	if f.IsStopped() {
		return f.End.Sub(*f.Start)
	}
	return time.Duration(0)
}

func (f *Frame) IsBefore(other *Frame) bool {
	return f.Start != nil && other.Start != nil && f.Start.Before(*other.Start)
}

func (f *Frame) IsAfter(other *Frame) bool {
	return !f.IsBefore(other) && f.Start != nil && other.Start != nil && f.Start.After(*other.Start)
}

func NewStartedFrame(project *Project) Frame {
	now := time.Now()
	return Frame{
		Id:        nextID(),
		ProjectId: project.ID,
		Start:     &now,
		Updated:   &now,
	}
}
