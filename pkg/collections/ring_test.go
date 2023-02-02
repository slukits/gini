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

package ring

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	. "github.com/slukits/gounit"
)

func init() {
	rand.Seed(time.Now().UnixMicro())
}

type ARing struct{ Suite }

func (s *ARing) SetUp(t *T) { t.Parallel() }

func (s *ARing) Capacity_defaults_to_twenty_elements(t *T) {
	t.Eq(20, (&Ring{}).Len())
}

func (s *ARing) Provides_no_element_if_empty(t *T) {
	r := &Ring{}
	providedElementsCount := 0
	r.For(func(i int, e interface{}) (stop bool) {
		providedElementsCount++
		return
	})
	t.Eq(0, providedElementsCount)
}

func elementsFX(n int) (ee []interface{}) {
	for i := 0; i < n; i++ {
		switch i {
		case 0:
			ee = append(ee, "1st")
		case 1:
			ee = append(ee, "2nd")
		case 2:
			ee = append(ee, "3rd")
		default:
			ee = append(ee, fmt.Sprintf("%dth", i+1))
		}
	}
	return ee
}

func (s *ARing) Below_its_capacity_has_added_elements(t *T) {
	r := &Ring{}
	n, got := rand.Intn(r.Len()-2), map[interface{}]bool{}
	exp := elementsFX(n + 1) // at least one elements

	r.Add(exp[0], exp[1:]...)
	r.For(func(_ int, e interface{}) (stop bool) {
		got[e] = true
		return
	})
	t.FatalIfNot(t.Eq(len(exp), len(got)))
	for _, e := range exp {
		t.True(got[e])
	}
}

func (s *ARing) Provides_elements_last_in_first_out(t *T) {
	r := &Ring{}
	ee := elementsFX(r.Len())
	r.Add(ee[0], ee[1:]...)
	providedElementsCount := 0
	r.For(func(i int, e interface{}) (stop bool) {
		i++
		t.Eq(e.(string), ee[r.Len()-i])
		providedElementsCount = i
		return
	})
	t.Eq(20, providedElementsCount)
}

func (s *ARing) Provides_latest_capacity_many_elements_on_overflow(
	t *T,
) {
	r := &Ring{}
	n, got := rand.Intn(r.Len())+r.Len()+1, make([]interface{}, r.Len())
	exp := elementsFX(n)
	r.Add(exp[0], exp[1:]...)
	r.For(func(i int, e interface{}) (stop bool) {
		got[i] = e
		return
	})
	t.Eq(r.Len(), len(got))
	expIdx := len(exp)
	for _, g := range got {
		expIdx--
		t.Eq(exp[expIdx], g)
	}
}

func (s *ARing) Reduces_its_elements_on_capacity_reduction(t *T) {
	r := &Ring{}
	smallerCapacity, ee := rand.Intn(r.Len()-2)+1, elementsFX(r.Len())
	r.Add(ee[0], ee[1:]...)
	r.Capacity = smallerCapacity
	got := make([]interface{}, smallerCapacity)
	r.For(func(i int, e interface{}) (stop bool) {
		got[i] = e
		return
	})
	eeIdx := len(ee)
	for _, e := range got {
		eeIdx--
		t.Eq(ee[eeIdx], e)
	}
}

func (s *ARing) Increases_its_elements_on_capacity_increase(t *T) {
	r, ee := &Ring{}, elementsFX(25)
	r.Add(ee[0], ee[1:20]...)
	r.Capacity = len(ee)
	r.Add(ee[20], ee[21:25]...)
	got := make([]interface{}, len(ee))
	r.For(func(i int, e interface{}) (stop bool) {
		got[i] = e
		return
	})
	eeIdx := len(ee)
	for _, e := range got {
		eeIdx--
		t.Eq(ee[eeIdx], e)
	}
}

func TestARing(t *testing.T) {
	t.Parallel()
	Run(&ARing{}, t)
}
