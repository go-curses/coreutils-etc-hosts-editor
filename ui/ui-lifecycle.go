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
	"fmt"

	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paths"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"

	editor "github.com/go-curses/coreutils-etc-hosts-editor"
)

func (c *CUI) startup(_ []interface{}, argv ...interface{}) enums.EventFlag {
	var err error
	var ok bool
	if c.App, c.Display, _, _, _, ok = ctk.ArgvApplicationSignalStartup(argv...); ok {

		if c.Display.App().GetContext().NArg() > 0 {
			c.SourceFile = c.Display.App().GetContext().Args().First()
		}
		if c.SourceFile == "" {
			c.SourceFile = "/etc/hosts"
		}

		if !paths.IsFile(c.SourceFile) {
			c.LastError = fmt.Errorf("%v not found or not a file", c.SourceFile)
			log.Error(c.LastError)
			return enums.EVENT_STOP
		}
		c.ReadOnlyMode = c.Display.App().GetContext().Bool("read-only")
		if !c.ReadOnlyMode && !paths.FileWritable(c.SourceFile) {
			log.WarnF("etc hosts file %v is not writable, read-only mode", c.SourceFile)
			c.ReadOnlyMode = true
		}

		if s := c.Display.Screen(); s != nil {
			s.EnableHostClipboard(true)
		}

		screenSize := ptypes.MakeRectangle(c.Display.Screen().Size())
		if screenSize.W < 60 || screenSize.H < 14 {
			c.LastError = fmt.Errorf("eheditor requires a terminal with at least 80x24 dimensions")
			log.Error(c.LastError)
			return enums.EVENT_STOP
		}

		title := fmt.Sprintf("%s - eheditor %v", c.SourceFile, c.App.Version())
		if c.ReadOnlyMode {
			title += " [read-only]"
		}

		if c.HostFile, err = editor.ParseFile(c.SourceFile); err != nil {
			c.LastError = fmt.Errorf("error parsing %v: %v", c.SourceFile, err)
			log.Error(c.LastError)
			return enums.EVENT_STOP
		}

		ctk.GetAccelMap().LoadFromString(eheditorAccelMap)

		c.Window = ctk.NewWindowWithTitle(title)
		c.Window.SetName("eheditor-window")
		c.Window.SetTheme(WindowTheme)
		// c.Window.SetDecorated(false)
		// _ = c.Window.SetBoolProperty(ctk.PropertyDebug, true)
		if err := c.Window.ImportStylesFromString(eheditorStyles); err != nil {
			c.Window.LogErr(err)
		}

		c.Window.AddAccelGroup(c.makeAccelmap())

		vbox := c.Window.GetVBox()
		vbox.SetSpacing(0)

		c.ContentsHBox = ctk.NewHBox(false, 0)
		c.ContentsHBox.Show()
		vbox.PackStart(c.ContentsHBox, true, true, 0)

		c.ContentsHBox.PackStart(c.makeEditor(), true, true, 0)

		vbox.PackEnd(c.makeActionButtonBox(), false, true, 0)

		c.switchToEditor()

		c.App.NotifyStartupComplete()
		c.Window.Show()

		return enums.EVENT_PASS
	}
	return enums.EVENT_STOP
}

func (c *CUI) shutdown(_ []interface{}, _ ...interface{}) enums.EventFlag {
	if c.LastError != nil {
		fmt.Printf("%v\n", c.LastError)
		log.InfoF("exiting (with error)")
	} else {
		log.InfoF("exiting (without error)")
	}
	return enums.EVENT_PASS
}
