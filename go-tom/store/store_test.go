package store_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/store"
	"github.com/jansorg/tom/go-tom/test_setup"
)

func Test_StoreNoDataDir(t *testing.T) {
	_, err := store.NewStore(filepath.Join(os.TempDir(), "tom-does-not-exist"), filepath.Join(os.TempDir(), "backup-dir"), 5)
	require.Error(t, err)
}

func Test_Store(t *testing.T) {
	dir, err := ioutil.TempDir("", "tom")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	backupDir, err := ioutil.TempDir("", "tom-backup")
	require.NoError(t, err)
	defer os.RemoveAll(backupDir)

	s, err := store.NewStore(dir, backupDir, 10)
	require.NoError(t, err)
	require.NoError(t, err)
	assert.Empty(t, s.Projects())
	assert.Empty(t, s.Tags())
	assert.Empty(t, s.Frames())

	// project
	newProject := model.Project{Name: "Project 42"}
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
	newTag := model.Tag{Name: "Tag 1"}
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
	addedProject, err = s.AddProject(model.Project{Name: "Project for Frame"})
	require.NoError(t, err)
	startTime, err := time.Parse(time.RFC822, "02 Jan 19 10:00 MST")
	require.NoError(t, err)
	endTime, err := time.Parse(time.RFC822, "02 Jan 19 10:00 MST")
	require.NoError(t, err)
	newFrame := model.Frame{
		ProjectId: addedProject.ID,
		Start:     &startTime,
		End:       &endTime,
	}
	addedFrame, err := s.AddFrame(newFrame)
	require.NoError(t, err)
	assert.NotEmpty(t, addedFrame.ID, "new ID must be added")

	err = s.RemoveFrame(addedFrame.ID)
	require.NoError(t, err)
	frames, err := s.FindFrames(func(f *model.Frame) (bool, error) {
		return f.ID == addedFrame.ID, nil
	})
	require.NoError(t, err)
	require.Empty(t, frames)

	// at this point save() must have been called and files must exists
	dataStore := s.(*store.DataStore)
	assert.FileExists(t, dataStore.ProjectFile)
	assert.FileExists(t, dataStore.TagFile)
	assert.FileExists(t, dataStore.FrameFile)

	// several backups have to exist at this points
	files, err := ioutil.ReadDir(backupDir)
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, name := range files {
		assert.True(t, strings.HasPrefix(name.Name(), "tom-"))
	}
}

func TestBackups(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	require.EqualValues(t, 0, countBackups(ctx.Store.BackupDirPath()))

	for i := 1; i <= 20; i++ {
		_, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames(fmt.Sprintf("project-%d", i))
		require.NoError(t, err)

		if i == 1 {
			require.EqualValues(t, 1, countBackups(ctx.Store.BackupDirPath()))
		} else if i <= ctx.Store.MaxBackups() {
			require.EqualValues(t, i, countBackups(ctx.Store.BackupDirPath()))
		} else {
			require.EqualValues(t, ctx.Store.MaxBackups(), countBackups(ctx.Store.BackupDirPath()))
		}
	}

	// open latest dir and check number of projects
	dirs, err := sortedBackupDirs(ctx.Store.BackupDirPath())
	require.NoError(t, err)
	newStore, err := store.NewStore(dirs[len(dirs)-1], "", 1)
	require.NoError(t, err)

	assert.EqualValues(t, 20, len(ctx.Store.Projects()), "expected backup to contain latest set of projects")
	assert.EqualValues(t, 19, len(newStore.Projects()), "expected backup to contain latest set of projects, 1 less than the live data")
}

func TestBackupsBatchMode(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	ctx.Store.StartBatch()

	require.EqualValues(t, 0, countBackups(ctx.Store.BackupDirPath()))

	for i := 1; i <= 20; i++ {
		_, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames(fmt.Sprintf("project-%d", i))
		require.NoError(t, err)
		require.EqualValues(t, 0, countBackups(ctx.Store.BackupDirPath()))
	}

	ctx.Store.StopBatch()
	require.EqualValues(t, 1, countBackups(ctx.Store.BackupDirPath()))

	_, _, err = ctx.StoreHelper.GetOrCreateNestedProjectNames(fmt.Sprintf("project-%d", 21))

	// open latest dir and check number of projects
	dirs, err := sortedBackupDirs(ctx.Store.BackupDirPath())
	require.NoError(t, err)
	newStore, err := store.NewStore(dirs[len(dirs)-1], "", 1)
	require.NoError(t, err)

	assert.EqualValues(t, 21, len(ctx.Store.Projects()), "expected backup to contain latest set of projects")
	assert.EqualValues(t, 20, len(newStore.Projects()), "expected backup to contain latest set of projects, 1 less than the live data")
}

func countBackups(path string) int {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return -1
	}

	c := 0
	for _, i := range infos {
		if strings.HasPrefix(i.Name(), "tom-") {
			c++
		}
	}

	return c
}

func sortedBackupDirs(path string) ([]string, error) {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, i := range infos {
		if strings.HasPrefix(i.Name(), "tom-") {
			dirs = append(dirs, filepath.Join(path, i.Name()))
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		return strings.Compare(dirs[i], dirs[j]) < 0
	})

	return dirs, nil
}
