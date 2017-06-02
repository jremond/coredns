// Package gurdil implements a middleware that returns details about the resolving
// querying it.
package gurdil

import (
	"net"
	"strings"

	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

// Gurdil is a middleware that returns your IP address, port and the protocol used for connecting
// to CoreDNS.
type Gurdil struct{
    Zone
}

// Zone with one domain
type Zone struct {
    domain string
}

// SetDomain puts the domain in the Zone
func (z *Zone) SetDomain(name string) {
    z.domain = name
}

// GetDomain returnt the domain name from Zone
func (z Zone) GetDomain() string {
    return z.domain
}

func extractIPv4(sIP string) (ip net.IP) {
    s := strings.Split(sIP, ".")
    for i := 0; i < len(s); i++ {
        p := net.ParseIP(strings.Join(s[i:i+4],"."))
        if p != nil {
            return p
        }
    }
    return nil
}


// ServeDNS implements the middleware.Handler interface.
func (wh Gurdil) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
    state := request.Request{W: w, Req: r}

    a := new(dns.Msg)
    a.SetReply(r)
    a.Compress = true
    a.Authoritative = true

    ip := state.IP()
    var rr dns.RR

    switch state.Family() {
    case 1:
        rr = new(dns.A)
        rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
        zoneDomain := wh.Zone.GetDomain()
        zoneDomain += "."
        if strings.HasSuffix(state.QName(), zoneDomain) {
            answerIP := extractIPv4(state.QName())
            rr.(*dns.A).A = answerIP
        }
    case 2:
        rr = new(dns.AAAA)
        rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass()}
        rr.(*dns.AAAA).AAAA = net.ParseIP(ip)
    }

    a.Answer = []dns.RR{rr}

    state.SizeAndDo(a)
    w.WriteMsg(a)

    return 0, nil
}

// Name implements the Handler interface.
func (wh Gurdil) Name() string { return "gurdil" }
