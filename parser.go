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
	"os"
	"path/filepath"
	"strings"
)

const (
	tagPrefix  = "<!--pvn"
	tagEndTmpl = "<!--pvnTmplEnd-->"
	tagEndEdit = "<!--pvnEditEnd"
	tagEndImpl = "<!--pvnImplEnd"
	tagBegLen  = len(tagEndImpl)
	tagBegEnd  = "-->"
)

type currentState int

const (
	stateText currentState = iota
	stateTmpl
	stateTmplId
	stateEdit
	stateEditId
	stateImpl
	stateImplId
)

func ParsePart(file *File, part Part, text []byte) {
	//fmt.Printf("ParsePart\n")

	var sb bytes.Buffer
	var pos, tagLen int
	var currentPart Part
	var currentEndTag string

	textLen := len(text)

	site := file.Site

	var state currentState = stateText

	for pos < textLen {
		switch state {
		case stateText:
			state = findBegTag(text, textLen, pos)
			if state == stateText {
				sb.WriteByte(text[pos])
				pos++
			} else {
				if sb.Len() > 0 {
					site.printDebugMessage("  add text")
					part.AddPart(NewTextPart(sb.Bytes()))
				}
				sb.Reset()
				pos += tagBegLen
			}

		case stateTmplId:
			tagLen = findEndTag(text, textLen, pos, tagBegEnd)
			if tagLen > 0 {
				newTemplateFileName := strings.TrimSpace(string(sb.Bytes()))
				if file.Tmpl == nil {
					file.Tmpl = &TmplPart{Id: newTemplateFileName}
					currentPart = file.Tmpl
				} else {
					sb.Reset()
					sb.WriteString("Multiple TMPL tag in ")
					sb.WriteString(file.Path)
					sb.WriteString(" (TMPL: ")
					if len(newTemplateFileName) > 0 {
						sb.WriteString(" file=")
						sb.WriteString(newTemplateFileName)
					} else {
						sb.WriteString(" root")
					}
					sb.WriteString(")")
					panic(PvntError{sb.String()})
				}

				sb.Reset()
				pos += tagLen
				state = stateTmpl
			} else {
				sb.WriteByte(text[pos])
				pos++
			}

		case stateTmpl:
			tagLen = findEndTag(text, textLen, pos, tagEndTmpl)
			if tagLen > 0 {
				site.printDebugMessage("  add tmpl")

				ParsePart(file, currentPart, sb.Bytes())

				part.AddPart(currentPart) //?

				templateFileName := currentPart.GetId()
				if len(templateFileName) > 0 {
					// We don't check result of getParent() for nil because file.getPath() return absolute path
					templatePath := filepath.Join(filepath.Dir(file.Path), templateFileName)

					if IsFileExist(templatePath) {
						site.printDebugMessage("  add tmpl file")
						parentFile := site.AddFile(templatePath)
						parentFile.AddChild(file)
					} else {
						panic(PvntError{"Template \"" + templateFileName + "\" declared in " + file.Path + " is not exists"})
					}
				}

				currentPart = nil
				sb.Reset()
				pos += tagLen
				state = stateText
			} else {
				sb.WriteByte(text[pos])
				pos++
			}

		case stateEditId:
			tagLen = findEndTag(text, textLen, pos, tagBegEnd)
			if tagLen > 0 {
				newId := strings.TrimSpace(string(sb.Bytes()))
				currentPart = &EditPart{Id: newId}
				sb.Reset()
				sb.WriteString(tagEndEdit)
				sb.WriteString(" ")
				sb.WriteString(currentPart.GetId())
				sb.WriteString(tagBegEnd)
				currentEndTag = sb.String()
				sb.Reset()
				pos += tagLen
				state = stateEdit
			} else {
				sb.WriteByte(text[pos])
				pos++
			}

		case stateEdit:
			tagLen = findEndTag(text, textLen, pos, currentEndTag)
			if tagLen > 0 {
				site.printDebugMessage("  add edit")

				ParsePart(file, currentPart, sb.Bytes())

				part.AddPart(currentPart)

				currentPart = nil
				sb.Reset()
				pos += tagLen
				state = stateText
			} else {
				sb.WriteByte(text[pos])
				pos++
			}

		case stateImplId:
			tagLen = findEndTag(text, textLen, pos, tagBegEnd)
			if tagLen > 0 {
				newId := strings.TrimSpace(string(sb.Bytes()))
				currentPart = &ImplPart{Id: newId}
				sb.Reset()
				sb.WriteString(tagEndImpl)
				sb.WriteString(" ")
				sb.WriteString(currentPart.GetId())
				sb.WriteString(tagBegEnd)
				currentEndTag = sb.String()
				sb.Reset()
				pos += tagLen
				state = stateImpl
			} else {
				sb.WriteByte(text[pos])
				pos++
			}

		case stateImpl:
			tagLen = findEndTag(text, textLen, pos, currentEndTag)
			if tagLen > 0 {
				site.printDebugMessage("  add impl")

				ParsePart(file, currentPart, sb.Bytes())

				part.AddPart(currentPart)
				file.ImplMap[currentPart.GetId()] = currentPart.(*ImplPart)

				currentPart = nil
				sb.Reset()
				pos += tagLen
				state = stateText
			} else {
				sb.WriteByte(text[pos])
				pos++
			}
		}

	}

	if state == stateText {
		if sb.Len() > 0 {
			site.printDebugMessage("  add text")
			part.AddPart(NewTextPart(sb.Bytes()))
		}
	} else {
		sb.Reset()
		sb.WriteString("End tag not found in ")
		sb.WriteString(file.Path)
		sb.WriteString(" (")
		sb.WriteString(part.ToString())
		if currentPart != nil {
			sb.WriteString(" -> ")
			sb.WriteString(currentPart.ToString())
		}
		sb.WriteString(")")
		panic(PvntError{sb.String()})
	}
}

//==============================================================
// Private methods.
//==============================================================

func findBegTag(t []byte, l int, p int) currentState {
	if p+tagBegLen < l {
		for i := 0; i < len(tagPrefix); i++ {
			if t[p+i] != tagPrefix[i] {
				return stateText
			}
		}

		s := stateText
		j := p + len(tagPrefix)

		if t[j] == 'T' && t[j+1] == 'm' && t[j+2] == 'p' && t[j+3] == 'l' {
			s = stateTmplId
		} else if t[j] == 'I' && t[j+1] == 'm' && t[j+2] == 'p' && t[j+3] == 'l' {
			s = stateImplId
		} else if t[j] == 'E' && t[j+1] == 'd' && t[j+2] == 'i' && t[j+3] == 't' {
			s = stateEditId
		} else {
			return stateText
		}

		if t[j+4] == 'B' && t[j+5] == 'e' && t[j+6] == 'g' {
			return s
		} else {
			return stateText
		}
	}
	return stateText
}

func findEndTag(t []byte, l int, p int, tag string) int {
	len := len(tag)
	if p+len <= l {
		for i := 0; i < len; i++ {
			if t[p+i] != tag[i] {
				return 0
			}
		}
		return len
	} else {
		return 0
	}
}

func IsFileExist(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
