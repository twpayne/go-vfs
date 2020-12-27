package vfst

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vfs "github.com/twpayne/go-vfs"
)

func TestBuilderBuild(t *testing.T) {
	for _, tc := range []struct {
		name  string
		umask os.FileMode
		root  interface{}
		tests interface{}
	}{
		{
			name:  "empty",
			umask: 0o22,
			tests: []Test{},
		},
		{
			name:  "dir",
			umask: 0o22,
			root: map[string]interface{}{
				"foo": &Dir{
					Perm: 0o755,
					Entries: map[string]interface{}{
						"bar": "baz",
					},
				},
			},
			tests: []Test{
				TestPath("/foo",
					TestIsDir,
					TestModePerm(0o755),
				),
				TestPath("/foo/bar",
					TestModeIsRegular,
					TestModePerm(0o644),
					TestContentsString("baz"),
				),
			},
		},
		{
			name:  "map_string_string",
			umask: 0o22,
			root: map[string]string{
				"foo": "bar",
			},
			tests: []Test{
				TestPath("/foo",
					TestModeIsRegular,
					TestModePerm(0o644),
					TestContentsString("bar"),
				),
			},
		},
		{
			name:  "map_string_empty_interface",
			umask: 0o22,
			root: map[string]interface{}{
				"foo": "bar",
				"baz": &File{Perm: 0o755, Contents: []byte("qux")},
				"dir": &Dir{Perm: 0o700},
			},
			tests: []Test{
				TestPath("/foo",
					TestModeIsRegular,
					TestModePerm(0o644),
					TestSize(3),
					TestContentsString("bar"),
				),
				TestPath("/baz",
					TestModeIsRegular,
					TestModePerm(0o755),
					TestSize(3),
					TestContentsString("qux"),
				),
				TestPath("/dir",
					TestIsDir,
					TestModePerm(0o700),
				),
			},
		},
		{
			name:  "long_paths",
			umask: 0o22,
			root: map[string]string{
				"/foo/bar": "baz",
			},
			tests: []Test{
				TestPath("/foo",
					TestIsDir,
					TestModePerm(0o755),
				),
				TestPath("/foo/bar",
					TestModeIsRegular,
					TestModePerm(0o644),
					TestSize(3),
					TestContentsString("baz"),
				),
			},
		},
		{
			name:  "symlink",
			umask: 0o22,
			root: map[string]interface{}{
				"foo": &Symlink{Target: "bar"},
			},
			tests: []Test{
				TestPath("/foo",
					TestModeType(os.ModeSymlink),
					TestSymlinkTarget("bar"),
				),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fs, cleanup, err := NewTestFS(tc.root, BuilderUmask(tc.umask), BuilderVerbose(true))
			require.NoError(t, err)
			defer cleanup()
			RunTests(t, fs, "", tc.tests)
		})
	}
}

// TestCoverage exercises as much functionality as possible to increase test
// coverage.
func TestCoverage(t *testing.T) {
	fs, cleanup, err := NewTestFS(map[string]interface{}{
		"/home/user/.bashrc": "# contents of user's .bashrc\n",
		"/home/user/empty":   []byte{},
		"/home/user/symlink": &Symlink{Target: "empty"},
		"/home/user/bin/hello.sh": &File{
			Perm:     0o755,
			Contents: []byte("echo hello\n"),
		},
		"/home/user/foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"baz": "qux",
			},
		},
		"/root": &Dir{
			Perm: 0o700,
			Entries: map[string]interface{}{
				".bashrc": "# contents of root's .bashrc\n",
			},
		},
	})
	require.NoError(t, err)
	defer cleanup()
	RunTests(t, fs, "", []interface{}{
		TestPath("/home",
			TestIsDir,
			TestModePerm(0o755),
		),
		TestPath("/notexist",
			TestDoesNotExist),
		map[string]Test{
			"home_user_bashrc": TestPath("/home/user/.bashrc",
				TestModeIsRegular,
				TestModePerm(0o644),
				TestContentsString("# contents of user's .bashrc\n"),
				TestMinSize(1),
				TestSysNlink(1),
			),
		},
		map[string]interface{}{
			"home_user_empty": TestPath("/home/user/empty",
				TestModeIsRegular,
				TestModePerm(0o644),
				TestSize(0),
			),
			"home_user_symlink": TestPath("/home/user/symlink",
				TestModeType(os.ModeSymlink),
				TestSymlinkTarget("empty"),
			),
			"foo_bar_baz": []Test{
				TestPath("/home/user/foo/bar/baz",
					TestModeIsRegular,
					TestModePerm(0o644),
					TestContentsString("qux"),
				),
			},
			"root": []interface{}{
				TestPath("/root",
					TestIsDir,
					TestModePerm(0o700),
				),
				TestPath("/root/.bashrc",
					TestModeIsRegular,
					TestModePerm(0o644),
					TestContentsString("# contents of root's .bashrc\n"),
				),
			},
		},
	})
}

func TestErrors(t *testing.T) {
	errSkip := errors.New("skip")
	for name, f := range map[string]func(*Builder, vfs.FS) error{
		"write_file_with_different_content": func(b *Builder, fs vfs.FS) error {
			return b.WriteFile(fs, "/home/user/.bashrc", nil, 0o644)
		},
		"write_file_with_different_perms": func(b *Builder, fs vfs.FS) error {
			if permEqual(0o644, 0o755) {
				return errSkip
			}
			return b.WriteFile(fs, "/home/user/.bashrc", []byte("# bashrc\n"), 0o755)
		},
		"write_file_to_existing_dir": func(b *Builder, fs vfs.FS) error {
			return b.WriteFile(fs, "/home/user", nil, 0o644)
		},
		"write_file_to_existing_symlink": func(b *Builder, fs vfs.FS) error {
			return b.WriteFile(fs, "/home/user/symlink", nil, 0o644)
		},
		"write_file_via_existing_dir": func(b *Builder, fs vfs.FS) error {
			return b.WriteFile(fs, "/home/user/empty/foo", nil, 0o644)
		},
		"write_file_via_existing_symlink": func(b *Builder, fs vfs.FS) error {
			return b.WriteFile(fs, "/home/user/symlink/foo", nil, 0o644)
		},
		"mkdir_existing_dir_with_different_perms": func(b *Builder, fs vfs.FS) error {
			if permEqual(0o755, 0o666) {
				return errSkip
			}
			return b.Mkdir(fs, "/home/user", 0o666)
		},
		"mkdir_to_existing_file": func(b *Builder, fs vfs.FS) error {
			return b.Mkdir(fs, "/home/user/empty", 0o755)
		},
		"mkdir_to_existing_symlink": func(b *Builder, fs vfs.FS) error {
			return b.Mkdir(fs, "/home/user/symlink", 0o755)
		},
		"mkdir_all_to_existing_file": func(b *Builder, fs vfs.FS) error {
			return b.Mkdir(fs, "/home/user/empty", 0o755)
		},
		"mkdir_all_via_existing_file": func(b *Builder, fs vfs.FS) error {
			return b.MkdirAll(fs, "/home/user/empty/foo", 0o755)
		},
		"mkdir_all_via_existing_symlink": func(b *Builder, fs vfs.FS) error {
			return b.MkdirAll(fs, "/home/user/symlink/foo", 0o755)
		},
	} {
		t.Run(name, func(t *testing.T) {
			fs, cleanup, err := newTestFS()
			require.NoError(t, err)
			defer cleanup()
			b := NewBuilder(BuilderVerbose(true))
			root := []interface{}{
				map[string]interface{}{
					"/home/user/.bashrc": "# bashrc\n",
					"/home/user/empty":   []byte{},
					"/home/user/foo":     &Dir{Perm: 0o755},
				},
				map[string]interface{}{
					"/home/user/symlink": &Symlink{Target: "empty"},
				},
			}
			require.NoError(t, b.Build(fs, root))
			assert.Error(t, f(b, fs))
		})
	}
}

func TestGlob(t *testing.T) {
	fs, cleanup, err := NewTestFS(map[string]interface{}{
		"/home/user/.bash_profile": "# contents of .bash_profile\n",
		"/home/user/.bashrc":       "# contents of .bashrc\n",
		"/home/user/.zshrc":        "# contents of .zshrc\n",
	})
	require.NoError(t, err)
	defer cleanup()
	for _, tc := range []struct {
		name            string
		pattern         string
		expectedMatches []string
	}{
		{
			name:    "all",
			pattern: "/home/user/*",
			expectedMatches: []string{
				"/home/user/.bash_profile",
				"/home/user/.bashrc",
				"/home/user/.zshrc",
			},
		},
		{
			name:    "star_rc",
			pattern: "/home/user/*rc",
			expectedMatches: []string{
				"/home/user/.bashrc",
				"/home/user/.zshrc",
			},
		},
		{
			name:    "all_subdir",
			pattern: "/home/*/*",
			expectedMatches: []string{
				"/home/user/.bash_profile",
				"/home/user/.bashrc",
				"/home/user/.zshrc",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			matches, err := fs.Glob(tc.pattern)
			require.NoError(t, err)
			require.Len(t, matches, len(tc.expectedMatches))
			for i, match := range matches {
				assert.True(t, filepath.IsAbs(match))
				expected := filepath.FromSlash(tc.expectedMatches[i])
				actual := strings.TrimPrefix(match, filepath.VolumeName(matches[i]))
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestIdempotency(t *testing.T) {
	for name, f := range map[string]func(*Builder, vfs.FS) error{
		"write_new_file": func(b *Builder, fs vfs.FS) error {
			return b.WriteFile(fs, "/home/user/empty", nil, 0o644)
		},
		"write_file_with_same_content_and_perms": func(b *Builder, fs vfs.FS) error {
			return b.WriteFile(fs, "/home/user/.bashrc", []byte("# bashrc\n"), 0o644)
		},
		"mkdir_existing_dir_with_same_perms": func(b *Builder, fs vfs.FS) error {
			return b.Mkdir(fs, "/home/user", 0o755)
		},
		"mkdir_new_dir": func(b *Builder, fs vfs.FS) error {
			return b.Mkdir(fs, "/home/user/foo", 0o755)
		},
		"mkdir_all_existing_dir": func(b *Builder, fs vfs.FS) error {
			return b.MkdirAll(fs, "/home/user", 0o755)
		},
		"mkdir_all_new_dir": func(b *Builder, fs vfs.FS) error {
			return b.MkdirAll(fs, "/usr/bin", 0o755)
		},
		"symlink_new_symlink": func(b *Builder, fs vfs.FS) error {
			return b.Symlink(fs, ".bashrc", "/home/user/symlink2")
		},
		"symlink_existing_symlink": func(b *Builder, fs vfs.FS) error {
			return b.Symlink(fs, ".bashrc", "/home/user/symlink")
		},
	} {
		t.Run(name, func(t *testing.T) {
			fs, cleanup, err := newTestFS()
			require.NoError(t, err)
			defer cleanup()
			b := NewBuilder(BuilderVerbose(true))
			root := map[string]interface{}{
				"/home/user/.bashrc": "# bashrc\n",
				"/home/user/symlink": &Symlink{Target: ".bashrc"},
			}
			require.NoError(t, b.Build(fs, root))
			assert.NoError(t, f(b, fs))
		})
	}
}
