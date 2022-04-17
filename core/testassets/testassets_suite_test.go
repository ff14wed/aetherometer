package testassets_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTestassets(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testassets Suite")
}
