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

package dir

import (
	"errors"
	"io/fs"
	"path"
	"runtime"
	"testing"

	"github.com/slukits/gini/pkg/env"
	"github.com/slukits/gini/pkg/lg"
	. "github.com/slukits/gounit"
)

const GINI_WEB_REPRO = "github.com/slukits/gini"

type _Dir struct{ Suite }

func (s *_Dir) SetUp(t *T) { t.Parallel() }

func (s *_Dir) Path_defaults_to_environments_working_directory(t *T) {
	var d Dir
	got := d.String() // this call initializes d.Log.Env
	t.Eq(d.Log.Env.WD(), got)
}

func (s *_Dir) Caller_provides_dir_of_caller_file(t *T) {
	_, f, _, ok := runtime.Caller(0)
	t.FatalIfNot(ok)
	t.Eq(path.Dir(f), (&Dir{}).Caller())
}

func (s *_Dir) Fatal_if_caller_cant_be_determined(t *T) {
	var d Dir
	d.Lib.Caller = func(i int) (
		pc uintptr, file string, line int, ok bool,
	) {
		return 0, "", 0, false
	}
	d.Log.Lib.Fatal = func(vv ...interface{}) {
		panic(vv[0])
	}
	defer func() {
		t.Contains(recover().(string), "can't determine caller")
	}()
	d.Caller()
}

func (s *_Dir) Repo_defaults_to_dir_of_current_binary(t *T) {
	// NOTE we can inject any path nested inside the repo path
	// to mock-up the "binary-path"
	d := &Dir{}
	d.Path = d.Caller().String()
	repo, ok := d.Repo()
	t.FatalIfNot(t.True(ok))
	t.True(repo.FileContains(".git/config", []byte(GINI_WEB_REPRO)))
}

func (s *_Dir) Repo_fails_if_not_inside_a_repo(t *T) {
	d := &Dir{Path: t.FS().Tmp().Path()}
	repo, ok := d.Repo()
	t.Not.True(ok)
	t.Eq((*Dir)(nil), repo)
}

func (s *_Dir) Doesnt_contain_a_dir_if_given_dir_not_readable(t *T) {
	d := &Dir{Log: lg.Logger{
		Env: (&env.Env{}).SetHome(t.FS().Tmp().Path())}}
	d.Lib.ReadDir = func(name string) ([]fs.DirEntry, error) {
		return nil, errors.New("read-dir error mock")
	}
	t.Not.True(d.Contains("blub"))
}

func logFX(t *T) lg.Logger {

	// make logger a in-memory logger by setting environments home to a
	// temporary directory.
	lgg := lg.Logger{Env: (&env.Env{}).SetHome(t.FS().Tmp().Path())}

	// make only the environment working directory change but the
	// executing binaries.
	lgg.Env.Lib.Chdir = func(path string) error { return nil }

	lgg.Env.ChWD(lgg.Env.Home())
	return lgg
}

func (s *_Dir) Reports_error_if_given_dir_is_not_readable(t *T) {
	d := &Dir{Log: logFX(t)}
	d.Lib.ReadDir = func(name string) ([]fs.DirEntry, error) {
		return nil, errors.New("read-dir error mock")
	}
	d.Contains("blub")
	t.Contains(d.Log.String(lg.ERR), "read-dir error mock")
}

func (s *_Dir) File_doesnt_contain_bytes_if_file_not_readable(t *T) {
	d := &Dir{Log: logFX(t)}
	t.Not.True(d.FileContains("blub", []byte("42")))
}

func (s *_Dir) Reports_error_if_queried_file_is_not_readable(t *T) {
	d := &Dir{Log: logFX(t)}
	d.FileContains("blub", []byte("42"))
	t.Contains(d.Log.String(lg.ERR), "file-contains")
	t.Contains(d.Log.String(lg.ERR), "no such file")
}

func TestDir(t *testing.T) {
	t.Parallel()
	Run(&_Dir{}, t)
}
