package datasheet

import (
	"io"
	"path/filepath"
)

// Collection encapsulates a collection of datasheets
type Collection struct {
	MapData      MapStore
	BNPCData     BNPCStore
	ActionData   ActionStore
	StatusData   StatusStore
	ClassJobData ClassJobStore
	RecipeData   RecipeStore
}

type dataTuple struct {
	path      string
	populator func(io.Reader) error
}

func (c *Collection) Populate(dataPath string) error {
	dataMapping := []dataTuple{
		{filepath.Join(dataPath, "Map.csv"), c.MapData.PopulateMaps},
		{filepath.Join(dataPath, "TerritoryType.csv"), c.MapData.PopulateTerritories},
		{filepath.Join(dataPath, "PlaceName.csv"), c.MapData.PopulatePlaceNames},

		{filepath.Join(dataPath, "BNpcName.csv"), c.BNPCData.PopulateBNPCNames},
		{filepath.Join(dataPath, "BNpcBase.csv"), c.BNPCData.PopulateBNPCBases},
		{filepath.Join(dataPath, "ModelChara.csv"), c.BNPCData.PopulateModelCharas},
		{filepath.Join(dataPath, "ModelSkeleton.csv"), c.BNPCData.PopulateModelSkeletons},

		{filepath.Join(dataPath, "Action.csv"), c.ActionData.PopulateActions},
		{filepath.Join(dataPath, "Omen.csv"), c.ActionData.PopulateOmens},
		{filepath.Join(dataPath, "CraftAction.csv"), c.ActionData.PopulateCraftActions},

		{filepath.Join(dataPath, "Status.csv"), c.StatusData.PopulateStatuses},

		{filepath.Join(dataPath, "ClassJob.csv"), c.ClassJobData.PopulateClassJobs},

		{filepath.Join(dataPath, "Recipe.csv"), c.RecipeData.PopulateRecipes},
		{filepath.Join(dataPath, "RecipeLevelTable.csv"), c.RecipeData.PopulateRecipeLevelTable},
		{filepath.Join(dataPath, "Item.csv"), c.RecipeData.PopulateItems},
	}
	fileReader := new(FileReader)
	for _, t := range dataMapping {
		fileReader.ReadFile(t.path, t.populator)
	}
	return fileReader.Error()
}
