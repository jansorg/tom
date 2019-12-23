package activity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/test_setup"
	"github.com/jansorg/tom/go-tom/util"
)

func Test_Activity(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.AmericanEnglish)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	project, err := ctx.Store.AddProject(model.Project{Name: "Project1"})
	require.NoError(t, err)

	control := NewActivityControl(ctx, false, false, time.Now())
	frame, err := control.Start(project.ID, "my new activity", []*model.Tag{})
	require.NoError(t, err)
	require.NotEmpty(t, frame.ID)
	require.EqualValues(t, "my new activity", frame.Notes)
	require.True(t, frame.IsActive())

	stoppedFrame, err := control.StopNewest("my updated notes", []*model.Tag{})
	require.NoError(t, err)
	require.EqualValues(t, frame.ID, stoppedFrame.ID)
	require.EqualValues(t, "my updated notes", stoppedFrame.Notes)
	require.False(t, stoppedFrame.IsActive())
}

func Test_ActivityNoteRequired(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.AmericanEnglish)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	project, err := ctx.Store.AddProject(model.Project{Name: "Project1"})
	require.NoError(t, err)
	project.SetNoteRequired(util.TrueP())
	project, err = ctx.Store.UpdateProject(*project)
	require.NoError(t, err)

	control := NewActivityControl(ctx, false, false, time.Now())
	frame, err := control.Start(project.ID, "my new activity", []*model.Tag{})
	require.NoError(t, err)
	require.NotEmpty(t, frame.ID)
	require.EqualValues(t, "my new activity", frame.Notes)
	require.True(t, frame.IsActive())

	stoppedFrame, err := control.StopNewest("", []*model.Tag{})
	require.EqualValues(t, NoteRequiredErr, err)

	stoppedFrame, err = control.StopNewest("note with content", []*model.Tag{})
	require.NoError(t, err)
	require.EqualValues(t, "note with content", stoppedFrame.Notes)
}
