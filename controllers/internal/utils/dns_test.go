package utils

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidDig(t *testing.T) {
	// arrange
	if !connected() {
		t.Skipf("no connectivity, skipping")
	}
	edgeDNSServer := "8.8.8.8"
	fqdn := "google.com"
	// act
	result, err := Dig(edgeDNSServer, fqdn)
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
	edgeDNSServer := "8.8.8.8"
	fqdn := ""
	// act
	result, err := Dig(edgeDNSServer, fqdn)
	// assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestEmptyEdgeDNS(t *testing.T) {
	// arrange
	edgeDNSServer := ""
	fqdn := "whatever"
	// act
	result, err := Dig(edgeDNSServer, fqdn)
	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestValidEdgeDNSButNonExistingFQDN(t *testing.T) {
	// arrange
	edgeDNSServer := "localhost"
	fqdn := "some-valid-ip-fqdn-123"
	// act
	result, err := Dig(edgeDNSServer, fqdn)
	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func connected() (ok bool) {
	res, err := http.Get("http://google.com")
	if err != nil {
		return false
	}
	return res.Body.Close() == nil
}
