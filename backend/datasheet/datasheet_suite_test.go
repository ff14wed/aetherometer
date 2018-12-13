package datasheet_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var InvalidJSON = `[{"Invalid"]`

func TestDatasheet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Datasheet Suite")
}
