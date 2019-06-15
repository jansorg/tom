package model

import (
	"sort"
	"time"

	"github.com/go-playground/locales"

	"github.com/jansorg/tom/go-tom/dateTime"
)

func NewFrameList(frames []*Frame) *FrameList {
	list := FrameList(frames)
	return &list
}

func NewSortedFrameList(frames []*Frame) *FrameList {
	list := NewFrameList(frames)
	list.Sort()
	return list
}

func NewEmptyFrameList() *FrameList {
	return &FrameList{}
}

type FrameList []*Frame

func (f *FrameList) Copy() *FrameList {
	return NewFrameList(*f)
}

func (f *FrameList) Empty() bool {
	return len(*f) == 0
}

func (f *FrameList) Frames() []*Frame {
	return *f
}

func (f *FrameList) Size() int {
	return len(*f)
}

func (f *FrameList) Append(value *Frame) {
	*f = append(*f, value)
}

func (f *FrameList) First() *Frame {
	if f.Empty() {
		return nil
	}
	return (*f)[0]
}

func (f *FrameList) Last() *Frame {
	if f.Empty() {
		return nil
	}
	return (*f)[len(*f)-1]
}

func (f *FrameList) Sort() {
	sort.Slice(*f, func(i, j int) bool {
		a := (*f)[i]
		b := (*f)[j]
		return a.IsBefore(b)
	})
}

func (f *FrameList) FilterByStartDate(minStartDate time.Time) {
	if minStartDate.IsZero() {
		return
	}

	f.Filter(func(frame *Frame) bool {
		return !frame.Start.Before(minStartDate)
	})
}

func (f *FrameList) FilterByEndDate(maxEndDate time.Time, acceptUnstopped bool) {
	if maxEndDate.IsZero() {
		return
	}

	f.Filter(func(frame *Frame) bool {
		return frame.End == nil && acceptUnstopped || frame.End != nil && !frame.End.After(maxEndDate)
	})
}

func (f *FrameList) FilterByDateRange(dateRange dateTime.DateRange, acceptUnstopped bool, keepOverlapping bool) {
	if keepOverlapping {
		f.Filter(func(frame *Frame) bool {
			return acceptUnstopped && frame.End == nil || dateRange.Intersects(frame.Start, frame.End)
		})
		return
	}

	f.FilterByStartDate(*dateRange.Start)
	if dateRange.IsClosed() {
		f.FilterByEndDate(*dateRange.End, acceptUnstopped)
	}
}

func (f *FrameList) FilterByDate(start time.Time, end time.Time, acceptUnstopped bool, keepOverlapping bool) {
	f.FilterByDatePtr(&start, &end, acceptUnstopped, keepOverlapping)
}

func (f *FrameList) FilterByDatePtr(start *time.Time, end *time.Time, acceptUnstopped bool, keepOverlapping bool) {
	if keepOverlapping {
		f.Filter(func(frame *Frame) bool {
			return acceptUnstopped && frame.End == nil ||
				start != nil && frame.Contains(start) ||
				end != nil && frame.Contains(end)
		})
		return
	}

	f.FilterByStartDate(*start)
	if end != nil {
		f.FilterByEndDate(*end, acceptUnstopped)
	}
}

func (f *FrameList) Filter(accepted func(frame *Frame) bool) {
	var result []*Frame
	for _, frame := range *f {
		if accepted(frame) {
			result = append(result, frame)
		}
	}

	*f = result
	f.Sort()
}

func (f *FrameList) ExcludeArchived() {
	f.Filter(func(frame *Frame) bool {
		return !frame.Archived
	})
}

// CutEntriesTo cuts entries which span more than start->end
func (f *FrameList) CutEntriesTo(start *time.Time, end *time.Time) {
	length := len(*f)
	for i := 0; i < length; i++ {
		frame := (*f)[i]

		// must not modify the original frame, it may still be referenced by other lists
		adaptStart := start != nil && frame.Start.Before(*start)
		adaptEnd := end != nil && frame.End != nil && frame.End.After(*end)

		if adaptStart || adaptEnd {
			copied := frame.copy()
			if adaptStart {
				copied.Start = start
			}
			if adaptEnd {
				copied.End = end
			}
			(*f)[i] = copied
		}
	}
	f.Sort()
}

// Split splits all frames into one ore more parts
// The part a frame belongs to is coputed by the key function
// because the distribution of keys is not always in order a map has to be used here
func (f *FrameList) Split(key func(f *Frame) interface{}) []*FrameList {
	if f.Empty() {
		return []*FrameList{}
	}

	mapping := map[interface{}][]*Frame{}
	for _, f := range *f {
		v := key(f)
		mapping[v] = append(mapping[v], f)
	}

	var parts []*FrameList
	for _, frames := range mapping {
		parts = append(parts, NewSortedFrameList(frames))
	}
	sort.SliceStable(parts, func(i, j int) bool {
		a := parts[i]
		b := parts[j]
		if a.Empty() {
			return b.Empty()
		}
		return a.First().IsBefore(b.First())
	})
	return parts
}

func (f *FrameList) MapByProject() map[string]*FrameList {
	result := make(map[string]*FrameList)

	for _, f := range *f {
		id := f.ProjectId
		if sub, ok := result[id]; !ok {
			result[id] = NewFrameList([]*Frame{f})
		} else {
			sub.Append(f)
		}
	}

	return result
}

func (f *FrameList) DateRange(locale locales.Translator) dateTime.DateRange {
	if f.Empty() {
		return dateTime.NewDateRange(nil, nil, locale)
	}

	var start *time.Time
	var end *time.Time

	// locate first closed item
	for _, f := range *f {
		if f.Start != nil {
			start = f.Start
			break
		}
	}

	// locate latest end value
	// the list is sorted by start value, the last item isn't necessarily having the latest end value
	// it should have it in normal circumstances, though
	// for now we're assuming this
	for i := f.Size() - 1; i >= 0; i-- {
		frameEnd := (*f)[i].End
		if frameEnd != nil {
			end = frameEnd
			break
		}
	}

	return dateTime.NewDateRange(start, end, locale)
}
