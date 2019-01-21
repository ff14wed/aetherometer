package update_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUpdate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Update Suite")
}
