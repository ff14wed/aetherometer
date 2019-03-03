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
			err := a.PopulateActions(bytes.NewReader([]byte(testassets.ActionCSV)))
			Expect(err).ToNot(HaveOccurred())
			Expect(a.Actions).To(HaveLen(len(testassets.ExpectedActionData)))
			for k, d := range a.Actions {
				Expect(d).To(Equal(testassets.ExpectedActionData[k]))
			}
		})

		It("returns an error if the datasheet is blank", func() {
			var a datasheet.ActionStore
			err := a.PopulateActions(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the datasheet is invalid", func() {
			var a datasheet.ActionStore
			err := a.PopulateActions(bytes.NewReader([]byte(InvalidCSV)))
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("PopulateOmens", func() {
		It("correctly populates the Omens store", func() {
			var a datasheet.ActionStore
			err := a.PopulateOmens(bytes.NewReader([]byte(testassets.OmenCSV)))
			Expect(err).ToNot(HaveOccurred())
			Expect(a.Omens).To(HaveLen(len(testassets.ExpectedOmenData)))
			for k, d := range a.Omens {
				Expect(d).To(Equal(testassets.ExpectedOmenData[k]))
			}
		})

		It("returns an error if the datasheet is blank", func() {
			var a datasheet.ActionStore
			err := a.PopulateOmens(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the datasheet is invalid", func() {
			var a datasheet.ActionStore
			err := a.PopulateOmens(bytes.NewReader([]byte(InvalidCSV)))
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetAction", func() {
		var a *datasheet.ActionStore

		BeforeEach(func() {
			a = new(datasheet.ActionStore)
			Expect(a.PopulateActions(bytes.NewReader([]byte(testassets.ActionCSV)))).To(Succeed())
			Expect(a.PopulateOmens(bytes.NewReader([]byte(testassets.OmenCSV)))).To(Succeed())
		})

		It("gets the action", func() {
			Expect(a.GetAction(2)).To(Equal(testassets.ExpectedActionData[2]))
		})

		It("returns an empty action if the action doesn't exist", func() {
			Expect(a.GetAction(1)).To(Equal(datasheet.Action{}))
		})

		Context("when the action has an omen", func() {
			It("returns an action with the Omen field set", func() {
				ac := a.GetAction(203)
				Expect(ac).ToNot(Equal(testassets.ExpectedActionData[203]))
				expectedAc := testassets.ExpectedActionData[203]
				expectedAc.Omen = testassets.ExpectedOmenData[1].Name
				Expect(ac).To(Equal(expectedAc))
			})
		})
	})
})
