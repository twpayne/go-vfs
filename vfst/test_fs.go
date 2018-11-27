package vfst

import (
	"io/ioutil"
	"os"

	"github.com/twpayne/go-vfs"
)

type TestFS struct {
	vfs.PathFS
	tempDir string
	keep    bool
}

func newTestFS() (*TestFS, func(), error) {
	tempDir, err := ioutil.TempDir("", "fs-fstest")
	if err != nil {
		return nil, func() {}, err
	}
	t := &TestFS{
		PathFS:  *vfs.NewPathFS(vfs.OSFS, tempDir),
		tempDir: tempDir,
		keep:    false,
	}
	return t, t.cleanup, nil
}

// NewTestFS returns a new *TestFS based in a temporary directory and a cleanup
// function, populated with root.
func NewTestFS(root interface{}, builderOptions ...BuilderOption) (*TestFS, func(), error) {
	fs, cleanup, err := newTestFS()
	if err != nil {
		cleanup()
		return nil, func() {}, err
	}
	if err := NewBuilder(builderOptions...).Build(fs, root); err != nil {
		cleanup()
		return nil, func() {}, err
	}
	return fs, cleanup, nil
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
		os.RemoveAll(t.tempDir)
	}
}
