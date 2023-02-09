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

package view

import (
	"testing"

	"github.com/slukits/gini/cmd/gini/view/internal/cnt"
	. "github.com/slukits/gounit"
	"github.com/slukits/lines"
)

type AView struct{ Suite }

func (s *AView) SetUp(t *T) { t.Parallel() }

func (s *AView) Has_a_context_component(t *T) {
	vw := &View{}
	lines.TermFixture(t.GoT(), 0, vw)
	_, ok := vw.Context().(*cnt.Context)
	t.True(ok)
}

func TestAView(t *testing.T) {
	t.Parallel()
	Run(&AView{}, t)
}
