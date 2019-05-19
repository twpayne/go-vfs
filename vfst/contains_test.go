package vfst

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vfs "github.com/twpayne/go-vfs"
)

func TestContains(t *testing.T) {
	type test struct {
		p         string
		prefix    string
		expectErr bool
		expected  bool
	}
	for _, tc := range []struct {
		name  string
		root  interface{}
		tests []test
	}{
		{
			name: "core",
			root: map[string]interface{}{
				"/home/user/file": "contents",
			},
			tests: []test{
				{
					p:        "/home/user",
					prefix:   "/home/user",
					expected: true,
				},
				{
					p:        "/home/user",
					prefix:   "/home",
					expected: true,
				},
				{
					p:        "/home/user",
					prefix:   "/",
					expected: true,
				},
				{
					p:        "/home/user/notexistpath",
					prefix:   "/home/user",
					expected: true,
				},
				{
					p:        "/home/user/notexistpath",
					prefix:   "/home",
					expected: true,
				},
				{
					p:        "/home/user/notexistpath",
					prefix:   "/",
					expected: true,
				},
				{
					p:        "/home/user/notexistdir/notexistpath",
					prefix:   "/home/user",
					expected: true,
				},
				{
					p:        "/home",
					prefix:   "/home/user",
					expected: false,
				},
				{
					p:        "/",
					prefix:   "/home/user",
					expected: false,
				},
				{
					p:        "/notexistpath",
					prefix:   "/home/user",
					expected: false,
				},
				{
					p:        "/notexistpath",
					prefix:   "/home",
					expected: false,
				},
				{
					p:        "/notexistpath",
					prefix:   "/",
					expected: true,
				},
			},
		},
		{
			name: "nonexistant_prefix",
			root: map[string]interface{}{
				"/home/user/file": "contents",
			},
			tests: []test{
				{
					p:         "/home/user",
					prefix:    "/notexistpath",
					expectErr: true,
				},
				{
					p:         "/home/user",
					prefix:    "/notexistdir/notexistpath",
					expectErr: true,
				},
			},
		},
		{
			name: "symlink_dir",
			root: []interface{}{
				map[string]interface{}{
					"/home/user/file": "contents",
				},
				map[string]interface{}{
					"/home/symlink": &Symlink{Target: "user"},
				},
			},
			tests: []test{
				{
					p:        "/home/symlink",
					prefix:   "/home/user",
					expected: true,
				},
				{
					p:        "/home/symlink",
					prefix:   "/home",
					expected: true,
				},
				{
					p:        "/home/symlink",
					prefix:   "/",
					expected: true,
				},
				{
					p:        "/home/symlink/notexistpath",
					prefix:   "/home/user",
					expected: true,
				},
				{
					p:        "/home/symlink/notexistpath",
					prefix:   "/home",
					expected: true,
				},
				{
					p:        "/home/symlink/notexistpath",
					prefix:   "/",
					expected: true,
				},
				{
					p:        "/home/symlink/notexistdir/notexistpath",
					prefix:   "/home/user",
					expected: true,
				},
				{
					p:        "/home/symlink/notexistdir/notexistpath",
					prefix:   "/home",
					expected: true,
				},
				{
					p:        "/home/symlink/notexistdir/notexistpath",
					prefix:   "/",
					expected: true,
				},
				{
					p:        "/home/symlink/notexistpath",
					prefix:   "/home/user",
					expected: true,
				},
			},
		},
		{
			name: "loop",
			root: map[string]interface{}{
				"/home/user": &Symlink{Target: "user"},
			},
			tests: []test{
				{
					p:         "/home/user",
					prefix:    "/home/user",
					expectErr: true,
				},
				{
					p:         "/home/user/notexistpath",
					prefix:    "/home/user",
					expectErr: true,
				},
				{
					p:         "/home/user/notexistdir/notexistpath",
					prefix:    "/home/user",
					expectErr: true,
				},
				{
					p:        "/home/user/notexistdir/notexistpath",
					prefix:   "/home",
					expected: true,
				},
			},
		},

		// Windows has a maximum path length of 260 chars ( - 12 if creating a directory)
		// so these tests are expected to error out on that platform.
		{
			name: "long_filename",
			root: map[string]interface{}{
				"/home/user": &Dir{Perm: 0755},
			},
			tests: []test{
				{
					p:         "/home/user/" + strings.Repeat("filename", 1024*1024), // 8MB filename
					prefix:    "/home/user",
					expectErr: runtime.GOOS == "windows",
					expected:  true,
				},
				{
					p:         "/home/user/" + strings.Repeat("filename", 1024*1024), // 8MB filename
					prefix:    "/home",
					expectErr: runtime.GOOS == "windows",
					expected:  true,
				},
				{
					p:         "/home/user/" + strings.Repeat("filename", 1024*1024), // 8MB filename
					prefix:    "/",
					expectErr: runtime.GOOS == "windows",
					expected:  true,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fs, cleanup, err := NewTestFS(tc.root)
			require.NoError(t, err)
			defer cleanup()
			for _, test := range tc.tests {
				actual, err := vfs.Contains(fs, test.p, test.prefix)
				if test.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, test.expected, actual)
				}
			}
		})
	}
}
