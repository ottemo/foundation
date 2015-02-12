package product_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestProduct(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Product Suite")
}
