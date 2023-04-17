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

	"github.com/go-curses/cdk/log"

	"github.com/go-curses/coreutils-etc-hosts-editor"
)

func (e *CEheditor) requestReload() {
	log.DebugF("reloading from: %v", e.SourceFile)
	var err error
	if e.HostFile, err = editor.ParseFile(e.SourceFile); err != nil {
		e.LastError = fmt.Errorf("error parsing %v: %v", e.SourceFile, err)
		log.Error(e.LastError)
		return
	}
	e.reloadViewer()
	e.reloadEditor()
	e.focusEditor(nil)
}

func (e *CEheditor) requestSave() {
	log.DebugF("saving to: %v", e.SourceFile)
	if e.HostFile != nil {
		if err := e.HostFile.Save(); err != nil {
			log.Error(err)
		}
	}
	e.requestReload()
	e.QuitButton.GrabFocus()
}

func (e *CEheditor) requestQuit() {
	e.Display.RequestQuit()
}