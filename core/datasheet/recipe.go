package datasheet

import (
	"fmt"
	"io"

	"github.com/ff14wed/aetherometer/core/models"
)

// RecipeStore stores all of the Recipe data.
type RecipeStore struct {
	Recipes          map[uint32]Recipe
	RecipeLevelTable map[uint16]RecipeLevel
	Items            map[uint32]Item
}

// Recipe stores some of the data for a game Recipe
type Recipe struct {
	Key              uint32 `datasheet:"key"`
	RecipeLevel      uint16 `datasheet:"RecipeLevelTable"`
	ItemID           uint32 `datasheet:"Item{Result}"`
	RecipeElement    byte   `datasheet:"RecipeElement"`
	DifficultyFactor uint16 `datasheet:"DifficultyFactor"`
	QualityFactor    uint16 `datasheet:"QualityFactor"`
	DurabilityFactor uint16 `datasheet:"DurabilityFactor"`
	CanHQ            bool   `datasheet:"CanHq"`
}

// RecipeLevel stores information about the difficulty of the recipe to
// improve quality and and increase progress
type RecipeLevel struct {
	Key        uint16 `datasheet:"key"`
	Difficulty uint16 `datasheet:"Difficulty"`
	Quality    uint16 `datasheet:"Quality"`
	Durability uint16 `datasheet:"Durability"`
}

// Item maps the ID of the item with the name
type Item struct {
	Key  uint32 `datasheet:"key"`
	Name string `datasheet:"Name"`
}

// PopulateRecipes will populate the RecipeStore with Recipe data provided a
// path to the data sheet for Recipes.
func (r *RecipeStore) PopulateRecipes(dataReader io.Reader) error {
	r.Recipes = make(map[uint32]Recipe)

	var rows []Recipe
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateRecipes: %s", err)
	}
	for _, recipe := range rows {
		r.Recipes[recipe.Key] = recipe
	}
	return nil
}

// PopulateRecipeLevelTable will populate the RecipeStore with RecipeLevel data
// provided a path to the data sheet for RecipeLevelTable.
func (r *RecipeStore) PopulateRecipeLevelTable(dataReader io.Reader) error {
	r.RecipeLevelTable = make(map[uint16]RecipeLevel)

	var rows []RecipeLevel
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateRecipeLevelTable: %s", err)
	}
	for _, recipeLevel := range rows {
		r.RecipeLevelTable[recipeLevel.Key] = recipeLevel
	}
	return nil
}

// PopulateItems will populate the RecipeStore with Item data
// provided a path to the data sheet for Items.
func (r *RecipeStore) PopulateItems(dataReader io.Reader) error {
	r.Items = make(map[uint32]Item)

	var rows []Item
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateItems: %s", err)
	}
	for _, item := range rows {
		r.Items[item.Key] = item
	}
	return nil
}

// GetInfo returns the normalized information about the recipe
func (r *RecipeStore) GetInfo(key uint32) models.RecipeInfo {
	recipe, found := r.Recipes[key]
	if !found {
		return models.RecipeInfo{}
	}
	info := models.RecipeInfo{
		ID:          int(key),
		RecipeLevel: int(recipe.RecipeLevel),
		Element:     int(recipe.RecipeElement),
		CanHQ:       recipe.CanHQ,
	}
	if recipeLevel, found := r.RecipeLevelTable[recipe.RecipeLevel]; found {
		info.Difficulty = int(float64(recipe.DifficultyFactor) * float64(recipeLevel.Difficulty) / 100)
		info.Quality = int(float64(recipe.QualityFactor) * float64(recipeLevel.Quality) / 100)
		info.Durability = int(float64(recipe.DurabilityFactor) * float64(recipeLevel.Durability) / 100)
	}
	if item, found := r.Items[recipe.ItemID]; found {
		info.Name = item.Name
	}

	return info
}
