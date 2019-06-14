package model

import (
	"fmt"
	"sort"
	"time"

	"github.com/jansorg/tom/go-tom/dateTime"
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
	Archived  bool       `json:"archived,omitempty"`
}

func (f *Frame) copy() *Frame {
	return &Frame{
		ID:        f.ID,
		ProjectId: f.ProjectId,
		Start:     f.Start,
		End:       f.End,
		Updated:   f.Updated,
		Notes:     f.Notes,
		TagIDs:    f.TagIDs,
		Archived:  f.Archived,
	}
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

func (f *Frame) ActiveDuration(end *time.Time) time.Duration {
	if f.IsActive() && end != nil {
		return end.Sub(*f.Start)
	}
	return f.Duration()
}

func (f *Frame) Intersection(activeEnd *time.Time, timeRange *dateTime.DateRange) time.Duration {
	var frameEnd *time.Time
	if f.IsActive() {
		if activeEnd == nil {
			return time.Duration(0)
		}
		frameEnd = activeEnd
	} else {
		frameEnd = f.End
	}

	if timeRange.ContainsP(f.Start) && timeRange.ContainsP(frameEnd) {
		return f.ActiveDuration(activeEnd)
	}

	// intersection at start of time range
	if !timeRange.ContainsP(f.Start) && timeRange.ContainsP(frameEnd) {
		return frameEnd.Sub(*timeRange.Start)
	}

	// intersection at end of time range
	if timeRange.ContainsP(f.Start) && !timeRange.ContainsP(frameEnd) {
		return timeRange.End.Sub(*f.Start)
	}

	// no intersection
	return time.Duration(0)
}

// Contains returns if this frame's time range contains this ref
// an unstopped frame contains ref if the start time is equal to it
func (f *Frame) Contains(ref *time.Time) bool {
	if ref == nil || ref.IsZero() {
		return false
	}

	if f.Start != nil && f.End != nil {
		debug := fmt.Sprintf("%s -> %s <- %s", f.Start.String(), ref.String(), f.End.String())
		fmt.Println(debug)

		return ref.Equal(*f.Start) ||
			ref.Equal(*f.End) ||
			ref.After(*f.Start) && ref.Before(*f.End)
	}
	return f.Start != nil && f.Start.Equal(*ref)
}

func (f *Frame) IsBefore(other *Frame) bool {
	return f.Start != nil && other.Start != nil && f.Start.Before(*other.Start)
}

func (f *Frame) IsAfter(other *Frame) bool {
	return !f.IsBefore(other) && f.Start != nil && other.Start != nil && f.Start.After(*other.Start)
}

func (f *Frame) Validate(requireID bool) error {
	if f.ID == "" && requireID {
		return fmt.Errorf("id of frame undefined")
	} else if f.Start == nil || f.Start.IsZero() {
		return fmt.Errorf("start time undefined")
	} else if f.ProjectId == "" {
		return fmt.Errorf("project id undefined")
	}
	return nil
}
