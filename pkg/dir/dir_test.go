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
	"path"
	"runtime"
	"testing"

	"github.com/slukits/gini/pkg/env"
	. "github.com/slukits/gounit"
)

const GINI_WEB_REPRO = "github.com/slukits/gini"

type _Dir struct{ Suite }

func (s *_Dir) SetUp(t *T) { t.Parallel() }

func (s *_Dir) Path_defaults_to_environments_working_directory(t *T) {
	t.Eq((&env.Env{}).WD(), (&Dir{}).String())
}

func (s *_Dir) Caller_provides_dir_of_caller_file(t *T) {
	_, f, _, ok := runtime.Caller(0)
	t.FatalIfNot(ok)
	t.Eq(path.Dir(f), (&Dir{}).Caller())
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

func TestDir(t *testing.T) {
	t.Parallel()
	Run(&_Dir{}, t)
}
