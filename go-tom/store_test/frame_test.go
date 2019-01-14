package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/test_setup"
)

func TestFrames(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	p, err := ctx.Store.AddProject(model.Project{Name: "p1"})
	require.NoError(t, err)

	frame, err := ctx.Store.AddFrame(model.NewStartedFrame(p))
	require.NoError(t, err)

	frame.ID = ""
	_, err = ctx.Store.UpdateFrame(*frame)
	require.Error(t, err)

	frame.ID = "my ID"
	frame.Start = nil
	_, err = ctx.Store.UpdateFrame(*frame)
	require.Error(t, err)
}
