package datasheet_test

import (
	"bytes"

	"github.com/ff14wed/sibyl/backend/datasheet"
	"github.com/ff14wed/sibyl/backend/testassets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map", func() {
	Describe("PopulateMaps", func() {
		It("correctly populates the map store", func() {
			m := new(datasheet.MapStore)
			err := m.PopulateMaps(bytes.NewReader([]byte(testassets.MapJSON)))
			Expect(err).ToNot(HaveOccurred())
			Expect(m.Maps).To(Equal(testassets.ExpectedMapInfo))
		})
		It("returns an error if the datasheet is blank", func() {
			m := new(datasheet.MapStore)
			err := m.PopulateMaps(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})
		It("returns an error if the datasheet is not valid JSON", func() {
			m := new(datasheet.MapStore)
			err := m.PopulateMaps(bytes.NewReader([]byte(InvalidJSON)))
			Expect(err).To(HaveOccurred())
		})
	})
	Describe("PopulateTerritories", func() {
		It("correctly populates the territory store", func() {
			m := new(datasheet.MapStore)
			err := m.PopulateTerritories(bytes.NewReader([]byte(testassets.TerritoryTypeJSON)))
			Expect(err).ToNot(HaveOccurred())
			Expect(m.Territories).To(Equal(testassets.ExpectedTerritoryInfo))
		})
		It("returns an error if the datasheet is blank", func() {
			m := new(datasheet.MapStore)
			err := m.PopulateTerritories(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})
		It("returns an error if the datasheet is not valid JSON", func() {
			m := new(datasheet.MapStore)
			err := m.PopulateTerritories(bytes.NewReader([]byte(InvalidJSON)))
			Expect(err).To(HaveOccurred())
		})
	})
	Describe("GetMaps", func() {
		var mapStore *datasheet.MapStore
		BeforeEach(func() {
			mapStore = new(datasheet.MapStore)
			mapStore.Maps = testassets.ExpectedMapInfo
			mapStore.Territories = testassets.ExpectedTerritoryInfo
		})
		It("correctly returns a map associated with territory ID", func() {
			Expect(mapStore.GetMaps(133)).To(Equal([]datasheet.MapInfo{
				datasheet.MapInfo{
					Key: 3, ID: "f1t2/00", SizeFactor: 200,
					PlaceName: "Old Gridania", TerritoryType: "f1t2",
				},
			}))
		})
		It("correctly returns the maps associated with territory ID", func() {
			Expect(mapStore.GetMaps(131)).To(ConsistOf(
				datasheet.MapInfo{
					Key: 14, ID: "w1t2/01", SizeFactor: 200, PlaceName: "Ul'dah - Steps of Thal",
					PlaceNameSub: "Merchant Strip", TerritoryType: "w1t2",
				},
				datasheet.MapInfo{
					Key: 73, ID: "w1t2/02", SizeFactor: 200, PlaceName: "Ul'dah - Steps of Thal",
					PlaceNameSub: "Hustings Strip", TerritoryType: "w1t2",
				},
			))
		})
		It("correctly returns maps in the case of a many to 1 (territory -> map ID) relation", func() {
			Expect(mapStore.GetMaps(206)).To(ConsistOf(
				datasheet.MapInfo{
					Key: 33, ID: "s1fa/00", SizeFactor: 400, PlaceName: "The Navel",
					TerritoryType: "s1fa",
				},
			))
			Expect(mapStore.GetMaps(293)).To(ConsistOf(
				datasheet.MapInfo{
					Key: 403, ID: "s1fa/00", SizeFactor: 400, PlaceName: "The Navel",
					TerritoryType: "s1fa_2",
				},
			))
		})
		It("returns empty array if the territory exists but the map does not", func() {
			Expect(mapStore.GetMaps(128)).To(BeNil())
		})
		It("returns nil if the territory does not exist", func() {
			Expect(mapStore.GetMaps(123)).To(BeNil())
		})
	})
})
