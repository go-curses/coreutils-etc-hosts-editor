package main

import (
	"github.com/go-curses/cdk/lib/ptypes"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/log"
	editor "github.com/go-curses/coreutils-etc-hosts-editor"
	"github.com/go-curses/ctk"
)

var (
	gViewerDomainLookup map[string]*editor.Host = nil
)

func reloadViewer() {
	gWindow.Freeze()

	gHostsVBox.Freeze()
	existing := gHostsVBox.GetChildren()
	for _, child := range existing {
		gHostsVBox.Remove(child)
		child.Destroy()
	}
	gHostsVBox.Thaw()

	gViewerDomainLookup = make(map[string]*editor.Host)
	var domains []string
	for _, host := range gEH.Hosts() {
		for _, domain := range host.Domains() {
			if !cstrings.StringSliceHasValue(domains, domain) {
				domains = append(domains, domain)
				gViewerDomainLookup[domain] = host
			}
		}
	}

	updateViewer()
	gWindow.Thaw()
	gWindow.Invalidate()
	gDisplay.RequestDraw()
	gDisplay.RequestSync()
}

func updateViewer() {
	var screenSize ptypes.Rectangle
	if screen := gDisplay.Screen(); screen != nil {
		screenSize = ptypes.MakeRectangle(screen.Size())
	} else {
		log.ErrorF("missing screen!")
		return
	}

	gWindow.Freeze()
	gHostsHBox.Freeze()
	gHostsVBox.Freeze()

	if screenSize.H < 30 {
		gSaveButton.SetSizeRequest(-1, 1)
		gReloadButton.SetSizeRequest(-1, 1)
		gQuitButton.SetSizeRequest(-1, 1)
	} else {
		gSaveButton.SetSizeRequest(-1, 3)
		gReloadButton.SetSizeRequest(-1, 3)
		gQuitButton.SetSizeRequest(-1, 3)
	}

	screenSize.Sub(2, 2)

	sepWidth := screenSize.W / 7
	if screenSize.W <= 100 {
		sepWidth = 1
	}

	gLeftSep.SetSizeRequest(sepWidth, -1)
	gRightSep.SetSizeRequest(sepWidth, -1)

	viewerWidth := screenSize.W - (sepWidth * 2)
	viewerWidth -= 2

	existing := gHostsVBox.GetChildren()
	totalExisting := len(existing)
	totalSize := ptypes.MakeRectangle(0, 0)

	for idx, host := range gEH.Hosts() {
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
			row := NewViewerRow(host, viewerWidth)
			row.Update(host, viewerWidth)
			gHostsVBox.PackStart(row.Frame, true, true, 0)

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

	gHostsVBox.SetSizeRequest(totalSize.W, totalSize.H)
	gHostsVBox.Thaw()
	gHostsHBox.Thaw()
	gHostsHBox.Resize()
	gWindow.ReApplyStyles()
	gWindow.Thaw()
	gWindow.Invalidate()
	gDisplay.RequestDraw()
	gDisplay.RequestShow()
}