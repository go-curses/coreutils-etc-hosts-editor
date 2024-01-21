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
	"regexp"
	"strconv"
	"strings"

	"github.com/go-corelibs/maps"
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
	gSidebarOuterWidth = 25
	gSidebarInnerWidth = 20
)

func (c *CUI) switchToEditor() {
	c.Window.Freeze()
	c.ContentsHBox.Freeze()
	for _, child := range c.ContentsHBox.GetChildren() {
		if child.ObjectID() == c.EditingHBox.ObjectID() {
			child.Show()
		} else {
			child.Hide()
		}
	}
	c.Window.Thaw()
	c.ContentsHBox.Thaw()
	c.focusEditor(c.SelectedHost)
}

func (c *CUI) makeEditor() ctk.Widget {
	c.EditingHBox = ctk.NewHBox(false, 1)
	c.EditingHBox.SetName("editing")
	c.EditingHBox.Show()

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
			eWidth = gSidebarInnerWidth - 5
			eLabel = "_Entry"
			eTheme = ActiveButtonTheme
			c.SidebarMode = ListByEntry
			sidebarCommentsFrame.Hide()
			sidebarLocalsFrame.Hide()
			sidebarCustomFrame.Hide()
			sidebarEntryFrame.Show()
		case ListByAddress:
			aWidth = gSidebarInnerWidth - 5
			aLabel = "_Address"
			aTheme = ActiveButtonTheme
			c.SidebarMode = ListByAddress
			sidebarCommentsFrame.Show()
			sidebarLocalsFrame.Show()
			sidebarCustomFrame.Show()
			sidebarEntryFrame.Hide()
		case ListByDomain:
			dWidth = gSidebarInnerWidth - 5
			dLabel = "_Domain"
			dTheme = ActiveButtonTheme
			c.SidebarMode = ListByDomain
			sidebarCommentsFrame.Show()
			sidebarLocalsFrame.Show()
			sidebarCustomFrame.Show()
			sidebarEntryFrame.Hide()
		}

		c.ByEntryButton.SetSizeRequest(eWidth, 1)
		c.ByAddressButton.SetSizeRequest(aWidth, 1)
		c.ByDomainsButton.SetSizeRequest(dWidth, 1)

		c.ByEntryButton.SetLabel(eLabel)
		c.ByAddressButton.SetLabel(aLabel)
		c.ByDomainsButton.SetLabel(dLabel)

		c.ByEntryButton.SetTheme(eTheme)
		c.ByAddressButton.SetTheme(aTheme)
		c.ByDomainsButton.SetTheme(dTheme)

		c.updateSidebarActionButtons()
	}

	c.ByDomainsButton = ctk.NewButtonWithMnemonic("_Domain")
	c.ByDomainsButton.Show()
	c.ByDomainsButton.Connect(ctk.SignalActivate, "by-domains-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		c.ByDomainsButton.LogDebug("clicked")
		changeSidebarMode(ListByDomain)
		c.reloadEditor()
		c.focusEditor(nil)
		return cenums.EVENT_STOP
	})
	c.ByDomainsButton.SetHasTooltip(true)
	c.ByDomainsButton.SetTooltipText("Click to list by domain names")
	sidebarViewCtrlHBox.PackStart(c.ByDomainsButton, false, false, 0)

	c.ByAddressButton = ctk.NewButtonWithMnemonic("_A")
	c.ByAddressButton.Show()
	c.ByAddressButton.Connect(ctk.SignalActivate, "by-address-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		c.ByAddressButton.LogDebug("clicked")
		changeSidebarMode(ListByAddress)
		c.reloadEditor()
		c.focusEditor(nil)
		return cenums.EVENT_STOP
	})
	c.ByAddressButton.SetHasTooltip(true)
	c.ByAddressButton.SetTooltipText("Click to list by IP addresses")
	sidebarViewCtrlHBox.PackStart(c.ByAddressButton, false, false, 0)

	c.ByEntryButton = ctk.NewButtonWithMnemonic("_E")
	c.ByEntryButton.Show()
	c.ByEntryButton.Connect(ctk.SignalActivate, "by-entry-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		c.ByEntryButton.LogDebug("clicked")
		changeSidebarMode(ListByEntry)
		c.reloadEditor()
		c.focusEditor(nil)
		return cenums.EVENT_STOP
	})
	c.ByEntryButton.SetHasTooltip(true)
	c.ByEntryButton.SetTooltipText("Click to list by hosts file entry")
	sidebarViewCtrlHBox.PackStart(c.ByEntryButton, false, false, 0)

	c.SidebarFrame = ctk.NewFrameWithWidget(sidebarViewCtrlHBox)
	c.SidebarFrame.Show()
	c.SidebarFrame.SetLabelAlign(0.0, 0.5)
	c.SidebarFrame.SetSizeRequest(gSidebarOuterWidth, -1)
	c.EditingHBox.PackStart(c.SidebarFrame, false, false, 0)

	sidebarVBox := ctk.NewVBox(false, 0)
	sidebarVBox.Show()
	c.SidebarFrame.Add(sidebarVBox)

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
		c.SidebarLocalsList.SetSizeRequest(gSidebarInnerWidth, -1)
		c.SidebarCustomList.SetSizeRequest(gSidebarInnerWidth, 0)
		c.SidebarCommentsList.SetSizeRequest(gSidebarInnerWidth, 0)
		c.Window.Resize()
		c.Window.ReApplyStyles()
		c.Display.RequestDraw()
		c.Display.RequestShow()
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
		c.SidebarLocalsList.SetSizeRequest(gSidebarInnerWidth, 0)
		c.SidebarCustomList.SetSizeRequest(gSidebarInnerWidth, -1)
		c.SidebarCommentsList.SetSizeRequest(gSidebarInnerWidth, 0)
		c.Window.Resize()
		c.Window.ReApplyStyles()
		c.Display.RequestDraw()
		c.Display.RequestShow()
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
		c.SidebarLocalsList.SetSizeRequest(gSidebarInnerWidth, 0)
		c.SidebarCustomList.SetSizeRequest(gSidebarInnerWidth, 0)
		c.SidebarCommentsList.SetSizeRequest(gSidebarInnerWidth, -1)
		c.Window.Resize()
		c.Window.ReApplyStyles()
		c.Display.RequestDraw()
		c.Display.RequestShow()
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

	c.SidebarLocalsList = ctk.NewVBox(false, 0)
	c.SidebarLocalsList.Show()
	c.SidebarLocalsList.SetSizeRequest(gSidebarInnerWidth, -1)
	sidebarLocalsScroll.Add(c.SidebarLocalsList)

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

	c.SidebarCustomList = ctk.NewVBox(false, 0)
	c.SidebarCustomList.Show()
	c.SidebarCustomList.SetSizeRequest(gSidebarInnerWidth, -1)
	sidebarCustomScroll.Add(c.SidebarCustomList)

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

	c.SidebarCommentsList = ctk.NewVBox(false, 0)
	c.SidebarCommentsList.Show()
	c.SidebarCommentsList.SetSizeRequest(gSidebarInnerWidth, -1)
	sidebarCommentsScroll.Add(c.SidebarCommentsList)

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

	c.SidebarEntryList = ctk.NewVBox(false, 0)
	c.SidebarEntryList.Show()
	c.SidebarEntryList.SetSizeRequest(gSidebarInnerWidth, -1)
	sidebarEntryScroll.Add(c.SidebarEntryList)

	// sidebar action buttons

	sidebarActionHBox := ctk.NewHBox(false, 1)
	sidebarActionHBox.Show()
	sidebarVBox.PackStart(sidebarActionHBox, false, false, 0)

	c.SidebarAddEntryButton = ctk.NewButtonWithLabel("+")
	c.SidebarAddEntryButton.Show()
	c.SidebarAddEntryButton.SetSizeRequest(-1, 1)
	c.SidebarAddEntryButton.SetHasTooltip(true)
	c.SidebarAddEntryButton.SetTooltipText("Add new entry")
	c.SidebarAddEntryButton.Connect(ctk.SignalActivate, gSidebarAddRowHandler, c.activateSidebarAddRowHandler)
	sidebarActionHBox.PackStart(c.SidebarAddEntryButton, true, true, 0)

	c.SidebarMoveEntryUpButton = ctk.NewButtonWithLabel(string(paint.RuneTriangleUp))
	// c.SidebarMoveEntryUpButton.Show()
	c.SidebarMoveEntryUpButton.SetSizeRequest(-1, 1)
	c.SidebarMoveEntryUpButton.SetHasTooltip(true)
	c.SidebarMoveEntryUpButton.SetTooltipText("Move selected entry up in the hosts file order")
	c.SidebarMoveEntryUpButton.Connect(ctk.SignalActivate, gSidebarMoveRowUpHandler, c.activateSidebarMoveRowUpHandler)
	sidebarActionHBox.PackStart(c.SidebarMoveEntryUpButton, true, true, 0)

	c.SidebarMoveEntryDownButton = ctk.NewButtonWithLabel(string(paint.RuneTriangleDown))
	// c.SidebarMoveEntryDownButton.Show()
	c.SidebarMoveEntryDownButton.SetSizeRequest(-1, 1)
	c.SidebarMoveEntryDownButton.SetHasTooltip(true)
	c.SidebarMoveEntryDownButton.SetTooltipText("Move selected entry down in the hosts file order")
	c.SidebarMoveEntryDownButton.Connect(ctk.SignalActivate, gSidebarMoveRowDownHandler, c.activateSidebarMoveRowDownHandler)
	sidebarActionHBox.PackStart(c.SidebarMoveEntryDownButton, true, true, 0)

	// nothing selected panel

	c.NothingSelectedFrame = ctk.NewFrame("")
	c.NothingSelectedFrame.Show()
	c.NothingSelectedFrame.SetLabelAlign(0.0, 0.5)
	c.EditingHBox.PackStart(c.NothingSelectedFrame, true, true, 0)

	nothingSelectedLabel := ctk.NewLabel("(please select a host)")
	nothingSelectedLabel.Show()
	nothingSelectedLabel.SetAlignment(0.5, 0.5)
	nothingSelectedLabel.SetJustify(cenums.JUSTIFY_CENTER)
	c.NothingSelectedFrame.Add(nothingSelectedLabel)

	// host entry panel

	c.HostSelectedFrame = ctk.NewFrame("")
	c.HostSelectedFrame.SetLabelAlign(0.0, 0.5)
	c.EditingHBox.PackStart(c.HostSelectedFrame, true, true, 0)

	panelVBox := ctk.NewVBox(false, 0)
	panelVBox.Show()
	c.HostSelectedFrame.Add(panelVBox)

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

	c.CommentsEntry = ctk.NewEntry("")
	c.CommentsEntry.SetName("editing-comment")
	c.CommentsEntry.Show()
	c.CommentsEntry.SetLineWrap(true)
	c.CommentsEntry.SetLineWrapMode(cenums.WRAP_NONE)
	c.CommentsEntry.SetSingleLineMode(false)
	c.CommentsEntry.SetSelectable(true)
	c.CommentsEntry.SetSizeRequest(-1, 3)
	panelVBox.PackStart(c.CommentsEntry, false, false, 0)

	c.HostEditVBox = ctk.NewVBox(false, 0)
	c.HostEditVBox.Show()
	panelVBox.PackStart(c.HostEditVBox, true, true, 0)

	addSeparator(c.HostEditVBox)
	addInstructions(c.HostEditVBox, "Enter an IP address (or nslookup domain):")

	addrBttnBox := ctk.NewHBox(true, 1)
	addrBttnBox.Show()
	addrBttnBox.SetSizeRequest(-1, 1)
	c.HostEditVBox.PackStart(addrBttnBox, false, false, 0)

	c.AddressEntry = ctk.NewEntry("")
	c.AddressEntry.Show()
	c.AddressEntry.SetSelectable(true)
	c.AddressEntry.SetLineWrap(false)
	c.AddressEntry.SetSizeRequest(-1, 1)
	c.AddressEntry.SetSingleLineMode(true)
	addrBttnBox.PackStart(c.AddressEntry, true, true, 0)

	c.AddressButton = ctk.NewButtonWithLabel("(address)")
	c.AddressButton.Show()
	c.AddressButton.SetSizeRequest(-1, 1)
	addrBttnBox.PackStart(c.AddressButton, true, true, 0)

	addSeparator(c.HostEditVBox)
	addInstructions(c.HostEditVBox, "Space separated list of domain names:")

	c.DomainsEntry = ctk.NewEntry("")
	c.DomainsEntry.Show()
	c.DomainsEntry.SetSelectable(true)
	c.DomainsEntry.SetSingleLineMode(false)
	// c.DomainsEntry.SetLineWrap(true)
	// c.DomainsEntry.SetLineWrapMode(cenums.WRAP_WORD)
	// c.DomainsEntry.SetJustify(cenums.JUSTIFY_LEFT)
	c.HostEditVBox.PackStart(c.DomainsEntry, true, true, 0)

	addSeparator(panelVBox)
	addInstructions(panelVBox, "Hosts file entry actions:")

	hostActionHBox := ctk.NewHBox(true, 1)
	hostActionHBox.Show()
	hostActionHBox.SetSizeRequest(-1, 1)
	panelVBox.PackStart(hostActionHBox, false, false, 0)

	c.ActivateButton = ctk.NewButtonWithLabel("")
	c.ActivateButton.Show()
	c.ActivateButton.SetSizeRequest(-1, 1)
	// panelVBox.PackStart(c.ActivateButton, false, false, 0)
	hostActionHBox.PackStart(c.ActivateButton, true, true, 0)

	c.DeleteButton = ctk.NewButtonWithLabel("click to delete")
	c.DeleteButton.Show()
	c.DeleteButton.SetSizeRequest(-1, 1)
	// panelVBox.PackStart(c.DeleteButton, true, true, 0)
	hostActionHBox.PackStart(c.DeleteButton, true, true, 0)

	changeSidebarMode(ListByDomain)

	return c.EditingHBox
}

func (c *CUI) reloadEditor() {
	c.Window.Freeze()
	defer func() {
		c.Window.Thaw()
		c.Window.Resize()
		c.Window.ReApplyStyles()
		c.Display.RequestDraw()
		c.Display.RequestSync()
	}()

	purgeContents := func(v ctk.VBox) {
		existing := v.GetChildren()
		for _, child := range existing {
			v.Remove(child)
			child.Destroy()
		}
	}

	c.EditorCommentList = []*editor.Host{}
	purgeContents(c.SidebarEntryList)
	purgeContents(c.SidebarLocalsList)
	purgeContents(c.SidebarCustomList)
	purgeContents(c.SidebarCommentsList)

	c.EditorAddressLookup = make(map[string]*editor.Host)
	c.EditorDomainsLookup = make(map[string]*editor.Host)

	var changed bool
	unique := make(map[string]int)
	for _, host := range c.HostFile.Hosts() {
		if !changed && host.Changed() {
			changed = true
		}
		if host.IsOnlyComment() {
			c.EditorCommentList = append(c.EditorCommentList, host)
			continue
		}
		c.EditorAddressLookup[host.Address()] = host
		if domains := host.Domains(); len(domains) == 0 {
			c.EditorDomainsLookup[""] = host
		} else {
			for _, domain := range domains {
				if _, found := unique[domain]; !found {
					unique[domain] = 1
				} else {
					unique[domain] += 1
				}
				if unique[domain] > 1 {
					c.EditorDomainsLookup[domain+" ("+strconv.Itoa(unique[domain])+")"] = host
				} else {
					c.EditorDomainsLookup[domain] = host
				}
			}
		}
	}

	c.SaveButton.SetSensitive(changed)
	c.ReloadButton.SetSensitive(changed)

	c.updateEditor()

}

func (c *CUI) updateEditor() {
	switch c.SidebarMode {
	case ListByDomain, ListByAddress:
		c.updateEditorByAddressOrDomain()
	case ListByEntry:
		c.updateEditorByEntry()
	}
}

func (c *CUI) updateEditorByAddressOrDomain() {
	var choices map[string]*editor.Host
	if c.SidebarMode == ListByAddress {
		choices = c.EditorAddressLookup
	} else {
		choices = c.EditorDomainsLookup
	}

	var localsCount, customCount int

	for idx, host := range c.EditorCommentList {
		key := fmt.Sprintf("Comment (%d)", idx+1)
		b := c.makeSidebarButton(key, host)
		c.SidebarCommentsList.PackStart(b, false, false, 0)
	}

	for _, key := range maps.SortedKeys(choices) {
		host := choices[key]
		b := c.makeSidebarButton(key, host)
		switch host.Importance() {
		case editor.HostIsLocalhostIPv4, editor.HostIsLocalhostIPv6:
			localsCount += 1
			c.SidebarLocalsList.PackStart(b, false, false, 0)
			c.SidebarLocalsList.SetSizeRequest(gSidebarInnerWidth, localsCount)
		default:
			customCount += 1
			c.SidebarCustomList.PackStart(b, false, false, 0)
			c.SidebarCustomList.SetSizeRequest(gSidebarInnerWidth, customCount)
		}
	}
}

func (c *CUI) updateEditorByEntry() {
	hosts := c.HostFile.Hosts()
	c.SidebarEntryList.SetSizeRequest(-1, len(hosts))
	var commentsCount int
	for idx, host := range hosts {
		key := strconv.Itoa(idx+1) + ". "
		if host.IsOnlyComment() {
			commentsCount += 1
			key += fmt.Sprintf("comment (%d)", commentsCount)
		} else {
			key += host.Address()
		}
		b := c.makeSidebarButton(key, host)
		c.SidebarEntryList.PackStart(b, false, false, 0)
	}
}

func (c *CUI) makeSidebarButton(key string, host *editor.Host) (b ctk.Button) {
	label := ctk.NewLabel(key)
	label.Show()
	label.SetJustify(cenums.JUSTIFY_LEFT)
	label.SetSizeRequest(-1, 1)
	label.SetSingleLineMode(true)

	b = ctk.NewButtonWithWidget(label)
	_ = b.InstallProperty(cdk.Property("host"), cdk.StructProperty, true, host)
	b.Show()
	b.SetSizeRequest(gSidebarInnerWidth, 1)
	b.Connect(ctk.SignalActivate, key+"-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		if h, ok := data[0].(*editor.Host); ok {
			c.focusEditor(h)
		}
		return cenums.EVENT_STOP
	}, host)

	var theme paint.Theme
	var name, tooltip string

	if c.SelectedHost != nil && c.SelectedHost.Equals(host) {
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

func (c *CUI) focusEditor(host *editor.Host) {
	c.Window.Freeze()
	defer func() {
		c.updateSidebarActionButtons()
		c.Window.Thaw()
		c.Window.Resize()
		c.Window.ReApplyStyles()
		c.Display.RequestDraw()
		c.Display.RequestShow()
		c.reloadEditor()
	}()

	if host == nil {
		c.Window.LogDebug("clearing editor focus")
		c.SelectedHost = nil
		c.HostSelectedFrame.Hide()
		c.NothingSelectedFrame.Show()
		return
	}

	c.Window.LogDebug("focusing editor on: %v", host.String())
	c.SelectedHost = host
	c.HostSelectedFrame.Show()
	c.NothingSelectedFrame.Hide()
	c.CommentsEntry.SetText(host.Comment())

	if v := host.Lookup(); v != "" {
		c.AddressEntry.SetText(v)
	} else {
		c.AddressEntry.SetText(host.Address())
	}

	_ = c.CommentsEntry.Disconnect(ctk.SignalChangedText, "comments-changed-handler")
	c.CommentsEntry.Connect(ctk.SignalChangedText, "comments-changed-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		h, _ := data[0].(*editor.Host)
		h.SetComment(c.CommentsEntry.GetText())
		c.CommentsEntry.LogDebug("updated host %v comment: %v", h.Address(), c.CommentsEntry.GetText())
		c.reloadEditor()
		return cenums.EVENT_STOP
	}, host)

	_ = c.DomainsEntry.Disconnect(ctk.SignalChangedText, "domains-changed-handler")
	// c.DomainsEntry.SetText(strings.Join(host.Domains(), " "))
	alloc := c.DomainsEntry.GetAllocation()
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
	c.DomainsEntry.SetText(strings.Join(domainLines, "\n"))
	c.DomainsEntry.Connect(ctk.SignalChangedText, "domains-changed-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		h, _ := data[0].(*editor.Host)
		h.SetDomains(c.DomainsEntry.GetText())
		c.reloadEditor()
		return cenums.EVENT_STOP
	}, host)

	handle := "activate-button-handler"
	_ = c.ActivateButton.Disconnect(ctk.SignalActivate, handle)
	if host.Importance() != editor.HostNotImportant {

		c.AddressEntry.SetSensitive(false)
		c.AddressButton.SetSensitive(false)
		c.ActivateButton.SetSensitive(false)

		c.ActivateButton.SetTheme(DefaultButtonTheme)
		c.ActivateButton.SetLabel("cannot deactivate host")

		c.DeleteButton.SetSensitive(false)
		c.DeleteButton.SetLabel("cannot delete host")
		_ = c.DeleteButton.Disconnect(ctk.SignalActivate, "delete-entry-handler")

	} else {

		c.AddressEntry.SetSensitive(true)
		addressEntryText := c.AddressEntry.GetText()
		c.AddressButton.SetSensitive(!cstrings.StringIsIP(addressEntryText) && cstrings.StringIsDomainName(addressEntryText))
		c.ActivateButton.SetSensitive(true)

		c.DeleteButton.SetSensitive(true)
		c.DeleteButton.SetLabel("click to delete")
		_ = c.DeleteButton.Disconnect(ctk.SignalActivate, "delete-entry-handler")
		c.DeleteButton.Connect(ctk.SignalActivate, "delete-entry-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
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
					for idx, hh := range c.HostFile.Hosts() {
						if h.Equals(hh) {
							log.DebugF("removing entry at index: %v", idx)
							c.HostFile.RemoveHost(idx)
							c.SelectedHost = nil
							break
						}
					}
					c.reloadEditor()
					c.focusEditor(nil)
				case enums.ResponseCancel, enums.ResponseClose, enums.ResponseNo:
					log.DebugF("user cancelled removal operation")
				}
			})
			return cenums.EVENT_STOP
		}, host)

		if host.Active() {

			c.ActivateButton.SetTheme(ActiveButtonTheme)
			c.ActivateButton.SetLabel("click to deactivate")
			c.ActivateButton.Connect(ctk.SignalActivate, handle, func(data []interface{}, argv ...interface{}) cenums.EventFlag {
				if h, ok := data[0].(*editor.Host); ok {
					h.SetActive(false)
					c.reloadEditor()
					c.focusEditor(h)
				}
				return cenums.EVENT_STOP
			}, host)

		} else {

			c.ActivateButton.SetTheme(DefaultButtonTheme)
			c.ActivateButton.SetLabel("click to activate")
			c.ActivateButton.Connect(ctk.SignalActivate, handle, func(data []interface{}, argv ...interface{}) cenums.EventFlag {
				if h, ok := data[0].(*editor.Host); ok {
					h.SetActive(true)
					c.reloadEditor()
					c.focusEditor(h)
				}
				return cenums.EVENT_STOP
			}, host)

		}
	}

	_ = c.AddressButton.Disconnect(ctk.SignalActivate, "address-activate-handler")
	c.AddressButton.Connect(
		ctk.SignalActivate,
		"address-activate-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			if err := c.newNsLookupDialog(host); err != nil {
				c.AddressButton.LogErr(err)
			}
			return cenums.EVENT_STOP
		},
		host,
	)

	actualLabel, actualTooltip := host.GetActualInfo()
	c.AddressButton.SetLabel(actualLabel)
	c.AddressButton.SetTooltipText(actualLabel + "\n" + actualTooltip)
	c.AddressButton.Resize()

	_ = c.AddressEntry.Disconnect(ctk.SignalChangedText, "address-text-changed-handler")
	c.AddressEntry.Connect(ctk.SignalChangedText, "address-text-changed-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		var h *editor.Host
		h, _ = data[0].(*editor.Host)
		text := h.Address()
		changed := c.AddressEntry.GetText()
		if cstrings.StringIsIP(changed) {
			if text != changed {
				c.AddressButton.SetLabel(fmt.Sprintf("(%v)", changed))
				c.AddressButton.SetTooltipText("is a valid IP address")
			} else {
				c.AddressButton.SetLabel(fmt.Sprintf("(%v)", text))
				c.AddressButton.SetTooltipText("is a valid IP address")
			}
			c.AddressButton.SetSensitive(false)
			h.SetAddress(changed)
		} else if cstrings.StringIsDomainName(changed) {
			if text != changed {
				c.AddressButton.SetLabel("(lookup changed)")
				c.AddressButton.SetTooltipText("click to perform domain lookup")
			}
			c.AddressButton.SetSensitive(true)
			h.SetAddress(changed)
			h.SetLookup(changed)
		} else {
			c.AddressButton.SetSensitive(false)
			c.AddressButton.SetLabel("(not ip or domain)")
			c.AddressButton.SetTooltipText("enter a valid address or domain name")
			h.SetAddress(changed)
			h.SetLookup("")
		}
		c.reloadEditor()
		return cenums.EVENT_PASS
	}, host)

	allEntries := append(
		append(
			c.SidebarLocalsList.GetChildren(),
			c.SidebarCustomList.GetChildren()...,
		),
		c.SidebarEntryList.GetChildren()...,
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
			switch c.SidebarMode {
			case ListByEntry:
				if c.SelectedHost != nil && c.SelectedHost.IsOnlyComment() {
					if h.Equals(c.SelectedHost) {
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
		c.CommentsEntry.SetSizeRequest(-1, -1)
		c.HostEditVBox.Hide()
		c.ActivateButton.Hide()
	} else {
		c.CommentsEntry.SetSizeRequest(-1, 3)
		c.HostEditVBox.Show()
		c.ActivateButton.Show()
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
