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

type Zone struct {
    domain string
}

func (z *Zone) SetDomain(name string) {
    z.domain = name
}

func (z Zone) GetDomain() string {
    return z.domain
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
        s := strings.Split(state.QName(), "." + wh.Zone.GetDomain())
        answer_ip := net.ParseIP(s[0]).To4()
        rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
        rr.(*dns.A).A = answer_ip
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
