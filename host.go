package editor

import (
	"fmt"
	"net"
	"strings"

	"github.com/go-curses/cdk/lib/paint"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/lib/sync"
)

type HostImportance string

const (
	HostNotImportant    HostImportance = "none"
	HostIsLocalhostIPv4 HostImportance = "ipv4"
	HostIsLocalhostIPv6 HostImportance = "ipv6"
)

type HostInfo struct {
	active  bool
	lookup  string
	address string
	comment string
	domains []string
}

func (h HostInfo) SameHostInfo(other HostInfo) (same bool) {
	same = h.active == other.active &&
		h.lookup == other.lookup &&
		h.address == other.address &&
		h.comment == other.comment &&
		cstrings.EqualStringSlices(h.domains, other.domains)
	return
}

type Host struct {
	HostInfo

	original    HostInfo
	onlyComment bool

	cache []net.IP

	sync.RWMutex
}

func NewComment(comment string) (host *Host) {
	host = new(Host)
	host.onlyComment = true
	host.active = false
	host.address = ""
	host.lookup = ""
	host.comment = comment
	host.domains = nil
	host.original = HostInfo{
		active:  false,
		address: "",
		lookup:  "",
		comment: comment,
		domains: nil,
	}
	return
}

func NewHostFromInfo(info HostInfo) (host *Host) {
	host = new(Host)
	host.onlyComment = false
	host.active = info.active
	host.address = info.address
	host.lookup = info.lookup
	host.comment = info.comment
	host.domains = info.domains
	host.original = info
	return
}

func (h *Host) Name() (ipOrName string) {
	h.RLock()
	defer h.RUnlock()
	if len(h.lookup) > 0 {
		return h.lookup
	}
	return h.address
}

func (h *Host) String() string {
	h.RLock()
	defer h.RUnlock()
	return h.address
}

func (h *Host) Equals(host *Host) bool {
	h.RLock()
	host.RLock()
	defer h.RUnlock()
	defer host.RUnlock()
	return h.address == host.address &&
		h.comment == host.comment &&
		h.active == host.active &&
		cstrings.EqualStringSlices(h.domains, host.domains)
}

func (h *Host) IsComment() bool {
	h.RLock()
	defer h.RUnlock()
	return h.onlyComment
}

func (h *Host) Changed() bool {
	h.RLock()
	defer h.RUnlock()
	return !h.SameHostInfo(h.original)
}

func (h *Host) Line() string {
	h.RLock()
	defer h.RUnlock()
	active := ""
	if !h.active {
		active = "#"
	}
	return fmt.Sprintf("%v%v\t%v\n", active, h.address, strings.Join(h.domains, " "))
}

func (h *Host) Empty() bool {
	if h.IsOnlyComment() {
		return h.comment == ""
	}
	return h.comment == "" &&
		h.address == "" &&
		h.lookup == "" &&
		len(h.domains) == 0
}

func (h *Host) Block() string {
	if h.Empty() {
		return ""
	}
	isComment := h.IsComment()
	h.RLock()
	defer h.RUnlock()
	out := ""
	if isComment {
		out += "###\n"
		for _, line := range rxNewlines.Split(h.comment, -1) {
			out += "# " + line + "\n"
		}
		out += "###\n"
		return out
	}
	if len(h.comment) > 0 {
		for _, line := range rxNewlines.Split(h.comment, -1) {
			out += "# " + line + "\n"
		}
	}
	if len(h.lookup) > 0 {
		out += fmt.Sprintf("#nslookup %v\n", h.lookup)
	}
	out += h.Line()
	return out
}

func (h *Host) PerformLookup() (found []net.IP, err error) {
	if len(h.cache) > 0 {
		found = h.cache
		return
	}
	if lookup := h.Lookup(); lookup != "" {
		if found, err = net.LookupIP(lookup); err == nil {
			h.Lock()
			h.cache = found
			h.Unlock()
		}
	} else {
		err = fmt.Errorf("missing domain to lookup")
	}
	return
}

func (h *Host) GetActualInfo() (label, tooltip string) {
	lookup := h.Lookup()
	address := h.Address()
	if cstrings.StringIsDomainName(lookup) {
		if found, err := h.PerformLookup(); err != nil {
			label = fmt.Sprintf("%v (!)", address)
			tooltip = err.Error()
			return
		} else {
			totalFound := len(found)
			if totalFound > 1 {
				for _, addr := range found {
					if address == addr.String() {
						label = fmt.Sprintf("%v (%v)", address, string(paint.RuneCheckbox))
						tooltip = fmt.Sprintf("is 1 of %d valid addresses", totalFound)
						return
					}
				}
			} else if totalFound == 1 && address == found[0].String() {
				label = fmt.Sprintf("%v (%v)", address, string(paint.RuneCheckbox))
				tooltip = "is the only valid\naddress for domain"
				return
			}
			label = fmt.Sprintf("%v (!)", address)
			tooltip = fmt.Sprintf("address not associated\nwith lookup domain")
			return
		}
	} else if cstrings.StringIsIP(address) {
		label = fmt.Sprintf("(%v)", address)
		tooltip = "is a valid IP address"
		return
	}
	label = "(not an address)"
	tooltip = "please enter a valid IP address"
	return
}

func (h *Host) SetActive(active bool) {
	h.Lock()
	h.active = active
	h.Unlock()
}

func (h *Host) Active() bool {
	h.RLock()
	defer h.RUnlock()
	return h.active
}

func (h *Host) SetLookup(value string) {
	h.Lock()
	h.lookup = value
	h.Unlock()
}

func (h *Host) Lookup() string {
	h.RLock()
	defer h.RUnlock()
	return h.lookup
}

func (h *Host) SetAddress(value string) {
	h.Lock()
	h.address = value
	h.Unlock()
}

func (h *Host) Address() string {
	h.RLock()
	defer h.RUnlock()
	return h.address
}

func (h *Host) SetComment(text string) {
	h.Lock()
	h.comment = strings.TrimSpace(text)
	h.Unlock()
}

func (h *Host) Comment() string {
	h.RLock()
	defer h.RUnlock()
	return h.comment
}

func (h *Host) IsOnlyComment() (onlyComment bool) {
	h.RLock()
	defer h.RUnlock()
	onlyComment = h.onlyComment
	return
}

func (h *Host) AppendComment(text string) {
	h.Lock()
	if len(h.comment) > 0 {
		h.comment += "\n"
	}
	h.comment += strings.TrimSpace(text)
	h.Unlock()
}

func (h *Host) SetDomains(text string) {
	h.Lock()
	h.domains = rxSpaceSep.Split(text, -1)
	h.Unlock()
}

func (h *Host) Domains() []string {
	h.RLock()
	defer h.RUnlock()
	return h.domains
}

func (h *Host) AddDomain(domain string) {
	h.Lock()
	h.domains = append(h.domains, domain)
	h.Unlock()
}

func (h *Host) RemoveDomain(domain string) {
	h.Lock()
	var domains []string
	for _, existing := range h.domains {
		if existing != domain {
			domains = append(domains, domain)
		}
	}
	h.domains = domains
	h.Unlock()
}

func (h *Host) HasDomain(needle string) (found bool) {
	h.RLock()
	defer h.RUnlock()
	for _, domain := range h.domains {
		if domain == needle {
			return true
		}
	}
	return
}

func (h *Host) Importance() HostImportance {
	for _, domain := range h.domains {
		switch domain {
		case "localhost":
			return HostIsLocalhostIPv4
		case "ip6-allrouters", "ip6-allnodes", "ip6-loopback", "ip6-localhost":
			return HostIsLocalhostIPv6
		}
	}
	return HostNotImportant
}