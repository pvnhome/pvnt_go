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
	"flag"
	"fmt"
	"os"
)

const PVNT_VERSION = "0.9.4 beta"

const (
	PARTS_INITIAL_CAPACITY = 2
	FILES_INITIAL_CAPACITY = 2
)

var (
	g_path    = flag.String("path", ".", "path to site")
	g_ext     = flag.String("ext", "html", "file extension")
	g_debug   = flag.Bool("d", false, "enable debug mode")
	g_version = flag.Bool("v", false, "print the version information and exit")
)

func show_usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-v] [-d] [-path=<path/to/site>] [-ext=<file extension>]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	fmt.Println("PVN Template engine (version", PVNT_VERSION, ")")

	flag.Usage = show_usage
	flag.Parse()

	if !*g_version {
		runPvnt(*g_path, *g_ext, *g_debug)
	}
}

func runPvnt(path string, ext string, debugMode bool) {
	// runPvnt will panic if there is any unrecoverable error.
	defer func() {
		if e := recover(); e != nil {
			err := e.(PvntError) // Will re-panic if not a PvntError.
			fmt.Println("ERROR: ", err.Description)
		}
	}()

	site := Site{Path: path, Ext: ext, DebugMode: debugMode, Files: make(map[string]*File)}

	fmt.Printf("INFO: Process all %s files in %s (debug=%t)\n", site.Ext, site.Path, site.DebugMode)

	site.Load()

	fmt.Println("INFO: files loaded")

	if site.DebugMode {
		site.PrintPartsTree()
	}

	site.Process()

	fmt.Println("INFO: files processed")

	site.Save()

	fmt.Println("INFO: files saved")
}
