package lbservice

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLbservice(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lbservice Suite")
}

var _ = Describe("Lbservice", func() {

})
