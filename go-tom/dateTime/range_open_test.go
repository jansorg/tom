package dateTime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntersectionHalfOpenRange(t *testing.T) {
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	start := ref.Add(time.Hour)
	end := ref.Add(3 * time.Hour)

	onlyStart := DateRange{Start: &ref}
	assert.Equal(t, 2*time.Hour, onlyStart.Intersection(&start, &end))

	cut := ref.Add(2 * time.Hour)
	onlyEnd := DateRange{End: &cut}
	assert.Equal(t, time.Hour, onlyEnd.Intersection(&start, &end))

	after := ref.Add(5 * time.Hour)
	assert.Equal(t, time.Duration(0), DateRange{Start: &after}.Intersection(&start, &end))
	assert.False(t, DateRange{Start: &after}.Intersects(&start, &end))
}
