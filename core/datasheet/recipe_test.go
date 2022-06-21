package datasheet_test

import (
	"bytes"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/models"
	"github.com/ff14wed/aetherometer/core/testassets"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Recipe", func() {
	Describe("PopulateRecipes", func() {
		It("correctly populates the Recipes store", func() {
			var r datasheet.RecipeStore
			err := r.PopulateRecipes(bytes.NewReader([]byte(testassets.RecipeCSV)))
			Expect(err).ToNot(HaveOccurred())
			Expect(r.Recipes).To(HaveLen(len(testassets.ExpectedRecipeData)))
			for k, d := range r.Recipes {
				Expect(d).To(Equal(testassets.ExpectedRecipeData[k]))
			}
		})

		It("returns an error if the datasheet is blank", func() {
			var r datasheet.RecipeStore
			err := r.PopulateRecipes(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the datasheet is invalid", func() {
			var r datasheet.RecipeStore
			err := r.PopulateRecipes(bytes.NewReader([]byte(InvalidCSV)))
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("PopulateRecipeLevelTable", func() {
		It("correctly populates the RecipeLevelTable store", func() {
			var r datasheet.RecipeStore
			err := r.PopulateRecipeLevelTable(bytes.NewReader([]byte(testassets.RecipeLevelTableCSV)))
			Expect(err).ToNot(HaveOccurred())
			Expect(r.RecipeLevelTable).To(HaveLen(len(testassets.ExpectedRecipeLevelTableData)))
			for k, d := range r.RecipeLevelTable {
				Expect(d).To(Equal(testassets.ExpectedRecipeLevelTableData[k]))
			}
		})

		It("returns an error if the datasheet is blank", func() {
			var r datasheet.RecipeStore
			err := r.PopulateRecipeLevelTable(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the datasheet is invalid", func() {
			var r datasheet.RecipeStore
			err := r.PopulateRecipeLevelTable(bytes.NewReader([]byte(InvalidCSV)))
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("PopulateItems", func() {
		It("correctly populates the Items store", func() {
			var r datasheet.RecipeStore
			err := r.PopulateItems(bytes.NewReader([]byte(testassets.ItemCSV)))
			Expect(err).ToNot(HaveOccurred())
			Expect(r.Items).To(HaveLen(len(testassets.ExpectedItemData)))
			for k, d := range r.Items {
				Expect(d).To(Equal(testassets.ExpectedItemData[k]))
			}
		})

		It("returns an error if the datasheet is blank", func() {
			var r datasheet.RecipeStore
			err := r.PopulateItems(bytes.NewReader([]byte("")))
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the datasheet is invalid", func() {
			var r datasheet.RecipeStore
			err := r.PopulateItems(bytes.NewReader([]byte(InvalidCSV)))
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetInfo", func() {
		var r *datasheet.RecipeStore

		BeforeEach(func() {
			r = new(datasheet.RecipeStore)
			Expect(r.PopulateRecipes(bytes.NewReader([]byte(testassets.RecipeCSV)))).To(Succeed())
			Expect(r.PopulateRecipeLevelTable(bytes.NewReader([]byte(testassets.RecipeLevelTableCSV)))).To(Succeed())
			Expect(r.PopulateItems(bytes.NewReader([]byte(testassets.ItemCSV)))).To(Succeed())
		})

		It("gets the full recipe details", func() {
			Expect(r.GetInfo(33074)).To(Equal(&models.RecipeInfo{
				ID:          33074,
				Name:        "Rakshasa Knuckles",
				RecipeLevel: 380,
				ItemID:      23769,
				Element:     0,
				CanHq:       true,
				Difficulty:  1500,
				Quality:     6100,
				Durability:  70,
			}))
		})

		It("returns incomplete information if there is no RecipeLevel for the recipe", func() {
			Expect(r.GetInfo(1)).To(Equal(&models.RecipeInfo{
				ID:          1,
				Name:        "Bronze Ingot",
				RecipeLevel: 1,
				ItemID:      5056,
				Element:     0,
				CanHq:       true,
			}))
		})

		It("returns incomplete information if there is no Item for the recipe", func() {
			Expect(r.GetInfo(33067)).To(Equal(&models.RecipeInfo{
				ID:          33067,
				RecipeLevel: 320,
				ItemID:      23002,
				Element:     0,
				CanHq:       true,
				Difficulty:  1200,
				Quality:     4800,
				Durability:  70,
			}))
		})
	})
})
