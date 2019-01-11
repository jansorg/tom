package model

import (
	"sort"
	"time"

	"github.com/jansorg/tom/go-tom/dateUtil"
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
	sort.SliceStable(*f, func(i, j int) bool {
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

func (f *FrameList) FilterByDateRange(dateRange dateUtil.DateRange, acceptUnstopped bool) {
	f.FilterByStartDate(*dateRange.Start)
	if dateRange.IsClosed() {
		f.FilterByEndDate(*dateRange.End, acceptUnstopped)
	}
}

func (f *FrameList) FilterByDate(start time.Time, end time.Time, acceptUnstopped bool) {
	f.FilterByStartDate(start)
	f.FilterByEndDate(end, acceptUnstopped)
}

func (f *FrameList) FilterByDatePtr(start *time.Time, end *time.Time, acceptUnstopped bool) {
	if start != nil && !start.IsZero() {
		f.FilterByStartDate(*start)
	}
	if end != nil && !end.IsZero() {
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

func (f *FrameList) SplitByProject() []*FrameList {
	return f.Split(func(frame *Frame) interface{} {
		return frame.ProjectId
	})
}

func (f *FrameList) SplitByYear() []*FrameList {
	return f.Split(func(frame *Frame) interface{} {
		return frame.Start.Year()
	})
}

func (f *FrameList) SplitByMonth() []*FrameList {
	return f.Split(func(frame *Frame) interface{} {
		y, m, _ := frame.Start.Date()
		return time.Date(y, m, 1, 0, 0, 0, 0, frame.Start.Location())
	})
}

func (f *FrameList) SplitByWeek() []*FrameList {
	return f.Split(func(frame *Frame) interface{} {
		y, week := frame.Start.ISOWeek()
		return y*1000 + week
	})
}

func (f *FrameList) SplitByDay() []*FrameList {
	return f.Split(func(frame *Frame) interface{} {
		y, m, d := frame.Start.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, frame.Start.Location())
	})
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
