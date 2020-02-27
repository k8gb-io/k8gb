// fake dns server that is used for external dns communication tests of ohmyglb

package gslb

import (
	"fmt"
	"strconv"
	"time"

	"github.com/miekg/dns"
)

func oldEdgeTimestamp() string {
	now := time.Now()
	duration, _ := time.ParseDuration("10m")
	before := now.Add(-duration)
	edgeTimestamp := fmt.Sprint(before.UTC().Format("2006-01-02T15:04:05"))
	return edgeTimestamp
}

var records = map[string][]string{
	"localtargets.app3.cloud.example.com.": []string{"10.1.0.1", "10.1.0.2", "10.1.0.3"},
	"test-gslb-ns-eu.example.com.":         []string{oldEdgeTimestamp()},
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
		case dns.TypeTXT:
			log.Info(fmt.Sprintf("Query for TXT %s\n", q.Name))
			TXTs := records[q.Name]
			log.Info(fmt.Sprintf("TXTs found: %s\n", TXTs))
			if len(TXTs) > 0 {
				for _, txt := range TXTs {
					rr, err := dns.NewRR(fmt.Sprintf("%s TXT %s", q.Name, txt))
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

	err := w.WriteMsg(m)
	if err != nil {
		log.Info(fmt.Sprintf("Failed to write message:%s", err))
	}
}

func fakedns() {
	// attach request handler func
	dns.HandleFunc("example.com.", handleDNSRequest)

	// start server
	port := 7753
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	go func() {
		log.Info(fmt.Sprintf("Starting at %d\n", port))
		err := server.ListenAndServe()
		defer func() {
			err := server.Shutdown()
			if err != nil {
				log.Error(err, "Failed to shutdown fakedns server")
			}

		}()
		if err != nil {
			log.Error(err, "Failed to start fakedns server")
		}
	}()
}
