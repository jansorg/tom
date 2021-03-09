package dateTime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/i18n"
)

func TestAdd(t *testing.T) {
	a := NewDurationSum()
	a.Add(10 * time.Minute)
	a.Add(15 * time.Minute)
	a.Add(1 * time.Minute)
	a.Add(30 * time.Second)
	a.Add(30 * time.Second)

	assert.EqualValues(t, 27*time.Minute, a.Get())
	assert.EqualValues(t, 27*time.Minute, a.GetExact())
}

func TestAddRounding(t *testing.T) {
	// one hour overlap
	start := time.Date(2018, time.January, 10, 12, 0, 0, 0, time.UTC)
	end := time.Date(2018, time.January, 10, 13, 0, 0, 0, time.UTC)
	dateRange := NewDateRange(&start, &end, i18n.FindLocale(language.English, false))

	a := NewDurationSumAll(RoundingUp(6*time.Minute), &dateRange, &end)
	a.Add(10 * time.Minute)
	a.Add(15 * time.Minute)
	a.Add(1 * time.Minute)
	a.Add(30 * time.Second)
	a.Add(30 * time.Second)

	assert.EqualValues(t, (12+18+6+6+6)*time.Minute, a.Get())
	assert.EqualValues(t, (10+15+1+0.5+0.5)*time.Minute, a.GetExact())

	// not add some date ranges
	before := time.Date(2018, time.January, 8, 12, 0, 0, 0, time.UTC)
	after := time.Date(2018, time.January, 20, 12, 0, 0, 0, time.UTC)

	a.AddStartEnd(before, after)
	assert.EqualValues(t, 1*time.Hour+(12+18+6+6+6)*time.Minute, a.Get())
	assert.EqualValues(t, 1*time.Hour+(10+15+1+0.5+0.5)*time.Minute, a.GetExact())
}
