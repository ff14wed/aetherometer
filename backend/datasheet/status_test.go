package datasheet_test

import (
	"bytes"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/testassets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Status", func() {
	Describe("PopulateStatuses", func() {
		It("correctly populates the Status store", func() {
			var s datasheet.StatusStore
			err := s.PopulateStatuses(bytes.NewReader([]byte(testassets.StatusCSV)))
			Expect(err).ToNot(HaveOccurred())
			Expect(s).To(HaveLen(len(testassets.ExpectedStatusData)))
			for k, d := range s {
				Expect(d).To(Equal(testassets.ExpectedStatusData[k]))
			}
		})

		It("returns an error if the datasheet is blank", func() {
			var s datasheet.StatusStore
			err := s.PopulateStatuses(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the datasheet is invalid", func() {
			var s datasheet.StatusStore
			err := s.PopulateStatuses(bytes.NewReader([]byte(InvalidCSV)))
			Expect(err).To(HaveOccurred())
		})
	})
})
