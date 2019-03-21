package vfs

import (
	"os"
	"path/filepath"
)

// A Stater implements Stat.
type Stater interface {
	Stat(string) (os.FileInfo, error)
}

// HasPrefix returns true if p or any parent of p is the same file as prefix.
// prefix must exist, but p may not. It is an expensive but accurate alternative
// to the deprecated filepath.HasPrefix.
func HasPrefix(fs Stater, p, prefix string) (bool, error) {
	prefixFI, err := fs.Stat(prefix)
	if err != nil {
		return false, err
	}
	for {
		fi, err := fs.Stat(p)
		switch {
		case err == nil:
			if os.SameFile(fi, prefixFI) {
				return true, nil
			}
		case os.IsNotExist(err):
			// Do nothing and skip ahead to trying p's parent directory.
		default:
			return false, err
		}
		parentDir := filepath.Dir(p)
		if parentDir == p {
			// Return when we stop making progress.
			return false, nil
		}
		p = parentDir
	}
}
