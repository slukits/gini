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
Package env provides access to GINI instances environment and is used as
identifier for environment specific type-instances of various types
making concurrency save parallel testing easier.
*/
package env

import (
	"os"
	"testing"

	. "github.com/slukits/gounit"
)

type env struct{ Suite }

func (s *env) SetUp(t *T) { t.Parallel() }

func (s *env) Home_dir_defaults_to_user_s_home_directory(t *T) {
	user, err := os.UserHomeDir()
	t.FatalOn(err)
	env := &Env{}
	t.Eq(user, env.Home())
	t.True(env.IsUser())
}

func (s *env) Is_not_associated_with_an_temp_dir_by_default(t *T) {
	t.Not.True((&Env{}).IsTemp())
}

func (s *env) Is_associated_with_an_temp_dir_if_root_set_to_one(t *T) {
	env := &Env{Root: t.FS().Tmp().Path()}
	t.True(env.IsTemp())
}

func TestEnv(t *testing.T) {
	t.Parallel()
	Run(&env{}, t)
}
