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
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk"
)

func (c *CUI) makeAccelmap() (ag ctk.AccelGroup) {
	ag = ctk.NewAccelGroup()
	ag.ConnectByPath(
		"<eheditor-window>/File/Quit",
		"quit-accel",
		func(argv ...interface{}) (handled bool) {
			ag.LogDebug("quit-accel called")
			c.requestQuit()
			return
		},
	)
	ag.ConnectByPath(
		"<eheditor-window>/File/Reload",
		"reload-accel",
		func(argv ...interface{}) (handled bool) {
			ag.LogDebug("reload-accel called")
			c.requestReload()
			return
		},
	)
	ag.ConnectByPath(
		"<eheditor-window>/File/Save",
		"save-accel",
		func(argv ...interface{}) (handled bool) {
			ag.LogDebug("save-accel called")
			c.requestSave()
			return
		},
	)
	return
}

func (c *CUI) makeActionButtonBox() ctk.HButtonBox {
	c.ActionHBox = ctk.NewHButtonBox(false, 1)
	c.ActionHBox.Show()
	c.ActionHBox.SetSizeRequest(-1, 1)

	actionSep := ctk.NewSeparator()
	actionSep.Show()
	c.ActionHBox.PackEnd(actionSep, true, true, 0)

	c.SaveButton = ctk.NewButtonWithMnemonic("_Save <F3>")
	c.SaveButton.Show()
	c.SaveButton.SetSizeRequest(-1, 1)
	if c.ReadOnlyMode {
		c.SaveButton.SetSensitive(false)
	} else {
		c.SaveButton.Connect(ctk.SignalActivate, "save-hosts", func(data []interface{}, argv ...interface{}) enums.EventFlag {
			c.requestSave()
			return enums.EVENT_STOP
		})
	}
	c.ActionHBox.PackEnd(c.SaveButton, false, false, 0)

	c.ReloadButton = ctk.NewButtonWithMnemonic("_Reload <F5>")
	c.ReloadButton.Show()
	c.ReloadButton.SetSizeRequest(-1, 1)
	c.ReloadButton.Connect(ctk.SignalActivate, "reload-hosts", func(data []interface{}, argv ...interface{}) enums.EventFlag {
		c.requestReload()
		return enums.EVENT_STOP
	})
	c.ActionHBox.PackEnd(c.ReloadButton, false, false, 0)

	c.QuitButton = ctk.NewButtonWithMnemonic("_Quit <F10>")
	c.QuitButton.Show()
	c.QuitButton.SetSizeRequest(-1, 1)
	c.QuitButton.Connect(ctk.SignalActivate, "quit-hosts", func(data []interface{}, argv ...interface{}) enums.EventFlag {
		c.requestQuit()
		return enums.EVENT_STOP
	})
	c.ActionHBox.PackEnd(c.QuitButton, false, false, 0)

	return c.ActionHBox
}
