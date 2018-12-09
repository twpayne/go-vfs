package vfsafero_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/twpayne/go-vfs/vfsafero"
	"github.com/twpayne/go-vfs/vfst"
)

func ExampleNewAferoFS() {

	Test := func(t *testing.T) {

		fs, cleanup, err := vfst.NewTestFS(map[string]interface{}{
			"/home/user/.bashrc": "# contents of .bashrc\n",
		})
		defer cleanup()
		if err != nil {
			t.Fatal(err)
		}

		aferoFS := vfsafero.NewAferoFS(fs)
		afero.WriteFile(aferoFS, "/home/user/foo", []byte("bar"), 0666)

		vfst.RunTests(t, fs, "",
			vfst.TestPath("/home/user/foo",
				vfst.TestContentsString("bar"),
			),
		)

	}

	Test(&testing.T{})
}
