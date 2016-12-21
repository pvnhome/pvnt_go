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
	"bytes"
	"fmt"
)

type TmplPart struct {
	Id       string
	Children []Part
}

//===========================================
// Part interface methods
//===========================================

var _ Part = (*TmplPart)(nil) // Verify that *TmplPart implements Part.

func (t *TmplPart) GetId() string {
	return t.Id
}

func (t *TmplPart) AddPart(p Part) {
	t.Children = append(t.Children, p)
}

func (t *TmplPart) GetChildren() []Part {
	return t.Children
}

func (t *TmplPart) ResetChildren() {
	if t.Children != nil {
		t.Children = t.Children[0:0]
	}
}

func (t *TmplPart) ToString() string {
	if t.isRoot() {
		return "TMPL: root"
	} else {
		return fmt.Sprintf("TMPL: file = %s", t.Id)
	}
}

func (t *TmplPart) Write(buf *bytes.Buffer) {
	if t.isRoot() {
		buf.WriteString("<!--pvnTmplBeg-->")
	} else {
		buf.WriteString("<!--pvnTmplBeg ")
		buf.WriteString(t.GetId())
		buf.WriteString("-->")
	}

	for _, part := range t.Children {
		part.Write(buf)
	}

	buf.WriteString("<!--pvnTmplEnd-->")
}

//===========================================
// Public methods
//===========================================

func (t *TmplPart) isRoot() bool {
	return len(t.Id) < 1
}
