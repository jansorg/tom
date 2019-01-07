package activity

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/gotime/go-tom/store"
	"github.com/jansorg/gotime/go-tom/testSetup"
)

func Test_Activity(t *testing.T) {
	ctx, err := testSetup.CreateTestContext(language.AmericanEnglish)
	require.NoError(t, err)
	defer testSetup.CleanupTestContext(ctx)

	project, err := ctx.Store.AddProject(store.Project{Name: "Project1"})
	require.NoError(t, err)

	control := NewActivityControl(ctx, false, false)
	frame, err := control.Start(project.ID, "my new activity", []string{})
	require.NoError(t, err)
	require.NotEmpty(t, frame.ID)
	require.EqualValues(t, "my new activity", frame.Notes)
	require.True(t, frame.IsActive())

	stoppedFrame, err := control.StopNewest("my updated notes", []string{})
	require.NoError(t, err)
	require.EqualValues(t, frame.ID, stoppedFrame.ID)
	require.EqualValues(t, "my updated notes", frame.Notes)
	require.False(t, frame.IsActive())
}
