package datasheet_test

import (
	"bytes"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/testassets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BNPC", func() {
	Describe("PopulateBNPCNames", func() {
		It("correctly populates the BNPCNames store", func() {
			b := new(datasheet.BNPCStore)
			err := b.PopulateBNPCNames(bytes.NewReader([]byte(testassets.BNPCNameJSON)))
			Expect(err).ToNot(HaveOccurred())
			Expect(b.BNPCNames).To(Equal(testassets.ExpectedBNPCNames))
		})
		It("returns an error if the datasheet is blank", func() {
			b := new(datasheet.BNPCStore)
			err := b.PopulateBNPCNames(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})
		It("returns an error if the datasheet is not valid JSON", func() {
			b := new(datasheet.BNPCStore)
			err := b.PopulateBNPCNames(bytes.NewReader([]byte(InvalidJSON)))
			Expect(err).To(HaveOccurred())
		})
	})
	Describe("PopulateBNPCBases", func() {
		It("correctly populates the BNPCBases store", func() {
			b := new(datasheet.BNPCStore)
			err := b.PopulateBNPCBases(bytes.NewReader([]byte(testassets.BNPCBaseJSON)))
			Expect(err).ToNot(HaveOccurred())
			Expect(b.BNPCBases).To(Equal(testassets.ExpectedBNPCBases))
		})
		It("returns an error if the datasheet is blank", func() {
			b := new(datasheet.BNPCStore)
			err := b.PopulateBNPCBases(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})
		It("returns an error if the datasheet is not valid JSON", func() {
			b := new(datasheet.BNPCStore)
			err := b.PopulateBNPCBases(bytes.NewReader([]byte(InvalidJSON)))
			Expect(err).To(HaveOccurred())
		})
	})
	Describe("PopulateModelCharas", func() {
		It("correctly populates the ModelCharas store", func() {
			b := new(datasheet.BNPCStore)
			err := b.PopulateModelCharas(bytes.NewReader([]byte(testassets.ModelCharaJSON)))
			Expect(err).ToNot(HaveOccurred())
			Expect(b.ModelCharas).To(Equal(testassets.ExpectedModelCharas))
		})
		It("returns an error if the datasheet is blank", func() {
			b := new(datasheet.BNPCStore)
			err := b.PopulateModelCharas(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})
		It("returns an error if the datasheet is not valid JSON", func() {
			b := new(datasheet.BNPCStore)
			err := b.PopulateModelCharas(bytes.NewReader([]byte(InvalidJSON)))
			Expect(err).To(HaveOccurred())
		})
	})
	Describe("PopulateModelSkeletons", func() {
		It("correctly populates the ModelSkeletons store", func() {
			b := new(datasheet.BNPCStore)
			err := b.PopulateModelSkeletons(bytes.NewReader([]byte(testassets.ModelSkeletonJSON)))
			Expect(err).ToNot(HaveOccurred())
			Expect(b.ModelSkeletons).To(Equal(testassets.ExpectedModelSkeletons))
		})
		It("returns an error if the datasheet is blank", func() {
			b := new(datasheet.BNPCStore)
			err := b.PopulateModelSkeletons(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})
		It("returns an error if the datasheet is not valid JSON", func() {
			b := new(datasheet.BNPCStore)
			err := b.PopulateModelSkeletons(bytes.NewReader([]byte(InvalidJSON)))
			Expect(err).To(HaveOccurred())
		})
	})
	Describe("GetBNPCInfo", func() {
		var bNPCStore *datasheet.BNPCStore
		BeforeEach(func() {
			bNPCStore = new(datasheet.BNPCStore)
			bNPCStore.BNPCBases = testassets.ExpectedBNPCBases
			bNPCStore.BNPCNames = testassets.ExpectedBNPCNames
			bNPCStore.ModelCharas = testassets.ExpectedModelCharas
			bNPCStore.ModelSkeletons = testassets.ExpectedModelSkeletons
		})
		It("correctly returns a BNPCInfo with the correct size", func() {
			Expect(bNPCStore.GetBNPCInfo(2, 3, 878)).To(Equal(&datasheet.BNPCInfo{
				Name: "Ruins Runner",
				Size: float32(1.2) * float32(0.2),
			}))
		})
		It("returns nil if the bNPCNameID does not exist", func() {
			Expect(bNPCStore.GetBNPCInfo(1337, 0, 0)).To(BeNil())
		})
		Context("if bNPCBaseID does not exist", func() {
			It("returns BNPCInfo with a name but default size, error = 1", func() {
				Expect(bNPCStore.GetBNPCInfo(2, 1337, 0)).To(Equal(&datasheet.BNPCInfo{
					Name:  "Ruins Runner",
					Size:  0.5,
					Error: 1,
				}))
			})
		})
		Context("if modelCharaID does not exist", func() {
			It("returns BNPCInfo with a name but scaled default size if the modelCharaID does not exist", func() {
				Expect(bNPCStore.GetBNPCInfo(2, 3, 0)).To(Equal(&datasheet.BNPCInfo{
					Name:  "Ruins Runner",
					Size:  float32(1.2) * float32(0.5),
					Error: 2,
				}))
			})
		})
		Context("if modelCharaID does not exist", func() {
			It("returns BNPCInfo with a name but scaled default size if the modelSkeleton does not exist", func() {
				Expect(bNPCStore.GetBNPCInfo(2, 3, 883)).To(Equal(&datasheet.BNPCInfo{
					Name:  "Ruins Runner",
					Size:  float32(1.2) * float32(0.5),
					Error: 3,
				}))
			})
		})
	})
})
