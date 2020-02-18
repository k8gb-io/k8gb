// fake dns server that is used for external dns communication tests of ohmyglb

package gslb

import (
	"fmt"
	"strconv"

	"github.com/miekg/dns"
)

var records = map[string][]string{
	"localtargets.app3.cloud.example.com.": []string{"10.1.0.1", "10.1.0.2", "10.1.0.3"},
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Info(fmt.Sprintf("Query for %s\n", q.Name))
			ips := records[q.Name]
			log.Info(fmt.Sprintf("IPs found: %s\n", ips))
			if len(ips) > 0 {
				for _, ip := range ips {
					rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
					if err == nil {
						m.Answer = append(m.Answer, rr)
					}
				}
			}
		}
	}
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func fakedns() {
	// attach request handler func
	dns.HandleFunc("cloud.example.com.", handleDNSRequest)

	// start server
	port := 7753
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	go func() {
		log.Info(fmt.Sprintf("Starting at %d\n", port))
		err := server.ListenAndServe()
		defer server.Shutdown()
		if err != nil {
			log.Error(err, "Failed to start fakedns server")
		}
	}()
}
