package frames

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/test_setup"
)

func TestArchiveCommand(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project")
	require.NoError(t, err)

	p2, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project", "subproject")
	require.NoError(t, err)

	now := time.Now()
	end := now.Add(10 * time.Minute)
	f1, err := ctx.Store.AddFrame(model.Frame{
		Start:     &now,
		End:       &end,
		ProjectId: p1.ID,
	})
	require.NoError(t, err)

	now = time.Now()
	end = now.Add(10 * time.Minute)
	f2, err := ctx.Store.AddFrame(model.Frame{
		Start:     &now,
		End:       &end,
		ProjectId: p2.ID,
	})
	require.NoError(t, err)

	err = archiveFrames(p1.GetFullName("/"), "/", false, ctx)
	require.NoError(t, err)

	f1New, err := ctx.Query.FrameByID(f1.ID)
	require.NoError(t, err)
	assert.True(t, f1New.Archived)

	f2New, err := ctx.Query.FrameByID(f2.ID)
	require.NoError(t, err)
	assert.False(t, f2New.Archived)

	// now archive frame of subproject p2
	err = archiveFrames(p2.ID, "/", false, ctx)
	f2New, err = ctx.Query.FrameByID(f2.ID)
	require.NoError(t, err)
	assert.True(t, f2New.Archived)
}

func TestArchiveCommandSubprojects(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project")
	require.NoError(t, err)

	p2, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project", "subproject")
	require.NoError(t, err)

	now := time.Now()
	end := now.Add(10 * time.Minute)
	f1, err := ctx.Store.AddFrame(model.Frame{
		Start:     &now,
		End:       &end,
		ProjectId: p1.ID,
	})
	require.NoError(t, err)

	now = time.Now()
	end = now.Add(10 * time.Minute)
	f2, err := ctx.Store.AddFrame(model.Frame{
		Start:     &now,
		End:       &end,
		ProjectId: p2.ID,
	})
	require.NoError(t, err)

	err = archiveFrames(p1.GetFullName("/"), "/", true, ctx)
	require.NoError(t, err)

	f1New, err := ctx.Query.FrameByID(f2.ID)
	require.NoError(t, err)

	f2New, err := ctx.Query.FrameByID(f2.ID)
	require.NoError(t, err)

	assert.True(t, f1New.Archived)
	assert.True(t, f2New.Archived)
}
