# Network Simulator in Go

Minimal network simulator written in Go. Inspired by [INetSim](https://www.inetsim.org/).

**Features:**
* Create DNS server on all addresses and respond to all queries with its own IP
* Create Web server (HTTP/HTTPs) responding to all requests based user configuration
* For HTTPs it generates on the fly server certificate (based on requested server name)

## Installing CA certificate

### Linux
Export the CA certificate somewhere with the "Save CA certificate" button. 
You can import the CA certificate directly in your browser settings, for example in Firefox go to Settings -> Certificates -> Authorities -> Import...
![firefox](https://user-images.githubusercontent.com/1062939/160276203-9eccd6bc-75ad-4bfb-a836-7d4907f82824.png)


### Windows
TBA

## Thanks 
Shamelessly used following resources:
* https://gist.github.com/walm/0d67b4fb2d5daf3edd4fad3e13b162cb
* https://shaneutt.com/blog/golang-ca-and-signed-cert-go/
