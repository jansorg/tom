package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	dirName, err := ioutil.TempDir("", "tom")
	require.NoError(t, err)

	source := filepath.Join(dirName, "a.txt")
	err = ioutil.WriteFile(source, []byte("hello world!"), 0600)
	require.NoError(t, err)
	defer os.Remove(source)

	target := filepath.Join(dirName, "b.txt")
	err = CopyFile(source, target, false)
	require.NoError(t, err)
	defer os.Remove(target)

	data, err := ioutil.ReadFile(target)
	require.NoError(t, err)
	assert.EqualValues(t, "hello world!", string(data))

	// overwrite
	err = ioutil.WriteFile(source, []byte("hello world, too!"), 0600)
	require.NoError(t, err)
	err = CopyFile(source, target, true)
	require.NoError(t, err)

	data, err = ioutil.ReadFile(target)
	require.NoError(t, err)
	assert.EqualValues(t, "hello world, too!", string(data))
}
