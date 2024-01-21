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

	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paths"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"

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

		e.Window.AddAccelGroup(e.makeAccelmap())

		vbox := e.Window.GetVBox()
		vbox.SetSpacing(0)

		e.ContentsHBox = ctk.NewHBox(false, 0)
		e.ContentsHBox.Show()
		vbox.PackStart(e.ContentsHBox, true, true, 0)

		e.ContentsHBox.PackStart(e.makeEditor(), true, true, 0)

		vbox.PackEnd(e.makeActionButtonBox(), false, true, 0)

		e.switchToEditor()

		e.App.NotifyStartupComplete()
		e.Window.Show()

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
