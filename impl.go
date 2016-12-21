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

type ImplPart struct {
	Id        string
	Children  []Part
	Processed bool
}

//===========================================
// Part interface methods
//===========================================

var _ Part = (*ImplPart)(nil) // Verify that *ImplPart implements Part.

func (i *ImplPart) GetId() string {
	return i.Id
}

func (i *ImplPart) AddPart(p Part) {
	i.Children = append(i.Children, p)
}

func (i *ImplPart) GetChildren() []Part {
	return i.Children
}

func (i *ImplPart) ResetChildren() {
	if i.Children != nil {
		i.Children = i.Children[0:0]
	}
}

func (i *ImplPart) ToString() string {
	return fmt.Sprintf("IMPL: %s", i.Id)
}

func (i *ImplPart) Write(buf *bytes.Buffer) {
	buf.WriteString("<!--pvnImplBeg ")
	buf.WriteString(i.GetId())
	buf.WriteString("-->")

	for _, part := range i.Children {
		part.Write(buf)
	}

	buf.WriteString("<!--pvnImplEnd ")
	buf.WriteString(i.GetId())
	buf.WriteString("-->")
}
