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
	"strings"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/cdk/lib/ptypes"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/ctk"

	editor "github.com/go-curses/coreutils-etc-hosts-editor"
)

/*

 Frame -[Controls]----------------------.
 | FrameVBox --------------------------.|
 | | Comment--------------------------.||
 | | |                                |||
 | | `--------------------------------'||
 | | FrameHBox -----------------------.||
 | | | InfoVBox ------. Domains -----.|||
 | | | | Address ---- | |            ||||
 | | | | Actual ----- | |            ||||
 | | | `--------------' `------------'|||
 | | `--------------------------------'||
 | `-----------------------------------'|
 `--------------------------------------'

*/

type ViewerRow struct {
	Host   *editor.Host
	Active bool

	Size *ptypes.Rectangle

	Frame     ctk.Frame
	FrameVBox ctk.VBox
	Comment   ctk.Entry
	FrameHBox ctk.HBox
	InfoVBox  ctk.VBox
	Actual    ctk.Button
	Address   ctk.Entry
	Domains   ctk.Entry

	Controls *ViewerRowControls

	Eheditor *CEheditor
}

func NewViewerRow(eheditor *CEheditor, host *editor.Host, viewerWidth int) (row *ViewerRow) {
	row = new(ViewerRow)
	row.Eheditor = eheditor

	row.Host = host
	row.Active = host.Active()

	row.Controls = newViewerRowControls(eheditor, row)
	row.Frame = ctk.NewFrameWithWidget(row.Controls.HBox)
	row.Frame.Show()
	row.Frame.SetLabelAlign(0.5, 0.5)
	_ = row.Frame.InstallProperty("ViewerRow", cdk.StructProperty, true, row)

	row.FrameVBox = ctk.NewVBox(true, 0)
	row.FrameVBox.Show()
	row.Frame.Add(row.FrameVBox)

	row.Comment = ctk.NewEntry("")
	row.Comment.SetName("comment")
	row.Comment.Show()
	row.Comment.SetSingleLineMode(false)
	row.Comment.SetLineWrapMode(cenums.WRAP_CHAR)
	row.Comment.SetTooltipText("(host comments)")
	row.Comment.SetHasTooltip(true)
	row.Comment.SetSelectable(true)
	row.Comment.SetSizeRequest(-1, 2)
	row.FrameVBox.PackStart(row.Comment, true, true, 0)

	row.FrameHBox = ctk.NewHBox(false, 0)
	row.FrameHBox.Show()
	row.FrameVBox.PackEnd(row.FrameHBox, false, false, 0)

	row.InfoVBox = ctk.NewVBox(false, 0)
	row.InfoVBox.Show()
	row.FrameHBox.PackStart(row.InfoVBox, true, true, 0)

	row.Address = ctk.NewEntry("")
	row.Address.SetName("address")
	row.Address.Show()
	row.Address.SetSizeRequest(-1, 1)
	row.Address.SetSingleLineMode(true)
	row.Address.SetLineWrapMode(cenums.WRAP_CHAR)
	row.Address.SetTooltipText("(address or domain to lookup)")
	row.InfoVBox.PackStart(row.Address, false, false, 0)

	row.Actual = ctk.NewButtonWithLabel("")
	row.Actual.SetName("actual")
	row.Actual.Show()
	row.Actual.SetHasTooltip(true)
	row.Actual.SetSizeRequest(-1, 1)
	row.InfoVBox.PackStart(row.Actual, false, false, 0)

	row.Domains = ctk.NewEntry("")
	row.Domains.SetName("domains")
	row.Domains.Show()
	row.Domains.SetSingleLineMode(false)
	row.Domains.SetLineWrap(true)
	row.Domains.SetLineWrapMode(cenums.WRAP_CHAR)
	row.Domains.SetJustify(cenums.JUSTIFY_NONE)
	row.Domains.SetTooltipText("(space or newline separated list of domain names)")
	row.FrameHBox.PackStart(row.Domains, true, true, 0)
	return
}

func (row *ViewerRow) IsCommentRow() bool {
	return row.Host.IsOnlyComment()
}

func (row *ViewerRow) Resize() {
	row.Frame.Resize()
	row.Controls.HBox.Resize()
}

func (row *ViewerRow) Update(host *editor.Host, viewerWidth int) {
	row.Frame.Freeze()
	if host.IsComment() {
		row.updateToComment(host, viewerWidth)
	} else {
		row.updateToHost(host, viewerWidth)
	}
	row.Frame.Thaw()
}

func (row *ViewerRow) updateToComment(host *editor.Host, viewerWidth int) {

	row.Host = host

	_, _, _, _, size := calcViewerSizes(viewerWidth)

	row.Size = size
	size = size.NewClone()

	row.Frame.SetSizeRequest(size.W, size.H)
	size.W -= 1
	size.H -= 1

	row.Frame.SetName("comment")
	row.Comment.SetSizeRequest(size.W, size.H)
	row.Comment.SetText(row.Host.Comment())

	row.Eheditor.Window.ApplyStylesTo(row.Comment)
	row.Comment.Invalidate()

	row.FrameHBox.Hide()

	_ = row.Comment.Disconnect(ctk.SignalChangedText, "comment-text-changed-handler")
	row.Comment.SetText(host.Comment())
	row.Comment.Connect(ctk.SignalChangedText, "comment-text-changed-handler", row.processCommentTextChanged)

	row.Controls.Update(row)
}

func (row *ViewerRow) processCommentTextChanged(data []interface{}, argv ...interface{}) cenums.EventFlag {
	row.Host.SetComment(row.Comment.GetText())
	row.Controls.UpdateToggle()
	row.Frame.Resize()
	return cenums.EVENT_PASS
}

func (row *ViewerRow) processDomainsTextChanged(data []interface{}, argv ...interface{}) cenums.EventFlag {
	row.Host.SetDomains(row.Domains.GetText())
	row.Controls.UpdateToggle()
	return cenums.EVENT_PASS
}

func (row *ViewerRow) updateToHost(host *editor.Host, viewerWidth int) {

	row.Host = host
	row.FrameHBox.Show()

	_, rightColumnWidth, actualWidth, addressWidth, size := calcViewerSizes(viewerWidth)

	row.Size = size
	size = size.NewClone()

	row.Frame.SetSizeRequest(size.W, size.H)
	size.W -= 1
	size.H -= 1

	row.FrameHBox.SetSizeRequest(size.W, 4)
	row.Comment.SetSizeRequest(-1, 2)
	row.InfoVBox.SetSizeRequest(rightColumnWidth, 4)
	row.Address.SetSizeRequest(addressWidth, 1)
	row.Actual.SetSizeRequest(actualWidth, 1)
	row.Domains.SetSizeRequest(rightColumnWidth, 2)

	if host.Active() {
		row.Frame.SetName("host-active")
	} else {
		row.Frame.SetName("host-inactive")
	}

	row.Eheditor.Window.ApplyStylesTo(row.Comment)
	row.Comment.Invalidate()

	name := host.Name()
	nameIsIP := cstrings.StringIsIP(name)
	row.Actual.SetSensitive(!nameIsIP)

	actualLabel, actualTooltip := host.GetActualInfo()
	row.Actual.SetLabel(actualLabel)
	row.Actual.SetTooltipText(actualTooltip)
	row.Actual.Resize()

	row.Address.SetText(name)

	_ = row.Actual.Disconnect(ctk.SignalActivate, "actual-activate-handler")
	row.Actual.Connect(ctk.SignalActivate, "actual-activate-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		if err := row.Eheditor.newNsLookupDialog(row.Host); err != nil {
			row.Actual.LogErr(err)
		}
		return cenums.EVENT_STOP
	})

	_ = row.Address.Disconnect(ctk.SignalChangedText, "address-text-changed-handler")
	row.Address.Connect(ctk.SignalChangedText, "address-text-changed-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		text := row.Host.Address()
		changed := row.Address.GetText()
		if row.Host.Active() {
			row.Frame.SetName("host-active")
		} else {
			row.Frame.SetName("host-inactive")
		}
		if cstrings.StringIsIP(changed) {
			if text != changed {
				row.Actual.SetLabel(fmt.Sprintf("(%v)", changed))
				row.Actual.SetTooltipText("is a valid IP address")
			} else {
				row.Actual.SetLabel(fmt.Sprintf("(%v)", text))
				row.Actual.SetTooltipText("is a valid IP address")
			}
			row.Actual.SetSensitive(false)
			row.Host.SetAddress(changed)
		} else if cstrings.StringIsDomainName(changed) {
			if text != changed {
				row.Actual.SetLabel("(lookup changed)")
				row.Actual.SetTooltipText("click to perform domain lookup")
			}
			row.Actual.SetSensitive(true)
			row.Host.SetAddress(changed)
			row.Host.SetLookup(changed)
		} else {
			row.Actual.SetSensitive(false)
			row.Actual.SetLabel("(not ip or domain)")
			row.Actual.SetTooltipText("enter a valid address or domain name")
			row.Host.SetAddress(changed)
			row.Host.SetLookup("")
		}
		row.Controls.UpdateToggle()
		return cenums.EVENT_PASS
	})

	_ = row.Comment.Disconnect(ctk.SignalChangedText, "comment-text-changed-handler")
	row.Comment.SetText(host.Comment())
	row.Comment.Connect(ctk.SignalChangedText, "comment-text-changed-handler", row.processCommentTextChanged)

	_ = row.Domains.Disconnect(ctk.SignalChangedText, "domains-text-changed-handler")
	row.Domains.SetText(strings.Join(host.Domains(), " "))
	row.Domains.Connect(ctk.SignalChangedText, "domains-text-changed-handler", row.processDomainsTextChanged)

	row.Controls.Update(row)
}

func calcViewerSizes(viewerWidth int) (left, right, actual, address int, size *ptypes.Rectangle) {
	left = cmath.FloorI(viewerWidth/2, 1)
	if left%2 != 0 {
		left -= 1
	}
	right = cmath.FloorI(viewerWidth-left-1, 1)
	actual = left / 2
	address = left - actual
	size = ptypes.NewRectangle(left+right+2, 6)
	return
}

func getViewerRowFromFrame(frame ctk.Frame) (row *ViewerRow) {
	if v, err := frame.GetStructProperty("ViewerRow"); err == nil {
		if r, ok := v.(*ViewerRow); ok {
			row = r
		}
	}
	return
}