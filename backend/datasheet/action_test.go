package datasheet_test

import (
	"bytes"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/testassets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Action", func() {
	Describe("PopulateActions", func() {
		It("correctly populates the Actions store", func() {
			var a datasheet.ActionStore
			err := a.PopulateActions(bytes.NewReader([]byte(testassets.ActionJSON)))
			Expect(err).ToNot(HaveOccurred())
			Expect(a).To(HaveLen(len(testassets.ExpectedActionData)))
			for k, d := range a {
				Expect(d).To(Equal(testassets.ExpectedActionData[k]))
			}
		})
		It("returns an error if the datasheet is blank", func() {
			var a datasheet.ActionStore
			err := a.PopulateActions(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})
		It("returns an error if the datasheet is not valid JSON", func() {
			var a datasheet.ActionStore
			err := a.PopulateActions(bytes.NewReader([]byte(InvalidJSON)))
			Expect(err).To(HaveOccurred())
		})
	})
})
