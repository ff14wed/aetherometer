package datasheet_test

import (
	"bytes"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/testassets"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map", func() {
	Describe("PopulateMaps", func() {
		It("correctly populates the map store", func() {
			m := new(datasheet.MapStore)
			err := m.PopulateMaps(bytes.NewReader([]byte(testassets.MapCSV)))
			Expect(err).ToNot(HaveOccurred())
			Expect(m.Maps).To(Equal(testassets.ExpectedMapInfo))
		})

		It("returns an error if the datasheet is blank", func() {
			m := new(datasheet.MapStore)
			err := m.PopulateMaps(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the datasheet is invalid", func() {
			m := new(datasheet.MapStore)
			err := m.PopulateMaps(bytes.NewReader([]byte(InvalidCSV)))
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("PopulateTerritories", func() {
		It("correctly populates the territory store", func() {
			m := new(datasheet.MapStore)
			err := m.PopulateTerritories(bytes.NewReader([]byte(testassets.TerritoryTypeCSV)))
			Expect(err).ToNot(HaveOccurred())
			Expect(m.Territories).To(Equal(testassets.ExpectedTerritoryInfo))
		})

		It("returns an error if the datasheet is blank", func() {
			m := new(datasheet.MapStore)
			err := m.PopulateTerritories(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the datasheet is invalid", func() {
			m := new(datasheet.MapStore)
			err := m.PopulateTerritories(bytes.NewReader([]byte(InvalidCSV)))
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("PopulatePlaceNames", func() {
		It("correctly populates the place names store", func() {
			m := new(datasheet.MapStore)
			err := m.PopulatePlaceNames(bytes.NewReader([]byte(testassets.PlaceNameCSV)))
			Expect(err).ToNot(HaveOccurred())
			Expect(m.PlaceNames).To(Equal(testassets.ExpectedPlaceNames))
		})

		It("returns an error if the datasheet is blank", func() {
			m := new(datasheet.MapStore)
			err := m.PopulatePlaceNames(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the datasheet is invalid", func() {
			m := new(datasheet.MapStore)
			err := m.PopulatePlaceNames(bytes.NewReader([]byte(InvalidCSV)))
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetMaps", func() {
		var mapStore *datasheet.MapStore

		BeforeEach(func() {
			mapStore = new(datasheet.MapStore)
			err := mapStore.PopulateMaps(bytes.NewReader([]byte(testassets.MapCSV)))
			Expect(err).ToNot(HaveOccurred())
			err = mapStore.PopulateTerritories(bytes.NewReader([]byte(testassets.TerritoryTypeCSV)))
			Expect(err).ToNot(HaveOccurred())
			err = mapStore.PopulatePlaceNames(bytes.NewReader([]byte(testassets.PlaceNameCSV)))
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly returns a map associated with territory ID and place names", func() {
			Expect(mapStore.GetMaps(133)).To(Equal([]models.MapInfo{
				{
					Key: 3, ID: "f1t2/00", SizeFactor: 200,
					PlaceName: "Old Gridania", TerritoryType: "f1t2",
				},
			}))
		})

		It("correctly returns (in sorted order) the maps associated with territory ID and place names", func() {
			Expect(mapStore.GetMaps(131)).To(Equal([]models.MapInfo{
				{
					Key: 14, ID: "w1t2/01", SizeFactor: 200, PlaceName: "Ul'dah - Steps of Thal",
					PlaceNameSub: "Merchant Strip", TerritoryType: "w1t2",
				},
				{
					Key: 73, ID: "w1t2/02", SizeFactor: 200, PlaceName: "Ul'dah - Steps of Thal",
					PlaceNameSub: "Hustings Strip", TerritoryType: "w1t2",
				},
			}))
		})

		It("correctly returns maps in the case of a many to 1 (territory -> map ID) relation", func() {
			Expect(mapStore.GetMaps(1046)).To(ConsistOf(
				models.MapInfo{
					Key: 33, ID: "s1fa/00", SizeFactor: 400, PlaceName: "The Navel",
					TerritoryType: "s1fa_re",
				},
			))
			Expect(mapStore.GetMaps(293)).To(ConsistOf(
				models.MapInfo{
					Key: 403, ID: "s1fa/00", SizeFactor: 400, PlaceName: "The Navel",
					TerritoryType: "s1fa_2",
				},
			))
		})

		It("returns the territory's Map if there are no maps associated the territory ID", func() {
			Expect(mapStore.GetMaps(296)).To(ConsistOf(
				models.MapInfo{
					Key: 403, ID: "s1fa/00", SizeFactor: 400, PlaceName: "The Navel",
					TerritoryType: "s1fa_2",
				},
			))
		})

		It("returns empty array if the territory exists but its map does not", func() {
			Expect(mapStore.GetMaps(128)).To(BeEmpty())
		})

		It("returns empty array if the territory does not exist", func() {
			Expect(mapStore.GetMaps(123)).To(BeEmpty())
		})
	})
})
