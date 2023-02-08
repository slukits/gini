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

package lg

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/slukits/gini/pkg/env"
	. "github.com/slukits/gounit"
)

const tst = "test"
const tstFl = tst + "." + FileSuffix

type logger struct{ Suite }

func (s *logger) For_fatal_defaults_to_log_fatal(t *T) {
	exp := fmt.Sprintf("%T::%[1]p", log.Fatal)
	got := fmt.Sprintf("%T::%[1]p", (&Logger{}).LibDefaults().Fatal)
	t.Eq(exp, got)
}

func (s *logger) For_logs_without_env_defaults_to_log_println(t *T) {
	exp := fmt.Sprintf("%T::%[1]p", log.Println)
	got := fmt.Sprintf("%T::%[1]p", (&Logger{}).LibDefaults().Println)
	t.Eq(exp, got)
}

func (s *logger) Logs_not_to_memory_if_nil(t *T) {
	t.Not.True((*Logger)(nil).InMemory(""))
}

func (s *logger) With_temp_env_logs_to_memory_by_default(t *T) {
	lgg := &Logger{Env: (&env.Env{}).SetHome(t.FS().Tmp().Path())}
	lgg.To("test", "test")
	t.True(lgg.InMemory("test"))
}

// memFX returns an in-memory Logger-fixture according to above test.
func memFX(t *T) *Logger {
	lgg := &Logger{Env: (&env.Env{}).SetHome(t.FS().Tmp().Path())}

	// make only the environment working directory change but the
	// executing binaries.
	lgg.Env.Lib.Chdir = func(path string) error { return nil }
	lgg.Env.ChWD(lgg.Env.Home())

	return lgg
}

func (s *logger) Logs_not_to_memory_if_log_not_exists(t *T) {
	lgg := memFX(t)
	t.Not.True(lgg.InMemory("test"))
	lgg.To("other", "42")
	t.Not.True(lgg.InMemory("test"))
}

func (s *logger) Logs_to_default_handler_if_no_env_set(t *T) {
	dflt, dfltF, gotMsg := "default", "default f", ""
	lgg := &Logger{}
	lgg.Lib.Println = func(vv ...interface{}) {
		gotMsg = vv[0].(string)
	}
	lgg.To("test", dflt)
	t.Contains(gotMsg, dflt)
	lgg.Tof("test", dfltF)
	t.Contains(gotMsg, dfltF)
}

func (s *logger) Logs_to_log_with_given_name(t *T) {
	lgg, msg, msgF := memFX(t), "msg", "msgF"

	t.Not.Contains(lgg.String("test"), msg)
	lgg.To("test", msg)
	t.Contains(lgg.String("test"), msg)

	t.Not.Contains(lgg.String("test"), msgF)
	lgg.Tof("test", msgF)
	t.Contains(lgg.String("test"), msgF)
}

func (s *logger) Fatal_logs_to_default_fatal_logger(t *T) {
	lgg, ftl := &Logger{}, "fatal f"
	lgg.Lib.Fatal = func(vv ...interface{}) { panic(vv[0].(string)) }
	defer func() { t.Contains(recover().(string), ftl) }()
	lgg.Fatal(ftl)
}

func (s *logger) Fatalf_logs_to_default_fatal_logger(t *T) {
	lgg, ftlF := &Logger{}, "fatal f"
	lgg.Lib.Fatal = func(vv ...interface{}) { panic(vv[0].(string)) }
	defer func() { t.Contains(recover().(string), ftlF) }()
	lgg.Fatalf("%s", ftlF)
}

func (s *logger) Log_writer_is_nil_if_log_doesnt_exist(t *T) {
	t.Eq((io.Writer)(nil), (&Logger{}).Writer("test"))
}

func (s *logger) Log_flags_are_zero_if_log_doesnt_exits(t *T) {
	t.Eq(0, (&Logger{}).Flags("test"))
}

func (s *logger) Logs_non_error_messages_with_default_flags(t *T) {
	lgg := memFX(t)
	lgg.To(tst, "test flagging")
	t.Eq(Flags, lgg.Flags(tst))
}

func (s *logger) Logs_error_messages_with_error_flags(t *T) {
	lgg := memFX(t)
	lgg.To(ERR, "test flagging")
	t.Eq(ErrFlags, lgg.Flags(ERR))
}

func (s *logger) Writes_tmp_log_files_on_write_temp_logs(t *T) {
	lgg := memFX(t)
	lgg.WriteTempLogs = true
	exp := filepath.Join(lgg.Env.Logging(), tstFl)
	_, err := os.Stat(exp)
	t.True(err != nil)
	lgg.To(tst, "written log")
	_, err = os.Stat(exp)
	t.FatalOn(err)
}

func (s *logger) Uses_std_lib_functions_for_failing_file_operations(
	t *T,
) {
	lgg := memFX(t)
	lgg.WriteTempLogs = true
	lgg.To(tst, "written log") // Lib is set to defaults now

	exp := fmt.Sprintf("%T::%[1]p", os.OpenFile)
	got := fmt.Sprintf("%T::%[1]p", lgg.Lib.OpenFile)
	t.Eq(exp, got)

	exp = fmt.Sprintf("%T::%[1]p", ioutil.ReadFile)
	got = fmt.Sprintf("%T::%[1]p", lgg.Lib.ReadFile)
	t.Eq(exp, got)
}

func (s *logger) Panics_if_handled_error_doesnt_stop_execution(t *T) {
	lgg := memFX(t)
	lgg.WriteTempLogs = true
	lgg.Lib.Fatal = func(vv ...interface{}) {}
	lgg.Env.Lib.MkdirAll = func(string, fs.FileMode) error {
		return errors.New("mk-dir-all failing")
	}
	t.Panics(func() { lgg.To(tst, "written log") })
}

func (s *logger) Panics_if_fatal_doesnt_stop_execution(t *T) {
	lgg := memFX(t)
	lgg.WriteTempLogs = true
	lgg.Lib.Fatal = func(vv ...interface{}) {}
	t.Panics(func() { lgg.Fatal("stop execution") })
}

func (s *logger) Reports_failing_log_dir_creation(t *T) {
	lgg, errMsg := memFX(t), ""
	lgg.WriteTempLogs = true
	lgg.Lib.Fatal = func(vv ...interface{}) {
		errMsg = vv[0].(string)
		panic("fatal mock")
	}
	lgg.Env.Lib.MkdirAll = func(string, fs.FileMode) error {
		return errors.New("mk-dir-all mock failing")
	}
	defer func() {
		t.FatalIfNot(t.Eq("fatal mock", recover().(string)))
		t.Contains(errMsg, "create dir")
	}()
	lgg.To(tst, "written log")
}

func (s *logger) Reports_failing_log_file_creation(t *T) {
	lgg, errMsg := memFX(t), ""
	lgg.WriteTempLogs = true
	lgg.Lib.Fatal = func(vv ...interface{}) {
		errMsg = vv[0].(string)
		panic("fatal mock")
	}
	lgg.Lib.OpenFile =
		func(n string, ff int, m fs.FileMode) (*os.File, error) {
			return nil, errors.New("open file mock failing")
		}
	defer func() {
		t.FatalIfNot(t.Eq("fatal mock", recover().(string)))
		t.Contains(errMsg, "open log-file")
	}()
	lgg.To(tst, "written log")
}

func (s *logger) Reports_zero_log_if_nil(t *T) {
	t.Eq("", (*Logger)(nil).String(tst))
}

func (s *logger) Reports_log_file_content(t *T) {
	lgg, content := memFX(t), "log-file content without prefix"
	lgg.WriteTempLogs = true
	lgg.To(tst, content)
	t.Contains(lgg.String(tst), content)
}

func (s *logger) Reports_zero_content_if_zero_log_name(t *T) {
	t.Eq("", memFX(t).String(""))
}

func (s *logger) Reports_failing_log_file_closing(t *T) {
	lgg, errMsg := memFX(t), ""
	lgg.WriteTempLogs = true
	lgg.Lib.Fatal = func(vv ...interface{}) {
		errMsg = vv[0].(string)
		panic("fatal mock")
	}
	defer func() {
		t.FatalIfNot(t.Eq("fatal mock", recover().(string)))
		t.Contains(errMsg, "close log-file")
	}()
	lgg.To(tst, "written log")
	t.FatalOn(lgg.Writer(tst).(*os.File).Close())
	lgg.String(tst)
}

func (s *logger) Reports_failing_log_file_read(t *T) {
	lgg, errMsg := memFX(t), ""
	lgg.WriteTempLogs = true
	lgg.Lib.Fatal = func(vv ...interface{}) {
		errMsg = vv[0].(string)
		panic("fatal mock")
	}
	defer func() {
		t.FatalIfNot(t.Eq("fatal mock", recover().(string)))
		t.Contains(errMsg, "read log-file")
	}()
	lgg.To(tst, "written log")
	lgg.Lib.ReadFile = func(name string) ([]byte, error) {
		return nil, errors.New("read-file error mock")
	}
	lgg.String(tst)
}

func (s *logger) Reports_failing_reopening_log_file(t *T) {
	lgg, errMsg, gotMsg := memFX(t), "open-file error mock", ""
	lgg.WriteTempLogs = true
	lgg.Lib.Fatal = func(vv ...interface{}) {
		gotMsg = vv[0].(string)
		panic("fatal mock")
	}
	defer func() {
		t.FatalIfNot(t.Eq("fatal mock", recover().(string)))
		t.Contains(gotMsg, errMsg)
		t.Contains(gotMsg, "reopen log-file")
	}()
	lgg.To(tst, "written log")
	lgg.Lib.OpenFile = func(
		name string, flags int, perm fs.FileMode,
	) (*os.File, error) {
		return nil, errors.New(errMsg)
	}
	lgg.String(tst)
}

func TestLogger(t *testing.T) {
	t.Parallel()
	Run(&logger{}, t)
}
