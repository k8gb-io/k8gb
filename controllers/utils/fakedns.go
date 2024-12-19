package utils

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/miekg/dns"
)

const periodDelay = 100 * time.Millisecond

type FakeDNSSettings struct {
	FakeDNSPort     int
	EdgeDNSZoneFQDN string
	DNSZoneFQDN     string
	Dump            bool
}

// DNSMock acts as DNS server but returns mock values
type DNSMock struct {
	// readinessProbe is the channel that is released when the dns server starts listening
	readinessProbe chan interface{}
	livenessProbe  chan interface{}
	settings       FakeDNSSettings
	records        map[uint16][]dns.RR
	server         *dns.Server
	err            error
}

type Result struct {
	Error error
}

func NewFakeDNS(settings FakeDNSSettings) *DNSMock {
	return &DNSMock{
		settings:       settings,
		readinessProbe: make(chan interface{}),
		livenessProbe:  make(chan interface{}),
		records:        make(map[uint16][]dns.RR),
		server:         &dns.Server{Addr: fmt.Sprintf("[::]:%v", settings.FakeDNSPort), Net: "udp", TsigSecret: nil, ReusePort: false},
	}
}

func (m *DNSMock) Start() *DNSMock {
	go m.startReadinessProbe()
	go func() {
		m.err = m.listen()
	}()
	<-m.readinessProbe
	fmt.Printf("FakeDNS listening on port %v \n", m.settings.FakeDNSPort)
	return m
}

func (m *DNSMock) RunTestFunc(f func()) *Result {
	if m.err == nil {
		f()
		go m.startLivenessProbe()
		m.err = m.server.Shutdown()
		<-m.livenessProbe
	}
	return &Result{
		Error: m.err,
	}
}

func (r *Result) RequireNoError(t *testing.T) {
	require.NoError(t, r.Error)
}

func (m *DNSMock) AddTXTRecord(fqdn string, strings ...string) *DNSMock {
	t := &dns.TXT{
		Hdr: dns.RR_Header{Name: fqdn, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0},
		Txt: strings,
	}
	m.records[dns.TypeTXT] = append(m.records[dns.TypeTXT], t)
	return m
}

func (m *DNSMock) AddNSRecord(fqdn, nsName string) *DNSMock {
	ns := &dns.NS{
		Hdr: dns.RR_Header{Name: fqdn, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 0},
		Ns:  nsName,
	}
	m.records[dns.TypeNS] = append(m.records[dns.TypeNS], ns)
	return m
}

func (m *DNSMock) AddARecord(fqdn string, ip net.IP) *DNSMock {
	rr := &dns.A{
		Hdr: dns.RR_Header{Name: fqdn, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
		A:   ip.To4(),
	}
	m.records[dns.TypeA] = append(m.records[dns.TypeA], rr)
	return m
}

func (m *DNSMock) AddAAAARecord(ip net.IP) *DNSMock {
	rr := &dns.A{
		Hdr: dns.RR_Header{Name: m.settings.DNSZoneFQDN, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 0},
		A:   ip.To16(),
	}
	m.records[dns.TypeAAAA] = append(m.records[dns.TypeAAAA], rr)
	return m
}

func (m *DNSMock) AddCNAMERecord(fqdn string, cname string) *DNSMock {
	rr := &dns.CNAME{
		Hdr:    dns.RR_Header{Name: fqdn, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: 0},
		Target: cname,
	}
	m.records[dns.TypeA] = append(m.records[dns.TypeA], rr)
	return m
}

func (m *DNSMock) listen() (err error) {
	dns.HandleFunc(m.settings.EdgeDNSZoneFQDN, m.handleReflect)
	for e := range m.serve() {
		if e != nil {
			err = fmt.Errorf("%s", e)
		}
	}
	return
}

func (m *DNSMock) hit() (err error) {
	g := new(dns.Msg)
	host := fmt.Sprintf("localhost:%v", m.settings.FakeDNSPort)
	g.SetQuestion(m.settings.DNSZoneFQDN, dns.TypeA)
	_, err = dns.Exchange(g, host)
	return
}

func (m *DNSMock) startReadinessProbe() {
	defer close(m.readinessProbe)
	var err error
	for i := 0; i < 5; i++ {
		err = m.hit()
		if err != nil {
			time.Sleep(periodDelay)
			continue
		}
		return
	}
	m.err = fmt.Errorf("readiness probe %s (%s)", err, m.err)
}

// liveness probe will be closed when the FakeDNS is not responding
func (m *DNSMock) startLivenessProbe() {
	defer close(m.livenessProbe)
	var err error
	for i := 0; i < 5; i++ {
		err = m.hit()
		if err != nil {
			return
		}
		time.Sleep(periodDelay)
	}
	m.err = fmt.Errorf("liveness probe %s (%s)", err, m.err)
}

func (m *DNSMock) serve() <-chan error {
	errors := make(chan error)
	go func() {
		defer close(errors)
		var err error
		if err = m.server.ListenAndServe(); err != nil {
			errors <- fmt.Errorf("failed to setup the server: %s", err.Error())
		}
	}()
	return errors
}

func (m *DNSMock) handleReflect(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Compress = false
	if m.records[r.Question[0].Qtype] != nil {
		for _, rr := range m.records[r.Question[0].Qtype] {
			fqdn := strings.Split(rr.String(), "\t")[0]
			if fqdn == r.Question[0].Name || m.settings.Dump {
				msg.Answer = append(msg.Answer, rr)
			}
		}
	}
	_ = w.WriteMsg(msg)
}
