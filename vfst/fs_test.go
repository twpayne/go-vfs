package vfst_test

import (
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert/v2"

	vfs "github.com/twpayne/go-vfs/v4"
	"github.com/twpayne/go-vfs/v4/vfst"
)

func TestWalk(t *testing.T) {
	fileSystem, cleanup, err := vfst.NewTestFS(map[string]interface{}{
		"/home/user/.bashrc":  "# .bashrc contents\n",
		"/home/user/skip/foo": "bar",
		"/home/user/symlink":  &vfst.Symlink{Target: "baz"},
	})
	assert.NoError(t, err)
	defer cleanup()
	pathTypeMap := make(map[string]fs.FileMode)
	assert.NoError(t, vfs.Walk(fileSystem, "/", func(path string, info fs.FileInfo, err error) error {
		assert.NoError(t, err)
		pathTypeMap[filepath.ToSlash(path)] = info.Mode() & fs.ModeType
		if filepath.Base(path) == "skip" {
			return vfs.SkipDir
		}
		return nil
	}))
	expectedPathTypeMap := map[string]fs.FileMode{
		"/":                  fs.ModeDir,
		"/home":              fs.ModeDir,
		"/home/user":         fs.ModeDir,
		"/home/user/.bashrc": 0,
		"/home/user/skip":    fs.ModeDir,
		"/home/user/symlink": fs.ModeSymlink,
	}
	assert.Equal(t, expectedPathTypeMap, pathTypeMap)
}
