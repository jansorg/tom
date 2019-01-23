package edit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/test_setup"
)

func Test_EditProject(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("parent", "child 1")
	require.NoError(t, err)

	p2, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("parent", "child 2")
	require.NoError(t, err)

	newParent, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("new parent", "new parent sub 1")
	require.NoError(t, err)

	parentName := newParent.GetFullName("/")
	err = doEditProjectCommand("my new project name", &parentName, "/", []string{p1.ID, p2.ID}, ctx)
	require.NoError(t, err)

	p, err := ctx.Query.ProjectByFullName([]string{"parent", "child 1"})
	require.Error(t, err)
	p, err = ctx.Query.ProjectByFullName([]string{"parent", "child 2"})
	require.Error(t, err)

	p, err = ctx.Query.ProjectByID(p1.ID)
	require.NoError(t, err)
	assert.EqualValues(t, newParent.ID, p.ParentID)
	assert.EqualValues(t, "my new project name", p.Name)

	p, err = ctx.Query.ProjectByID(p2.ID)
	require.NoError(t, err)
	assert.EqualValues(t, newParent.ID, p.ParentID)
	assert.EqualValues(t, "my new project name", p.Name)
}
