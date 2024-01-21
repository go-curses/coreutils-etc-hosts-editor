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

package editor

import (
	"fmt"

	cpaths "github.com/go-curses/cdk/lib/paths"
	"github.com/go-curses/cdk/lib/sync"
	"github.com/go-curses/cdk/log"
)

var (
	HostValidations = map[string]string{
		"localhost":      "IPv4 (localhost)",
		"ip6-localhost":  "IPv6 (localhost)",
		"ip6-loopback":   "IPv6 (loopback)",
		"ip6-allnodes":   "IPv6 (all nodes)",
		"ip6-allrouters": "IPv6 (all routers)",
	}
)

type Hostfile struct {
	Path    string
	hosts   []*Host
	Comment string

	sync.RWMutex
}

func (eh *Hostfile) Hosts() []*Host {
	eh.RLock()
	defer eh.RUnlock()
	return eh.hosts
}

func (eh *Hostfile) Save() (err error) {
	if cpaths.FileWritable(eh.Path) {
		content := eheditorFileHeading + "\n"
		eh.RLock()
		defer eh.RUnlock()
		for _, host := range eh.hosts {
			content += "\n"
			content += host.Block()
		}
		err = cpaths.WriteFile(eh.Path, content)
	} else {
		err = fmt.Errorf("%v is not writable", eh.Path)
	}
	return
}

func (eh *Hostfile) Validate() (errs []error) {
	for required, label := range HostValidations {
		found := false
		for _, host := range eh.hosts {
			if host.HasDomain(required) {
				found = true
				break
			}
		}
		if !found {
			errs = append(errs, fmt.Errorf("missing required entry for %v", label))
		}
	}
	return
}

func (eh *Hostfile) Len() int {
	eh.RLock()
	defer eh.RUnlock()
	return len(eh.hosts)
}

func (eh *Hostfile) IndexOf(host *Host) int {
	eh.RLock()
	defer eh.RUnlock()
	if host != nil {
		for i, h := range eh.hosts {
			if h.Equals(host) {
				return i
			}
		}
	}
	return -1
}

func (eh *Hostfile) InsertHost(host *Host, idx int) {
	eh.Lock()
	eh.hosts = eh.insertHost(eh.hosts, host, idx)
	eh.Unlock()
}

func (eh *Hostfile) insertHost(hosts []*Host, host *Host, idx int) []*Host {
	temp := append([]*Host{}, hosts...)
	if idx >= 0 && idx < len(temp) {
		return append(temp[:idx], append([]*Host{host}, temp[idx:]...)...)
	}
	return append(temp, host)
}

func (eh *Hostfile) RemoveHost(idx int) {
	eh.Lock()
	eh.hosts = eh.removeHost(eh.hosts, idx)
	eh.Unlock()
}

func (eh *Hostfile) removeHost(hosts []*Host, idx int) []*Host {
	temp := append([]*Host{}, hosts...)
	if idx >= 0 && idx < len(temp) {
		return append(temp[:idx], temp[idx+1:]...)
	}
	return temp
}

func (eh *Hostfile) MoveHost(from, to int) {
	eh.Lock()
	eh.hosts = eh.moveHost(eh.hosts, from, to)
	eh.Unlock()
}

func (eh *Hostfile) moveHost(hosts []*Host, from, to int) (updated []*Host) {
	temp := append([]*Host{}, hosts...)
	total := len(temp)
	if from >= 0 && from < total {
		tgt := temp[from]
		if to >= 0 && to < total {
			modified := eh.removeHost(temp, from)
			updated = eh.insertHost(modified, tgt, to)
			log.DebugF("tgt=%v, from=%v, to=%v,\ntemp=%v\nmodified=%v\nupdated=%v", tgt, from, to, temp, modified, updated)
			return
		}
	}
	updated = temp
	return
}
