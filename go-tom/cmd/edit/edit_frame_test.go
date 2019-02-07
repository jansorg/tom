package edit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/test_setup"
	"github.com/jansorg/tom/go-tom/util"
)

func TestEditFrame(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top", "p1")
	require.NoError(t, err)

	p2, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("p2")
	require.NoError(t, err)

	start := time.Date(2018, time.January, 10, 0, 0, 0, 0, time.UTC)
	end := start.Add(6 * time.Hour)

	f1, err := ctx.Store.AddFrame(model.Frame{Start: &start, End: &end, ProjectId: p1.ID})
	require.NoError(t, err)

	start2 := time.Date(2018, time.January, 10, 1, 0, 0, 0, time.UTC)
	f2, err := ctx.Store.AddFrame(model.Frame{Start: &start2, ProjectId: p1.ID})
	require.NoError(t, err)

	newEnd := end.Add(6 * time.Hour)
	newEndString := newEnd.Format(time.RFC3339)
	newNotes := "my new notes"
	err = doEditFrameCommand(ctx, []string{f1.ID, f2.ID}, nil, &newEndString, &newNotes, &(p2.ID), "/", nil)
	require.NoError(t, err)

	newF1, err := ctx.Query.FrameByID(f1.ID)
	require.NoError(t, err)
	assert.EqualValues(t, "my new notes", newF1.Notes)
	assert.EqualValues(t, p2.ID, newF1.ProjectId)
	assert.EqualValues(t, &start, newF1.Start, "start time must not be modified if no new value was passed")
	assert.EqualValues(t, &newEnd, newF1.End)

	newF2, err := ctx.Query.FrameByID(f2.ID)
	require.NoError(t, err)
	assert.EqualValues(t, "my new notes", newF2.Notes)
	assert.EqualValues(t, p2.ID, newF2.ProjectId)
	assert.EqualValues(t, &start2, newF2.Start, "start time must not be modified if no new value was passed")
	assert.EqualValues(t, &newEnd, newF2.End)

	// update f2 to use p1
	projectName1 := p1.GetFullName("/")
	err = doEditFrameCommand(ctx, []string{f2.ID}, nil, nil, nil, &projectName1, "/", nil)
	require.NoError(t, err)
	newF2, err = ctx.Query.FrameByID(f2.ID)
	require.NoError(t, err)
	assert.EqualValues(t, p1.ID, newF2.ProjectId)

	// update f2 to be archived
	err = doEditFrameCommand(ctx, []string{f2.ID}, nil, nil, nil, nil, "/", util.TrueP())
	require.NoError(t, err)
	assert.EqualValues(t, true, newF2.Archived)
}

func TestEditFrameErrors(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p1, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames("top", "p1")
	require.NoError(t, err)

	start := time.Date(2018, time.January, 10, 1, 0, 0, 0, time.UTC)
	f, err := ctx.Store.AddFrame(model.Frame{Start: &start, ProjectId: p1.ID})
	require.NoError(t, err)

	empty := ""
	err = doEditFrameCommand(ctx, []string{f.ID}, nil, nil, nil, &empty, "/", nil)
	require.Error(t, err, "empty project must not be accepted")

	name := "Invalid/project/name"
	err = doEditFrameCommand(ctx, []string{f.ID}, nil, nil, nil, &name, "/", nil)
	require.Error(t, err, "not existing project must not be accepted")

	err = doEditFrameCommand(ctx, []string{"does not exist"}, nil, nil, nil, &empty, "/", nil)
	require.Error(t, err, "invalid frame id must not be accepted")
}
