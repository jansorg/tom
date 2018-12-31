package store

import "time"

type Frame struct {
	ID        string     `json:"id"`
	ProjectId string     `json:"project"`
	Start     *time.Time `json:"start,omitempty"`
	End       *time.Time `json:"end,omitempty"`
	Updated   *time.Time `json:"updated,omitempty"`
	Notes     string     `json:"notes,omitempty"`
}

func (f *Frame) IsStopped() bool {
	return f.End != nil && !f.End.IsZero()
}

func (f *Frame) IsSingleDay() bool {
	if f.IsActive() {
		return true
	}

	y1,m1,d1 := f.Start.Date()
	y2,m2,d2 := f.Start.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
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
		ID:        nextID(),
		ProjectId: project.ID,
		Start:     &now,
		Updated:   &now,
	}
}
