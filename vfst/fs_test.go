package vfst_test

import (
	"errors"
	"io/fs"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alecthomas/assert/v2"

	vfs "github.com/twpayne/go-vfs/v5"
	"github.com/twpayne/go-vfs/v5/vfst"
)

func TestWalk(t *testing.T) {
	fileSystem, cleanup, err := vfst.NewTestFS(map[string]any{
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

func TestWalkErrors(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test uses UNIX file permissions")
	}
	for _, tc := range []struct {
		name          string
		root          any
		postFunc      func(vfs.FS) error
		expectedPaths []string
	}{
		{
			name: "empty",
			expectedPaths: []string{
				"/",
			},
		},
		{
			name: "simple",
			root: map[string]any{
				"/dir/subdir/subsubdir/file": "",
			},
			expectedPaths: []string{
				"/",
				"/dir",
				"/dir/subdir",
				"/dir/subdir/subsubdir",
				"/dir/subdir/subsubdir/file",
			},
		},
		{
			name: "private_subdir",
			root: map[string]any{
				"/dir/subdir/subsubdir/file": "",
			},
			postFunc: func(fileSystem vfs.FS) error {
				return fileSystem.Chmod("/dir/subdir", 0)
			},
			expectedPaths: []string{
				"/",
				"/dir",
				"/dir/subdir",
			},
		},
		{
			name: "private_subdir_keep_going",
			root: map[string]any{
				"/dir/subdir/subsubdir/file": "",
				"/dir/subdir2/file":          "",
			},
			postFunc: func(fileSystem vfs.FS) error {
				return fileSystem.Chmod("/dir/subdir", 0)
			},
			expectedPaths: []string{
				"/",
				"/dir",
				"/dir/subdir",
				"/dir/subdir2",
				"/dir/subdir2/file",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fileSystem, cleanup, err := vfst.NewTestFS(tc.root)
			assert.NoError(t, err)
			if tc.postFunc != nil {
				assert.NoError(t, tc.postFunc(fileSystem))
			}
			defer cleanup()
			var actualPaths []string
			assert.NoError(t, vfs.Walk(fileSystem, "/", func(path string, info fs.FileInfo, err error) error {
				switch {
				case errors.Is(err, fs.ErrPermission):
					if info.IsDir() {
						return vfs.SkipDir
					}
					return nil
				case err != nil:
					return err
				default:
					actualPaths = append(actualPaths, path)
					return nil
				}
			}))
			assert.Equal(t, tc.expectedPaths, actualPaths)
		})
	}
}
