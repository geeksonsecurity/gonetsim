package dns

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"strconv"
)

type Server struct {
	Port      int
	tcpServer *dns.Server
	udpServer *dns.Server
}

func (d *Server) Stop() {
	log.Printf("Stopping DNS servers\n")
	if d.tcpServer != nil {
		log.Printf("Stopping TCP DNS server\n")
		d.tcpServer.Shutdown()
	}
	if d.udpServer != nil {
		log.Printf("Stopping UDP DNS server\n")
		d.udpServer.Shutdown()
	}
}

func (d *Server) Start(callback func(domain string)) {
	log.Printf("Starting DNS server at port %d\n", d.Port)

	// attach request handler func
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Compress = false
		switch r.Opcode {
		case dns.OpcodeQuery:
			for _, q := range m.Question {
				switch q.Qtype {
				case dns.TypeA:
					log.Printf("Query for %s\n", q.Name)
					callback(q.Name)
					rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, "127.0.0.1"))
					if err == nil {
						m.Answer = append(m.Answer, rr)
					}
				}
			}
		}
		w.WriteMsg(m)
	})

	go func() {
		d.tcpServer = &dns.Server{Addr: "127.0.0.1:" + strconv.Itoa(d.Port), Net: "udp"}

		if err := d.tcpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
		defer d.tcpServer.Shutdown()
	}()

	go func() {
		d.udpServer = &dns.Server{Addr: ":" + strconv.Itoa(d.Port), Net: "tcp"}
		if err := d.udpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set tcp listener %s\n", err.Error())
		}
		defer d.udpServer.Shutdown()
	}()
}
