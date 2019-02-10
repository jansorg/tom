package properties

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrencyType(t *testing.T) {
	assertParsed(t, "100 EUR", "EUR", 100*100)
	assertParsed(t, "EUR 100", "EUR", 100*100)

	assertParsed(t, "100.20 EUR", "EUR", 100*100+20)
	assertParsed(t, "EUR 100.20", "EUR", 100*100+20)
}

func assertParsed(t *testing.T, v string, unit string, amount int64) {
	value, err := CurrencyType.Parse(v, "id")
	require.NoError(t, err)
	assert.EqualValues(t, "id", value.PropertyID())
	c, err := value.AsCurrency()
	require.NoError(t, err)
	assert.EqualValues(t, amount, c.Amount())
	assert.EqualValues(t, unit, c.Currency().Code)
}
