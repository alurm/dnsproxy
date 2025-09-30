# DNS blocklist proxy

## About

This is a simple DNS proxy accepting a blocklist made as a test task for a company.

There are two programs: dnsproxy and extra-recursive-resolver.

Documentation is available for both of them in Go comments, suitable for godoc.

## Direct dependencies

- Go.
- Parts of Go's standard library.
- DNS library: <https://github.com/miekg/dns>.
- GNU make.
- Entr command for rerunning commands when files change.

## Make targets

Make check does some test queries if dnsproxy is running. It also performs some checks on extra-recursive-resolver program.

Make configuration.json copies example-configuration.json as the configuration.

Make build builds both programs.

Make clean removes some of the build artifacts.

Make watch provides a simple compile-edit loop using entr.

Make README.html generates this document in HTML form.

## Resources used during development

<https://jvns.ca/blog/2022/02/01/a-dns-resolver-in-80-lines-of-go/>

<https://reintech.io/blog/implementing-a-dns-server-in-go/>
