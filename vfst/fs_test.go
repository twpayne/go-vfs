package vfst

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	vfs "github.com/twpayne/go-vfs"
)

func TestWalk(t *testing.T) {
	fs, cleanup, err := NewTestFS(map[string]interface{}{
		"/home/user/.bashrc":  "# .bashrc contents\n",
		"/home/user/skip/foo": "bar",
		"/home/user/symlink":  &Symlink{Target: "baz"},
	})
	defer cleanup()
	if err != nil {
		t.Fatal(err)
	}
	pathTypeMap := make(map[string]os.FileMode)
	if err := vfs.Walk(fs, "/", func(path string, info os.FileInfo, err error) error {
		pathTypeMap[filepath.ToSlash(path)] = info.Mode() & os.ModeType
		if err != nil {
			t.Errorf("walkFn(%q, %v, %v) called, want err == <nil>", path, info, err)
		}
		if filepath.Base(path) == "skip" {
			return filepath.SkipDir
		}
		return nil

	}); err != nil {
		t.Errorf("vfs.Walk(...) == %v, want <nil>", err)
	}
	wantPathTypeMap := map[string]os.FileMode{
		"/":                  os.ModeDir,
		"/home":              os.ModeDir,
		"/home/user":         os.ModeDir,
		"/home/user/.bashrc": 0,
		"/home/user/skip":    os.ModeDir,
		"/home/user/symlink": os.ModeSymlink,
	}
	if !reflect.DeepEqual(pathTypeMap, wantPathTypeMap) {
		t.Errorf("pathTypeMap == %+v, want %+v", pathTypeMap, wantPathTypeMap)
	}
}
