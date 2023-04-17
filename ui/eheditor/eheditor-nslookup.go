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
	"net"

	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
	"github.com/go-curses/ctk/lib/enums"

	"github.com/go-curses/coreutils-etc-hosts-editor"
)

func (e *CEheditor) newNsLookupDialog(host *editor.Host) (err error) {
	e.HostsVBox.Freeze()
	e.EditingHBox.Freeze()

	var found []net.IP
	if found, err = host.PerformLookup(); err != nil {
		e.HostsVBox.Thaw()
		e.EditingHBox.Thaw()
		return err
	}

	numFound := len(found)
	if numFound == 0 {
		ctk.NewMessageDialog("nslookup", fmt.Sprintf("No hosts found for domain:\n%v", host.Lookup()))
		e.HostsVBox.Thaw()
		e.EditingHBox.Thaw()
		return fmt.Errorf("domain hosts not found")
	}

	var options []interface{}

	for idx, ip := range found {
		options = append(options, ip.String())
		options = append(options, enums.ResponseType(idx+1))
	}

	dialog := ctk.NewButtonMenuDialog(
		"Select an IP for: "+host.Lookup(),
		"",
		options...,
	)
	dialog.SetSizeRequest(42, 10)
	dialog.RunFunc(func(response enums.ResponseType, argv ...interface{}) {
		e.HostsVBox.Thaw()
		e.EditingHBox.Thaw()
		if len(argv) >= 2 {
			h, _ := argv[0].(*editor.Host)
			if available, ok := argv[1].([]net.IP); ok {
				if idx := int(response); idx > 0 {
					ip := available[idx-1]
					log.DebugF("selected ip: %v (idx=%v,found=%v)", ip.String(), idx, found)
					h.SetAddress(ip.String())
					e.reloadContents()
					e.focusEditor(h)
				} else {
					log.DebugF("ip selection cancelled")
				}
			}

		}
	}, host, found)

	return
}