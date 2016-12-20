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
	"io/ioutil"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type Site struct {
	Path      string // Path to site
	Ext       string // Desired file extension
	DebugMode bool   // Debug mode flag

	Files map[string]*File
}

//===========================================
// Main methods.
//===========================================

func (site *Site) Load() {
	files, err := ioutil.ReadDir(site.Path)
	if err != nil {
		fmt.Println("ERROR: ", err)
		panic(PvntError{"Can't read directory: " + site.Path})
	}

	ext := "." + site.Ext

	for _, file := range files {
		if !file.IsDir() && path.Ext(file.Name()) == ext {
			site.AddFile(path.Join(site.Path, file.Name()))
		}
	}
}

func (site *Site) Process() {
	fmt.Println("process")
}

func (site *Site) Save() {
	fmt.Println("save")
	//panic(PvntError{"Can't save file (panic)"})
}

func (site *Site) AddFile(path string) *File {
	var file *File
	var ok bool

	fullName, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("ERROR: file =", path, ",", err)
	} else {
		file, ok = site.Files[fullName]
		if ok {
			site.printDebugMessage("file from cache: %s", fullName)
		} else {
			file = NewFile(site, fullName)

			if text, err := file.Load(); err == nil {
				site.printDebugMessage("add file:        %s", fullName)
				ParsePart(file, file, text)
				site.Files[fullName] = file
			} else {
				fmt.Println("ERROR: add file=", file.Path, ", error=", err)
			}
		}
	}

	return file
}

func (site *Site) PrintPartsTree() {
	fmt.Println("PARTS TREE:")

	// Sort files
	fileNames := make([]string, 0, len(site.Files))
	for name := range site.Files {
		fileNames = append(fileNames, name)
	}
	sort.Strings(fileNames)

	// Print tree of parts
	for _, name := range fileNames {
		file := site.Files[name]
		fmt.Println("  ", name)

		implIds := file.GetImplIds()

		if len(implIds) > 0 {
			fmt.Println("      IMPL IDS:", implIds)
		}
		site.printFilePartsTree(file)
	}
}

func (site *Site) printDebugMessage(format string, a ...interface{}) {
	if site.DebugMode {
		fmt.Println(fmt.Sprintf(format, a...))
	}
}

//===========================================
// Private methods.
//===========================================

func (site *Site) printFilePartsTree(file *File) {
	for _, part := range file.GetChildren() {
		site.printFilePartTree(part, 6)
	}
}

func (site *Site) printFilePartTree(part Part, depth int) {
	ident := strings.Repeat(" ", depth)

	if len(part.GetChildren()) > 0 {
		fmt.Printf("%s%s {\n", ident, part.ToString())
		for _, child := range part.GetChildren() {
			site.printFilePartTree(child, depth+3)
		}
		fmt.Printf("%s}\n", ident)
	} else {
		fmt.Printf("%s%s\n", ident, part.ToString())
	}
}
