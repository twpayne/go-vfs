package vfst

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vfs "github.com/twpayne/go-vfs"
)

func TestHasPrefix(t *testing.T) {
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
			root: map[string]interface{}{
				"/home/user/file": "contents",
				"/home/symlink":   &Symlink{Target: "user"},
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
	} {
		t.Run(tc.name, func(t *testing.T) {
			fs, cleanup, err := NewTestFS(tc.root)
			require.NoError(t, err)
			defer cleanup()
			for _, test := range tc.tests {
				actual, err := vfs.HasPrefix(fs, test.p, test.prefix)
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
