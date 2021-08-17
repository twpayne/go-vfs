package vfst_test

import (
	"testing"

	"github.com/twpayne/go-vfs/v4/vfst"
)

func ExampleNewTestFS_complex() {
	Test := func(t *testing.T) {
		t.Helper()

		// Describe the structure of the filesystem using a map from filenames to
		// file or directory contents.
		root := map[string]interface{}{
			// A string or []byte is sets a file's contents.
			"/home/user/.bashrc": "# contents of user's .bashrc\n",
			"/home/user/empty":   []byte{},
			// To set non-default permissions on a file, create an &vfst.File.
			"/home/user/bin/hello.sh": &vfst.File{
				Perm:     0o755,
				Contents: []byte("echo hello\n"),
			},
			// Directories can be nested.
			"/home/user/foo": map[string]interface{}{
				"bar": map[string]interface{}{
					"baz": "qux",
				},
			},
			// To set non-default permissions on a directory, create an
			// &vfst.Dir.
			"/root": &vfst.Dir{
				Perm: 0o700,
				Entries: map[string]interface{}{
					".bashrc": "# contents of root's .bashrc\n",
				},
			},
		}

		// Create and populate an *vfst.TestFS
		fileSystem, cleanup, err := vfst.NewTestFS(root)
		if err != nil {
			t.Fatal(err)
		}
		defer cleanup()

		// Create tests by creating data structures containing Tests.
		tests := []interface{}{
			// Test multiple properties of a single path with TestPath.
			vfst.TestPath("/home",
				vfst.TestIsDir,
				vfst.TestModePerm(0o755),
			),
			vfst.TestPath("/home/user",
				vfst.TestIsDir,
				vfst.TestModePerm(0o755),
			),
			vfst.TestPath("/home/user/.bashrc",
				vfst.TestModeIsRegular,
				vfst.TestModePerm(0o644),
				vfst.TestContentsString("# contents of user's .bashrc\n"),
			),
			// Maps with string keys create sub tests with testing.T.Run. The key
			// is used as the test name.
			map[string]interface{}{
				"home_user_empty": vfst.TestPath("/home/user/empty",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o644),
					vfst.TestSize(0),
				),
				"foo_bar_baz": vfst.TestPath("/home/user/foo/bar/baz",
					vfst.TestModeIsRegular,
					vfst.TestModePerm(0o644),
					vfst.TestContentsString("qux"),
				),
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
		}

		// RunTests traverses the data structure and running all Tests.
		vfst.RunTests(t, fileSystem, "", tests)

		// Optionally, calling fileSystem.Keep() prevents the cleanup function
		// from removing the temporary directory, so you can inspect it later.
		// The directory itself is returned by fileSystem.TempDir().
		// fileSystem.Keep()
		t.Logf("fs.TempDir() == %s", fileSystem.TempDir())
	}

	Test(&testing.T{})
}

func ExampleNewTestFS() {
	Test := func(t *testing.T) {
		t.Helper()

		fileSystem, cleanup, err := vfst.NewTestFS(map[string]interface{}{
			"/home/user/.bashrc": "# contents of user's .bashrc\n",
		})
		if err != nil {
			t.Fatal(err)
		}
		defer cleanup()

		vfst.RunTests(t, fileSystem, "bashrc",
			vfst.TestPath("/home/user/.bashrc",
				vfst.TestModeIsRegular,
				vfst.TestContentsString("# contents of user's .bashrc\n"),
			),
		)
	}

	Test(&testing.T{})
}
