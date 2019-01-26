package remove

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/test_setup"
)

func Test_RemoveProjects(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	top, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top")
	require.NoError(t, err)

	child1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top", "child1")
	require.NoError(t, err)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top", "child1", "child1.1")
	require.NoError(t, err)

	p2, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top", "child1", "child1.2")
	require.NoError(t, err)

	now := time.Now()
	end := now.Add(10 * time.Minute)

	// frame in child1.1
	f1, err := ctx.Store.AddFrame(model.Frame{Start: &now, End: &end, ProjectId: p1.ID})
	require.NoError(t, err)

	// frame in child1.2
	f2, err := ctx.Store.AddFrame(model.Frame{Start: &now, End: &end, ProjectId: p2.ID})
	require.NoError(t, err)

	// 2 frames in top
	fTop1, err := ctx.Store.AddFrame(model.Frame{Start: &now, End: &end, ProjectId: top.ID})
	require.NoError(t, err)
	_, err = ctx.Store.AddFrame(model.Frame{Start: &now, End: &end, ProjectId: top.ID})
	require.NoError(t, err)

	// 2 frames in top/child1
	_, err = ctx.Store.AddFrame(model.Frame{Start: &now, End: &end, ProjectId: child1.ID})
	require.NoError(t, err)
	_, err = ctx.Store.AddFrame(model.Frame{Start: &now, End: &end, ProjectId: child1.ID})
	require.NoError(t, err)

	projects, frames, err := ctx.StoreHelper.RemoveProject(p1)
	require.NoError(t, err)
	require.EqualValues(t, 1, projects)
	require.EqualValues(t, 1, frames)
	_, err = ctx.Query.FrameByID(f1.ID)
	require.Error(t, err, "the frame must ahve been deleted")

	// remove top and hierarchy below
	projects, frames, err = ctx.StoreHelper.RemoveProject(top)
	require.NoError(t, err)
	_, err = ctx.Store.ProjectByID(top.ID)
	require.Error(t, err)

	require.EqualValues(t, 3, projects, "top, top/child1, and top/child1/child1.2 must have been removed")
	require.EqualValues(t, 5, frames, "2 frames of top, 2 of top/child1 and 1 frame of child1.2 must have been removed")
	_, err = ctx.Query.FrameByID(f2.ID)
	require.Error(t, err)
	_, err = ctx.Query.FrameByID(fTop1.ID)
	require.Error(t, err)
}
