package visitor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestVisitor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Visitor Suite")
}
