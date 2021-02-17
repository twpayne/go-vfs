package vfst

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vfs "github.com/twpayne/go-vfs/v2"
)

func TestWalk(t *testing.T) {
	fs, cleanup, err := NewTestFS(map[string]interface{}{
		"/home/user/.bashrc":  "# .bashrc contents\n",
		"/home/user/skip/foo": "bar",
		"/home/user/symlink":  &Symlink{Target: "baz"},
	})
	require.NoError(t, err)
	defer cleanup()
	pathTypeMap := make(map[string]os.FileMode)
	require.NoError(t, vfs.Walk(fs, "/", func(path string, info os.FileInfo, err error) error {
		assert.NoError(t, err)
		pathTypeMap[filepath.ToSlash(path)] = info.Mode() & os.ModeType
		if filepath.Base(path) == "skip" {
			return vfs.SkipDir
		}
		return nil
	}))
	expectedPathTypeMap := map[string]os.FileMode{
		"/":                  os.ModeDir,
		"/home":              os.ModeDir,
		"/home/user":         os.ModeDir,
		"/home/user/.bashrc": 0,
		"/home/user/skip":    os.ModeDir,
		"/home/user/symlink": os.ModeSymlink,
	}
	assert.Equal(t, expectedPathTypeMap, pathTypeMap)
}
