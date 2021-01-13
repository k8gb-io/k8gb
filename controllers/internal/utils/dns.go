package utils

import (
	"fmt"
	"sort"

	"github.com/lixiangzhong/dnsutil"
)

// Dig retrieves list of tuple <IP address, A record > from edge DNS server for specific FQDN
func Dig(edgeDNSServer, fqdn string) ([]string, error) {
	var dig dnsutil.Dig
	if edgeDNSServer == "" {
		return nil, fmt.Errorf("empty edgeDNSServer")
	}
	err := dig.SetDNS(edgeDNSServer)
	if err != nil {
		err = fmt.Errorf("dig error: can't set query dns (%s) with error(%s)", edgeDNSServer, err)
		return nil, err
	}
	a, err := dig.A(fqdn)
	if err != nil {
		err = fmt.Errorf("dig error: can't dig fqdn(%s) with error(%s)", fqdn, err)
		return nil, err
	}
	var IPs []string
	for _, ip := range a {
		IPs = append(IPs, fmt.Sprint(ip.A))
	}
	sort.Strings(IPs)
	return IPs, nil
}
