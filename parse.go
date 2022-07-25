package editor

import (
	"regexp"
	"strings"

	"github.com/go-curses/cdk/lib/paths"
	"github.com/go-curses/cdk/log"
)

const (
	eheditorFileHeading = "## eheditor wrote this file"
)

var (
	rxCommentLine = regexp.MustCompile(`^\s*#+\s*([^#].+?)\s*$`)
	rxLookupLine  = regexp.MustCompile(`^\s*#nslookup ([a-zA-Z][-_.a-zA-Z\d]+?)\s*$`)
	rxUnHostLine  = regexp.MustCompile(`^\s*#+\s*([:a-f\d][:.a-fA-F\d]+?)\s+(.+?)\s*$`)
	rxHostLine    = regexp.MustCompile(`^\s*([^#][:.a-fA-F\d]+?)\s+(.+?)\s*$`)
	rxEmptyLine   = regexp.MustCompile(`^\s*$`)
	rxSpaceSep    = regexp.MustCompile(`\s+`)
	rxNewlines    = regexp.MustCompile(`\r??\n`)
)

func ParseFile(path string) (eh *Hostfile, err error) {
	var contents string
	if contents, err = paths.ReadFile(path); err != nil {
		return
	}
	eh = new(Hostfile)
	eh.Path = path
	eh.hosts = make([]*Host, 0)
	lines := strings.Split(contents, "\n")
	if lines[0] == eheditorFileHeading {
		if err := parseOwnFile(lines[1:], eh); err != nil {
			return nil, err
		}
	} else if err := parseOtherFile(lines, eh); err != nil {
		return nil, err
	}
	return eh, nil
}

func processCommentBlocks(eh *Hostfile) (err error) {
	for _, host := range eh.Hosts() {
		if host.address == "" && len(host.domains) == 0 {
			host.onlyComment = true
		}
	}
	return
}

func parseOtherFile(lines []string, eh *Hostfile) (err error) {
	var current *Host

	for _, line := range lines {

		if m := rxHostLine.FindAllStringSubmatch(line, -1); m != nil {
			current = nil
			host := HostInfo{address: m[0][1]}
			host.domains = rxSpaceSep.Split(m[0][2], -1)
			host.active = true
			eh.hosts = append(eh.hosts, NewHostFromInfo(host))
			log.DebugF("line: \"%v\", host: %v", line, host)
			continue
		}

		if m := rxUnHostLine.FindAllStringSubmatch(line, -1); m != nil {
			current = nil
			host := HostInfo{address: m[0][1]}
			host.domains = rxSpaceSep.Split(m[0][2], -1)
			host.active = false
			eh.hosts = append(eh.hosts, NewHostFromInfo(host))
			log.DebugF("line: \"%v\", unhost: %v", line, host)
			continue
		}

		if m := rxCommentLine.FindAllStringSubmatch(line, -1); m != nil {
			if current == nil {
				current = NewComment(m[0][1])
				eh.hosts = append(eh.hosts, current)
			} else {
				current.AppendComment(m[0][1])
			}
			log.DebugF("line: \"%v\", comment: %v", line, current)
			continue
		}

		if line != "" {
			log.DebugF("skipping line: \"%v\"", line)
		}
	}

	return processCommentBlocks(eh)
}

func parseOwnFile(lines []string, eh *Hostfile) (err error) {
	var current *HostInfo = nil

	for _, line := range lines {
		if rxEmptyLine.MatchString(line) {
			if current != nil {
				eh.hosts = append(eh.hosts, NewHostFromInfo(*current))
				current = nil
			}
			log.DebugF("empty: \"%v\", current: %v", line, current)
			continue
		}

		if m := rxHostLine.FindAllStringSubmatch(line, -1); m != nil {
			if current == nil {
				current = &HostInfo{address: m[0][1]}
			} else {
				current.address = m[0][1]
			}
			current.domains = rxSpaceSep.Split(m[0][2], -1)
			current.active = true
			log.DebugF("host: \"%v\", current: %v", line, current)
			continue
		}

		if m := rxLookupLine.FindAllStringSubmatch(line, -1); m != nil {
			if current == nil {
				current = &HostInfo{}
			}
			current.lookup = m[0][1]
			log.DebugF("lookup: \"%v\", current: %v", line, current)
			continue
		}

		if m := rxUnHostLine.FindAllStringSubmatch(line, -1); m != nil {
			if current == nil {
				current = &HostInfo{address: m[0][1]}
			} else {
				current.address = m[0][1]
			}
			current.domains = rxSpaceSep.Split(m[0][2], -1)
			current.active = false
			log.DebugF("inactive: \"%v\", current: %v", line, current)
			continue
		}

		if m := rxCommentLine.FindAllStringSubmatch(line, -1); m != nil {
			if current == nil {
				current = &HostInfo{}
			}
			if len(current.comment) > 0 {
				current.comment += "\n"
			}
			current.comment += strings.TrimSpace(m[0][1])
			log.DebugF("comment: \"%v\", current: %v", line, current)
			continue
		}

		log.WarnF("unknown line: \"%v\", current: %v\n", line, current)
	}

	return processCommentBlocks(eh)
}