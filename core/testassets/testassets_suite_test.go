package testassets_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTestassets(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testassets Suite")
}
