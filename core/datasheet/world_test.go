package datasheet_test

import (
	"bytes"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/testassets"
	. "github.com/onsi/ginkgo/v2"
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

	Describe("Lookup", func() {
		var s datasheet.WorldStore
		BeforeEach(func() {
			err := s.PopulateWorlds(bytes.NewReader([]byte(testassets.WorldCSV)))
			Expect(err).ToNot(HaveOccurred())
			Expect(s).To(HaveLen(len(testassets.ExpectedWorldData)))
		})

		It("returns the World with the requested world ID", func() {
			Expect(s.Lookup(5)).To(Equal(models.World{
				ID:   5,
				Name: "c-contents2",
			}))
		})

		It("returns an Unknown world if the requested world does not exist", func() {
			Expect(s.Lookup(123)).To(Equal(models.World{
				ID:   123,
				Name: "Unknown_123",
			}))
		})

	})
})
