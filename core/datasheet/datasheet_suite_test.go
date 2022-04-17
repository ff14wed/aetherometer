package datasheet_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

const InvalidCSV = "key,0,1\n#,Singular"

func TestDatasheet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Datasheet Suite")
}
