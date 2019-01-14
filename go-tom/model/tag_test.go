package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateTag(t *testing.T) {
	v := &Tag{}
	assert.Error(t, v.Validate())

	v.ID = "id"
	assert.Error(t, v.Validate())

	v.Name = "name"
	assert.NoError(t, v.Validate())
}


