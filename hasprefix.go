package vfs

import (
	"os"
	"path/filepath"
	"syscall"
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
			goto TryParent
		case os.IsNotExist(err):
			goto TryParent
		case os.IsPermission(err):
			goto TryParent
		default:
			// Remove any os.PathError or os.SyscallError wrapping, if present.
			for {
				if pathError, ok := err.(*os.PathError); ok {
					err = pathError.Err
				} else if syscallError, ok := err.(*os.SyscallError); ok {
					err = syscallError.Err
				} else {
					break
				}
			}
			// Ignore some syscall.Errnos.
			if errno, ok := err.(syscall.Errno); ok {
				switch errno {
				case syscall.ELOOP:
					goto TryParent
				case syscall.EMLINK:
					goto TryParent
				case syscall.ENAMETOOLONG:
					goto TryParent
				case syscall.ENOENT:
					goto TryParent
				case syscall.EOVERFLOW:
					goto TryParent
				}
			}
			return false, err
		}
	TryParent:
		parentDir := filepath.Dir(p)
		if parentDir == p {
			// Return when we stop making progress.
			return false, nil
		}
		p = parentDir
	}
}
