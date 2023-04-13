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
	"fmt"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paths"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
	enums2 "github.com/go-curses/ctk/lib/enums"

	"github.com/go-curses/coreutils-etc-hosts-editor"
)

func (e *CEheditor) startup(data []interface{}, argv ...interface{}) enums.EventFlag {
	var err error
	var ok bool
	if e.App, e.Display, _, _, _, ok = ctk.ArgvApplicationSignalStartup(argv...); ok {

		if e.Display.App().GetContext().NArg() > 0 {
			e.SourceFile = e.Display.App().GetContext().Args().First()
		}
		if e.SourceFile == "" {
			e.SourceFile = "/etc/hosts"
		}

		if !paths.IsFile(e.SourceFile) {
			e.LastError = fmt.Errorf("%v not found or not a file", e.SourceFile)
			log.Error(e.LastError)
			return enums.EVENT_STOP
		}
		e.ReadOnlyMode = e.Display.App().GetContext().Bool("read-only")
		if !e.ReadOnlyMode && !paths.FileWritable(e.SourceFile) {
			log.WarnF("etc hosts file %v is not writable, read-only mode", e.SourceFile)
			e.ReadOnlyMode = true
		}

		// e.Display.CaptureCtrlC()
		if s := e.Display.Screen(); s != nil {
			s.EnableHostClipboard(true)
		}

		screenSize := ptypes.MakeRectangle(e.Display.Screen().Size())
		if screenSize.W < 60 || screenSize.H < 14 {
			e.LastError = fmt.Errorf("eheditor requires a terminal with at least 80x24 dimensions")
			log.Error(e.LastError)
			return enums.EVENT_STOP
		}

		title := fmt.Sprintf("%s - eheditor %v", e.SourceFile, e.App.Version())
		if e.ReadOnlyMode {
			title += " [read-only]"
		}

		if e.HostFile, err = editor.ParseFile(e.SourceFile); err != nil {
			e.LastError = fmt.Errorf("error parsing %v: %v", e.SourceFile, err)
			log.Error(e.LastError)
			return enums.EVENT_STOP
		}

		ctk.GetAccelMap().LoadFromString(eheditorAccelMap)

		e.Window = ctk.NewWindowWithTitle(title)
		e.Window.SetName("eheditor-window")
		// e.Window.Show()
		e.Window.SetTheme(WindowTheme)
		// e.Window.SetDecorated(false)
		// _ = e.Window.SetBoolProperty(ctk.PropertyDebug, true)
		if err := e.Window.ImportStylesFromString(eheditorStyles); err != nil {
			e.Window.LogErr(err)
		}

		ag := ctk.NewAccelGroup()
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
		e.Window.AddAccelGroup(ag)

		vbox := e.Window.GetVBox()
		vbox.SetSpacing(1)

		e.HostsHBox = ctk.NewHBox(false, 0)
		e.HostsHBox.Show()
		// _ = e.HostsHBox.SetBoolProperty(cdk.PropertyDebug, true)
		// _ = e.HostsHBox.SetBoolProperty(ctk.PropertyDebugChildren, true)
		vbox.PackStart(e.HostsHBox, true, true, 0)

		e.LeftSep = ctk.NewSeparator()
		e.LeftSep.Show()
		e.HostsHBox.PackStart(e.LeftSep, false, true, 0)

		e.HostsViewport = ctk.NewScrolledViewport()
		e.HostsViewport.SetTheme(ViewerTheme)
		e.HostsViewport.Show()
		e.HostsViewport.SetPolicy(enums2.PolicyAutomatic, enums2.PolicyAutomatic)
		e.HostsVBox = ctk.NewVBox(false, 1)
		e.HostsVBox.Show()
		// _ = e.HostsVBox.SetBoolProperty(ctk.PropertyDebug, true)
		e.HostsViewport.Add(e.HostsVBox)
		// e.HostsVBox.SetBoolProperty(cdk.PropertyDebug, true)
		// e.HostsViewport.SetBoolProperty(cdk.PropertyDebug, true)
		e.HostsHBox.PackStart(e.HostsViewport, true, true, 0)

		e.RightSep = ctk.NewSeparator()
		e.RightSep.Show()
		e.HostsHBox.PackStart(e.RightSep, false, false, 0)

		e.ActionHBox = ctk.NewHBox(false, 1)
		e.ActionHBox.Show()
		vbox.PackEnd(e.ActionHBox, false, true, 0)

		actionSep := ctk.NewSeparator()
		actionSep.Show()
		e.ActionHBox.PackStart(actionSep, true, true, 0)

		e.SaveButton = ctk.NewButtonWithMnemonic("_Save <F3>")
		e.SaveButton.Show()
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
		e.ReloadButton.Connect(ctk.SignalActivate, "reload-hosts", func(data []interface{}, argv ...interface{}) enums.EventFlag {
			e.requestReload()
			return enums.EVENT_STOP
		})
		e.ActionHBox.PackEnd(e.ReloadButton, false, false, 0)

		e.QuitButton = ctk.NewButtonWithMnemonic("_Quit <F10>")
		e.QuitButton.Show()
		e.QuitButton.Connect(ctk.SignalActivate, "quit-hosts", func(data []interface{}, argv ...interface{}) enums.EventFlag {
			e.requestQuit()
			return enums.EVENT_STOP
		})
		e.ActionHBox.PackEnd(e.QuitButton, false, false, 0)

		e.App.NotifyStartupComplete()
		e.Window.Show()
		e.HostsViewport.GrabFocus()
		e.reloadViewer()
		e.Display.Connect(cdk.SignalEventResize, "display-resize-handler", func(data []interface{}, argv ...interface{}) enums.EventFlag {
			if e.App.StartupCompleted() {
				e.updateViewer()
			}
			return enums.EVENT_PASS
		})
		return enums.EVENT_PASS
	}
	return enums.EVENT_STOP
}

func (e *CEheditor) shutdown(_ []interface{}, _ ...interface{}) enums.EventFlag {
	if e.LastError != nil {
		fmt.Printf("%v\n", e.LastError)
		log.InfoF("exiting (with error)")
	} else {
		log.InfoF("exiting (without error)")
	}
	return enums.EVENT_PASS
}