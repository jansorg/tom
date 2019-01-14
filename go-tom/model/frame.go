package model

import (
	"fmt"
	"sort"
	"time"
)

func NewStartedFrame(project *Project) Frame {
	now := time.Now()
	return Frame{
		ID:        NextID(),
		ProjectId: project.ID,
		Start:     &now,
		Updated:   &now,
	}
}

type Frame struct {
	ID        string     `json:"id"`
	ProjectId string     `json:"project"`
	Start     *time.Time `json:"start,omitempty"`
	End       *time.Time `json:"end,omitempty"`
	Updated   *time.Time `json:"updated,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	TagIDs    []string   `json:"tags,omitempty"`
}

func (f *Frame) sortTagIDs() {
	sort.Strings(f.TagIDs)
}

func (f *Frame) AddTagID(id string) {
	i := sort.SearchStrings(f.TagIDs, id)
	if i >= len(f.TagIDs) || f.TagIDs[i] != id {
		// fixme user insertion index
		f.TagIDs = append(f.TagIDs, id)
		f.sortTagIDs()
	}
}

func (f *Frame) RemoveTagID(id string) {
	i := sort.SearchStrings(f.TagIDs, id)
	if i < len(f.TagIDs) && f.TagIDs[i] == id {
		f.TagIDs = append(f.TagIDs[:i], f.TagIDs[i+1:]...)
	}
}

func (f *Frame) AddTags(newTags ...*Tag) {
	for _, t := range newTags {
		f.AddTagID(t.ID)
	}
}

func (f *Frame) HasTag(tag *Tag) bool {
	if tag == nil {
		return false
	}

	i := sort.SearchStrings(f.TagIDs, tag.ID)
	return i < len(f.TagIDs) && f.TagIDs[i] == tag.ID
}

func (f *Frame) IsStopped() bool {
	return f.End != nil && !f.End.IsZero()
}

func (f *Frame) IsSingleDay() bool {
	if f.IsActive() {
		return true
	}

	y1, m1, d1 := f.Start.Date()
	y2, m2, d2 := f.Start.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

func (f *Frame) IsActive() bool {
	return f.End == nil || f.End.IsZero()
}

func (f *Frame) Stop() {
	f.StopAt(time.Now())
}

func (f *Frame) StopAt(time time.Time) {
	f.End = &time
	f.Updated = &time
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

func (f *Frame) Validate() error {
	if f.ID == "" {
		return fmt.Errorf("id of frame undefined")
	} else if f.Start == nil || f.Start.IsZero() {
		return fmt.Errorf("start time undefined")
	}
	return nil
}
