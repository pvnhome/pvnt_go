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
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type File struct {
	Site *Site  // Site
	Path string // Path to file

	fileLen int         // File length
	perm    os.FileMode // // File mode bits

	Parent   *File   // Parent file
	Children []*File // File children
	Parts    []Part  // File parts

	Tmpl    *TmplPart            // Template part
	ImplMap map[string]*ImplPart // Implementations map
}

//===========================================
// Part interface methods
//===========================================

var _ Part = (*File)(nil) // Verify that *File implements Part.

func (f *File) GetId() string {
	return f.Path
}

func (f *File) AddPart(p Part) {
	f.Parts = append(f.Parts, p)
}

func (f *File) GetChildren() []Part {
	return f.Parts
}

func (f *File) ResetChildren() {
	if f.Parts != nil {
		f.Parts = f.Parts[0:0]
	}
}

func (f *File) ToString() string {
	return fmt.Sprintf("FILE: %s", f.Path)
}

func (f *File) Write(buf *bytes.Buffer) {
}

//===========================================
// File methods
//===========================================

func NewFile(site *Site, name string) *File {
	newFile := &File{Site: site, Path: name}
	newFile.Parts = make([]Part, 0, PARTS_INITIAL_CAPACITY)
	newFile.Children = make([]*File, 0, FILES_INITIAL_CAPACITY)
	newFile.ImplMap = make(map[string]*ImplPart)
	return newFile
}

func (f *File) Load() ([]byte, error) {
	info, err := os.Stat(f.Path)
	if err != nil {
		return nil, err
	}

	f.perm = info.Mode()

	content, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return content, err
	}

	f.fileLen = len(content)

	return content, nil
}

func (f *File) Save() {
	buf := bytes.NewBuffer(make([]byte, 0, f.fileLen+f.fileLen/2))

	for _, part := range f.Parts {
		part.Write(buf)
	}

	if err := ioutil.WriteFile(f.Path, buf.Bytes(), f.perm); err != nil {
		fmt.Println("ERROR: save file =", f.Path, ",", err)
	}
}

func (f *File) GetImplIds() string {
	ids := make([]string, len(f.ImplMap))

	i := 0
	for key := range f.ImplMap {
		ids[i] = key
		i++
	}
	sort.Strings(ids)

	return strings.Join(ids, ", ")
}

func (f *File) IsRoot() bool {
	return f.Parent == nil
}

func (f *File) AddChild(child *File) {
	child.Parent = f
	f.Children = append(f.Children, child)
}

func (f *File) Process() {
	if !f.IsRoot() {
		f.Site.printDebugMessage("process file: %s", f.Path)
		f.ResetChildren()
		f.process(f.Parent.GetChildren(), f)

		ids := make([]string, 0, len(f.ImplMap))
		for key, impl := range f.ImplMap {
			if !impl.Processed {
				ids = append(ids, key)
			}
		}

		if len(ids) > 0 {
			sort.Strings(ids)

			var sb bytes.Buffer

			sb.WriteString("We have IMPL tags [")
			sb.WriteString(strings.Join(ids, ", "))
			sb.WriteString("] without EDIT tags in template in ")
			sb.WriteString(f.Path)

			panic(PvntError{sb.String()})
		}
	}

	for _, file := range f.Children {
		file.Process()
	}
}

//===========================================
// Private methods.
//===========================================

func (f *File) process(templateParts []Part, currentPart Part) {
	for _, part := range templateParts {
		if _, ok := part.(*TextPart); ok {
			currentPart.AddPart(part)
		} else {
			_, isEdit := part.(*EditPart)
			implPart, isImplFound := f.ImplMap[part.GetId()]
			if isEdit && isImplFound {
				implPart.Processed = true
				currentPart.AddPart(implPart)
			} else {
				if _, ok := part.(*ImplPart); ok && isImplFound {
					implPart.Processed = true
					// TODO We need to compare part with implPart and we need to print message if the difference will be detected
				}

				clonedPart := clonePart(f, part)
				currentPart.AddPart(clonedPart)
				if len(part.GetChildren()) > 0 {
					f.process(part.GetChildren(), clonedPart)
				}
			}
		}
	}
}

func clonePart(f *File, part Part) Part {
	switch t := part.(type) {
	case *EditPart:
		return &EditPart{Id: t.GetId()}
	case *TmplPart:
		return &TmplPart{Id: f.Tmpl.Id}
	case *ImplPart:
		return &ImplPart{Id: t.GetId()}
	default:
		return part
	}
}
