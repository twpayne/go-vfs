package vfs

//nolint:godox
// FIXME implement path/filepath.WalkDir

import (
	"os"
	"path/filepath"
	"sort"
)

// SkipDir is filepath.SkipDir.
//nolint:gochecknoglobals
var SkipDir = filepath.SkipDir

// A LstatReadDirer implements all the functionality needed by Walk.
type LstatReadDirer interface {
	Lstat(name string) (os.FileInfo, error)
	ReadDir(name string) ([]os.DirEntry, error)
}

type dirEntriesByName []os.DirEntry

func (is dirEntriesByName) Len() int           { return len(is) }
func (is dirEntriesByName) Less(i, j int) bool { return is[i].Name() < is[j].Name() }
func (is dirEntriesByName) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }

// walk recursively walks fs from path.
func walk(fs LstatReadDirer, path string, walkFn filepath.WalkFunc, info os.FileInfo, err error) error {
	if err != nil {
		return walkFn(path, info, err)
	}
	err = walkFn(path, info, nil)
	if !info.IsDir() {
		return err
	}
	if err == filepath.SkipDir {
		return nil
	}
	dirEntries, err := fs.ReadDir(path)
	if err != nil {
		return err
	}
	sort.Sort(dirEntriesByName(dirEntries))
	for _, dirEntry := range dirEntries {
		name := dirEntry.Name()
		if name == "." || name == ".." {
			continue
		}
		info, err := dirEntry.Info()
		if err != nil {
			return err
		}
		if err := walk(fs, filepath.Join(path, dirEntry.Name()), walkFn, info, nil); err != nil {
			return err
		}
	}
	return nil
}

// Walk is the equivalent of filepath.Walk but operates on fs. Entries are
// returned in lexicographical order.
func Walk(fs LstatReadDirer, path string, walkFn filepath.WalkFunc) error {
	info, err := fs.Lstat(path)
	return walk(fs, path, walkFn, info, err)
}

// WalkSlash is the equivalent of Walk but all paths are converted to use
// forward slashes with filepath.ToSlash.
func WalkSlash(fs LstatReadDirer, path string, walkFn filepath.WalkFunc) error {
	return Walk(fs, path, func(path string, info os.FileInfo, err error) error {
		return walkFn(filepath.ToSlash(path), info, err)
	})
}
