package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRounding(t *testing.T) {
	for value, expected := range map[string]int64{
		"19.99 EUR": 1999,
		"0.29 EUR":  29,
		"USD 1.15":  115,
		"100 EUR":   10000,
	} {
		m, err := Parse(value)
		require.NoError(t, err, value)
		assert.EqualValues(t, expected, m.Amount(), value)
	}
}

func TestParseInvalid(t *testing.T) {
	for _, value := range []string{"100", "", "100 ABCQ", "abc EUR"} {
		_, err := Parse(value)
		assert.Error(t, err, value)
	}
}

func TestMultipleRounding(t *testing.T) {
	m := NewMoney(1999, "EUR")
	assert.EqualValues(t, 3998, m.Multiple(2.0).Amount())
	assert.EqualValues(t, 1000, NewMoney(6000, "EUR").Multiple(1.0/6.0).Amount())
}
