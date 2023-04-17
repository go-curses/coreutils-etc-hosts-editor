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
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
	"github.com/go-curses/ctk/lib/enums"

	editor "github.com/go-curses/coreutils-etc-hosts-editor"
)

type ViewerRowControls struct {
	Row      *ViewerRow
	HBox     ctk.HBox
	Toggle   ctk.Button
	MoveUp   ctk.Button
	MoveDn   ctk.Button
	AddEntry ctk.Button
	DelEntry ctk.Button

	Eheditor *CEheditor
}

func newViewerRowControls(eheditor *CEheditor, row *ViewerRow) (ctrls *ViewerRowControls) {
	ctrls = new(ViewerRowControls)
	ctrls.Eheditor = eheditor

	ctrls.Row = row

	ctrls.HBox = ctk.NewHBox(false, 1)
	ctrls.HBox.Show()
	theme := ctrls.HBox.GetTheme()
	theme.Content.FillRune = rune(0)
	ctrls.HBox.SetTheme(theme)
	ctrls.HBox.SetSizeRequest(20, 1)
	// _ = ctrls.FrameHBox.SetBoolProperty(ctk.PropertyDebug, true)

	arrowUp := ctk.NewArrow(enums.ArrowUp)
	arrows, _ := paint.GetArrows(paint.WideArrow)
	arrowUp.SetArrowRuneSet(arrows)
	arrowUp.Show()
	ctrls.MoveUp = ctk.NewButtonWithWidget(arrowUp)
	ctrls.MoveUp.Show()
	ctrls.MoveUp.SetSizeRequest(-1, 1)
	ctrls.MoveUp.SetTooltipText("Click to move host up in the list")
	ctrls.HBox.PackStart(ctrls.MoveUp, false, false, 0)

	arrowDown := ctk.NewArrow(enums.ArrowDown)
	arrowDown.SetArrowRuneSet(arrows)
	arrowDown.Show()
	ctrls.MoveDn = ctk.NewButtonWithWidget(arrowDown)
	ctrls.MoveDn.Show()
	ctrls.MoveDn.SetSizeRequest(-1, 1)
	ctrls.MoveUp.SetTooltipText("Click to move host down in the list")
	ctrls.HBox.PackStart(ctrls.MoveDn, false, false, 0)

	ctrls.Toggle = ctk.NewButtonWithLabel("")
	ctrls.Toggle.Show()
	ctrls.HBox.PackStart(ctrls.Toggle, false, false, 0)

	ctrls.AddEntry = ctk.NewButtonWithLabel("+")
	ctrls.AddEntry.Show()
	ctrls.AddEntry.SetTooltipText("click to add a new host or comment entry above")
	ctrls.AddEntry.SetHasTooltip(true)
	ctrls.HBox.PackEnd(ctrls.AddEntry, false, false, 0)

	ctrls.DelEntry = ctk.NewButtonWithLabel("-")
	ctrls.DelEntry.Show()
	if ctrls.Row.Host.IsOnlyComment() {
		ctrls.DelEntry.SetTooltipText("click to delete this comment entry")
	} else {
		ctrls.DelEntry.SetTooltipText("click to delete this host entry")
	}
	ctrls.DelEntry.SetHasTooltip(true)
	ctrls.HBox.PackEnd(ctrls.DelEntry, false, false, 0)

	ctrls.Update(row)
	return
}

func (ctrls *ViewerRowControls) UpdateToggle() {
	if ctrls.Row.IsCommentRow() {
		commentLabel := "Comment Only"
		if ctrls.Row.Host.Changed() {
			commentLabel = "* " + commentLabel
		}
		ctrls.Toggle.SetLabel(commentLabel)
		ctrls.Toggle.SetHasTooltip(false)
		ctrls.Toggle.SetSizeRequest(len(commentLabel)+2, 1)
	} else {
		stateLabel := "Host (inactive)"
		stateTooltip := "Click to set active (uncomment)"
		if ctrls.Row.Host.Active() {
			stateLabel = "Host (active)"
			stateTooltip = "Click to set inactive (comment out)"
		}
		if ctrls.Row.Host.Changed() {
			stateLabel = "* " + stateLabel
		}
		ctrls.Toggle.SetLabel(stateLabel)
		ctrls.Toggle.SetTooltipText(stateTooltip)
		ctrls.Toggle.SetHasTooltip(true)
		ctrls.Toggle.SetSizeRequest(len(stateLabel)+2, 1)
	}
}

func (ctrls *ViewerRowControls) Update(row *ViewerRow) {
	ctrls.Row = row

	idx := ctrls.Eheditor.HostFile.IndexOf(row.Host)
	last := ctrls.Eheditor.HostFile.Len() - 1

	ctrls.UpdateToggle()

	if row.Host.IsComment() {
		ctrls.Toggle.SetSensitive(false)
	} else {
		ctrls.Toggle.SetSensitive(true)

		_ = ctrls.Toggle.Disconnect(ctk.SignalActivate, "toggle-row-state-handler")
		ctrls.Toggle.Connect(ctk.SignalActivate, "toggle-row-state-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			row.Host.SetActive(!row.Host.Active())
			ctrls.UpdateToggle()
			ctrls.Eheditor.updateViewer()
			return cenums.EVENT_STOP
		})
	}

	if idx == 0 {
		ctrls.MoveUp.SetSensitive(false)
		ctrls.MoveUp.SetHasTooltip(false)
	} else if !ctrls.MoveUp.IsSensitive() {
		ctrls.MoveUp.SetSensitive(true)
		ctrls.MoveUp.SetHasTooltip(true)
	}

	_ = ctrls.MoveUp.Disconnect(ctk.SignalActivate, "move-row-up-handler")
	ctrls.MoveUp.Connect(ctk.SignalActivate, "move-row-up-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		ctrls.Eheditor.HostFile.MoveHost(idx, idx-1)
		ctrls.Eheditor.updateViewer()
		return cenums.EVENT_STOP
	})

	if idx == last {
		ctrls.MoveDn.SetSensitive(false)
		ctrls.MoveDn.SetHasTooltip(false)
	} else if !ctrls.MoveDn.IsSensitive() {
		ctrls.MoveDn.SetSensitive(true)
		ctrls.MoveDn.SetHasTooltip(true)
	}

	_ = ctrls.MoveDn.Disconnect(ctk.SignalActivate, "move-row-down-handler")
	ctrls.MoveDn.Connect(ctk.SignalActivate, "move-row-down-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		ctrls.Eheditor.HostFile.MoveHost(idx, idx+1)
		ctrls.Eheditor.updateViewer()
		return cenums.EVENT_STOP
	})

	_ = ctrls.AddEntry.Disconnect(ctk.SignalActivate, "add-row-handler")
	ctrls.AddEntry.Connect(ctk.SignalActivate, "add-row-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		// ctrls.Eheditor.HostFile.InsertHost(editor.NewHostFromInfo(editor.HostInfo{}), idx-1)

		var entries []interface{}
		entries = append(entries, "Comment Entry", 1)
		entries = append(entries, "Host Entry", 2)
		// add presets that don't already exist, ie localhost

		dialog := ctk.NewButtonMenuDialog(
			"Add Entry",
			"Select a type of entry to add:",
			entries...,
		)
		dialog.SetSizeRequest(30, 8)
		dialog.RunFunc(func(response enums.ResponseType, argv ...interface{}) {
			switch response {
			case 1: // add comment
				ctrls.AddEntry.LogDebug("add comment at index: %v", idx)
				ctrls.Eheditor.HostFile.InsertHost(editor.NewComment(""), idx)
				ctrls.Eheditor.reloadViewer()
			case 2: // add host
				ctrls.AddEntry.LogDebug("add host at index: %v", idx)
				ctrls.Eheditor.HostFile.InsertHost(editor.NewHostFromInfo(editor.HostInfo{}), idx)
				ctrls.Eheditor.reloadViewer()
			default:
				if idx := int(response); idx < 0 {
					ctrls.AddEntry.LogDebug("new entry action cancelled")
				} else {
					log.ErrorF("unhandled dialog response: %v", response)
				}
			}
		})

		return cenums.EVENT_STOP
	})

	_ = ctrls.DelEntry.Disconnect(ctk.SignalActivate, "del-row-handler")
	ctrls.DelEntry.Connect(ctk.SignalActivate, "del-row-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		message := ""
		if ctrls.Row.Host.Empty() {
			if ctrls.Row.Host.IsOnlyComment() {
				message = "(empty comment entry)"
			} else {
				message = "(empty host entry)"
			}
		} else {
			message = ctrls.Row.Host.Block()
		}
		ctk.NewYesNoDialog("Remove Entry?", message, true).
			RunFunc(func(response enums.ResponseType, argv ...interface{}) {
				switch response {
				case enums.ResponseYes:
					log.DebugF("removing entry at index: %v", idx)
					ctrls.Eheditor.HostFile.RemoveHost(idx)
					ctrls.Eheditor.reloadViewer()
				case enums.ResponseCancel, enums.ResponseClose, enums.ResponseNo:
					log.DebugF("user cancelled removal operation")
				}
			})
		return cenums.EVENT_STOP
	})
}