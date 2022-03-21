# NetworkSimulator in Go

Minimal network simulator written in Go. Inspired by [INetSim](https://www.inetsim.org/).

**Features:**
* Generate a CA certificate on start (you can import that as [trusted certificate authority](https://asu.secure.force.com/kb/articles/FAQ/How-Do-I-Add-Certificates-to-the-Trusted-Root-Certification-Authorities-Store-for-a-Local-Computer))
* Create DNS server on all addresses and respond to all queries with its own IP
* Create HTTP server responding to all requests `200 OK`
* Create HTTPs server responding to all requests `200 OK`. Generates on the fly server certificate (based on requested server name)

## Thanks 
Shamelessly used following resources:
* https://gist.github.com/walm/0d67b4fb2d5daf3edd4fad3e13b162cb
* https://shaneutt.com/blog/golang-ca-and-signed-cert-go/