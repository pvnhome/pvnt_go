// PVN Template utility
// Copyright (C) 2016,  Victor Pyankov
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
)

type EditPart struct {
	Id       string
	Children []Part
}

//===========================================
// Part interface methods
//===========================================

var _ Part = (*EditPart)(nil) // Verify that *EditPart implements Part.

func (e *EditPart) GetId() string {
	return e.Id
}

func (e *EditPart) AddPart(p Part) {
	e.Children = append(e.Children, p)
}

func (e *EditPart) GetChildren() []Part {
	return e.Children
}

func (e *EditPart) ResetChildren() {
	if e.Children != nil {
		e.Children = e.Children[0:0]
	}
}

func (e *EditPart) ToString() string {
	return fmt.Sprintf("EDIT: id = %s", e.Id)
}
