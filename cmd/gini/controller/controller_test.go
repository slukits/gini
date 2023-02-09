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
controller tests are typically end to end tests.  Hence being able to
inject a lines terminal fixture for ui creation which allows to simulate
user input and retrieve screen states should be sufficient for testing.
In case we actually do need a controller instance to test something we
can extend the Init-type by a "Listener func(*Controller)"-property
which gets right after (!) the call "ll.WaitForQuit()" the tested
controller instance reported.
*/

package controller

import (
	"fmt"
	"testing"

	"github.com/slukits/gini/cmd/gini/view"
	. "github.com/slukits/gounit"
	"github.com/slukits/lines"
)

type GINI struct{ Suite }

func (s *GINI) Ui_factory_defaults_to_lines_terminal(t *T) {
	exp := fmt.Sprintf("%T::%[1]p", lines.Term)
	got := fmt.Sprintf("%T::%[1]p", (&Init{}).UIFactory())
	t.Eq(exp, got)
}

func (s *GINI) Exits_if_ui_cannot_be_obtained(t *T) {
	var init Init
	init.Log.Lib.Fatal = func(vv ...interface{}) {
		panic(vv[0])
	}
	init.Lines = func(c lines.Componenter) *lines.Lines {
		panic("terminal screen can't be obtained")
	}
	defer func() {
		t.Contains(recover().(string), "GINI: controller: panic")
	}()
	New(init)
}

func (s *GINI) Passes_non_blocking_on_lines_fixture_for_ui(t *T) {
	var init Init
	init.Log.Lib.Fatal = func(vv ...interface{}) {
		panic(vv[0])
	}
	init.Lines = func(c lines.Componenter) *lines.Lines {
		return lines.TermFixture(t.GoT(), 0, &view.View{}).Lines
	}
	blockingChan := make(chan struct{})
	blockingFunc := func(c chan struct{}) {
		New(init)
		close(c)
	}
	blockingFunc(blockingChan)
	select {
	case <-blockingChan:
	case <-t.Timeout(0):
		t.Error("controller initialization timed out.")
	}
}

func TestGINI(t *testing.T) {
	t.Parallel()
	Run(&GINI{}, t)
}
