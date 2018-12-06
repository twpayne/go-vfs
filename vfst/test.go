// +build !windows

package vfst

import (
	"os"
	"syscall"
	"testing"

	vfs "github.com/twpayne/go-vfs"
)

// permEqual returns if perm1 and perm2 represent the same permissions. On
// Windows, it always returns true.
func permEqual(perm1, perm2 os.FileMode) bool {
	return perm1&os.ModePerm == perm2&os.ModePerm
}

// TestSysNlink returns a PathTest that verifies that the the path's
// Sys().(*syscall.Stat_t).Nlink is equal to wantNlink. If path's Sys() cannot
// be converted to a *syscall.Stat_t, it does nothing.
func TestSysNlink(wantNlink int) PathTest {
	return func(t *testing.T, fs vfs.FS, path string) {
		info, ok, err := fs.LstatIfPossible(path)
		if err != nil {
			t.Errorf("fs.LstatIfPossible(%q) == %+v, %v, %v, want !<nil>, _, <nil>", path, info, ok, err)
			return
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok && int(stat.Nlink) != wantNlink {
			t.Errorf("fs.LstatIfPossible(%q).Sys().(*syscall.Stat_t).Nlink == %d, want %d", path, stat.Nlink, wantNlink)
		}
	}
}
