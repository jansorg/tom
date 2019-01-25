package query_test

import (
	"testing"
	"time"

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

func Test_RecentlyTracked(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, err := ctx.Store.AddProject(model.Project{Name: "Project1"})
	require.NoError(t, err)

	p2, err := ctx.Store.AddProject(model.Project{Name: "Project2"})
	require.NoError(t, err)

	p3, err := ctx.Store.AddProject(model.Project{Name: "Project3"})
	require.NoError(t, err)

	recent, err := ctx.Query.FindRecentlyTrackedProjects(5)
	require.NoError(t, err)
	assert.Empty(t, recent, "no projects without frames")

	start := time.Now()
	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p1.ID, Start: &start})
	require.NoError(t, err)
	recent, err = ctx.Query.FindRecentlyTrackedProjects(5)
	require.NoError(t, err)
	assert.EqualValues(t, 1, recent.Size(), "expected the only project")
	assert.EqualValues(t, p1.ID, recent.First().ID, "expected the only project")

	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p2.ID, Start: &start})
	require.NoError(t, err)

	_, err = ctx.Store.AddFrame(model.Frame{ProjectId: p3.ID, Start: &start})
	require.NoError(t, err)

	recent, err = ctx.Query.FindRecentlyTrackedProjects(5)
	require.NoError(t, err)
	assert.EqualValues(t, 3, recent.Size(), "expected the only project")
	assert.EqualValues(t, p3.ID, recent[0].ID, "expected the only project")
	assert.EqualValues(t, p2.ID, recent[1].ID, "expected the only project")
	assert.EqualValues(t, p1.ID, recent[2].ID, "expected the only project")
}
