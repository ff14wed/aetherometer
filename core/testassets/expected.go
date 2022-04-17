package testassets

import "github.com/ff14wed/aetherometer/core/datasheet"

// ExpectedActionData derives from ActionCSV
var ExpectedActionData = map[uint32]datasheet.Action{
	0:    {Key: 0},
	2:    {Key: 2, Name: "Interaction", CastType: 1, Range: 3},
	3:    {Key: 3, Name: "Sprint", CastType: 1},
	4:    {Key: 4, Name: "Mount", CastType: 1},
	5:    {Key: 5, Name: "Teleport", CastType: 1},
	7:    {Key: 7, Name: "attack", Range: -1, CastType: 1},
	9:    {Key: 9, Name: "Fast Blade", Range: -1, CastType: 1},
	11:   {Key: 11, Name: "Savage Blade", Range: -1, CastType: 1},
	26:   {Key: 26, Name: "Sword Oath", CastType: 1},
	50:   {Key: 50, Name: "Unchained", CastType: 1},
	102:  {Key: 102, Name: "Flaming Arrow", Range: -1, TargetArea: true, CastType: 7, EffectRange: 5},
	203:  {Key: 203, Name: "Skyshard", Range: 25, TargetArea: true, CastType: 2, EffectRange: 8, OmenID: 1},
	4238: {Key: 4238, Name: "Big Shot", Range: 30, CastType: 4, EffectRange: 30, XAxisModifier: 4, OmenID: 2},
}

// ExpectedOmenData derives from OmenCSV
var ExpectedOmenData = map[uint16]datasheet.Omen{
	0: {Key: 0, Name: ""},
	1: {Key: 1, Name: "general_1bf"},
	2: {Key: 2, Name: "general02f"},
}

// ExpectedCraftActionData derives from CraftActionCSV
var ExpectedCraftActionData = map[uint32]datasheet.CraftAction{
	100000: {Key: 100000, Name: ""},
	100001: {Key: 100001, Name: "Basic Synthesis"},
	100002: {Key: 100002, Name: "Basic Touch"},
}

// ExpectedBNPCBases derives from BNPCBaseCSV
var ExpectedBNPCBases = map[uint32]datasheet.BNPCBase{
	0: {Key: 0, Scale: 1},
	1: {Key: 1, Scale: 1},
	2: {Key: 2, Scale: 1},
	3: {Key: 3, Scale: 1.2},
}

// ExpectedBNPCNames derives from BNPCNameCSV
var ExpectedBNPCNames = map[uint32]datasheet.BNPCName{
	0: {Key: 0, Name: ""},
	1: {Key: 1, Name: ""},
	2: {Key: 2, Name: "ruins runner"},
	3: {Key: 3, Name: "antelope doe"},
}

// ExpectedModelCharas derives from ModelCharaCSV
var ExpectedModelCharas = map[uint32]datasheet.ModelChara{
	878: {Key: 878, Model: 8094},
	879: {Key: 879, Model: 8095},
	880: {Key: 880, Model: 8096},
	881: {Key: 881, Model: 8097},
	882: {Key: 882, Model: 8098},
	883: {Key: 883, Model: 8099},
}

// ExpectedModelSkeletons derives from ModelSkeletonCSV
var ExpectedModelSkeletons = map[uint32]datasheet.ModelSkeleton{
	8094: {Key: 8094, Radius: 0.2},
	8095: {Key: 8095, Radius: 0.2},
	8096: {Key: 8096, Radius: 0.2},
	8097: {Key: 8097, Radius: 0.2},
	8098: {Key: 8098, Radius: 0.2},
}

// ExpectedMapInfo derives from MapCSV
var ExpectedMapInfo = map[uint16]datasheet.MapInfo{
	0:   {SizeFactor: 100},
	1:   {Key: 1, ID: "default/00", SizeFactor: 100, PlaceName: 21, TerritoryType: 1},
	2:   {Key: 2, ID: "f1t1/00", SizeFactor: 200, PlaceName: 52, TerritoryType: 132},
	3:   {Key: 3, ID: "f1t2/00", SizeFactor: 200, PlaceName: 53, TerritoryType: 133},
	14:  {Key: 14, ID: "w1t2/01", SizeFactor: 200, PlaceName: 41, PlaceNameSub: 373, TerritoryType: 131},
	73:  {Key: 73, ID: "w1t2/02", SizeFactor: 200, PlaceName: 41, PlaceNameSub: 698, TerritoryType: 131},
	178: {Key: 178, ID: "w1b4/00", SizeFactor: 200, OffsetX: -448, OffsetY: 0, PlaceName: 1409, TerritoryType: 196},
	33:  {Key: 33, ID: "s1fa/00", SizeFactor: 400, PlaceName: 359, TerritoryType: 1046},
	403: {Key: 403, ID: "s1fa/00", SizeFactor: 400, PlaceName: 359, PlaceNameSub: 19, TerritoryType: 293},
}

// ExpectedPlaceNames derives from PlaceNameCSV
var ExpectedPlaceNames = map[uint16]datasheet.PlaceName{
	0:    {Key: 0, Name: ""},
	19:   {Key: 19, Name: ""},
	21:   {Key: 21, Name: "Eorzea"},
	41:   {Key: 41, Name: "Ul'dah - Steps of Thal"},
	52:   {Key: 52, Name: "New Gridania"},
	53:   {Key: 53, Name: "Old Gridania"},
	359:  {Key: 359, Name: "The Navel"},
	373:  {Key: 373, Name: "Merchant Strip"},
	698:  {Key: 698, Name: "Hustings Strip"},
	1409: {Key: 1409, Name: "The Burning Heart"},
}

// ExpectedTerritoryInfo derives from TerritoryTypeCSV
var ExpectedTerritoryInfo = map[uint16]datasheet.TerritoryInfo{
	1:    {Key: 1, Name: "", Map: 0},
	128:  {Key: 128, Name: "s1t1", Map: 11},
	129:  {Key: 129, Name: "s1t2", Map: 12},
	130:  {Key: 130, Name: "w1t1", Map: 13},
	131:  {Key: 131, Name: "w1t2", Map: 14},
	132:  {Key: 132, Name: "f1t1", Map: 2},
	133:  {Key: 133, Name: "f1t2", Map: 3},
	196:  {Key: 196, Name: "w1b4", Map: 178},
	293:  {Key: 293, Name: "s1fa_2", Map: 403},
	296:  {Key: 296, Name: "s1fa_3", Map: 403},
	1046: {Key: 1046, Name: "s1fa_re", Map: 33},
}

// ExpectedStatusData derives from StatusCSV
var ExpectedStatusData = map[uint32]datasheet.Status{
	0: {Key: 0},
	1: {Key: 1, Name: "Petrification", Description: "Stone-like rigidity is preventing the execution of actions."},
	2: {Key: 2, Name: "Stun", Description: "Unable to execute actions."},
}

// ExpectedClassJobData derives from ClassJobCSV
var ExpectedClassJobData = map[byte]datasheet.ClassJob{
	0: {Key: 0, Name: "adventurer", Abbreviation: "ADV"},
	1: {Key: 1, Name: "gladiator", Abbreviation: "GLA"},
	2: {Key: 2, Name: "pugilist", Abbreviation: "PGL"},
}

// ExpectedRecipeData derives from RecipeCSV
var ExpectedRecipeData = map[uint32]datasheet.Recipe{
	1: {
		Key:              1,
		RecipeLevel:      1,
		ItemID:           5056,
		RecipeElement:    0,
		DifficultyFactor: 50,
		QualityFactor:    80,
		DurabilityFactor: 67,
		CanHQ:            true,
	},
	33067: {
		Key:              33067,
		RecipeLevel:      320,
		ItemID:           23002,
		RecipeElement:    0,
		DifficultyFactor: 100,
		QualityFactor:    100,
		DurabilityFactor: 100,
		CanHQ:            true,
	},
	33068: {
		Key:              33068,
		RecipeLevel:      320,
		ItemID:           23374,
		RecipeElement:    0,
		DifficultyFactor: 100,
		QualityFactor:    100,
		DurabilityFactor: 100,
		CanHQ:            true,
	},
	33073: {
		Key:              33073,
		RecipeLevel:      380,
		ItemID:           23768,
		RecipeElement:    0,
		DifficultyFactor: 100,
		QualityFactor:    100,
		DurabilityFactor: 100,
		CanHQ:            true,
	},
	33074: {
		Key:              33074,
		RecipeLevel:      380,
		ItemID:           23769,
		RecipeElement:    0,
		DifficultyFactor: 100,
		QualityFactor:    100,
		DurabilityFactor: 100,
		CanHQ:            true,
	},
}

// ExpectedRecipeLevelTableData derives from RecipeLevelTableCSV
var ExpectedRecipeLevelTableData = map[uint16]datasheet.RecipeLevel{
	320: {Key: 320, Difficulty: 1200, Quality: 4800, Durability: 70},
	380: {Key: 380, Difficulty: 1500, Quality: 6100, Durability: 70},
}

// ExpectedItemData derives from ItemCSV
var ExpectedItemData = map[uint32]datasheet.Item{
	5056:  {Key: 5056, Name: "Bronze Ingot"},
	23374: {Key: 23374, Name: "Quaintrelle's Dress Shoes"},
	23768: {Key: 23768, Name: "Rakshasa Blade"},
	23769: {Key: 23769, Name: "Rakshasa Knuckles"},
}

// ExpectedWorldData derives from WorldCSV
var ExpectedWorldData = map[uint32]datasheet.World{
	0: {Key: 0, Name: "crossworld"},
	1: {Key: 1, Name: "reserved1"},
	2: {Key: 2, Name: "c-contents"},
	3: {Key: 3, Name: "c-whiteae"},
	4: {Key: 4, Name: "c-baudinii"},
	5: {Key: 5, Name: "c-contents2"},
}
