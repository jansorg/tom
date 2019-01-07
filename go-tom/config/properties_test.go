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

	_, ok := InvoiceDescriptionProperty.Get(data)
	assert.False(t, ok)

	InvoiceDescriptionProperty.Set("mine", data)
	s, ok := InvoiceDescriptionProperty.Get(data)
	assert.True(t, ok)
	assert.EqualValues(t, "mine", s)

	_, ok = InvoiceAddressProperty.Get(data)
	assert.False(t, ok)

	InvoiceAddressProperty.Set("my address", data)
	s, ok = InvoiceAddressProperty.Get(data)
	assert.True(t, ok)
	assert.EqualValues(t, "my address", s)
}
