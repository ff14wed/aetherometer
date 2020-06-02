package datasheet_test

import (
	"bytes"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/testassets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("World", func() {
	Describe("PopulateWorlds", func() {
		It("correctly populates the World store", func() {
			var s datasheet.WorldStore
			err := s.PopulateWorlds(bytes.NewReader([]byte(testassets.WorldCSV)))
			Expect(err).ToNot(HaveOccurred())
			Expect(s).To(HaveLen(len(testassets.ExpectedWorldData)))
			for k, d := range s {
				Expect(d).To(Equal(testassets.ExpectedWorldData[k]))
			}
		})

		It("returns an error if the datasheet is blank", func() {
			var s datasheet.WorldStore
			err := s.PopulateWorlds(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the datasheet is invalid", func() {
			var s datasheet.WorldStore
			err := s.PopulateWorlds(bytes.NewReader([]byte(InvalidCSV)))
			Expect(err).To(HaveOccurred())
		})
	})
})
