// Package proxy provides a Serve function for starting the DNS proxy.
//
// Conf variable is used for storing configuration.
//
// Proxy doesn't claim itself to be Authoritative, but maybe it should.
//
// Proxy checks domain names for blockage only for the requests, not responses.
//
// Proxy doesn't set UDPSize, leaving it at the minumum, but maybe it should.
package proxy

import (
	"github.com/miekg/dns"
	"net"
)

func Query(name string, server net.IP, question uint16) *dns.Msg {
	var m dns.Msg
	m.SetQuestion(name, question)
	var c dns.Client
	reply, _, err := c.Exchange(&m, server.String()+":53")
	if err != nil {
		panic(err)
	}
	return reply
}

func Serve() {
	server := dns.Server{
		Addr: ":53",
		Net:  "udp",
		// UDPSize: 1 << 16,
		Handler: dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
			var m dns.Msg
			m.SetReply(r)
			// m.Authoritative = true
			for _, q := range r.Question {
				if Conf.Blocklist[q.Name] {
					m.SetRcode(r, dns.RcodeRefused)
					err := w.WriteMsg(&m)
					if err != nil {
						panic(err)
					}
					return
				}
				rrs := Query(q.Name, Conf.Server, q.Qtype).Answer
				// for _, rr := range rrs {
				// 	header := rr.Header()
				// 	fmt.Println(header.Name)
				// 	if conf.blocklist[header.Name] {
				// 		m.SetRcode(r, dns.RcodeRefused)
				// 		err := w.WriteMsg(&m)
				// 		if err != nil {
				// 			panic(err)
				// 		}
				// 		return
				// 	}
				// }
				m.Answer = append(m.Answer, rrs...)
			}
			err := w.WriteMsg(&m)
			if err != nil {
				panic(err)
			}
		}),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
