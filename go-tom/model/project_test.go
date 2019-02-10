package model

import (
	"encoding/json"
	"testing"

	"github.com/rhymond/go-money"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jansorg/tom/go-tom/properties"
)

func Test_ValidateProject(t *testing.T) {
	p := &Project{}
	assert.Error(t, p.Validate())

	p.ID = "id"
	assert.Error(t, p.Validate())

	p.Name = "name"
	assert.NoError(t, p.Validate())
}

func TestJson(t *testing.T) {
	p := &Project{
		Properties: properties.PropertyValues{
			properties.NewCurrency("myProp", money.New(1000, "EUR")),
		},
	}

	bytes, err := json.Marshal(p)
	require.NoError(t, err)

	var p2 Project
	err = json.Unmarshal(bytes, &p2)
	require.NoError(t, err)

	assert.EqualValues(t, 1, len(p2.Properties))
}

func TestJsonEmpty(t *testing.T) {
	p := &Project{}

	bytes, err := json.Marshal(p)
	require.NoError(t, err)

	var p2 Project
	err = json.Unmarshal(bytes, &p2)
	require.NoError(t, err)

	assert.EqualValues(t, 0, len(p2.Properties))
}
