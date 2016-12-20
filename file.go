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
	"sort"
	"strings"
)

type File struct {
	Site *Site  // Site
	Path string // Path to file

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
	//fmt.Printf("File.AddPart: len=%d cap=%d (%s)\n", len(f.Parts), cap(f.Parts), f.Path)
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
	return ioutil.ReadFile(f.Path)
}

func (f *File) Save() {
	fmt.Println("save file: ", f.Path)

	//content := []byte("temporary file's content")
	/*
		dir, err := ioutil.TempDir("", "example")
		if err != nil {
			log.Fatal(err)
		}

		defer os.RemoveAll(dir) // clean up

		tmpfn := filepath.Join(dir, "tmpfile")
		if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
			log.Fatal(err)
		}
	*/
	/*
	   try (BufferedWriter writer = Files.newBufferedWriter(path, StandardCharsets.UTF_8)) {
	      for (Part part : getChildren()) {
	         part.write(writer);
	      }
	   }
	*/
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
	//fmt.Printf("File.AddChild: len=%d cap=%d (%s)\n", len(f.Children), cap(f.Children), f.Path)
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
	/*
	   for (Part part : templateParts) {
	      if (part instanceof TextPart) {
	         currentPart.addPart(part);
	      } else {
	         if (part instanceof EditPart && implMap.containsKey(part.getId())) {
	            ImplPart implPart = implMap.get(part.getId());
	            implPart.setProcessed(true);
	            currentPart.addPart(implPart);
	         } else {
	            if (part instanceof ImplPart) {
	               ImplPart implPart = implMap.get(part.getId());
	               if (implPart != null) {
	                  implPart.setProcessed(true);
	                  // TODO We need to compare part with implPart and we need to print message if the difference will be detected
	               }
	            }

	            Part clonedPart = clone(part);
	            currentPart.addPart(clonedPart);
	            if (!part.getChildren().isEmpty()) {
	               process(part.getChildren(), clonedPart);
	            }
	         }
	      }
	   }
	*/
}

func clonePart(part Part) Part {
	/*
		switch t := part.(type) {
		case *TmplPart:
			fmt.Printf("type TmplPart\n")
		case *File:
			fmt.Printf("type File\n")
		default:
			fmt.Printf("unexpected type %T\n", t)
		}
	*/
	/*
		      if (part instanceof EditPart) {
		         newPart = new EditPart(part);
		      } else if (part instanceof TmplPart) {
		         newPart = new TmplPart(getTmpl().getId());
		      } else if (part instanceof ImplPart) {
		         newPart = new ImplPart(part);
		      } else {
				return part
			}
	*/
	return nil
}
