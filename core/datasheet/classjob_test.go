package datasheet_test

import (
	"bytes"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/testassets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClassJob", func() {
	Describe("PopulateClassJob", func() {
		It("correctly populates the ClassJob store", func() {
			var c datasheet.ClassJobStore
			err := c.PopulateClassJobs(bytes.NewReader([]byte(testassets.ClassJobCSV)))
			Expect(err).ToNot(HaveOccurred())
			Expect(c).To(HaveLen(len(testassets.ExpectedClassJobData)))
			for k, d := range c {
				Expect(d).To(Equal(testassets.ExpectedClassJobData[k]))
			}
		})

		It("returns an error if the datasheet is blank", func() {
			var c datasheet.ClassJobStore
			err := c.PopulateClassJobs(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the datasheet is invalid", func() {
			var c datasheet.ClassJobStore
			err := c.PopulateClassJobs(bytes.NewReader([]byte(InvalidCSV)))
			Expect(err).To(HaveOccurred())
		})
	})
})
