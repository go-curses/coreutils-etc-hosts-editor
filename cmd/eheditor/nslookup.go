package main

import (
	"fmt"
	"net"

	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
	"github.com/go-curses/ctk/lib/enums"

	editor "github.com/go-curses/coreutils-etc-hosts-editor"
)

func newNsLookupDialog(host *editor.Host) (err error) {
	gHostsVBox.Freeze()

	var found []net.IP
	if found, err = host.PerformLookup(); err != nil {
		gHostsVBox.Thaw()
		return err
	}

	numFound := len(found)
	if numFound == 0 {
		ctk.NewMessageDialog("nslookup", fmt.Sprintf("No hosts found for domain:\n%v", host.Lookup()))
		return fmt.Errorf("domain hosts not found")
	}

	var options []interface{}

	for idx, ip := range found {
		options = append(options, ip.String())
		options = append(options, enums.ResponseType(idx+1))
	}

	dialog := ctk.NewButtonMenuDialog(
		"Select an address",
		fmt.Sprintf("%d addresses found for:\n%v", numFound, host.Lookup()),
		options...,
	)

	dialog.RunFunc(func(response enums.ResponseType, argv ...interface{}) {
		if len(argv) >= 1 {
			if available, ok := argv[0].([]net.IP); ok {
				gHostsVBox.Thaw()
				if idx := int(response); idx > 0 {
					ip := available[idx-1]
					log.DebugF("selected ip: %v (idx=%v,found=%v)", ip.String(), idx, found)
					host.SetAddress(ip.String())
					updateViewer()
				} else {
					log.DebugF("ip selection cancelled")
				}
			}
		}
	}, found)

	return
}