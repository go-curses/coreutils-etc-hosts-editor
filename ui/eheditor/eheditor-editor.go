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
	"regexp"
	"strconv"
	"strings"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
	"github.com/go-curses/ctk/lib/enums"

	editor "github.com/go-curses/coreutils-etc-hosts-editor"
)

const (
	gSidebarOuterWidth = 21
	gSidebarInnerWidth = 16
)

func (e *CEheditor) switchToEditor() {
	e.Window.Freeze()
	e.ContentsHBox.Freeze()
	for _, child := range e.ContentsHBox.GetChildren() {
		if child.ObjectID() == e.EditingHBox.ObjectID() {
			child.Show()
		} else {
			child.Hide()
		}
	}
	e.Window.Thaw()
	e.ContentsHBox.Thaw()
	e.EditorButton.SetTheme(ActiveButtonTheme)
	e.EditorButton.GrabFocus()
	e.ViewerButton.SetTheme(DefaultButtonTheme)
	e.ContentsHBox.Thaw()
	e.focusEditor(e.SelectedHost)
}

func (e *CEheditor) makeEditor() ctk.Widget {
	e.EditingHBox = ctk.NewHBox(false, 1)
	e.EditingHBox.SetName("editing")
	e.EditingHBox.Show()

	sidebarViewCtrlHBox := ctk.NewHBox(false, 0)
	sidebarViewCtrlHBox.Show()
	sidebarViewCtrlHBox.SetSizeRequest(-1, 1)

	var sidebarEntryFrame ctk.Frame
	var sidebarLocalsFrame ctk.Frame
	var sidebarCustomFrame ctk.Frame
	var sidebarCommentsFrame ctk.Frame

	changeSidebarMode := func(mode SidebarListMode) {
		dWidth, aWidth, eWidth := 3, 3, 3
		dLabel, aLabel, eLabel := "_D", "_A", "_E"
		dTheme, aTheme, eTheme := DefaultButtonTheme, DefaultButtonTheme, DefaultButtonTheme

		switch mode {
		case ListByEntry:
			eWidth = 11
			eLabel = "_Entry"
			eTheme = ActiveButtonTheme
			e.SidebarMode = ListByEntry
			sidebarCommentsFrame.Hide()
			sidebarLocalsFrame.Hide()
			sidebarCustomFrame.Hide()
			sidebarEntryFrame.Show()
		case ListByAddress:
			aWidth = 11
			aLabel = "_Address"
			aTheme = ActiveButtonTheme
			e.SidebarMode = ListByAddress
			sidebarCommentsFrame.Show()
			sidebarLocalsFrame.Show()
			sidebarCustomFrame.Show()
			sidebarEntryFrame.Hide()
		case ListByDomain:
			dWidth = 11
			dLabel = "_Domain"
			dTheme = ActiveButtonTheme
			e.SidebarMode = ListByDomain
			sidebarCommentsFrame.Show()
			sidebarLocalsFrame.Show()
			sidebarCustomFrame.Show()
			sidebarEntryFrame.Hide()
		}

		e.ByEntryButton.SetSizeRequest(eWidth, 1)
		e.ByAddressButton.SetSizeRequest(aWidth, 1)
		e.ByDomainsButton.SetSizeRequest(dWidth, 1)

		e.ByEntryButton.SetLabel(eLabel)
		e.ByAddressButton.SetLabel(aLabel)
		e.ByDomainsButton.SetLabel(dLabel)

		e.ByEntryButton.SetTheme(eTheme)
		e.ByAddressButton.SetTheme(aTheme)
		e.ByDomainsButton.SetTheme(dTheme)

		e.updateSidebarActionButtons()
	}

	e.ByDomainsButton = ctk.NewButtonWithMnemonic("_Domain")
	e.ByDomainsButton.Show()
	e.ByDomainsButton.Connect(ctk.SignalActivate, "by-domains-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		e.ByDomainsButton.LogDebug("clicked")
		changeSidebarMode(ListByDomain)
		e.reloadEditor()
		e.focusEditor(nil)
		return cenums.EVENT_STOP
	})
	e.ByDomainsButton.SetHasTooltip(true)
	e.ByDomainsButton.SetTooltipText("Click to list by domain names")
	sidebarViewCtrlHBox.PackStart(e.ByDomainsButton, false, false, 0)

	e.ByAddressButton = ctk.NewButtonWithMnemonic("_A")
	e.ByAddressButton.Show()
	e.ByAddressButton.Connect(ctk.SignalActivate, "by-address-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		e.ByAddressButton.LogDebug("clicked")
		changeSidebarMode(ListByAddress)
		e.reloadEditor()
		e.focusEditor(nil)
		return cenums.EVENT_STOP
	})
	e.ByAddressButton.SetHasTooltip(true)
	e.ByAddressButton.SetTooltipText("Click to list by IP addresses")
	sidebarViewCtrlHBox.PackStart(e.ByAddressButton, false, false, 0)

	e.ByEntryButton = ctk.NewButtonWithMnemonic("_E")
	e.ByEntryButton.Show()
	e.ByEntryButton.Connect(ctk.SignalActivate, "by-entry-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		e.ByEntryButton.LogDebug("clicked")
		changeSidebarMode(ListByEntry)
		e.reloadEditor()
		e.focusEditor(nil)
		return cenums.EVENT_STOP
	})
	e.ByEntryButton.SetHasTooltip(true)
	e.ByEntryButton.SetTooltipText("Click to list by hosts file entry")
	sidebarViewCtrlHBox.PackStart(e.ByEntryButton, false, false, 0)

	e.SidebarFrame = ctk.NewFrameWithWidget(sidebarViewCtrlHBox)
	e.SidebarFrame.Show()
	e.SidebarFrame.SetLabelAlign(0.0, 0.5)
	e.SidebarFrame.SetSizeRequest(gSidebarOuterWidth, -1)
	e.EditingHBox.PackStart(e.SidebarFrame, false, false, 0)

	sidebarVBox := ctk.NewVBox(false, 0)
	sidebarVBox.Show()
	e.SidebarFrame.Add(sidebarVBox)

	// list toggles

	toggleLocals := ctk.NewButtonWithLabel(string(paint.RuneTriangleRight) + " locals")
	toggleLocals.Show()
	toggleLocals.SetTheme(SidebarHeaderTheme)
	toggleCustom := ctk.NewButtonWithLabel(string(paint.RuneTriangleDown) + " custom")
	toggleCustom.Show()
	toggleCustom.SetTheme(SidebarHeaderTheme)
	toggleComments := ctk.NewButtonWithLabel(string(paint.RuneTriangleRight) + " comments")
	toggleComments.Show()
	toggleComments.SetTheme(SidebarHeaderTheme)

	toggleLocals.Connect(ctk.SignalActivate, "locals-toggle-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		toggleLocals.SetLabel(string(paint.RuneTriangleDown) + " locals")
		toggleCustom.SetLabel(string(paint.RuneTriangleRight) + " custom")
		toggleComments.SetLabel(string(paint.RuneTriangleRight) + " comments")
		sidebarLocalsFrame.SetSizeRequest(-1, -1)
		sidebarCustomFrame.SetSizeRequest(-1, 1)
		sidebarCommentsFrame.SetSizeRequest(-1, 1)
		e.SidebarLocalsList.SetSizeRequest(gSidebarInnerWidth, -1)
		e.SidebarCustomList.SetSizeRequest(gSidebarInnerWidth, 0)
		e.SidebarCommentsList.SetSizeRequest(gSidebarInnerWidth, 0)
		e.Window.Resize()
		e.Window.ReApplyStyles()
		e.Display.RequestDraw()
		e.Display.RequestShow()
		// toggleListSection("locals")
		return cenums.EVENT_STOP
	})

	toggleCustom.Connect(ctk.SignalActivate, "custom-toggle-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		toggleLocals.SetLabel(string(paint.RuneRArrow) + " locals")
		toggleCustom.SetLabel(string(paint.RuneDArrow) + " custom")
		toggleComments.SetLabel(string(paint.RuneRArrow) + " comments")
		sidebarLocalsFrame.SetSizeRequest(-1, 1)
		sidebarCustomFrame.SetSizeRequest(-1, -1)
		sidebarCommentsFrame.SetSizeRequest(-1, 1)
		e.SidebarLocalsList.SetSizeRequest(gSidebarInnerWidth, 0)
		e.SidebarCustomList.SetSizeRequest(gSidebarInnerWidth, -1)
		e.SidebarCommentsList.SetSizeRequest(gSidebarInnerWidth, 0)
		e.Window.Resize()
		e.Window.ReApplyStyles()
		e.Display.RequestDraw()
		e.Display.RequestShow()
		// toggleListSection("custom")
		return cenums.EVENT_STOP
	})

	toggleComments.Connect(ctk.SignalActivate, "custom-toggle-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		toggleLocals.SetLabel(string(paint.RuneRArrow) + " locals")
		toggleCustom.SetLabel(string(paint.RuneRArrow) + " custom")
		toggleComments.SetLabel(string(paint.RuneDArrow) + " comments")
		sidebarLocalsFrame.SetSizeRequest(-1, 1)
		sidebarCustomFrame.SetSizeRequest(-1, 1)
		sidebarCommentsFrame.SetSizeRequest(-1, -1)
		e.SidebarLocalsList.SetSizeRequest(gSidebarInnerWidth, 0)
		e.SidebarCustomList.SetSizeRequest(gSidebarInnerWidth, 0)
		e.SidebarCommentsList.SetSizeRequest(gSidebarInnerWidth, -1)
		e.Window.Resize()
		e.Window.ReApplyStyles()
		e.Display.RequestDraw()
		e.Display.RequestShow()
		// toggleListSection("custom")
		return cenums.EVENT_STOP
	})

	// localhost list

	sidebarLocalsFrame = ctk.NewFrameWithWidget(toggleLocals)
	sidebarLocalsFrame.Show()
	sidebarLocalsFrame.SetTheme(SidebarFrameTheme)
	sidebarLocalsFrame.SetSizeRequest(-1, 1)
	sidebarLocalsFrame.SetLabelAlign(0.0, 0.5)
	sidebarVBox.PackStart(sidebarLocalsFrame, false, false, 0)

	sidebarLocalsScroll := ctk.NewScrolledViewport()
	sidebarLocalsScroll.Show()
	sidebarLocalsScroll.SetPolicy(enums.PolicyAutomatic, enums.PolicyNever)
	sidebarLocalsFrame.Add(sidebarLocalsScroll)

	e.SidebarLocalsList = ctk.NewVBox(false, 0)
	e.SidebarLocalsList.Show()
	e.SidebarLocalsList.SetSizeRequest(gSidebarInnerWidth, -1)
	sidebarLocalsScroll.Add(e.SidebarLocalsList)

	// custom list

	sidebarCustomFrame = ctk.NewFrameWithWidget(toggleCustom)
	sidebarCustomFrame.Show()
	sidebarCustomFrame.SetSizeRequest(-1, -1)
	sidebarCustomFrame.SetLabelAlign(0.0, 0.5)
	sidebarCustomFrame.SetTheme(SidebarFrameTheme)
	sidebarVBox.PackStart(sidebarCustomFrame, false, false, 0)

	sidebarCustomScroll := ctk.NewScrolledViewport()
	sidebarCustomScroll.Show()
	sidebarCustomScroll.SetPolicy(enums.PolicyAutomatic, enums.PolicyNever)
	sidebarCustomFrame.Add(sidebarCustomScroll)

	e.SidebarCustomList = ctk.NewVBox(false, 0)
	e.SidebarCustomList.Show()
	e.SidebarCustomList.SetSizeRequest(gSidebarInnerWidth, -1)
	sidebarCustomScroll.Add(e.SidebarCustomList)

	// comments list

	sidebarCommentsFrame = ctk.NewFrameWithWidget(toggleComments)
	sidebarCommentsFrame.Show()
	sidebarCommentsFrame.SetSizeRequest(-1, 1)
	sidebarCommentsFrame.SetLabelAlign(0.0, 0.5)
	sidebarCommentsFrame.SetTheme(SidebarFrameTheme)
	sidebarVBox.PackStart(sidebarCommentsFrame, false, false, 0)

	sidebarCommentsScroll := ctk.NewScrolledViewport()
	sidebarCommentsScroll.Show()
	sidebarCommentsScroll.SetPolicy(enums.PolicyAutomatic, enums.PolicyNever)
	sidebarCommentsFrame.Add(sidebarCommentsScroll)

	e.SidebarCommentsList = ctk.NewVBox(false, 0)
	e.SidebarCommentsList.Show()
	e.SidebarCommentsList.SetSizeRequest(gSidebarInnerWidth, -1)
	sidebarCommentsScroll.Add(e.SidebarCommentsList)

	// entry list

	sidebarEntryFrame = ctk.NewFrame("entries")
	// sidebarEntryFrame.Show()
	sidebarEntryFrame.SetSizeRequest(-1, -1)
	sidebarEntryFrame.SetLabelAlign(0.0, 0.5)
	sidebarEntryFrame.SetTheme(SidebarFrameTheme)
	sidebarVBox.PackStart(sidebarEntryFrame, false, false, 0)

	sidebarEntryScroll := ctk.NewScrolledViewport()
	sidebarEntryScroll.Show()
	sidebarEntryScroll.SetPolicy(enums.PolicyAutomatic, enums.PolicyNever)
	sidebarEntryFrame.Add(sidebarEntryScroll)

	e.SidebarEntryList = ctk.NewVBox(false, 0)
	e.SidebarEntryList.Show()
	e.SidebarEntryList.SetSizeRequest(gSidebarInnerWidth, -1)
	sidebarEntryScroll.Add(e.SidebarEntryList)

	// sidebar action buttons

	sidebarActionHBox := ctk.NewHBox(false, 1)
	sidebarActionHBox.Show()
	sidebarVBox.PackStart(sidebarActionHBox, false, false, 0)

	e.SidebarAddEntryButton = ctk.NewButtonWithLabel("+")
	e.SidebarAddEntryButton.Show()
	e.SidebarAddEntryButton.SetSizeRequest(-1, 1)
	e.SidebarAddEntryButton.SetHasTooltip(true)
	e.SidebarAddEntryButton.SetTooltipText("Add new entry")
	e.SidebarAddEntryButton.Connect(ctk.SignalActivate, gSidebarAddRowHandler, e.activateSidebarAddRowHandler)
	sidebarActionHBox.PackStart(e.SidebarAddEntryButton, true, true, 0)

	e.SidebarMoveEntryUpButton = ctk.NewButtonWithLabel(string(paint.RuneTriangleUp))
	// e.SidebarMoveEntryUpButton.Show()
	e.SidebarMoveEntryUpButton.SetSizeRequest(-1, 1)
	e.SidebarMoveEntryUpButton.SetHasTooltip(true)
	e.SidebarMoveEntryUpButton.SetTooltipText("Move selected entry up in the hosts file order")
	e.SidebarMoveEntryUpButton.Connect(ctk.SignalActivate, gSidebarMoveRowUpHandler, e.activateSidebarMoveRowUpHandler)
	sidebarActionHBox.PackStart(e.SidebarMoveEntryUpButton, true, true, 0)

	e.SidebarMoveEntryDownButton = ctk.NewButtonWithLabel(string(paint.RuneTriangleDown))
	// e.SidebarMoveEntryDownButton.Show()
	e.SidebarMoveEntryDownButton.SetSizeRequest(-1, 1)
	e.SidebarMoveEntryDownButton.SetHasTooltip(true)
	e.SidebarMoveEntryDownButton.SetTooltipText("Move selected entry down in the hosts file order")
	e.SidebarMoveEntryDownButton.Connect(ctk.SignalActivate, gSidebarMoveRowDownHandler, e.activateSidebarMoveRowDownHandler)
	sidebarActionHBox.PackStart(e.SidebarMoveEntryDownButton, true, true, 0)

	// nothing selected panel

	e.NothingSelectedFrame = ctk.NewFrame("")
	e.NothingSelectedFrame.Show()
	e.NothingSelectedFrame.SetLabelAlign(0.0, 0.5)
	e.EditingHBox.PackStart(e.NothingSelectedFrame, true, true, 0)

	nothingSelectedLabel := ctk.NewLabel("(please select a host)")
	nothingSelectedLabel.Show()
	nothingSelectedLabel.SetAlignment(0.5, 0.5)
	nothingSelectedLabel.SetJustify(cenums.JUSTIFY_CENTER)
	e.NothingSelectedFrame.Add(nothingSelectedLabel)

	// host entry panel

	e.HostSelectedFrame = ctk.NewFrame("")
	e.HostSelectedFrame.SetLabelAlign(0.0, 0.5)
	e.EditingHBox.PackStart(e.HostSelectedFrame, true, true, 0)

	panelVBox := ctk.NewVBox(false, 0)
	panelVBox.Show()
	e.HostSelectedFrame.Add(panelVBox)

	addSeparator := func(vbox ctk.VBox) {
		sep := ctk.NewSeparator()
		sep.Show()
		sep.SetSizeRequest(-1, 1)
		vbox.PackStart(sep, false, false, 0)
	}

	addInstructions := func(vbox ctk.VBox, format string, argv ...interface{}) {
		label, _ := ctk.NewLabelWithMarkup(fmt.Sprintf(format, argv...))
		label.Show()
		label.SetSingleLineMode(true)
		label.SetSizeRequest(-1, 1)
		vbox.PackStart(label, false, false, 0)
	}

	addInstructions(panelVBox, "Comments for this entry:")

	e.CommentsEntry = ctk.NewEntry("")
	e.CommentsEntry.SetName("editing-comment")
	e.CommentsEntry.Show()
	e.CommentsEntry.SetLineWrap(true)
	e.CommentsEntry.SetLineWrapMode(cenums.WRAP_NONE)
	e.CommentsEntry.SetSingleLineMode(false)
	e.CommentsEntry.SetSelectable(true)
	e.CommentsEntry.SetSizeRequest(-1, 3)
	panelVBox.PackStart(e.CommentsEntry, false, false, 0)

	e.HostEditVBox = ctk.NewVBox(false, 0)
	e.HostEditVBox.Show()
	panelVBox.PackStart(e.HostEditVBox, true, true, 0)

	addSeparator(e.HostEditVBox)
	addInstructions(e.HostEditVBox, "Enter an IP address (or nslookup domain):")

	addrBttnBox := ctk.NewHBox(true, 1)
	addrBttnBox.Show()
	addrBttnBox.SetSizeRequest(-1, 1)
	e.HostEditVBox.PackStart(addrBttnBox, false, false, 0)

	e.AddressEntry = ctk.NewEntry("")
	e.AddressEntry.Show()
	e.AddressEntry.SetSelectable(true)
	e.AddressEntry.SetLineWrap(false)
	e.AddressEntry.SetSizeRequest(-1, 1)
	e.AddressEntry.SetSingleLineMode(true)
	addrBttnBox.PackStart(e.AddressEntry, true, true, 0)

	e.AddressButton = ctk.NewButtonWithLabel("(address)")
	e.AddressButton.Show()
	e.AddressButton.SetSizeRequest(-1, 1)
	addrBttnBox.PackStart(e.AddressButton, true, true, 0)

	addSeparator(e.HostEditVBox)
	addInstructions(e.HostEditVBox, "Space separated list of domain names:")

	e.DomainsEntry = ctk.NewEntry("")
	e.DomainsEntry.Show()
	e.DomainsEntry.SetSelectable(true)
	e.DomainsEntry.SetSingleLineMode(false)
	// e.DomainsEntry.SetLineWrap(true)
	// e.DomainsEntry.SetLineWrapMode(cenums.WRAP_WORD)
	// e.DomainsEntry.SetJustify(cenums.JUSTIFY_LEFT)
	e.HostEditVBox.PackStart(e.DomainsEntry, true, true, 0)

	addSeparator(panelVBox)
	addInstructions(panelVBox, "Hosts file entry actions:")

	hostActionHBox := ctk.NewHBox(true, 1)
	hostActionHBox.Show()
	hostActionHBox.SetSizeRequest(-1, 1)
	panelVBox.PackStart(hostActionHBox, false, false, 0)

	e.ActivateButton = ctk.NewButtonWithLabel("")
	e.ActivateButton.Show()
	e.ActivateButton.SetSizeRequest(-1, 1)
	// panelVBox.PackStart(e.ActivateButton, false, false, 0)
	hostActionHBox.PackStart(e.ActivateButton, true, true, 0)

	e.DeleteButton = ctk.NewButtonWithLabel("click to delete")
	e.DeleteButton.Show()
	e.DeleteButton.SetSizeRequest(-1, 1)
	// panelVBox.PackStart(e.DeleteButton, true, true, 0)
	hostActionHBox.PackStart(e.DeleteButton, true, true, 0)

	changeSidebarMode(ListByDomain)

	return e.EditingHBox
}

func (e *CEheditor) reloadEditor() {
	e.Window.Freeze()
	defer func() {
		e.Window.Thaw()
		e.Window.Resize()
		e.Window.ReApplyStyles()
		e.Display.RequestDraw()
		e.Display.RequestSync()
	}()

	purgeContents := func(v ctk.VBox) {
		existing := v.GetChildren()
		for _, child := range existing {
			v.Remove(child)
			child.Destroy()
		}
	}

	e.EditorCommentList = []*editor.Host{}
	purgeContents(e.SidebarEntryList)
	purgeContents(e.SidebarLocalsList)
	purgeContents(e.SidebarCustomList)
	purgeContents(e.SidebarCommentsList)

	e.EditorAddressLookup = make(map[string]*editor.Host)
	e.EditorDomainsLookup = make(map[string]*editor.Host)

	var changed bool
	unique := make(map[string]int)
	for _, host := range e.HostFile.Hosts() {
		if !changed && host.Changed() {
			changed = true
		}
		if host.IsOnlyComment() {
			e.EditorCommentList = append(e.EditorCommentList, host)
			continue
		}
		e.EditorAddressLookup[host.Address()] = host
		if domains := host.Domains(); len(domains) == 0 {
			e.EditorDomainsLookup[""] = host
		} else {
			for _, domain := range domains {
				if _, found := unique[domain]; !found {
					unique[domain] = 1
				} else {
					unique[domain] += 1
				}
				if unique[domain] > 1 {
					e.EditorDomainsLookup[domain+" ("+strconv.Itoa(unique[domain])+")"] = host
				} else {
					e.EditorDomainsLookup[domain] = host
				}
			}
		}
	}

	e.SaveButton.SetSensitive(changed)
	e.ReloadButton.SetSensitive(changed)

	e.updateEditor()

}

func (e *CEheditor) updateEditor() {
	switch e.SidebarMode {
	case ListByDomain, ListByAddress:
		e.updateEditorByAddressOrDomain()
	case ListByEntry:
		e.updateEditorByEntry()
	}
}

func (e *CEheditor) updateEditorByAddressOrDomain() {
	var choices map[string]*editor.Host
	if e.SidebarMode == ListByAddress {
		choices = e.EditorAddressLookup
	} else {
		choices = e.EditorDomainsLookup
	}

	var localsCount, customCount int

	for idx, host := range e.EditorCommentList {
		key := fmt.Sprintf("Comment (%d)", idx+1)
		b := e.makeSidebarButton(key, host)
		e.SidebarCommentsList.PackStart(b, false, false, 0)
	}

	for _, key := range SortedKeys(choices) {
		host := choices[key]
		b := e.makeSidebarButton(key, host)
		switch host.Importance() {
		case editor.HostIsLocalhostIPv4:
			localsCount += 1
			e.SidebarLocalsList.PackStart(b, false, false, 0)
			e.SidebarLocalsList.SetSizeRequest(gSidebarInnerWidth, localsCount)
		case editor.HostIsLocalhostIPv6:
			localsCount += 1
			e.SidebarLocalsList.PackStart(b, false, false, 0)
			e.SidebarLocalsList.SetSizeRequest(gSidebarInnerWidth, localsCount)
		default:
			customCount += 1
			e.SidebarCustomList.PackStart(b, false, false, 0)
			e.SidebarCustomList.SetSizeRequest(gSidebarInnerWidth, customCount)
		}
	}
}

func (e *CEheditor) updateEditorByEntry() {
	hosts := e.HostFile.Hosts()
	e.SidebarEntryList.SetSizeRequest(-1, len(hosts))
	var commentsCount int
	for idx, host := range hosts {
		key := strconv.Itoa(idx+1) + ". "
		if host.IsOnlyComment() {
			commentsCount += 1
			key += fmt.Sprintf("comment (%d)", commentsCount)
		} else {
			key += host.Address()
		}
		b := e.makeSidebarButton(key, host)
		e.SidebarEntryList.PackStart(b, false, false, 0)
	}
}

func (e *CEheditor) makeSidebarButton(key string, host *editor.Host) (b ctk.Button) {
	label := ctk.NewLabel(key)
	label.Show()
	label.SetJustify(cenums.JUSTIFY_LEFT)
	label.SetSizeRequest(-1, 1)
	label.SetSingleLineMode(true)

	b = ctk.NewButtonWithWidget(label)
	_ = b.InstallProperty(cdk.Property("host"), cdk.StructProperty, true, host)
	b.Show()
	b.SetSizeRequest(-1, 1)
	b.Connect(ctk.SignalActivate, key+"-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		if h, ok := data[0].(*editor.Host); ok {
			e.focusEditor(h)
		}
		return cenums.EVENT_STOP
	}, host)

	var theme paint.Theme
	var name, tooltip string

	if e.SelectedHost != nil && e.SelectedHost.Equals(host) {
		name = "editing-list-selected"
	} else {
		name = "editing-list-unselected"
	}

	if host.Active() {
		theme = SidebarActiveTheme
		tooltip = key + " is active"
	} else {
		theme = SidebarButtonTheme
		tooltip = key + " is inactive"
	}

	switch host.Importance() {
	case editor.HostIsLocalhostIPv4:
		tooltip += "\n" + key + " points to an IPv4 localhost address"
	case editor.HostIsLocalhostIPv6:
		tooltip += "\n" + key + " points to an IPv6 localhost address"
	default:
		if lookup := host.Lookup(); lookup == "" {
			tooltip += "\n" + key + " points to a static ip address"
		} else {
			tooltip += "\n" + key + " points to a dynamic ip address"
			tooltip += "\n" + key + " gets the ip from " + lookup
		}
	}

	b.SetName(name)
	b.SetTheme(theme)
	b.SetHasTooltip(true)
	b.SetTooltipText(tooltip)
	return
}

func (e *CEheditor) focusEditor(host *editor.Host) {
	e.Window.Freeze()
	defer func() {
		e.updateSidebarActionButtons()
		e.Window.Thaw()
		e.Window.Resize()
		e.Window.ReApplyStyles()
		e.Display.RequestDraw()
		e.Display.RequestShow()
		e.reloadEditor()
	}()

	if host == nil {
		e.Window.LogDebug("clearing editor focus")
		e.SelectedHost = nil
		e.HostSelectedFrame.Hide()
		e.NothingSelectedFrame.Show()
		return
	}

	e.Window.LogDebug("focusing editor on: %v", host.String())
	e.SelectedHost = host
	e.HostSelectedFrame.Show()
	e.NothingSelectedFrame.Hide()
	e.CommentsEntry.SetText(host.Comment())

	if v := host.Lookup(); v != "" {
		e.AddressEntry.SetText(v)
	} else {
		e.AddressEntry.SetText(host.Address())
	}

	_ = e.CommentsEntry.Disconnect(ctk.SignalChangedText, "comments-changed-handler")
	e.CommentsEntry.Connect(ctk.SignalChangedText, "comments-changed-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		h, _ := data[0].(*editor.Host)
		h.SetComment(e.CommentsEntry.GetText())
		e.CommentsEntry.LogDebug("updated host %v comment: %v", h.Address(), e.CommentsEntry.GetText())
		e.reloadEditor()
		return cenums.EVENT_STOP
	}, host)

	_ = e.DomainsEntry.Disconnect(ctk.SignalChangedText, "domains-changed-handler")
	// e.DomainsEntry.SetText(strings.Join(host.Domains(), " "))
	alloc := e.DomainsEntry.GetAllocation()
	var domainLines []string
	var current string
	for _, domain := range host.Domains() {
		currentLen := len(current)
		if currentLen > 0 {
			domainLen := len(domain)
			if currentLen+1+domainLen >= alloc.W {
				domainLines = append(domainLines, current)
				current = ""
			} else {
				current += " "
			}
		}
		current += domain
	}
	if current != "" {
		domainLines = append(domainLines, current)
	}
	e.DomainsEntry.SetText(strings.Join(domainLines, "\n"))
	e.DomainsEntry.Connect(ctk.SignalChangedText, "domains-changed-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		h, _ := data[0].(*editor.Host)
		h.SetDomains(e.DomainsEntry.GetText())
		e.reloadEditor()
		return cenums.EVENT_STOP
	}, host)

	handle := "activate-button-handler"
	_ = e.ActivateButton.Disconnect(ctk.SignalActivate, handle)
	if host.Importance() != editor.HostNotImportant {

		e.AddressEntry.SetSensitive(false)
		e.AddressButton.SetSensitive(false)
		e.ActivateButton.SetSensitive(false)

		e.ActivateButton.SetTheme(DefaultButtonTheme)
		e.ActivateButton.SetLabel("cannot deactivate host")

		e.DeleteButton.SetSensitive(false)
		e.DeleteButton.SetLabel("cannot delete host")
		_ = e.DeleteButton.Disconnect(ctk.SignalActivate, "delete-entry-handler")

	} else {

		e.AddressEntry.SetSensitive(true)
		addressEntryText := e.AddressEntry.GetText()
		e.AddressButton.SetSensitive(!cstrings.StringIsIP(addressEntryText) && cstrings.StringIsDomainName(addressEntryText))
		e.ActivateButton.SetSensitive(true)

		e.DeleteButton.SetSensitive(true)
		e.DeleteButton.SetLabel("click to delete")
		_ = e.DeleteButton.Disconnect(ctk.SignalActivate, "delete-entry-handler")
		e.DeleteButton.Connect(ctk.SignalActivate, "delete-entry-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			h, _ := data[0].(*editor.Host)
			message := ""
			if h.Empty() {
				if h.IsOnlyComment() {
					message = "(empty comment entry)"
				} else {
					message = "(empty host entry)"
				}
			} else {
				message = h.Block()
			}
			d := ctk.NewYesNoDialog("Remove Entry?", message, true)
			d.SetSizeRequest(54, 10)
			d.RunFunc(func(response enums.ResponseType, argv ...interface{}) {
				switch response {
				case enums.ResponseYes:
					for idx, hh := range e.HostFile.Hosts() {
						if h.Equals(hh) {
							log.DebugF("removing entry at index: %v", idx)
							e.HostFile.RemoveHost(idx)
							e.SelectedHost = nil
							break
						}
					}
					e.reloadEditor()
					e.focusEditor(nil)
				case enums.ResponseCancel, enums.ResponseClose, enums.ResponseNo:
					log.DebugF("user cancelled removal operation")
				}
			})
			return cenums.EVENT_STOP
		}, host)

		if host.Active() {

			e.ActivateButton.SetTheme(ActiveButtonTheme)
			e.ActivateButton.SetLabel("click to deactivate")
			e.ActivateButton.Connect(ctk.SignalActivate, handle, func(data []interface{}, argv ...interface{}) cenums.EventFlag {
				if h, ok := data[0].(*editor.Host); ok {
					h.SetActive(false)
					e.reloadEditor()
					e.focusEditor(h)
				}
				return cenums.EVENT_STOP
			}, host)

		} else {

			e.ActivateButton.SetTheme(DefaultButtonTheme)
			e.ActivateButton.SetLabel("click to activate")
			e.ActivateButton.Connect(ctk.SignalActivate, handle, func(data []interface{}, argv ...interface{}) cenums.EventFlag {
				if h, ok := data[0].(*editor.Host); ok {
					h.SetActive(true)
					e.reloadEditor()
					e.focusEditor(h)
				}
				return cenums.EVENT_STOP
			}, host)

		}
	}

	_ = e.AddressButton.Disconnect(ctk.SignalActivate, "address-activate-handler")
	e.AddressButton.Connect(
		ctk.SignalActivate,
		"address-activate-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			if err := e.newNsLookupDialog(host); err != nil {
				e.AddressButton.LogErr(err)
			}
			return cenums.EVENT_STOP
		},
		host,
	)

	actualLabel, actualTooltip := host.GetActualInfo()
	e.AddressButton.SetLabel(actualLabel)
	e.AddressButton.SetTooltipText(actualLabel + "\n" + actualTooltip)
	e.AddressButton.Resize()

	_ = e.AddressEntry.Disconnect(ctk.SignalChangedText, "address-text-changed-handler")
	e.AddressEntry.Connect(ctk.SignalChangedText, "address-text-changed-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		var h *editor.Host
		h, _ = data[0].(*editor.Host)
		text := h.Address()
		changed := e.AddressEntry.GetText()
		if cstrings.StringIsIP(changed) {
			if text != changed {
				e.AddressButton.SetLabel(fmt.Sprintf("(%v)", changed))
				e.AddressButton.SetTooltipText("is a valid IP address")
			} else {
				e.AddressButton.SetLabel(fmt.Sprintf("(%v)", text))
				e.AddressButton.SetTooltipText("is a valid IP address")
			}
			e.AddressButton.SetSensitive(false)
			h.SetAddress(changed)
		} else if cstrings.StringIsDomainName(changed) {
			if text != changed {
				e.AddressButton.SetLabel("(lookup changed)")
				e.AddressButton.SetTooltipText("click to perform domain lookup")
			}
			e.AddressButton.SetSensitive(true)
			h.SetAddress(changed)
			h.SetLookup(changed)
		} else {
			e.AddressButton.SetSensitive(false)
			e.AddressButton.SetLabel("(not ip or domain)")
			e.AddressButton.SetTooltipText("enter a valid address or domain name")
			h.SetAddress(changed)
			h.SetLookup("")
		}
		e.reloadEditor()
		return cenums.EVENT_PASS
	}, host)

	allEntries := append(
		append(
			e.SidebarLocalsList.GetChildren(),
			e.SidebarCustomList.GetChildren()...,
		),
		e.SidebarEntryList.GetChildren()...,
	)
	for _, child := range allEntries {
		if b, ok := child.(ctk.Button); ok {
			var h *editor.Host
			if v, err := b.GetStructProperty(cdk.Property("host")); err != nil {
				b.LogError("error getting button host: %v", err)
				continue
			} else if h, ok = v.(*editor.Host); !ok {
				b.LogError("error button host is not an *editor.Host")
				continue
			}
			switch e.SidebarMode {
			case ListByEntry:
				if e.SelectedHost != nil && e.SelectedHost.IsOnlyComment() {
					if h.Equals(e.SelectedHost) {
						b.SetName("editing-list-selected")
					} else {
						b.SetName("editing-list-unselected")
					}
					break
				}
				fallthrough
			case ListByAddress:
				if h.Address() == host.Address() {
					b.SetName("editing-list-selected")
				} else {
					b.SetName("editing-list-unselected")
				}
			case ListByDomain:
				_, key, _ := getSidebarButtonInfo(b)
				if host.Address() == h.Address() && cstrings.StringSliceHasValue(host.Domains(), key) {
					b.SetName("editing-list-selected")
				} else {
					b.SetName("editing-list-unselected")
				}
			}
		}
	}

	if host.IsOnlyComment() {
		e.CommentsEntry.SetSizeRequest(-1, -1)
		e.HostEditVBox.Hide()
		e.ActivateButton.Hide()
	} else {
		e.CommentsEntry.SetSizeRequest(-1, 3)
		e.HostEditVBox.Show()
		e.ActivateButton.Show()
	}
}

var rxSidebarButtonLabel = regexp.MustCompile(`^(\d+\.)??\s??(\S+?)\s??(\(\d+\))??$`)

func getSidebarButtonInfo(b ctk.Button) (idx int, key, extra string) {
	key = b.GetLabel()
	if rxSidebarButtonLabel.MatchString(key) {
		m := rxSidebarButtonLabel.FindAllStringSubmatch(key, 1)
		var err error
		if idx, err = strconv.Atoi(m[0][1]); err == nil {
			idx -= 1
		}
		key = m[0][2]
		extra = m[0][3]
		return
	}
	// fallback?
	parts := strings.Split(key, " ")
	if len(parts) <= 1 {
		return
	}
	key = parts[0]
	extra = strings.Join(parts[1:], " ")
	return
}