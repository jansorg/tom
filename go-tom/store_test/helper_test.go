package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/test_setup"
)

func Test_RenameProject(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.German)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	_, _, err = ctx.StoreHelper.GetOrCreateNestedProjectNames("top1", "childExisting")
	require.NoError(t, err)

	_, _, err = ctx.StoreHelper.GetOrCreateNestedProjectNames("top2", "childExisting")
	require.NoError(t, err)

	c1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top1", "child")
	require.NoError(t, err)

	// top1/child -> childNewName
	renamed, err := ctx.StoreHelper.RenameProject(c1, []string{"childNewName"})
	require.NoError(t, err)
	require.EqualValues(t, "childNewName", renamed.Name)
	require.EqualValues(t, "childNewName", renamed.GetFullName("/"))

	// childNewName -> top1/childNewName
	renamed, err = ctx.StoreHelper.RenameProject(renamed, []string{"top1", "childNewName"})
	require.NoError(t, err)
	require.EqualValues(t, "childNewName", renamed.Name)
	require.EqualValues(t, "top1/childNewName", renamed.GetFullName("/"))

	// top1/childNewName -> top1/child
	renamed, err = ctx.StoreHelper.RenameProjectByName("top1/childNewName", "top1/child")
	require.NoError(t, err, "renaming under same parent must succeed")
	require.EqualValues(t, "child", renamed.Name)
	require.EqualValues(t, "top1/child", renamed.GetFullName("/"))
	require.EqualValues(t, "top1", renamed.Parent().GetFullName("/"))

	// top1/child -> top1/childExisting
	renamed, err = ctx.StoreHelper.RenameProjectByName("top1/child", "top1/childExisting")
	require.Error(t, err, "renaming like an existing child must fail")

	// top1/child -> top2/childExisting
	renamed, err = ctx.StoreHelper.RenameProjectByName("top1/child", "top2/childExisting")
	require.Error(t, err, "renaming like an existing child must fail")

	// top1/child -> top2/childNewName
	renamed, err = ctx.StoreHelper.RenameProjectByName("top1/child", "top2/childNewName")
	require.NoError(t, err, "moving a child from one parent to another must succeed")
	require.EqualValues(t, "childNewName", renamed.Name)
	require.EqualValues(t, "top2/childNewName", renamed.GetFullName("/"))

	// top2/childNewName -> newParent/sub/childName
	renamed, err = ctx.StoreHelper.RenameProjectByName("top2/childNewName", "newParent/sub/childName")
	require.NoError(t, err, "moving a child from one parent to another must succeed")
	require.EqualValues(t, "childName", renamed.Name)
	require.EqualValues(t, "newParent/sub/childName", renamed.GetFullName("/"))

	newParent, err := ctx.Query.ProjectByFullName([]string{"newParent", "sub"})
	require.NoError(t, err)
	require.EqualValues(t, newParent.ID, renamed.ParentID, "parent ids must match after rename")

	// fixme test that a rename of a parent will also update the fullnames of all child projects
}
