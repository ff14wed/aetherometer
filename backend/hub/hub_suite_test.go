package hub_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hub Suite")
}
