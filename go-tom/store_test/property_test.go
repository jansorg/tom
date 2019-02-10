package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/properties"
	"github.com/jansorg/tom/go-tom/test_setup"
)

func TestProperties(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	ctx.StoreHelper.GetOrCreateNestedProjectNames("top")
	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top", "p1")
	p1_sub, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top", "p1", "sub1")
	p2, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top", "p2")

	hourlyRate, err := ctx.Store.AddProperty(&properties.Property{
		Name:               "hourlyRate",
		ApplyToSubprojects: true,
		TypeID:             properties.CurrencyType.ID(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, hourlyRate.ID)

	_, err = ctx.Store.AddProperty(&properties.Property{
		Name:   "hourlyRate",
		TypeID: properties.CurrencyType.ID(),
	})
	require.Error(t, err)

	err = p1.SetPropertyValue(hourlyRate.ID, "70.0")
	require.Error(t, err)

	err = p1.SetPropertyValue(hourlyRate.ID, "70 EUR")
	require.NoError(t, err)
	p1, err = ctx.Store.UpdateProject(*p1)
	require.NoError(t, err)

	value, err := ctx.Query.FindPropertyValue(hourlyRate.ID, p1.ID)
	require.NoError(t, err, "expected a property value for p1")
	currencyValue, err := value.AsCurrency()
	require.NoError(t, err)
	assert.EqualValues(t, "€70.00", currencyValue.Display())

	value, err = ctx.Query.FindPropertyValue(hourlyRate.ID, p1_sub.ID)
	require.NoError(t, err)
	currencyValue, err = value.AsCurrency()
	require.NoError(t, err)
	assert.EqualValues(t, "€70.00", currencyValue.Display())

	value, err = ctx.Query.FindPropertyValue(hourlyRate.ID, p2.ID)
	require.Error(t, err, "p2 must not inherit property value from p1")
}
