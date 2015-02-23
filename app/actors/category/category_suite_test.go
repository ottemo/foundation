package category_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCategory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Category Suite")
}
