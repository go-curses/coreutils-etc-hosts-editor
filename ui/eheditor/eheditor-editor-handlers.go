// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eheditor

import (
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/log"
	editor "github.com/go-curses/coreutils-etc-hosts-editor"
	"github.com/go-curses/ctk"
	"github.com/go-curses/ctk/lib/enums"
)

const gSidebarAddRowHandler = "editor-add-row-handler"

func (e *CEheditor) activateSidebarAddRowHandler(data []interface{}, argv ...interface{}) cenums.EventFlag {
	idx := e.HostFile.Len()

	var entries []interface{}
	entries = append(entries, "Comment Entry", 1)
	entries = append(entries, "Host Entry", 2)
	// add presets that don't already exist, ie localhost

	dialog := ctk.NewButtonMenuDialog(
		"Add Entry",
		"Select a type of entry to add:",
		entries...,
	)
	dialog.SetSizeRequest(32, 9)
	dialog.RunFunc(func(response enums.ResponseType, argv ...interface{}) {
		switch response {
		case 1: // add comment
			e.SidebarAddEntryButton.LogDebug("add comment at index: %v", idx)
			h := editor.NewComment("")
			e.HostFile.InsertHost(h, idx)
			e.reloadContents()
			e.focusEditor(h)
		case 2: // add host
			e.SidebarAddEntryButton.LogDebug("add host at index: %v", idx)
			h := editor.NewHostFromInfo(editor.HostInfo{})
			e.HostFile.InsertHost(h, idx)
			e.reloadContents()
			e.focusEditor(h)
		default:
			if responseId := int(response); responseId < 0 {
				e.SidebarAddEntryButton.LogDebug("new entry action cancelled")
			} else {
				log.ErrorF("unhandled dialog response: %v", response)
			}
		}
	})

	return cenums.EVENT_STOP
}

const gSidebarMoveRowUpHandler = "editor-move-row-up-handler"

func (e *CEheditor) activateSidebarMoveRowUpHandler(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if e.SelectedHost == nil {
		return cenums.EVENT_STOP
	}
	thisIdx := e.HostFile.IndexOf(e.SelectedHost)
	nextIdx := thisIdx - 1
	if nextIdx < 0 {
		return cenums.EVENT_STOP
	}
	e.HostFile.MoveHost(thisIdx, nextIdx)
	e.reloadEditor()
	e.focusEditor(e.SelectedHost)
	return cenums.EVENT_STOP
}

const gSidebarMoveRowDownHandler = "editor-move-row-down-handler"

func (e *CEheditor) activateSidebarMoveRowDownHandler(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if e.SelectedHost == nil {
		return cenums.EVENT_STOP
	}
	lastIdx := e.HostFile.Len() - 1
	thisIdx := e.HostFile.IndexOf(e.SelectedHost)
	nextIdx := thisIdx + 1
	if nextIdx > lastIdx {
		return cenums.EVENT_STOP
	}
	e.HostFile.MoveHost(thisIdx, nextIdx)
	e.reloadEditor()
	e.focusEditor(e.SelectedHost)
	return cenums.EVENT_STOP
}

func (e *CEheditor) updateSidebarActionButtons() {
	if e.SidebarMode != ListByEntry {
		e.SidebarMoveEntryUpButton.Hide()
		e.SidebarMoveEntryDownButton.Hide()
		return
	}
	e.SidebarMoveEntryUpButton.Show()
	e.SidebarMoveEntryDownButton.Show()
	var thisIdx int
	if e.SelectedHost == nil {
		thisIdx = -1
	} else {
		thisIdx = e.HostFile.IndexOf(e.SelectedHost)
	}
	lastIdx := e.HostFile.Len() - 1
	switch {
	case thisIdx == 0:
		e.SidebarMoveEntryUpButton.SetSensitive(false)
		e.SidebarMoveEntryDownButton.SetSensitive(true)
	case thisIdx == lastIdx:
		e.SidebarMoveEntryUpButton.SetSensitive(true)
		e.SidebarMoveEntryDownButton.SetSensitive(false)
	case thisIdx > 0 && thisIdx < lastIdx:
		e.SidebarMoveEntryUpButton.SetSensitive(true)
		e.SidebarMoveEntryDownButton.SetSensitive(true)
	default:
		e.SidebarMoveEntryUpButton.SetSensitive(false)
		e.SidebarMoveEntryDownButton.SetSensitive(false)
	}
}