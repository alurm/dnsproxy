/*
Extra-recursive-resolver tries to resolve domains to their A records.

Effectively, this is a simplistic DNS resolver.

Program takes a domain and returns an IPv4 address if it finds one.
*/
package main

import (
	"fmt"
	"github.com/alurm/dnsproxy/proxy"
	"github.com/miekg/dns"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func resolve(name string) *net.IP {
	get := func(in []dns.RR, _type uint16) dns.RR {
		for _, record := range in {
			if record.Header().Rrtype == _type {
				return record
			}
		}
		return nil
	}

	getCanonicalName := func(msg dns.Msg) *string {
		if record := get(msg.Answer, dns.TypeCNAME); record != nil {
			return &record.(*dns.CNAME).Target
		}
		return nil
	}

	getRedirectDomain := func(msg dns.Msg) *string {
		if record := get(msg.Ns, dns.TypeNS); record != nil {
			return &record.(*dns.NS).Ns
		}
		return nil
	}

	getRedirect := func(msg dns.Msg) net.IP {
		if record := get(msg.Extra, dns.TypeA); record != nil {
			return record.(*dns.A).A
		}
		return nil
	}

	getIP := func(msg dns.Msg) *dns.A {
		if record := get(msg.Answer, dns.TypeA); record != nil {
			return record.(*dns.A)
		}
		return nil
	}

	currentServer := proxy.Conf.Server

	for {
		reply := proxy.Query(name, currentServer, dns.TypeA)
		if ip := getIP(*reply); ip != nil {
			return &ip.A
		} else if canonical := getCanonicalName(*reply); canonical != nil {
			return resolve(*canonical)
		} else if secondServer := getRedirect(*reply); secondServer != nil {
			currentServer = secondServer
		} else if domain := getRedirectDomain(*reply); domain != nil {
			resolution := resolve(*domain)
			if resolution == nil {
				return nil
			}
			currentServer = *resolution
		} else {
			return nil
		}
	}
}

func main() {
	proxy.Conf = proxy.ReadConfiguration(
		filepath.Dir(proxy.Here()) + "/configuration.json",
	)
	if len(os.Args) != 2 {
		fmt.Println("usage: ./extra-recursive-resolver domain")
		os.Exit(1)
	}
	name := os.Args[1]
	if !strings.HasSuffix(name, ".") {
		name += "."
	}
	ip := resolve(name)
	if ip == nil {
		fmt.Println("couldn't resolve domain")
		os.Exit(1)
	}
	fmt.Println(ip.String())
}
