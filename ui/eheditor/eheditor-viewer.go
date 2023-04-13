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
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"

	"github.com/go-curses/coreutils-etc-hosts-editor"
)

func (e *CEheditor) reloadViewer() {
	e.Window.Freeze()

	e.HostsVBox.Freeze()
	existing := e.HostsVBox.GetChildren()
	for _, child := range existing {
		e.HostsVBox.Remove(child)
		child.Destroy()
	}
	e.HostsVBox.Thaw()

	e.ViewerDomainLookup = make(map[string]*editor.Host)
	var domains []string
	for _, host := range e.HostFile.Hosts() {
		for _, domain := range host.Domains() {
			if !strings.StringSliceHasValue(domains, domain) {
				domains = append(domains, domain)
				e.ViewerDomainLookup[domain] = host
			}
		}
	}

	e.updateViewer()
	e.Window.Thaw()
	e.Window.Invalidate()
	e.Display.RequestDraw()
	e.Display.RequestSync()
}

func (e *CEheditor) updateViewer() {
	var screenSize ptypes.Rectangle
	if screen := e.Display.Screen(); screen != nil {
		screenSize = ptypes.MakeRectangle(screen.Size())
	} else {
		log.ErrorF("missing screen!")
		return
	}

	e.Window.Freeze()
	e.HostsHBox.Freeze()
	e.HostsVBox.Freeze()

	if screenSize.H < 30 {
		e.SaveButton.SetSizeRequest(-1, 1)
		e.ReloadButton.SetSizeRequest(-1, 1)
		e.QuitButton.SetSizeRequest(-1, 1)
	} else {
		e.SaveButton.SetSizeRequest(-1, 3)
		e.ReloadButton.SetSizeRequest(-1, 3)
		e.QuitButton.SetSizeRequest(-1, 3)
	}

	screenSize.Sub(2, 2)

	sepWidth := screenSize.W / 7
	if screenSize.W <= 100 {
		sepWidth = 1
	}

	e.LeftSep.SetSizeRequest(sepWidth, -1)
	e.RightSep.SetSizeRequest(sepWidth, -1)

	viewerWidth := screenSize.W - (sepWidth * 2)
	viewerWidth -= 2

	existing := e.HostsVBox.GetChildren()
	totalExisting := len(existing)
	totalSize := ptypes.MakeRectangle(0, 0)

	for idx, host := range e.HostFile.Hosts() {
		var child ctk.Widget
		if idx < totalExisting {
			child = existing[idx]
			if frame, ok := child.Self().(ctk.Frame); ok {
				if row := getViewerRowFromFrame(frame); row != nil {
					row.Update(host, viewerWidth)
					rw, rh := row.Frame.GetSizeRequest()
					if totalSize.W < rw {
						totalSize.W = rw
					}
					if idx > 0 {
						totalSize.H += 1
					}
					totalSize.H += rh
					log.DebugF("updated %v with %v", frame.ObjectInfo(), host.Name())
				} else {
					log.ErrorF("row not found for child: %v", child.ObjectInfo())
				}
			}
		} else {
			row := NewViewerRow(e, host, viewerWidth)
			row.Update(host, viewerWidth)
			e.HostsVBox.PackStart(row.Frame, true, true, 0)

			rw, rh := row.Frame.GetSizeRequest()
			if totalSize.W < rw {
				totalSize.W = rw
			}
			if idx > 0 {
				totalSize.H += 1
			}
			totalSize.H += rh
			log.DebugF("created %v with %v", row.Frame.ObjectInfo(), host.Name())
		}
	}

	e.HostsVBox.SetSizeRequest(totalSize.W, totalSize.H)
	e.HostsVBox.Thaw()
	e.HostsHBox.Thaw()
	e.HostsHBox.Resize()
	e.Window.ReApplyStyles()
	e.Window.Thaw()
	e.Window.Invalidate()
	e.Display.RequestDraw()
	e.Display.RequestShow()
}