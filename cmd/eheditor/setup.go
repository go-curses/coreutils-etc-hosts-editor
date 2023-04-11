package main

import (
	"fmt"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paths"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/log"
	editor "github.com/go-curses/coreutils-etc-hosts-editor"
	"github.com/go-curses/ctk"
	"github.com/go-curses/ctk/lib/enums"
)

var (
	gEH            *editor.Hostfile
	gSrc                                = ""
	gLastError     error                = nil
	gReadOnlyMode                       = false
	gLeftSep       ctk.Separator        = nil
	gRightSep      ctk.Separator        = nil
	gHostsViewport ctk.ScrolledViewport = nil
	gHostsVBox     ctk.VBox             = nil
	gHostsHBox     ctk.HBox             = nil
	gActionHBox    ctk.HBox             = nil
	gApp           ctk.Application      = nil
	gDisplay       cdk.Display          = nil
	gWindow        ctk.Window           = nil
	gSaveButton    ctk.Button           = nil
	gReloadButton  ctk.Button           = nil
	gQuitButton    ctk.Button           = nil
)

func startup(data []interface{}, argv ...interface{}) cenums.EventFlag {
	var err error
	var ok bool
	if gApp, gDisplay, _, _, _, ok = ctk.ArgvApplicationSignalStartup(argv...); ok {

		if gDisplay.App().GetContext().NArg() > 0 {
			gSrc = gDisplay.App().GetContext().Args().First()
		}
		if gSrc == "" {
			gSrc = "/etc/hosts"
		}

		if !paths.IsFile(gSrc) {
			gLastError = fmt.Errorf("%v not found or not a file", gSrc)
			log.Error(gLastError)
			return cenums.EVENT_STOP
		}
		gReadOnlyMode = gDisplay.App().GetContext().Bool("read-only")
		if !gReadOnlyMode && !paths.FileWritable(gSrc) {
			log.WarnF("etc hosts file %v is not writable, read-only mode", gSrc)
			gReadOnlyMode = true
		}

		// gDisplay.CaptureCtrlC()
		if s := gDisplay.Screen(); s != nil {
			s.EnableHostClipboard(true)
		}

		screenSize := ptypes.MakeRectangle(gDisplay.Screen().Size())
		if screenSize.W < 60 || screenSize.H < 14 {
			gLastError = fmt.Errorf("eheditor requires a terminal with at least 80x24 dimensions")
			log.Error(gLastError)
			return cenums.EVENT_STOP
		}

		title := fmt.Sprintf("%s - eheditor (v%v)", gSrc, gApp.Version())
		if gReadOnlyMode {
			title += " [read-only]"
		}

		if gEH, err = editor.ParseFile(gSrc); err != nil {
			gLastError = fmt.Errorf("error parsing %v: %v", gSrc, err)
			log.Error(gLastError)
			return cenums.EVENT_STOP
		}

		ctk.GetAccelMap().LoadFromString(eheditorAccelMap)

		gWindow = ctk.NewWindowWithTitle(title)
		gWindow.SetName("eheditor-window")
		// gWindow.Show()
		gWindow.SetTheme(WindowTheme)
		// gWindow.SetDecorated(false)
		// _ = gWindow.SetBoolProperty(ctk.PropertyDebug, true)
		if err := gWindow.ImportStylesFromString(eheditorStyles); err != nil {
			gWindow.LogErr(err)
		}

		ag := ctk.NewAccelGroup()
		ag.ConnectByPath(
			"<eheditor-window>/File/Quit",
			"quit-accel",
			func(argv ...interface{}) (handled bool) {
				ag.LogDebug("quit-accel called")
				requestQuit()
				return
			},
		)
		ag.ConnectByPath(
			"<eheditor-window>/File/Reload",
			"reload-accel",
			func(argv ...interface{}) (handled bool) {
				ag.LogDebug("reload-accel called")
				requestReload()
				return
			},
		)
		ag.ConnectByPath(
			"<eheditor-window>/File/Save",
			"save-accel",
			func(argv ...interface{}) (handled bool) {
				ag.LogDebug("save-accel called")
				requestSave()
				return
			},
		)
		gWindow.AddAccelGroup(ag)

		vbox := gWindow.GetVBox()
		vbox.SetSpacing(1)

		gHostsHBox = ctk.NewHBox(false, 0)
		gHostsHBox.Show()
		// _ = gHostsHBox.SetBoolProperty(cdk.PropertyDebug, true)
		// _ = gHostsHBox.SetBoolProperty(ctk.PropertyDebugChildren, true)
		vbox.PackStart(gHostsHBox, true, true, 0)

		gLeftSep = ctk.NewSeparator()
		gLeftSep.Show()
		gHostsHBox.PackStart(gLeftSep, false, true, 0)

		gHostsViewport = ctk.NewScrolledViewport()
		gHostsViewport.SetTheme(ViewerTheme)
		gHostsViewport.Show()
		gHostsViewport.SetPolicy(enums.PolicyAutomatic, enums.PolicyAutomatic)
		gHostsVBox = ctk.NewVBox(false, 1)
		gHostsVBox.Show()
		// _ = gHostsVBox.SetBoolProperty(ctk.PropertyDebug, true)
		gHostsViewport.Add(gHostsVBox)
		// gHostsVBox.SetBoolProperty(cdk.PropertyDebug, true)
		// gHostsViewport.SetBoolProperty(cdk.PropertyDebug, true)
		gHostsHBox.PackStart(gHostsViewport, true, true, 0)

		gRightSep = ctk.NewSeparator()
		gRightSep.Show()
		gHostsHBox.PackStart(gRightSep, false, false, 0)

		gActionHBox = ctk.NewHBox(false, 1)
		gActionHBox.Show()
		vbox.PackEnd(gActionHBox, false, true, 0)

		actionSep := ctk.NewSeparator()
		actionSep.Show()
		gActionHBox.PackStart(actionSep, true, true, 0)

		gSaveButton = ctk.NewButtonWithMnemonic("_Save <F3>")
		gSaveButton.Show()
		if gReadOnlyMode {
			gSaveButton.SetSensitive(false)
		} else {
			gSaveButton.Connect(ctk.SignalActivate, "save-hosts", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
				requestSave()
				return cenums.EVENT_STOP
			})
		}
		gActionHBox.PackEnd(gSaveButton, false, false, 0)

		gReloadButton = ctk.NewButtonWithMnemonic("_Reload <F5>")
		gReloadButton.Show()
		gReloadButton.Connect(ctk.SignalActivate, "reload-hosts", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			requestReload()
			return cenums.EVENT_STOP
		})
		gActionHBox.PackEnd(gReloadButton, false, false, 0)

		gQuitButton = ctk.NewButtonWithMnemonic("_Quit <F10>")
		gQuitButton.Show()
		gQuitButton.Connect(ctk.SignalActivate, "quit-hosts", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			requestQuit()
			return cenums.EVENT_STOP
		})
		gActionHBox.PackEnd(gQuitButton, false, false, 0)

		gApp.NotifyStartupComplete()
		gWindow.Show()
		gHostsViewport.GrabFocus()
		reloadViewer()
		gDisplay.Connect(cdk.SignalEventResize, "display-resize-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			if gApp.StartupCompleted() {
				updateViewer()
			}
			return cenums.EVENT_PASS
		})
		return cenums.EVENT_PASS
	}
	return cenums.EVENT_STOP
}

func shutdown(_ []interface{}, _ ...interface{}) cenums.EventFlag {
	if gLastError != nil {
		fmt.Printf("%v\n", gLastError)
		log.InfoF("exiting (with error)")
	} else {
		log.InfoF("exiting (without error)")
	}
	return cenums.EVENT_PASS
}

func requestReload() {
	log.DebugF("reloading from: %v", gSrc)
	var err error
	if gEH, err = editor.ParseFile(gSrc); err != nil {
		gLastError = fmt.Errorf("error parsing %v: %v", gSrc, err)
		log.Error(gLastError)
		return
	}
	reloadViewer()
}

func requestSave() {
	log.DebugF("saving to: %v", gSrc)
	if gEH != nil {
		if err := gEH.Save(); err != nil {
			log.Error(err)
		}
	}
	requestReload()
	gQuitButton.GrabFocus()
}

func requestQuit() {
	gDisplay.RequestQuit()
}