package gurdil

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/middleware"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("gurdil", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	c.Next() // 'gurdil'
	c.Next() // 'gurdil'
	if c.NextArg() {
		return middleware.Error("gurdil", c.ArgErr())
	}
    z := Zone{}
    z.SetDomain(c.Val())
	dnsserver.GetConfig(c).AddMiddleware(func(next middleware.Handler) middleware.Handler {
        return Gurdil{ Zone: z}
	})

	return nil
}
