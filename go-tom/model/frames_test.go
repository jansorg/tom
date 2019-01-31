package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSplitByDay(t *testing.T) {
	buckets := NewFrameList([]*Frame{
		{
			Start: newDay(10, 12, 0),
			End:   newDay(10, 13, 30),
		},
		{
			Start: newDay(10, 14, 0),
			End:   newDay(10, 16, 0),
		},
		{
			Start: newDay(11, 12, 0),
			End:   newDay(11, 13, 30),
		},
		{
			Start: newDay(13, 1, 0),
			End:   newDay(13, 3, 0),
		},
	}).SplitByDay(time.Local)

	assert.EqualValues(t, 3, len(buckets))
	assert.EqualValues(t, 2, buckets[0].Size())
	assert.EqualValues(t, 1, buckets[1].Size())
	assert.EqualValues(t, 1, buckets[1].Size())
}

func TestSplitByMonth(t *testing.T) {
	buckets := NewFrameList([]*Frame{
		{
			Start: newDate(2017, time.February, 10, 12, 0),
			End:   newDate(2017, time.February, 11, 20, 0),
		},
		{
			Start: newDate(2018, time.February, 10, 12, 0),
			End:   newDate(2018, time.February, 13, 12, 0),
		},
		{
			Start: newDate(2018, time.February, 14, 12, 0),
			End:   newDate(2018, time.February, 16, 12, 0),
		},
	}).SplitByMonth(time.Local)

	assert.EqualValues(t, 2, len(buckets))
	assert.EqualValues(t, 1, buckets[0].Size())
	assert.EqualValues(t, 2, buckets[1].Size())
}

func TestSplitByYear(t *testing.T) {
	buckets := NewFrameList([]*Frame{
		{
			Start: newDate(2017, time.February, 10, 12, 0),
			End:   newDate(2017, time.February, 11, 20, 0),
		},
		{
			Start: newDate(2018, time.February, 10, 12, 0),
			End:   newDate(2018, time.February, 13, 12, 0),
		},
		{
			Start: newDate(2018, time.February, 14, 12, 0),
			End:   newDate(2018, time.February, 16, 12, 0),
		},
		{
			Start: newDate(2018, time.March, 14, 12, 0),
			End:   newDate(2018, time.April, 16, 12, 0),
		},
	}).SplitByYear(time.Local)

	assert.EqualValues(t, 2, len(buckets))
	assert.EqualValues(t, 1, buckets[0].Size())
	assert.EqualValues(t, 3, buckets[1].Size())
}

func newDay(day, hour, minute int) *time.Time {
	date := time.Date(2018, time.December, day, hour, minute, 0, 0, time.Local)
	return &date
}

func newDate(year int, month time.Month, day, hour, minute int) *time.Time {
	date := time.Date(year, month, day, hour, minute, 0, 0, time.Local)
	return &date
}
