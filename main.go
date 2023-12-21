/*
Dnsproxy proxies a DNS and blocks excluded domains.

If a request is made for a blocked domain, proxy returns "refused" response cod.e.

Configuration file for proxy is named configuration.json and must lie in the same real (not symbolically linked) directory as dnsproxy program.

Example file example-configuration.json is provided.

There are no options.

UDP is used for transport.
*/
package main

import (
	"github.com/alurm/dnsproxy/proxy"
)

func main() {
	proxy.Conf = proxy.ReadConfiguration(proxy.Here() + "/configuration.json")
	proxy.Serve()
}
