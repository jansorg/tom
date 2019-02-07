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

func TestShowExcluded(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project")
	require.NoError(t, err)

	p2, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("project", "subproject")
	require.NoError(t, err)

	now := time.Now()
	end := now.Add(10 * time.Minute)
	_, err = ctx.Store.AddFrame(model.Frame{
		Start:     &now,
		End:       &end,
		ProjectId: p1.ID,
	})
	require.NoError(t, err)

	_, err = ctx.Store.AddFrame(model.Frame{
		Start:     &now,
		End:       &end,
		ProjectId: p1.ID,
		Archived:  true,
	})
	require.NoError(t, err)

	_, err = ctx.Store.AddFrame(model.Frame{
		Start:     &now,
		End:       &end,
		ProjectId: p2.ID,
		Archived:  true,
	})
	require.NoError(t, err)

	frames := filterFrames("", ctx, true, true)
	assert.EqualValues(t, frames.Size(), 3)

	frames = filterFrames("", ctx, false, true)
	assert.EqualValues(t, frames.Size(), 3)

	frames = filterFrames(p1.ID, ctx, true, true)
	assert.EqualValues(t, frames.Size(), 3)

	frames = filterFrames("", ctx, true, false)
	assert.EqualValues(t, frames.Size(), 1)
}
