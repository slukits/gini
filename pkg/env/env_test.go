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
	"errors"
	"os"
	"strings"
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

func (s *env) Panics_if_no_home_dir_and_no_fatal_handler(t *T) {
	env := &Env{}
	env.Lib.UserHomeDir = func() (string, error) {
		return "", errors.New("err: home-dir mock")
	}
	t.Panics(func() { env.Home() })
}

func (s *env) Panics_if_fatal_handler_isnt_stopping_execution(t *T) {
	env := &Env{FatalHandler: &fatalerMock{fatal: func(i ...interface{}) {}}}
	env.Lib.UserHomeDir = func() (string, error) {
		return "", errors.New("err: home-dir mock")
	}
	t.Panics(func() { env.Home() })
}

func (s *env) Fatales_if_no_home_dir_and_fatal_handler(t *T) {
	ftlExp := "home-dir fatal mock"
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("fatal wasn't called")
		}
		t.Eq(r, ftlExp)
	}()
	fatal := &fatalerMock{fatal: func(i ...interface{}) {
		panic(ftlExp)
	}}
	env := &Env{FatalHandler: fatal}
	env.Lib.UserHomeDir = func() (string, error) {
		return "", errors.New("err: home-dir mock")
	}
	env.Home()
}

func (e *env) Sets_its_home_directory_to_given(t *T) {
	path := t.FS().Tmp().Path()
	env := (&Env{}).SetHome(path)
	t.Eq(path, env.Home())
}

func (e *env) Setting_home_dir_is_noop_if_zero_str_given(t *T) {
	var env Env
	home := env.Home()
	env.SetHome("")
	t.Eq(home, env.Home())
}

func (s *env) Is_not_user_env_if_nil(t *T) {
	t.Not.True((*Env)(nil).IsUser())
}

func (s *env) Fatal_if_is_user_cant_get_home(t *T) {
	first := true
	ftl, rcv := fatalMockRecover(t, "is-user home-dir fatal mock")
	defer rcv()
	env := &Env{FatalHandler: ftl}
	env.Lib.UserHomeDir = func() (string, error) {
		if first {
			first = false
			home, err := os.UserHomeDir()
			t.FatalOn(err)
			return home, nil
		}
		return "", errors.New("err: home-dir mock")
	}
	env.IsUser()
}

func (s *env) Is_not_associated_with_tmp_if_nil(t *T) {
	t.Not.True((*Env)(nil).IsTemp())
}

func (s *env) Is_not_associated_with_an_temp_dir_by_default(t *T) {
	t.Not.True((&Env{}).IsTemp())
}

func (s *env) Is_associated_with_an_temp_dir_if_home_set_to_one(t *T) {
	env := (&Env{}).SetHome(t.FS().Tmp().Path())
	t.True(env.IsTemp())
}

func (s *env) Fatal_if_no_working_dir_and_fatal_handler(t *T) {
	ftl, rcv := fatalMockRecover(t, "working-dir fatal mock")
	defer rcv()
	env := &Env{FatalHandler: ftl}
	env.Lib.Getwd = func() (string, error) {
		return "", errors.New("err: working-dir mock")
	}
	env.WD()
}

func (s *env) Working_directory_is_os_working_directory(t *T) {
	wd, err := os.Getwd()
	t.FatalOn(err)
	t.Eq(wd, (&Env{}).WD())
}

func (s *env) Fatal_if_config_dir_cant_be_determined(t *T) {
	ftl, rcv := fatalMockRecover(t, "config-dir fatal mock")
	defer rcv()
	env := &Env{FatalHandler: ftl}
	env.Lib.UserConfigDir = func() (string, error) {
		return "", errors.New("err: config-dir mock")
	}
	env.Lib.UserHomeDir = func() (string, error) {
		return "", errors.New("err: working-dir mock")
	}
	env.Conf()
}

func (e *env) Config_dir_is_in_home_dir_if_no_os_config(t *T) {
	env := &Env{}
	env.Lib.UserConfigDir = func() (string, error) {
		return "", errors.New("err: config-dir mock")
	}
	t.True(strings.HasPrefix(env.Conf(), env.Home()))
}

func (e *env) Config_dir_is_in_home_dir_if_not_user_home(t *T) {
	env := (&Env{}).SetHome(t.FS().Tmp().Path())
	t.True(strings.HasPrefix(env.Conf(), env.Home()))
}

func (e *env) Config_dir_is_in_os_s_config_dir_by_default(t *T) {
	cnf, err := os.UserConfigDir()
	t.FatalOn(err)
	t.True(strings.HasPrefix((&Env{}).Conf(), cnf))
}

func (e *env) Logging_dir_is_in_config_dir(t *T) {
	env := &Env{}
	t.True(strings.HasPrefix(env.Logging(), env.Conf()))
}

func (e *env) Creates_logging_directory(t *T) {
	env := (&Env{}).SetHome(t.FS().Tmp().Path())
	t.FatalIfNot(t.True(strings.HasPrefix(env.Logging(), env.Home())))
	_, err := os.Stat(env.Logging())
	t.True(err != nil)
	t.FatalOn(env.MkLogging())
	_, err = os.Stat(env.Logging())
	t.True(err == nil)
}

func (e *env) Errors_if_working_directory_cant_be_changed(t *T) {
	env := (&Env{}).SetHome(t.FS().Tmp().Path())
	env.Lib.Chdir = func(path string) error {
		return errors.New("chdir error mock")
	}
	err := env.ChWD("test")
	t.Contains(err.Error(), "chdir error mock")
}

func (e *env) Changes_working_directory(t *T) {
	home := t.FS().Tmp().Path()
	env := (&Env{}).SetHome(home)
	env.Lib.Chdir = func(path string) error {
		t.Eq(path, home)
		return nil
	}
	t.FatalOn(env.ChWD(env.Home()))
	t.Eq(env.WD(), env.Home())
}

func fatalMockRecover(t *T, exp string) (Fataler, func()) {

	// fatalerMock must panic to end the execution of the function
	// calling fatal
	return &fatalerMock{func(i ...interface{}) { panic(exp) }},

		// since fataler should panic that panic needs to be recovered
		// to continue with the tests.
		func() {
			r := recover()
			if r == nil {
				t.Fatal("fatal wasn't called")
			}
			t.Eq(r, exp)
		}
}

func TestEnv(t *testing.T) {
	t.Parallel()
	Run(&env{}, t)
}

type fatalerMock struct {
	fatal func(...interface{})
}

func (m *fatalerMock) Fatal(vv ...interface{}) {
	m.fatal(vv...)
}
