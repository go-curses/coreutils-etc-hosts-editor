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
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk"
)

func (e *CEheditor) makeAccelmap() (ag ctk.AccelGroup) {
	ag = ctk.NewAccelGroup()
	ag.ConnectByPath(
		"<eheditor-window>/File/Quit",
		"quit-accel",
		func(argv ...interface{}) (handled bool) {
			ag.LogDebug("quit-accel called")
			e.requestQuit()
			return
		},
	)
	ag.ConnectByPath(
		"<eheditor-window>/File/Reload",
		"reload-accel",
		func(argv ...interface{}) (handled bool) {
			ag.LogDebug("reload-accel called")
			e.requestReload()
			return
		},
	)
	ag.ConnectByPath(
		"<eheditor-window>/File/Save",
		"save-accel",
		func(argv ...interface{}) (handled bool) {
			ag.LogDebug("save-accel called")
			e.requestSave()
			return
		},
	)
	return
}

func (e *CEheditor) makeActionButtonBox() ctk.HButtonBox {
	e.ActionHBox = ctk.NewHButtonBox(false, 1)
	e.ActionHBox.Show()
	e.ActionHBox.SetSizeRequest(-1, 1)

	actionSep := ctk.NewSeparator()
	actionSep.Show()
	e.ActionHBox.PackEnd(actionSep, true, true, 0)

	e.SaveButton = ctk.NewButtonWithMnemonic("_Save <F3>")
	e.SaveButton.Show()
	e.SaveButton.SetSizeRequest(-1, 1)
	if e.ReadOnlyMode {
		e.SaveButton.SetSensitive(false)
	} else {
		e.SaveButton.Connect(ctk.SignalActivate, "save-hosts", func(data []interface{}, argv ...interface{}) enums.EventFlag {
			e.requestSave()
			return enums.EVENT_STOP
		})
	}
	e.ActionHBox.PackEnd(e.SaveButton, false, false, 0)

	e.ReloadButton = ctk.NewButtonWithMnemonic("_Reload <F5>")
	e.ReloadButton.Show()
	e.ReloadButton.SetSizeRequest(-1, 1)
	e.ReloadButton.Connect(ctk.SignalActivate, "reload-hosts", func(data []interface{}, argv ...interface{}) enums.EventFlag {
		e.requestReload()
		return enums.EVENT_STOP
	})
	e.ActionHBox.PackEnd(e.ReloadButton, false, false, 0)

	e.QuitButton = ctk.NewButtonWithMnemonic("_Quit <F10>")
	e.QuitButton.Show()
	e.QuitButton.SetSizeRequest(-1, 1)
	e.QuitButton.Connect(ctk.SignalActivate, "quit-hosts", func(data []interface{}, argv ...interface{}) enums.EventFlag {
		e.requestQuit()
		return enums.EVENT_STOP
	})
	e.ActionHBox.PackEnd(e.QuitButton, false, false, 0)

	return e.ActionHBox
}
