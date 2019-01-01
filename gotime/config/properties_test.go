package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Data struct {
	p map[string]string
}

func (d Data) GetProperties() map[string]string {
	return d.p
}

func Test_PropertyTests(t *testing.T) {
	data := Data{
		p: map[string]string{},
	}

	_, ok := DescriptionProperty.Get(data)
	assert.False(t, ok)

	DescriptionProperty.Set("mine", data)
	s, ok := DescriptionProperty.Get(data)
	assert.True(t, ok)
	assert.EqualValues(t, "mine", s)

	_, ok = AddressProperty.Get(data)
	assert.False(t, ok)

	AddressProperty.Set("my address", data)
	s, ok = AddressProperty.Get(data)
	assert.True(t, ok)
	assert.EqualValues(t, "my address", s)
}
