package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/config"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/test_setup"
)

func Test_InheritedProps(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	parent, err := ctx.Store.AddProject(model.Project{
		Name: "Top",
	})
	require.NoError(t, err)

	child, err := ctx.Store.AddProject(model.Project{
		Name:     "Top > Sub",
		ParentID: parent.ID,
	})
	require.NoError(t, err)

	_, ok := ctx.Query.GetInheritedStringProp(parent.ID, config.InvoiceDescriptionProperty)
	assert.False(t, ok)

	_, ok = ctx.Query.GetInheritedStringProp(child.ID, config.InvoiceDescriptionProperty)
	assert.False(t, ok)

	config.InvoiceDescriptionProperty.Set("top description", parent)
	v, ok := ctx.Query.GetInheritedStringProp(parent.ID, config.InvoiceDescriptionProperty)
	assert.EqualValues(t, "top description", v)
	assert.True(t, ok)

	v, ok = ctx.Query.GetInheritedStringProp(child.ID, config.InvoiceDescriptionProperty)
	assert.EqualValues(t, "top description", v)
	assert.True(t, ok)
}
