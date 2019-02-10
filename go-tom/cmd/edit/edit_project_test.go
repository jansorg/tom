package edit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/money"
	"github.com/jansorg/tom/go-tom/test_setup"
	"github.com/jansorg/tom/go-tom/util"
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
	err = doEditProjectCommand("my new project name", &parentName, "/", util.StringP("10.50 USD"), []string{p1.ID, p2.ID}, ctx)
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

	assert.EqualValues(t, money.NewMoney(1050, "USD").CurrencyCode(), p2.HourlyRate().CurrencyCode())
	assert.EqualValues(t, money.NewMoney(1050, "USD").Amount(), p2.HourlyRate().Amount())
}

func Test_EditProjectMoveToTop(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("parent", "child 1")
	require.NoError(t, err)

	emptyID := ""
	err = doEditProjectCommand("", &emptyID, "/", nil, []string{p1.ID}, ctx)
	require.NoError(t, err)

	p, err := ctx.Query.ProjectByFullName([]string{"child 1"})
	require.NoError(t, err)
	require.EqualValues(t, p1.ID, p.ID)

	p, err = ctx.Query.ProjectByID(p.ID)
	require.NoError(t, err)
	require.EqualValues(t, "child 1", p.Name)
	require.EqualValues(t, "child 1", p.GetFullName("/"))
}

func Test_EditProjectMoveToOwnChild(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	top, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("parent")
	require.NoError(t, err)

	child, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("parent", "child 1")
	require.NoError(t, err)

	err = doEditProjectCommand("", &child.ID, "/", nil, []string{top.ID}, ctx)
	require.Error(t, err, "moving a project into it's own child scope must fail")

	err = doEditProjectCommand("", &top.ID, "/", nil, []string{top.ID}, ctx)
	require.Error(t, err, "making a project its own child must fail")
}

func Test_EditProjectHourlyRate(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	top, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("parent")
	require.NoError(t, err)

	err = doEditProjectCommand("", &top.ID, "/", util.StringP("100.50 EUR"), []string{top.ID}, ctx)
	top, _, _ = ctx.StoreHelper.GetOrCreateNestedProjectNames("parent")
	assert.EqualValues(t, money.NewMoney(10050, "EUR"), top.HourlyRate())

	err = doEditProjectCommand("", &top.ID, "/", util.StringP("10.75 USD"), []string{top.ID}, ctx)
	top, _, _ = ctx.StoreHelper.GetOrCreateNestedProjectNames("parent")
	assert.EqualValues(t, money.NewMoney(1075, "USD"), top.HourlyRate())

	err = doEditProjectCommand("", &top.ID, "/", util.StringP(""), []string{top.ID}, ctx)
	top, _, _ = ctx.StoreHelper.GetOrCreateNestedProjectNames("parent")
	assert.Nil(t, top.HourlyRate())
}
