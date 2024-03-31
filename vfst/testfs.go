package vfst

import (
	"errors"
	"io/fs"
	"os"

	vfs "github.com/twpayne/go-vfs/v5"
)

// A TestFS is a virtual filesystem based in a temporary directory.
type TestFS struct {
	vfs.PathFS
	tempDir string
	keep    bool
}

// NewEmptyTestFS returns a new empty TestFS and a cleanup function.
func NewEmptyTestFS() (*TestFS, func(), error) {
	tempDir, err := os.MkdirTemp("", "go-vfs-vfst")
	if err != nil {
		return nil, nil, err
	}
	t := &TestFS{
		PathFS:  *vfs.NewPathFS(vfs.OSFS, tempDir),
		tempDir: tempDir,
		keep:    false,
	}
	return t, t.cleanup, nil
}

// NewTestFS returns a new *TestFS populated with root and a cleanup function.
func NewTestFS(root any, builderOptions ...BuilderOption) (*TestFS, func(), error) {
	fileSystem, cleanup, err := NewEmptyTestFS()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	if err := NewBuilder(builderOptions...).Build(fileSystem, root); err != nil {
		cleanup()
		return nil, nil, err
	}
	return fileSystem, cleanup, nil
}

// Keep prevents t's cleanup function from removing the temporary directory. It
// has no effect if cleanup has already been called.
func (t *TestFS) Keep() {
	t.keep = true
}

// TempDir returns t's temporary directory.
func (t *TestFS) TempDir() string {
	return t.tempDir
}

func (t *TestFS) cleanup() {
	if !t.keep {
		for {
			// Remove t.tempDir but try to recover from permission denied errors
			// by chmod'ing the path that causes the error.
			err := os.RemoveAll(t.tempDir)
			if err == nil {
				break
			}
			if !errors.Is(err, fs.ErrPermission) {
				break
			}
			var pathErr *os.PathError
			if !errors.As(err, &pathErr) {
				break
			}
			if err := os.Chmod(pathErr.Path, 0o777); err != nil {
				break
			}
		}
	}
}
