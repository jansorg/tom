package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/model"
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
	renamed, err := ctx.StoreHelper.RenameProject(c1, []string{"childNewName"}, true)
	require.NoError(t, err)
	require.EqualValues(t, "childNewName", renamed.Name)
	require.EqualValues(t, "childNewName", renamed.GetFullName("/"))

	// childNewName -> top1/childNewName without allowd hierarchy change must fail
	_, err = ctx.StoreHelper.RenameProject(renamed, []string{"top1", "childNewName"}, false)
	require.Error(t, err)

	// childNewName -> top1/childNewName
	renamed, err = ctx.StoreHelper.RenameProject(renamed, []string{"top1", "childNewName"}, true)
	require.NoError(t, err)
	require.EqualValues(t, "childNewName", renamed.Name)
	require.EqualValues(t, "top1/childNewName", renamed.GetFullName("/"))

	// top1/childNewName -> top1/childNewName2 (with hierarchyUpdate disabled)
	renamed, err = ctx.StoreHelper.RenameProject(renamed, []string{"childNewName2"}, false)
	require.NoError(t, err)
	require.EqualValues(t, "childNewName2", renamed.Name)
	require.EqualValues(t, "top1/childNewName2", renamed.GetFullName("/"))

	// top1/childNewName -> top1/child
	renamed, err = ctx.StoreHelper.RenameProjectByIDOrName("top1/childNewName2", "top1/child")
	require.NoError(t, err, "renaming under same parent must succeed")
	require.EqualValues(t, "child", renamed.Name)
	require.EqualValues(t, "top1/child", renamed.GetFullName("/"))
	require.EqualValues(t, "top1", renamed.Parent().GetFullName("/"))

	// top1/child -> top1/childExisting
	renamed, err = ctx.StoreHelper.RenameProjectByIDOrName("top1/child", "top1/childExisting")
	require.Error(t, err, "renaming like an existing child must fail")

	// top1/child -> top2/childExisting
	renamed, err = ctx.StoreHelper.RenameProjectByIDOrName("top1/child", "top2/childExisting")
	require.Error(t, err, "renaming like an existing child must fail")

	// top1/child -> top2/childNewName
	renamed, err = ctx.StoreHelper.RenameProjectByIDOrName("top1/child", "top2/childNewName")
	require.NoError(t, err, "moving a child from one parent to another must succeed")
	require.EqualValues(t, "childNewName", renamed.Name)
	require.EqualValues(t, "top2/childNewName", renamed.GetFullName("/"))

	// top2/childNewName -> newParent/sub/childName
	renamed, err = ctx.StoreHelper.RenameProjectByIDOrName("top2/childNewName", "newParent/sub/childName")
	require.NoError(t, err, "moving a child from one parent to another must succeed")
	require.EqualValues(t, "childName", renamed.Name)
	require.EqualValues(t, "newParent/sub/childName", renamed.GetFullName("/"))

	newParent, err := ctx.Query.ProjectByFullName([]string{"newParent", "sub"})
	require.NoError(t, err)
	require.EqualValues(t, newParent.ID, renamed.ParentID, "parent ids must match after rename")

	// fixme test that a rename of a parent will also update the fullnames of all child projects
}

func Test_RenameProjectEmpty(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p, err := ctx.Store.AddProject(model.Project{Name: "my project"})
	require.NoError(t, err)

	p, err = ctx.StoreHelper.RenameProjectByIDOrName(p.ID, "my new project")
	require.NoError(t, err)
	assert.EqualValues(t, "my new project", p.Name)

	_, err = ctx.StoreHelper.RenameProjectByIDOrName(p.ID, "")
	require.Errorf(t, err, "renaming to an empty name must fail")

	_, err = ctx.StoreHelper.RenameProjectByIDOrName(p.ID, "   ")
	require.Errorf(t, err, "renaming to a blank name must fail")
}
