package vfst_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"

	vfs "github.com/twpayne/go-vfs/v4"
	"github.com/twpayne/go-vfs/v4/vfst"
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
					"/home/symlink": &vfst.Symlink{Target: "user"},
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
				"/home/user": &vfst.Symlink{Target: "user"},
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
	} {
		t.Run(tc.name, func(t *testing.T) {
			fileSystem, cleanup, err := vfst.NewTestFS(tc.root)
			assert.NoError(t, err)
			defer cleanup()
			for _, test := range tc.tests {
				actual, err := vfs.Contains(fileSystem, test.p, test.prefix)
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
