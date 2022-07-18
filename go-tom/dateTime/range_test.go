package dateTime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/i18n"
)

func TestDurationIntersection(t *testing.T) {
	// 2018-01-01, 12am
	start := time.Date(2018, time.January, 1, 12, 0, 0, 0, time.UTC)
	// 2018-01-02, 4pm --> 28 hours
	end := time.Date(2018, time.January, 2, 16, 0, 0, 0, time.UTC)
	// 2018-01-02, 4pm --> 28 hours

	wayBefore := time.Date(2016, time.December, 24, 12, 0, 0, 0, time.UTC)
	before := time.Date(2017, time.December, 24, 12, 0, 0, 0, time.UTC)
	between := time.Date(2018, time.January, 1, 20, 0, 0, 0, time.UTC)
	after := time.Date(2018, time.February, 1, 12, 0, 0, 0, time.UTC)
	wayAfter := time.Date(2019, time.February, 1, 12, 0, 0, 0, time.UTC)

	r := NewDateRange(&start, &end, i18n.FindLocale(language.English, false))

	assert.EqualValues(t, 28*time.Hour, r.Intersection(&start, &end))
	assert.EqualValues(t, 28*time.Hour, r.Intersection(&before, &end))
	assert.EqualValues(t, 28*time.Hour, r.Intersection(&start, &after))
	assert.EqualValues(t, 28*time.Hour, r.Intersection(&before, &after))

	assert.EqualValues(t, 8*time.Hour, r.Intersection(&before, &between))
	assert.EqualValues(t, 20*time.Hour, r.Intersection(&between, &end))
	assert.EqualValues(t, 20*time.Hour, r.Intersection(&between, &after))

	assert.EqualValues(t, 0, r.Intersection(&wayBefore, &before))
	assert.EqualValues(t, 0, r.Intersection(&after, &wayAfter))
}

func TestOpenRange(t *testing.T) {
	// 2018-01-01, 12am
	start := time.Date(2018, time.January, 5, 12, 0, 0, 0, time.UTC)
	// 2018-01-02, 4pm --> 28 hours
	end := time.Date(2018, time.January, 6, 16, 0, 0, 0, time.UTC)
	// 2018-01-02, 4pm --> 28 hours

	r := NewDateRange(&start, nil, i18n.FindLocale(language.English, false))
	require.EqualValues(t, "1/5/18 –", r.MinimalString())
	require.EqualValues(t, "2018-01-05 –", r.ShortString())

	r = NewDateRange(nil, &end, i18n.FindLocale(language.English, false))
	require.EqualValues(t, "– 1/6/18", r.MinimalString())
	require.EqualValues(t, "– 2018-01-06", r.ShortString())

	r = NewDateRange(&start, &end, i18n.FindLocale(language.English, false))
	require.EqualValues(t, "1/5/18 – 1/6/18", r.MinimalString())
	require.EqualValues(t, "2018-01-05 – 2018-01-06", r.ShortString())
}
