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

package ui

import (
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/log"
	editor "github.com/go-curses/coreutils-etc-hosts-editor"
	"github.com/go-curses/ctk"
	"github.com/go-curses/ctk/lib/enums"
)

const gSidebarAddRowHandler = "editor-add-row-handler"

func (c *CUI) activateSidebarAddRowHandler(data []interface{}, argv ...interface{}) cenums.EventFlag {
	idx := c.HostFile.Len()

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
			c.SidebarAddEntryButton.LogDebug("add comment at index: %v", idx)
			h := editor.NewComment("")
			c.HostFile.InsertHost(h, idx)
			c.requestReloadContents()
			c.focusEditor(h)
		case 2: // add host
			c.SidebarAddEntryButton.LogDebug("add host at index: %v", idx)
			h := editor.NewHostFromInfo(editor.HostInfo{})
			c.HostFile.InsertHost(h, idx)
			c.requestReloadContents()
			c.focusEditor(h)
		default:
			if responseId := int(response); responseId < 0 {
				c.SidebarAddEntryButton.LogDebug("new entry action cancelled")
			} else {
				log.ErrorF("unhandled dialog response: %v", response)
			}
		}
	})

	return cenums.EVENT_STOP
}

const gSidebarMoveRowUpHandler = "editor-move-row-up-handler"

func (c *CUI) activateSidebarMoveRowUpHandler(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if c.SelectedHost == nil {
		return cenums.EVENT_STOP
	}
	thisIdx := c.HostFile.IndexOf(c.SelectedHost)
	nextIdx := thisIdx - 1
	if nextIdx < 0 {
		return cenums.EVENT_STOP
	}
	c.HostFile.MoveHost(thisIdx, nextIdx)
	c.reloadEditor()
	c.focusEditor(c.SelectedHost)
	return cenums.EVENT_STOP
}

const gSidebarMoveRowDownHandler = "editor-move-row-down-handler"

func (c *CUI) activateSidebarMoveRowDownHandler(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if c.SelectedHost == nil {
		return cenums.EVENT_STOP
	}
	lastIdx := c.HostFile.Len() - 1
	thisIdx := c.HostFile.IndexOf(c.SelectedHost)
	nextIdx := thisIdx + 1
	if nextIdx > lastIdx {
		return cenums.EVENT_STOP
	}
	c.HostFile.MoveHost(thisIdx, nextIdx)
	c.reloadEditor()
	c.focusEditor(c.SelectedHost)
	return cenums.EVENT_STOP
}

func (c *CUI) updateSidebarActionButtons() {
	if c.SidebarMode != ListByEntry {
		c.SidebarMoveEntryUpButton.Hide()
		c.SidebarMoveEntryDownButton.Hide()
		return
	}
	c.SidebarMoveEntryUpButton.Show()
	c.SidebarMoveEntryDownButton.Show()
	var thisIdx int
	if c.SelectedHost == nil {
		thisIdx = -1
	} else {
		thisIdx = c.HostFile.IndexOf(c.SelectedHost)
	}
	lastIdx := c.HostFile.Len() - 1
	switch {
	case thisIdx == 0:
		c.SidebarMoveEntryUpButton.SetSensitive(false)
		c.SidebarMoveEntryDownButton.SetSensitive(true)
	case thisIdx == lastIdx:
		c.SidebarMoveEntryUpButton.SetSensitive(true)
		c.SidebarMoveEntryDownButton.SetSensitive(false)
	case thisIdx > 0 && thisIdx < lastIdx:
		c.SidebarMoveEntryUpButton.SetSensitive(true)
		c.SidebarMoveEntryDownButton.SetSensitive(true)
	default:
		c.SidebarMoveEntryUpButton.SetSensitive(false)
		c.SidebarMoveEntryDownButton.SetSensitive(false)
	}
}
