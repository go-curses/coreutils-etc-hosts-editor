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
	_ "embed"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/sync"
	"github.com/go-curses/ctk"

	editor "github.com/go-curses/coreutils-etc-hosts-editor"
)

//go:embed eheditor.accelmap
var eheditorAccelMap string

//go:embed eheditor.styles
var eheditorStyles string

type SidebarListMode uint8

const (
	ListByDomain SidebarListMode = iota
	ListByAddress
	ListByEntry
)

type CUI struct {
	App ctk.Application

	HostFile     *editor.Hostfile
	SourceFile   string
	LastError    error
	ReadOnlyMode bool

	ContentsHBox ctk.HBox
	ActionHBox   ctk.HButtonBox
	Display      cdk.Display
	Window       ctk.Window
	SaveButton   ctk.Button
	ReloadButton ctk.Button
	QuitButton   ctk.Button

	EditingHBox ctk.HBox

	ByDomainsButton ctk.Button
	ByAddressButton ctk.Button
	ByEntryButton   ctk.Button

	SidebarFrame        ctk.Frame
	SidebarEntryList    ctk.VBox
	SidebarLocalsList   ctk.VBox
	SidebarCustomList   ctk.VBox
	SidebarCommentsList ctk.VBox

	SidebarAddEntryButton      ctk.Button
	SidebarMoveEntryUpButton   ctk.Button
	SidebarMoveEntryDownButton ctk.Button

	CommentsEntry  ctk.Entry
	HostEditVBox   ctk.VBox
	AddressEntry   ctk.Entry
	AddressButton  ctk.Button
	DomainsEntry   ctk.Entry
	ActivateButton ctk.Button
	DeleteButton   ctk.Button

	HostSelectedFrame    ctk.Frame
	NothingSelectedFrame ctk.Frame
	CommentSelectedFrame ctk.Frame

	SidebarMode  SidebarListMode
	SelectedHost *editor.Host

	EditorCommentList   []*editor.Host
	EditorAddressLookup map[string]*editor.Host
	EditorDomainsLookup map[string]*editor.Host

	sync.RWMutex
}

func NewUI(name string, usage string, description string, version string, tag string, title string, ttyPath string) (e *CUI) {
	e = &CUI{
		App: ctk.NewApplication(name, usage, description, version, tag, title, ttyPath),
	}
	e.App.Connect(cdk.SignalStartup, "eheditor-startup-handler", e.startup)
	// e.App.Connect(cdk.SignalStartupComplete, "eheditor-startup-complete-handler", startupComplete)
	e.App.Connect(cdk.SignalShutdown, "eheditor-quit-handler", e.shutdown)
	return
}

func (c *CUI) Run(argv []string) (err error) {
	err = c.App.Run(argv)
	return
}
