package testassets_test

import (
	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/testassets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test Asset", func() {
	var collection *datasheet.Collection
	BeforeEach(func() {
		collection = new(datasheet.Collection)
		Expect(collection.Populate("../../resources/datasheets")).To(Succeed())
	})

	It("is up to date with the action CSV", func() {
		for k, v := range testassets.ExpectedActionData {
			Expect(collection.ActionData.Actions).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the omen CSV", func() {
		for k, v := range testassets.ExpectedOmenData {
			Expect(collection.ActionData.Omens).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the craft action CSV", func() {
		for k, v := range testassets.ExpectedCraftActionData {
			Expect(collection.ActionData.CraftActions).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the BNPCBase CSV", func() {
		for k, v := range testassets.ExpectedBNPCBases {
			Expect(collection.BNPCData.BNPCBases).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the BNPCName CSV", func() {
		for k, v := range testassets.ExpectedBNPCNames {
			Expect(collection.BNPCData.BNPCNames).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the ModelChara CSV", func() {
		for k, v := range testassets.ExpectedModelCharas {
			Expect(collection.BNPCData.ModelCharas).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the ModelSkeleton CSV", func() {
		for k, v := range testassets.ExpectedModelSkeletons {
			Expect(collection.BNPCData.ModelSkeletons).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the MapInfo CSV", func() {
		for k, v := range testassets.ExpectedMapInfo {
			Expect(collection.MapData.Maps).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the PlaceNames CSV", func() {
		for k, v := range testassets.ExpectedPlaceNames {
			Expect(collection.MapData.PlaceNames).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the TerritoryType CSV", func() {
		for k, v := range testassets.ExpectedTerritoryInfo {
			Expect(collection.MapData.Territories).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the Status CSV", func() {
		for k, v := range testassets.ExpectedStatusData {
			Expect(collection.StatusData).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the ClassJob CSV", func() {
		for k, v := range testassets.ExpectedClassJobData {
			Expect(collection.ClassJobData).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the Recipe CSV", func() {
		for k, v := range testassets.ExpectedRecipeData {
			Expect(collection.RecipeData.Recipes).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the RecipeLevelTable CSV", func() {
		for k, v := range testassets.ExpectedRecipeLevelTableData {
			Expect(collection.RecipeData.RecipeLevelTable).To(HaveKeyWithValue(k, v))
		}
	})

	It("is up to date with the Item CSV", func() {
		for k, v := range testassets.ExpectedItemData {
			Expect(collection.RecipeData.Items).To(HaveKeyWithValue(k, v))
		}
	})
})
