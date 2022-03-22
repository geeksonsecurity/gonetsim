package dns

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"strconv"
)

func StartDNSServer(port int, callback func(domain string)) {
	log.Printf("Starting DNS server at port %d\n", port)

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
		srv := &dns.Server{Addr: "127.0.0.1:" + strconv.Itoa(port), Net: "udp"}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
		defer srv.Shutdown()
	}()

	go func() {
		srv := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "tcp"}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set tcp listener %s\n", err.Error())
		}
		defer srv.Shutdown()
	}()
}
