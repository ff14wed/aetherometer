package datasheet_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ff14wed/aetherometer/core/datasheet"
	"github.com/ff14wed/aetherometer/core/testassets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Collection", func() {
	Describe("Populate", func() {
		var tmpDir string
		BeforeEach(func() {
			var err error
			tmpDir, err = ioutil.TempDir("", "collection-test")
			Expect(err).ToNot(HaveOccurred())

			fileToTestAssets := map[string]string{
				"Map.csv":           testassets.MapCSV,
				"TerritoryType.csv": testassets.TerritoryTypeCSV,
				"PlaceName.csv":     testassets.PlaceNameCSV,

				"BNpcName.csv":      testassets.BNPCNameCSV,
				"BNpcBase.csv":      testassets.BNPCBaseCSV,
				"ModelChara.csv":    testassets.ModelCharaCSV,
				"ModelSkeleton.csv": testassets.ModelSkeletonCSV,

				"Action.csv":      testassets.ActionCSV,
				"Omen.csv":        testassets.OmenCSV,
				"CraftAction.csv": testassets.CraftActionCSV,

				"Status.csv": testassets.StatusCSV,

				"ClassJob.csv": testassets.ClassJobCSV,

				"Recipe.csv":           testassets.RecipeCSV,
				"RecipeLevelTable.csv": testassets.RecipeLevelTableCSV,
				"Item.csv":             testassets.ItemCSV,
			}

			for name, contents := range fileToTestAssets {
				err := ioutil.WriteFile(filepath.Join(tmpDir, name), []byte(contents), 0777)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		AfterEach(func() {
			Expect(os.RemoveAll(tmpDir)).To(Succeed())
		})

		It("successfully reads in the into the Collection", func() {
			collection := new(datasheet.Collection)
			err := collection.Populate(tmpDir)
			Expect(err).ToNot(HaveOccurred())
			Expect(collection.MapData.Maps).ToNot(BeEmpty())
		})

		Context("when there is an error reading in a file", func() {
			BeforeEach(func() {
				err := os.Remove(filepath.Join(tmpDir, "Map.csv"))
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns an error and does not read in further files", func() {
				collection := new(datasheet.Collection)
				err := collection.Populate(tmpDir)
				Expect(err).To(MatchError(MatchRegexp("Map.csv: no such file or directory")))
				Expect(collection.MapData.Territories).To(BeEmpty())
			})
		})
	})
})
