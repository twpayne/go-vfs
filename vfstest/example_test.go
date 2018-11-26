package vfstest_test

import (
	"testing"

	"github.com/twpayne/go-vfs/vfstest"
)

func ExampleNewTempFS_complex() {

	Test := func(t *testing.T) {
		// Describe the structure of the filesystem using a map from filenames to
		// file or directory contents.
		root := map[string]interface{}{
			// A string or []byte is sets a file's contents.
			"/home/user/.bashrc": "# contents of user's .bashrc\n",
			"/home/user/empty":   []byte{},
			// To set non-default permissions on a file, create an &vfstest.File.
			"/home/user/bin/hello.sh": &vfstest.File{
				Perm:     0755,
				Contents: []byte("echo hello\n"),
			},
			// Directories can be nested.
			"/home/user/foo": map[string]interface{}{
				"bar": map[string]interface{}{
					"baz": "qux",
				},
			},
			// To set non-default permissions on a directory, create an
			// &vfstest.Dir.
			"/root": &vfstest.Dir{
				Perm: 0700,
				Entries: map[string]interface{}{
					".bashrc": "# contents of root's .bashrc\n",
				},
			},
		}

		// Create and populate an *vfs.FS
		fs, cleanup, err := vfstest.NewTempFS(root)
		defer cleanup()
		if err != nil {
			t.Fatal(err)
		}

		// Create tests by creating data structures containing Tests.
		tests := []interface{}{
			// Test multiple properties of a single path with TestPath.
			vfstest.TestPath("/home",
				vfstest.TestIsDir,
				vfstest.TestModePerm(0755)),
			vfstest.TestPath("/home/user",
				vfstest.TestIsDir,
				vfstest.TestModePerm(0755)),
			vfstest.TestPath("/home/user/.bashrc",
				vfstest.TestModeIsRegular,
				vfstest.TestModePerm(0644),
				vfstest.TestContentsString("# contents of user's .bashrc\n")),
			// Maps with string keys create sub tests with testing.T.Run. The key
			// is used as the test name.
			map[string]interface{}{
				"home_user_empty": vfstest.TestPath("/home/user/empty",
					vfstest.TestModeIsRegular,
					vfstest.TestModePerm(0644),
					vfstest.TestSize(0)),
				"foo_bar_baz": vfstest.TestPath("/home/user/foo/bar/baz",
					vfstest.TestModeIsRegular,
					vfstest.TestModePerm(0644),
					vfstest.TestContentsString("qux")),
				"root": []interface{}{
					vfstest.TestPath("/root",
						vfstest.TestIsDir,
						vfstest.TestModePerm(0700)),
					vfstest.TestPath("/root/.bashrc",
						vfstest.TestModeIsRegular,
						vfstest.TestModePerm(0644),
						vfstest.TestContentsString("# contents of root's .bashrc\n")),
				},
			},
		}

		// RunTests traverses the data structure and running all Tests.
		vfstest.RunTests(t, fs, "", tests)
	}

	Test(&testing.T{})
}

func ExampleNewTempFS() {

	Test := func(t *testing.T) {
		fs, cleanup, err := vfstest.NewTempFS(map[string]string{
			"/home/user/.bashrc": "# contents of user's .bashrc\n",
		})
		defer cleanup()
		if err != nil {
			t.Fatal(err)
		}

		vfstest.RunTests(t, fs, "",
			vfstest.TestPath("/home/user/.bashrc",
				vfstest.TestContentsString("# contents of user's .bashrc\n")),
		)
	}

	Test(&testing.T{})
}
