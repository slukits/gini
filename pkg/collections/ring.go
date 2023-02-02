/*
Copyright 2023 - present Stephan Lukits. All rights reserved.
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
Package ring provides a general purpose ring data structure.
*/
package ring

// Ring a can hold Ring.Capacity many elements whereas the oldest
// element is removed to make space for a new added element if the
// capacity is reached.  Stored elements are provide last in first out.
// The capacity may be adapted by the client and a Ring's zero value is
// ready to use.
type Ring struct {

	// Capacity represents the number of elements a ring can hold.  Is
	// the capacity not positive it defaults to 20.
	Capacity int

	data        []interface{}
	nextElement int
}

// Len returns the capacity of a ring which defaults to 20 elements.
func (r *Ring) Len() int {
	r.ensureCapacityConsistency()
	return r.Capacity
}

func (r *Ring) ensureCapacityConsistency() {
	if r.Capacity > 0 && r.Capacity == len(r.data) {
		return
	}
	if r.Capacity <= 0 {
		r.Capacity = 20
		r.data = make([]interface{}, 20)
		return
	}
	if r.Capacity < len(r.data) {
		r.data = r.data[len(r.data)-r.Capacity:]
		if r.nextElement > r.Capacity {
			r.nextElement = r.Capacity
		}
		return
	}
	r.data = append(
		r.data,
		make([]interface{}, r.Capacity-len(r.data))...,
	)
}

// Add adds given element e (and elements ee) to given ring r whereas
// the oldest element is removed to make space for a new element if r
// holds capacity many elements.
func (r *Ring) Add(e interface{}, ee ...interface{}) {
	r.add(e)
	for _, e := range ee {
		r.add(e)
	}
}

func (r *Ring) add(e interface{}) {
	if r.nextElement < r.Len() {
		r.data[r.nextElement] = e
		r.nextElement++
		return
	}
	copy(r.data, r.data[1:])
	r.data[r.nextElement-1] = e
}

// For provides to given receiver rcv stored elements of given ring
// r last in first out until r runs out of elements or the receiver
// stops the callback by returning true.
func (r *Ring) For(rcv func(int, interface{}) (stop bool)) {
	r.ensureCapacityConsistency()
	for i := r.nextElement - 1; i >= 0; i-- {
		rcv(r.nextElement-1-i, r.data[i])
	}
}
