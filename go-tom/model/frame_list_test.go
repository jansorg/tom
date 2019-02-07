package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArchived(t *testing.T) {
	f := NewEmptyFrameList()
	now := time.Now()
	end := time.Now().Add(10 * time.Minute)
	f.Append(&Frame{
		Start:    &now,
		End:      &end,
		Archived: false,
	})
	f.Append(&Frame{
		Start:    &now,
		End:      &end,
		Archived: true,
	})

	assert.EqualValues(t, 2, f.Size())

	f.ExcludeArchived()
	assert.EqualValues(t, 1, f.Size())
}
