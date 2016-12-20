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

type TextPart struct {
	Text []byte
}

//===========================================
// Main methods
//===========================================

func (t *TextPart) SetText(newText []byte) {
	t.Text = make([]byte, len(newText))
	copy(t.Text, newText)
}

//===========================================
// Part interface methods
//===========================================

var _ Part = (*TextPart)(nil) // Verify that *TextPart implements Part.

func (t *TextPart) GetId() string {
	return ""
}

func (t *TextPart) AddPart(p Part) {
}

func (t *TextPart) GetChildren() []Part {
	return make([]Part, 0)
}

func (t *TextPart) ResetChildren() {
}

func (t *TextPart) ToString() string {
	return "TEXT"
}
