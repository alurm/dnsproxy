package main

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"io"
	"net"
	"os"
	"path/filepath"
)

type configurationJSON struct {
	Server    string
	Blocklist []string
}

type configuration struct {
	server net.IP
	blocklist map[string]bool
}

func readConfiguration() configuration {
	var bytes []byte
	{
		executable, err := os.Executable()
		if err != nil {
			panic(err)
		}
		executable, err = filepath.EvalSymlinks(executable)
		if err != nil {
			panic(err)
		}
		path := filepath.Dir(executable) + "/configuration.json"
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		bytes, err = io.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}

	var cJSON configurationJSON
	err := json.Unmarshal(
		bytes,
		&cJSON,
	)
	if err != nil {
		panic(err)
	}

	var c configuration

	server := net.ParseIP(cJSON.Server)
	if server == nil {
		panic("server configuration is incorrect")
	}

	c.server = server

	c.blocklist = map[string]bool{}
	for _, domain := range cJSON.Blocklist {
		c.blocklist[domain] = true
	}

	return c
}

var conf configuration

func resolve(name string) dns.A {
	server := conf.server
	for {
		reply := dnsQuery(name, server)
		if ip := getAnswer(*reply); ip != nil {
			return *ip
		} else if canonical := getCanonicalName(*reply); canonical != nil {
			return resolve(*canonical)
		} else if secondServer := getRedirect(*reply); secondServer != nil {
			server = secondServer
		} else if domain := getRedirectDomain(*reply); domain != nil {
			server = resolve(*domain).A
		} else {
			panic("can't resolve domain")
		}
	}
}

func getCanonicalName(msg dns.Msg) *string {
	for _, record := range msg.Answer {
		if record.Header().Rrtype == dns.TypeCNAME {
			return &record.(*dns.CNAME).Target
		}
	}
	return nil
}

func getRedirectDomain(msg dns.Msg) *string {
	for _, record := range msg.Ns {
		if record.Header().Rrtype == dns.TypeNS {
			return &record.(*dns.NS).Ns
		}
	}
	return nil
}

func getRedirect(msg dns.Msg) net.IP {
	for _, record := range msg.Extra {
		if record.Header().Rrtype == dns.TypeA {
			return record.(*dns.A).A
		}
	}
	return nil
}

func getAnswer(msg dns.Msg) *dns.A {
	for _, record := range msg.Answer {
		if record.Header().Rrtype == dns.TypeA {
			return record.(*dns.A)
		}
	}
	return nil
}

func dnsQuery(name string, server net.IP) *dns.Msg {
	var msg dns.Msg
	msg.SetQuestion(name, dns.TypeA)
	var c dns.Client
	reply, _, err := c.Exchange(&msg, server.String()+":53")
	if err != nil {
		panic(err)
	}
	return reply
}

func serve() {
	server := dns.Server{
		Addr: ":53",
		Net: "udp",
		// UDPSize: 1 << 16,
		Handler: dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
			var m dns.Msg
			m.SetReply(r)
			// m.Authoritative = true
			for _, q := range r.Question {
				if q.Qtype != dns.TypeA {
					m.SetRcode(r, dns.RcodeNotImplemented)
					err := w.WriteMsg(&m)
					if err != nil {
						panic(err)
					}
					return
				}
				a := resolve(q.Name)
				m.Answer = append(m.Answer, &a)
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

func main() {
	conf = readConfiguration()
	fmt.Printf("%#v\n", conf)
	fmt.Printf("%#v\n", resolve("example.com.").A.String())
	serve()
}
