package address_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAddress(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Address Suite")
}
