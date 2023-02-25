package vfst_test

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vfs "github.com/twpayne/go-vfs/v4"
	"github.com/twpayne/go-vfs/v4/vfst"
)

func TestBuilderBuild(t *testing.T) {
	for _, tc := range []struct {
		name  string
		umask fs.FileMode
		root  interface{}
		tests interface{}
	}{
		{
			name:  "empty",
			umask: 0o22,
			tests: []vfst.Test{},
		},
		{
			name:  "dir",
			umask: 0o22,
			root: map[string]interface{}{
				"foo": &vfst.Dir{
					Perm: 0o755,
					Entries: map[string]interface{}{
						"bar": "baz",
					},
				},
			},
			tests: []vfst.Test{
				vfst.TestPath("/foo",
					vfst.TestIsDir,
					vfst.TestModePerm(0o755),
				),
				vfst.TestPath("/foo/bar",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o644),
					vfst.TestContentsString("baz"),
				),
			},
		},
		{
			name:  "map_string_string",
			umask: 0o22,
			root: map[string]string{
				"foo": "bar",
			},
			tests: []vfst.Test{
				vfst.TestPath("/foo",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o644),
					vfst.TestContentsString("bar"),
				),
			},
		},
		{
			name:  "map_string_empty_interface",
			umask: 0o22,
			root: map[string]interface{}{
				"foo": "bar",
				"baz": &vfst.File{Perm: 0o755, Contents: []byte("qux")},
				"dir": &vfst.Dir{Perm: 0o700},
			},
			tests: []vfst.Test{
				vfst.TestPath("/foo",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o644),
					vfst.TestSize(3),
					vfst.TestContentsString("bar"),
				),
				vfst.TestPath("/baz",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o755),
					vfst.TestSize(3),
					vfst.TestContentsString("qux"),
				),
				vfst.TestPath("/dir",
					vfst.TestIsDir,
					vfst.TestModePerm(0o700),
				),
			},
		},
		{
			name:  "long_paths",
			umask: 0o22,
			root: map[string]string{
				"/foo/bar": "baz",
			},
			tests: []vfst.Test{
				vfst.TestPath("/foo",
					vfst.TestIsDir,
					vfst.TestModePerm(0o755),
				),
				vfst.TestPath("/foo/bar",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o644),
					vfst.TestSize(3),
					vfst.TestContentsString("baz"),
				),
			},
		},
		{
			name:  "symlink",
			umask: 0o22,
			root: map[string]interface{}{
				"foo": &vfst.Symlink{Target: "bar"},
			},
			tests: []vfst.Test{
				vfst.TestPath("/foo",
					vfst.TestModeType(fs.ModeSymlink),
					vfst.TestSymlinkTarget("bar"),
				),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fileSystem, cleanup, err := vfst.NewTestFS(tc.root, vfst.BuilderUmask(tc.umask), vfst.BuilderVerbose(true))
			require.NoError(t, err)
			defer cleanup()
			vfst.RunTests(t, fileSystem, "", tc.tests)
		})
	}
}

// TestCoverage exercises as much functionality as possible to increase test
// coverage.
func TestCoverage(t *testing.T) {
	fileSystem, cleanup, err := vfst.NewTestFS(map[string]interface{}{
		"/home/user/.bashrc": "# contents of user's .bashrc\n",
		"/home/user/empty":   []byte{},
		"/home/user/symlink": &vfst.Symlink{Target: "empty"},
		"/home/user/bin/hello.sh": &vfst.File{
			Perm:     0o755,
			Contents: []byte("echo hello\n"),
		},
		"/home/user/foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"baz": "qux",
			},
		},
		"/root": &vfst.Dir{
			Perm: 0o700,
			Entries: map[string]interface{}{
				".bashrc": "# contents of root's .bashrc\n",
			},
		},
	})
	require.NoError(t, err)
	defer cleanup()
	vfst.RunTests(t, fileSystem, "", []interface{}{
		vfst.TestPath("/home",
			vfst.TestIsDir,
			vfst.TestModePerm(0o755),
		),
		vfst.TestPath("/notexist",
			vfst.TestDoesNotExist),
		map[string]vfst.Test{
			"home_user_bashrc": vfst.TestPath("/home/user/.bashrc",
				vfst.TestModeIsRegular,
				vfst.TestModePerm(0o644),
				vfst.TestContentsString("# contents of user's .bashrc\n"),
				vfst.TestMinSize(1),
				vfst.TestSysNlink(1),
			),
		},
		map[string]interface{}{
			"home_user_empty": vfst.TestPath("/home/user/empty",
				vfst.TestModeIsRegular,
				vfst.TestModePerm(0o644),
				vfst.TestSize(0),
			),
			"home_user_symlink": vfst.TestPath("/home/user/symlink",
				vfst.TestModeType(fs.ModeSymlink),
				vfst.TestSymlinkTarget("empty"),
			),
			"foo_bar_baz": []vfst.Test{
				vfst.TestPath("/home/user/foo/bar/baz",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o644),
					vfst.TestContentsString("qux"),
				),
			},
			"root": []interface{}{
				vfst.TestPath("/root",
					vfst.TestIsDir,
					vfst.TestModePerm(0o700),
				),
				vfst.TestPath("/root/.bashrc",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o644),
					vfst.TestContentsString("# contents of root's .bashrc\n"),
				),
			},
		},
	})
}

func TestErrors(t *testing.T) {
	errSkip := errors.New("skip")
	for name, f := range map[string]func(*vfst.Builder, vfs.FS) error{
		"write_file_with_different_content": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.WriteFile(fileSystem, "/home/user/.bashrc", nil, 0o644)
		},
		"write_file_with_different_perms": func(b *vfst.Builder, fileSystem vfs.FS) error {
			if vfst.PermEqual(0o644, 0o755) {
				return errSkip
			}
			return b.WriteFile(fileSystem, "/home/user/.bashrc", []byte("# bashrc\n"), 0o755)
		},
		"write_file_to_existing_dir": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.WriteFile(fileSystem, "/home/user", nil, 0o644)
		},
		"write_file_to_existing_symlink": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.WriteFile(fileSystem, "/home/user/symlink", nil, 0o644)
		},
		"write_file_via_existing_dir": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.WriteFile(fileSystem, "/home/user/empty/foo", nil, 0o644)
		},
		"write_file_via_existing_symlink": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.WriteFile(fileSystem, "/home/user/symlink/foo", nil, 0o644)
		},
		"mkdir_existing_dir_with_different_perms": func(b *vfst.Builder, fileSystem vfs.FS) error {
			if vfst.PermEqual(0o755, 0o666) {
				return errSkip
			}
			return b.Mkdir(fileSystem, "/home/user", 0o666)
		},
		"mkdir_to_existing_file": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.Mkdir(fileSystem, "/home/user/empty", 0o755)
		},
		"mkdir_to_existing_symlink": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.Mkdir(fileSystem, "/home/user/symlink", 0o755)
		},
		"mkdir_all_to_existing_file": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.Mkdir(fileSystem, "/home/user/empty", 0o755)
		},
		"mkdir_all_via_existing_file": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.MkdirAll(fileSystem, "/home/user/empty/foo", 0o755)
		},
		"mkdir_all_via_existing_symlink": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.MkdirAll(fileSystem, "/home/user/symlink/foo", 0o755)
		},
	} {
		t.Run(name, func(t *testing.T) {
			fileSystem, cleanup, err := vfst.NewEmptyTestFS()
			require.NoError(t, err)
			defer cleanup()
			b := vfst.NewBuilder(vfst.BuilderVerbose(true))
			root := []interface{}{
				map[string]interface{}{
					"/home/user/.bashrc": "# bashrc\n",
					"/home/user/empty":   []byte{},
					"/home/user/foo":     &vfst.Dir{Perm: 0o755},
				},
				map[string]interface{}{
					"/home/user/symlink": &vfst.Symlink{Target: "empty"},
				},
			}
			require.NoError(t, b.Build(fileSystem, root))
			assert.Error(t, f(b, fileSystem))
		})
	}
}

func TestGlob(t *testing.T) {
	fileSystem, cleanup, err := vfst.NewTestFS(map[string]interface{}{
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
			matches, err := fileSystem.Glob(tc.pattern)
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
	for name, f := range map[string]func(*vfst.Builder, vfs.FS) error{
		"write_new_file": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.WriteFile(fileSystem, "/home/user/empty", nil, 0o644)
		},
		"write_file_with_same_content_and_perms": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.WriteFile(fileSystem, "/home/user/.bashrc", []byte("# bashrc\n"), 0o644)
		},
		"mkdir_existing_dir_with_same_perms": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.Mkdir(fileSystem, "/home/user", 0o755)
		},
		"mkdir_new_dir": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.Mkdir(fileSystem, "/home/user/foo", 0o755)
		},
		"mkdir_all_existing_dir": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.MkdirAll(fileSystem, "/home/user", 0o755)
		},
		"mkdir_all_new_dir": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.MkdirAll(fileSystem, "/usr/bin", 0o755)
		},
		"symlink_new_symlink": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.Symlink(fileSystem, ".bashrc", "/home/user/symlink2")
		},
		"symlink_existing_symlink": func(b *vfst.Builder, fileSystem vfs.FS) error {
			return b.Symlink(fileSystem, ".bashrc", "/home/user/symlink")
		},
	} {
		t.Run(name, func(t *testing.T) {
			fileSystem, cleanup, err := vfst.NewEmptyTestFS()
			require.NoError(t, err)
			defer cleanup()
			b := vfst.NewBuilder(vfst.BuilderVerbose(true))
			root := map[string]interface{}{
				"/home/user/.bashrc": "# bashrc\n",
				"/home/user/symlink": &vfst.Symlink{Target: ".bashrc"},
			}
			require.NoError(t, b.Build(fileSystem, root))
			assert.NoError(t, f(b, fileSystem))
		})
	}
}
