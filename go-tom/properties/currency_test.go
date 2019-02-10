package properties

import (
	"encoding/json"
	"testing"

	"github.com/rhymond/go-money"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJson(t *testing.T) {
	v := CurrencyValue{value: money.New(100, "EUR")}
	assert.True(t, v.IsCurrency())
	assert.False(t, v.IsFloat())
	assert.False(t, v.IsString())

	bytes, err := json.Marshal(v)
	require.NoError(t, err)

	newValue := CurrencyValue{}
	err = json.Unmarshal(bytes, &newValue)
	require.NoError(t, err)

	equal, err := v.value.Equals(newValue.value)
	require.NoError(t, err)
	assert.True(t, equal)

	m, err := newValue.AsCurrency()
	assert.EqualValues(t, 100, m.Amount())
	assert.EqualValues(t, "EUR", m.Currency().Code)
}
