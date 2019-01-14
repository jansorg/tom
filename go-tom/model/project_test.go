package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateProject(t *testing.T) {
	p := &Project{}
	assert.Error(t, p.Validate())

	p.ID = "id"
	assert.Error(t, p.Validate())

	p.Name = "name"
	assert.NoError(t, p.Validate())
}

