package store_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jansorg/gotime/go-tom/store"
)

func Test_StoreNoDataDir(t *testing.T) {
	_, err := store.NewStore(filepath.Join(os.TempDir(), "gotime-does-not-exist"))
	require.Error(t, err)
}

func Test_Store(t *testing.T) {
	dir, err := ioutil.TempDir("", "gotime")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	s, err := store.NewStore(dir)
	require.NoError(t, err)
	require.NoError(t, err)
	assert.Empty(t, s.Projects())
	assert.Empty(t, s.Tags())
	assert.Empty(t, s.Frames())

	// project
	newProject := store.Project{Name: "Project 42"}
	addedProject, err := s.AddProject(newProject)
	require.NoError(t, err)
	assert.EqualValues(t, newProject.Name, addedProject.Name)
	assert.Empty(t, newProject.ID)
	assert.NotEmpty(t, addedProject.ID, "a new ID must be added to the newly created project")

	// 1st removal must succeed
	err = s.RemoveProject(addedProject.ID)
	require.NoError(t, err)
	assert.Empty(t, s.Projects())

	// 2nd removal must fail
	err = s.RemoveProject(addedProject.ID)
	require.Error(t, err)

	// tag
	newTag := store.Tag{Name: "Tag 1"}
	addedTag, err := s.AddTag(newTag)
	require.NoError(t, err)
	assert.EqualValues(t, newTag.Name, addedTag.Name)
	assert.Empty(t, newTag.ID)
	assert.NotEmpty(t, addedTag.ID, "a new ID must be added to the newly created tag")

	// 1st removal must succeed
	err = s.RemoveTag(addedTag.ID)
	require.NoError(t, err)
	assert.Empty(t, s.Tags())

	// 2nd removal must fail
	err = s.RemoveProject(addedTag.ID)
	require.Error(t, err)

	// frames
	addedProject, err = s.AddProject(store.Project{Name: "Project for Frame"})
	require.NoError(t, err)
	startTime, err := time.Parse(time.RFC822, "02 Jan 19 10:00 MST")
	require.NoError(t, err)
	endTime, err := time.Parse(time.RFC822, "02 Jan 19 10:00 MST")
	require.NoError(t, err)
	newFrame := store.Frame{
		ProjectId: addedProject.ID,
		Start:     &startTime,
		End:       &endTime,
	}
	addedFrame, err := s.AddFrame(newFrame)
	require.NoError(t, err)
	assert.NotEmpty(t, addedFrame.ID, "new ID must be added")

	err = s.RemoveFrame(addedFrame.ID)
	require.NoError(t, err)

	// at this point save() must have been called and files must exists
	dataStore := s.(*store.DataStore)
	assert.FileExists(t, dataStore.ProjectFile)
	assert.FileExists(t, dataStore.TagFile)
	assert.FileExists(t, dataStore.FrameFile)
}
