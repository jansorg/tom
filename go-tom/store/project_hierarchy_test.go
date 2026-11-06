package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/test_setup"
)

func TestRenameParentUpdatesSubprojects(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	child, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("a", "b", "c")
	require.NoError(t, err)

	_, err = ctx.StoreHelper.RenameProjectByIDOrName("a", "x")
	require.NoError(t, err)

	assert.EqualValues(t, []string{"x", "b", "c"}, child.FullName)
	_, err = ctx.Query.ProjectByFullName([]string{"x", "b", "c"})
	assert.NoError(t, err)
	_, err = ctx.Query.ProjectByFullName([]string{"a", "b", "c"})
	assert.Error(t, err)
}

func TestRejectProjectCycles(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	b, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("a", "b")
	require.NoError(t, err)
	a := b.Parent()
	require.NotNil(t, a)

	// rename into itself or into its own subtree
	_, err = ctx.StoreHelper.RenameProjectByIDOrName("a", "a/new")
	assert.Error(t, err)
	_, err = ctx.StoreHelper.RenameProjectByIDOrName("a", "a/b/new")
	assert.Error(t, err)
	assert.Len(t, ctx.Store.Projects(), 2, "no intermediate projects must be created")

	// direct updates
	updated := *a
	updated.ParentID = a.ID
	_, err = ctx.Store.UpdateProject(updated)
	assert.Error(t, err)

	updated.ParentID = b.ID
	_, err = ctx.Store.UpdateProject(updated)
	assert.Error(t, err)

	updated.ParentID = "does-not-exist"
	_, err = ctx.Store.UpdateProject(updated)
	assert.Error(t, err)

	assert.EqualValues(t, "", a.ParentID)
	assert.EqualValues(t, []string{"a", "b"}, b.FullName)
}
