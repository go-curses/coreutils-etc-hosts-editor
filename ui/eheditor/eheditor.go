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
	_ "embed"

	"github.com/go-curses/cdk"
	"github.com/go-curses/ctk"

	editor "github.com/go-curses/coreutils-etc-hosts-editor"
)

//go:embed eheditor.accelmap
var eheditorAccelMap string

//go:embed eheditor.styles
var eheditorStyles string

type CEheditor struct {
	App ctk.Application

	HostFile     *editor.Hostfile
	SourceFile   string
	LastError    error
	ReadOnlyMode bool

	LeftSep       ctk.Separator
	RightSep      ctk.Separator
	HostsViewport ctk.ScrolledViewport
	HostsVBox     ctk.VBox
	HostsHBox     ctk.HBox
	ActionHBox    ctk.HBox
	Display       cdk.Display
	Window        ctk.Window
	SaveButton    ctk.Button
	ReloadButton  ctk.Button
	QuitButton    ctk.Button

	ViewerDomainLookup map[string]*editor.Host
}

func NewEheditor(name string, usage string, description string, version string, tag string, title string, ttyPath string) (e *CEheditor) {
	e = &CEheditor{
		App: ctk.NewApplication(name, usage, description, version, tag, title, ttyPath),
	}
	e.App.Connect(cdk.SignalStartup, "eheditor-startup-handler", e.startup)
	// e.App.Connect(cdk.SignalStartupComplete, "eheditor-startup-complete-handler", startupComplete)
	e.App.Connect(cdk.SignalShutdown, "eheditor-quit-handler", e.shutdown)
	return
}

func (e *CEheditor) Run(argv []string) (err error) {
	err = e.App.Run(argv)
	return
}