# `go-vfs`

[![GoDoc](https://godoc.org/github.com/twpayne/go-vfs?status.svg)](https://godoc.org/github.com/twpayne/go-vfs)
[![Build Status](https://travis-ci.org/twpayne/go-vfs.svg?branch=master)](https://travis-ci.org/twpayne/go-vfs)
[![Build status](https://ci.appveyor.com/api/projects/status/m0nup45u310krjah?svg=true)](https://ci.appveyor.com/project/twpayne/go-vfs)
[![Report Card](https://goreportcard.com/badge/github.com/twpayne/go-vfs)](https://goreportcard.com/report/github.com/twpayne/go-vfs)

Package `go-vfs` provides an abstraction of the `os` and `ioutil` packages that
is easy to test.

## Key features

 * File system abstraction layer for commonly-used `os` and `ioutil` functions
   from the standard library.

 * Powerful testing framework, `vfst`. For a quick tour of `vfst`'s features,
   see [the examples in the
documentation](https://godoc.org/github.com/twpayne/go-vfs/vfst#pkg-examples).

## Quick start

`go-vfs` provides implementations of the `FS` interface:

```go
// An FS is an abstraction over commonly-used functions in the os and ioutil
// packages.
type FS interface {
    Chmod(name string, mode os.FileMode) error
    Chown(name string, uid, git int) error
    Chtimes(name string, atime, mtime time.Time) error
    Lchown(name string, uid, git int) error
    Lstat(name string) (os.FileInfo, error)
    Mkdir(name string, perm os.FileMode) error
    Open(name string) (*os.File, error)
    ReadDir(dirname string) ([]os.FileInfo, error)
    ReadFile(filename string) ([]byte, error)
    Readlink(name string) (string, error)
    Remove(name string) error
    RemoveAll(name string) error
    Rename(oldpath, newpath string) error
    Stat(name string) (os.FileInfo, error)
    Symlink(oldname, newname string) error
    Truncate(name string, size int64) error
    WriteFile(filename string, data []byte, perm os.FileMode) error
}
```

To use `go-vfs`, you write your code to use the `FS` interface, and then use
`vfst` to test it.

`go-vfs` also provides functions `MkdirAll` (equivalent to `os.MkdirAll`) and
`Walk` (equivalent to `filepath.Walk`) that operate on an `FS`.

The implementations of `FS` provided are:

 * `OSFS` which calls the underlying `os` and `ioutil` functions directly.

 * `PathFS` which transforms all paths to provide a poor-man's `chroot`.

 * `ReadOnlyFS` which prevents modification of the underlying FS.

 * `TestFS` which assists running tests on a real filesystem but in a temporary
   directory that is easily cleaned up.

Example usage:

```go
// writeConfigFile is the function we're going to test. It can make arbitrary
// changes to the filesystem through fs.
func writeConfigFile(fs vfs.FS) error {
    return fs.WriteFile("/home/user/app.conf", []byte(`app config`), 0644)
}

// TestWriteConfigFile is our test function.
func TestWriteConfigFile(t *testing.T) {
    // Create and populate an temporary directory with a home directory.
    fs, cleanup, err := vfst.NewTestFS(map[string]string{
        "/home/user/.bashrc": "# contents of user's .bashrc\n",
    })

    // Ensure that the temporary directory is removed.
    defer cleanup()

    // Check that the directory was populated successfully.
    if err != nil {
        t.Fatalf("vfsTest.NewTestFS(_) == _, _, %v, want _, _, <nil>", err)
    }

    // Call the function we want to test.
    if err := writeConfigFile(fs); err != nil {
        t.Error(err)
    }

    // Check properties of the filesystem after our function has modified it.
    vfst.RunTest(t, fs, "",
        vfst.PathTest("/home/user/app.conf",
            vfst.TestModeIsRegular,
            vfst.TestModePerm(0644),
            vfst.TestContentsString("app config"))),
}
```


## Motivation

`go-vfs` was inspired by
[`github.com/spf13/afero`](https://github.com/spf13/afero) and
[`github.com/twpayne/aferot`](https://github.com/twpayne/aferot). So, why not
use these?

 * `afero` has several critical bugs in its in-memory mock filesystem
   implementation `MemMapFs`, to the point that it is unusable for non-trivial
test cases.  `vfs` does not attempt to implent an in-memory mock filesystem,
and instead only provides a thin layer around the standard libary's `os` and
`ioutil` packages, and as such should have fewer bugs.

 * `afero` does not support creating or reading symbolic links, and its
   `LstatIfPossible` interface is clumsy to use as it is not part of the
`afero.Fs` interface. `vfs` provides out-of-the-box support for symbolic links
with all methods in the `FS` interface.

 * `afero` has been effectively abandoned by its author, and a "friendly fork"
   ([`github.com/absfs/afero`](https://github.com/absfs/afero)) has not seen
much activity. `vfs`, by providing much less functionality than `afero`, should
be smaller and easier to maintain.

## License

The MIT License (MIT)

Copyright (c) 2018 Tom Payne

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
