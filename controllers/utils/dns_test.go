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
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	defaultFqdn = "google.com"
)

var defaultEdgeDNSServer = DNSServer{
	Host: "8.8.8.8",
	Port: 53,
}

func TestValidDigFQDNWithDot(t *testing.T) {
	// arrange
	if !connected() {
		t.Skipf("no connectivity, skipping")
	}
	fqdn := defaultFqdn
	// act
	result, err := Dig(fqdn+".", defaultEdgeDNSServer)
	// assert
	if strings.Contains(fmt.Sprintf("%v", err), "timeout") {

		t.Skipf("%s timeouts", defaultEdgeDNSServer)
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.NotEmpty(t, result[0])
}

func TestValidDig(t *testing.T) {
	// arrange
	if !connected() {
		t.Skipf("no connectivity, skipping")
	}
	fqdn := defaultFqdn
	// act
	result, err := Dig(fqdn, defaultEdgeDNSServer)
	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.NotEmpty(t, result[0])
}

func TestEmptyFQDNButValidEdgeDNS(t *testing.T) {
	// arrange
	if !connected() {
		t.Skipf("no connectivity, skipping")
	}
	fqdn := ""
	// act
	result, err := Dig(fqdn, defaultEdgeDNSServer)
	// assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestEmptyEdgeDNS(t *testing.T) {
	// arrange
	fqdn := "whatever"
	// act
	result, err := Dig(fqdn, DNSServer{Host: "", Port: 53})
	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestEmptyDNSList(t *testing.T) {
	// arrange
	fqdn := "whatever"
	// act
	result, err := Dig(fqdn, []DNSServer{}...)
	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestOneValidEdgeDNSInTheList(t *testing.T) {
	if !connected() {
		t.Skipf("no connectivity, skipping")
	}
	// arrange
	edgeDNSServers := []DNSServer{
		{Host: "127.1.2.3", Port: 53}, // wrong
		{Host: "8.8.8.8", Port: 153},  // wrong
		{Host: "8.8.8.8", Port: 53},   // ok
		{Host: "8.8.8.8", Port: 253},  // wrong
	}
	fqdn := defaultFqdn
	// act
	result, err := Dig(fqdn, edgeDNSServers...)
	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.NotEmpty(t, result[0])
}

func TestNoValidEdgeDNSInTheList(t *testing.T) {
	// arrange
	edgeDNSServers := []DNSServer{
		{Host: "", Port: 53},         // wrong
		{Host: "8.8.8.8", Port: 153}, // wrong
		{Host: "8.8.4.4", Port: 253}, // wrong
	}
	fqdn := defaultFqdn
	// act
	result, err := Dig(fqdn, edgeDNSServers...)
	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestEmptyEdgeDNSInTheList(t *testing.T) {
	// arrange
	edgeDNSServers := []DNSServer{
		{Host: "", Port: 53},        // wrong
		{Host: "8.8.8.8", Port: 53}, // ok
		{Host: "8.8.4.4", Port: 53}, // ok
	}
	fqdn := defaultFqdn
	// act
	result, err := Dig(fqdn, edgeDNSServers...)
	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestMultipleValidEdgeDNSInTheList(t *testing.T) {
	if !connected() {
		t.Skipf("no connectivity, skipping")
	}
	// arrange
	edgeDNSServers := []DNSServer{
		{Host: "1.1.1.1", Port: 53}, // ok
		{Host: "8.8.8.8", Port: 53}, // ok
		{Host: "8.8.4.4", Port: 53}, // ok
	}
	fqdn := defaultFqdn
	// act
	result, err := Dig(fqdn, edgeDNSServers...)
	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.NotEmpty(t, result[0])
}

func TestValidEdgeDNSButNonExistingFQDN(t *testing.T) {
	// arrange
	edgeDNSServer := "localhost"
	fqdn := "some-valid-ip-fqdn-123"
	// act
	result, err := Dig(fqdn, DNSServer{Host: edgeDNSServer, Port: 53})
	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func connected() (ok bool) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://google.com", nil)
	client := &http.Client{}
	go func() {
		time.Sleep(time.Second * 2)
		cancel()
	}()
	res, err := client.Do(req)
	if err != nil {
		return false
	}
	return res.Body.Close() == nil
}
