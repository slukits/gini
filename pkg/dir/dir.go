/*
Copyright 2022 - present Stephan Lukits. All rights reserved.
Use of this source code is governed by the GNU GPLv3 that can
be found in the LICENSE file.

This file is part of GINI.

GINI is free software: you can redistribute it and/or modify it
under the terms of the GNU General Public License as published
by the Free Software Foundation, either version 3 of the License,
or (at your option) any later version.

GINI is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with GINI. If not, see <https://www.gnu.org/licenses/#GPL>.
*/

/*
Package dir provides features around the file system which exceed
provided features by the os-package.
*/
package dir

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/slukits/gini/pkg/env"
	"github.com/slukits/gini/pkg/lg"
)

// Dir provides file-system operation compounded of os-file-system
// operations.  The zero-type is ready to use.
type Dir struct {

	// Log is a logger for reporting errors it defaults to the
	// zero-logger.
	Log *lg.Logger

	// Lib provides the std-lib functions Dir needs to provide its
	// features
	Lib Lib

	// Path of an Dir instance defaulting to Log.Env's working directory.
	Path string

	initLib bool
}

func (d *Dir) env() *env.Env {
	if d.lg().Env == nil {
		d.Log.Env = &env.Env{}
	}
	return d.Log.Env
}

func (d *Dir) path() string {
	if d.Path == "" {
		d.Path = d.env().WD()
	}
	return d.Path
}

func (d *Dir) lg() *lg.Logger {
	if d.Log == nil {
		d.Log = &lg.Logger{}
	}
	return d.Log
}

func (d *Dir) lib() Lib {
	if !d.initLib {
		d.initLib = true
		if d.Lib.Caller == nil {
			d.Lib.Caller = runtime.Caller
		}
		if d.Lib.ReadDir == nil {
			d.Lib.ReadDir = os.ReadDir
		}
		if d.Lib.ReadFile == nil {
			d.Lib.ReadFile = ioutil.ReadFile
		}
	}
	return d.Lib
}

// WalkUp returns a closure up which goes with each call one directory up
// starting at given directory d.  up returns nil and false if the last
// provided directory has no parent directory.
func (d *Dir) WalkUp() (up func() (*Dir, bool)) {
	p := d.path()
	return func() (*Dir, bool) {
		if p == filepath.Dir(p) {
			return nil, false
		}
		p = filepath.Dir(p)
		return &Dir{Log: d.lg(), Path: p}, true
	}
}

// Repo walks given Dir d's path up until a directory d is found
// containing a .git directory and returns d along with a true value.
// If no such directory is found nil and false is returned.
func (d *Dir) Repo() (*Dir, bool) {
	up, _d, next := d.WalkUp(), d, true
	for next && !_d.Contains(".git") {
		_d, next = up()
	}
	if _d == nil {
		return nil, false
	}
	return _d, true
}

// Caller returns the directory of the file containing Caller calling
// function.  Caller calls Log.Fatal if runtime-caller can't be
// determined.
func (d *Dir) Caller() *Dir {
	_, f, _, ok := d.lib().Caller(1)
	if !ok {
		d.lg().Fatal(
			"gini: pkg: dir: Caller: can't determine caller stack")
	}
	return &Dir{
		Log:  d.lg(),
		Path: filepath.Dir(f),
	}
}

// Contains returns true if given directory d contains a directory with
// given name dirName;  otherwise false is returned.  If given directory
// d is not readable an error is logged to lg.ERR.
func (d *Dir) Contains(dirName string) bool {
	ee, err := d.lib().ReadDir(d.path())
	if err != nil {
		d.lg().Tof(lg.ERR, "gini: pkg: dir: contains: %v", err)
		return false
	}
	for _, e := range ee {
		if !e.IsDir() {
			continue
		}
		if e.Name() != dirName {
			continue
		}
		return true
	}
	return false
}

// FileContains returns ture if given file fl in given directory d
// contains given bytes bb; otherwise false is returned.  If fl can't be
// read an error is logged to lg.ERR.
func (d *Dir) FileContains(fl string, bb []byte) bool {
	fbb, err := d.lib().ReadFile(filepath.Join(d.path(), fl))
	if err != nil {
		d.lg().Tof(lg.ERR, "gini: pkg: dir: file-contains: %v", err)
		return false
	}
	return bytes.Contains(fbb, bb)
}

// String returns given directories d path.
func (d *Dir) String() string {
	return d.path()
}

// Lib provides std-lib functions which may fail.
type Lib struct {

	// Caller defaults to runtime.Caller and its semantics
	Caller func(int) (pc uintptr, file string, line int, ok bool)

	// ReadFile defaults to ioutil.ReadFile and its semantics
	ReadFile func(name string) ([]byte, error)

	// ReadDir defaults to os.ReadDir and its semantics
	ReadDir func(name string) ([]fs.DirEntry, error)
}
