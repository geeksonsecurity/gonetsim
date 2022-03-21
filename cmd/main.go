package main

import (
	"gonetsim/internal/dns"
	"gonetsim/internal/web"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dns.StartDNSServer(53)
	web.StartHttpServer(8080)
	web.StartHttpsServer(8443)
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Fatalf("Signal (%v) received, stopping\n", s)
}
