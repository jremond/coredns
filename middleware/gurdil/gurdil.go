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
    if len(s) > 4 {
        for i := 0; i < len(s); i++ {
            p := net.ParseIP(strings.Join(s[i:i+4],"."))
            if p != nil {
                return p
            }
        }
    }
    return nil
}

func extractIPv6(sIP string) (ip net.IP) {
    // remove anyname before the IPv6 address
    w := strings.Split(sIP, ".")
    s := strings.Split(w[len(w)-1], ":")
    for i := 0; i < len(s); i++ {
        p := net.ParseIP(strings.Join(s[i:],":"))
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

    zoneDomain := "."
    zoneDomain += wh.Zone.GetDomain()
    zoneDomain += "."

    var rr dns.RR

    rr = new(dns.A)
    rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}

    if strings.HasSuffix(state.QName(), zoneDomain) {
        answerIPv4 := extractIPv4(strings.TrimSuffix(state.QName(), zoneDomain))
        rr.(*dns.A).A = answerIPv4
        a.Answer = []dns.RR{rr}
    }

    rr = new(dns.AAAA)
    rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass()}
    if strings.HasSuffix(state.QName(), zoneDomain) {
        answerIPv6 := extractIPv6(strings.TrimSuffix(state.QName(), zoneDomain))
        rr.(*dns.AAAA).AAAA = answerIPv6
        a.Answer = []dns.RR{rr}
    }

    state.SizeAndDo(a)
    w.WriteMsg(a)

    return 0, nil
}

// Name implements the Handler interface.
func (wh Gurdil) Name() string { return "gurdil" }
