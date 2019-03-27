package testassets

import "github.com/ff14wed/sibyl/backend/datasheet"

// ExpectedActionData derives from ActionCSV
var ExpectedActionData = map[uint32]datasheet.Action{
	0:    datasheet.Action{Key: 0},
	2:    datasheet.Action{Key: 2, Name: "Interaction", CastType: 1, Range: 3},
	3:    datasheet.Action{Key: 3, Name: "Sprint", CastType: 1},
	4:    datasheet.Action{Key: 4, Name: "Mount", CastType: 1},
	5:    datasheet.Action{Key: 5, Name: "Teleport", CastType: 1},
	7:    datasheet.Action{Key: 7, Name: "attack", Range: -1, CastType: 1},
	9:    datasheet.Action{Key: 9, Name: "Fast Blade", Range: -1, CastType: 1},
	11:   datasheet.Action{Key: 11, Name: "Savage Blade", Range: -1, CastType: 1},
	26:   datasheet.Action{Key: 26, Name: "Sword Oath", CastType: 1},
	50:   datasheet.Action{Key: 50, Name: "Unchained", CastType: 1},
	102:  datasheet.Action{Key: 102, Name: "Flaming Arrow", Range: -1, TargetArea: true, CastType: 7, EffectRange: 5},
	203:  datasheet.Action{Key: 203, Name: "Skyshard", Range: 25, TargetArea: true, CastType: 2, EffectRange: 8, OmenID: 1},
	4238: datasheet.Action{Key: 4238, Name: "Big Shot", Range: 30, CastType: 4, EffectRange: 30, XAxisModifier: 4, OmenID: 2},
}

// ExpectedOmenData derives from OmenCSV
var ExpectedOmenData = map[uint16]datasheet.Omen{
	0: datasheet.Omen{Key: 0, Name: ""},
	1: datasheet.Omen{Key: 1, Name: "general_1bf"},
	2: datasheet.Omen{Key: 2, Name: "general02f"},
}

// ExpectedBNPCBases derives from BNPCBaseCSV
var ExpectedBNPCBases = map[uint32]datasheet.BNPCBase{
	0: datasheet.BNPCBase{Key: 0, Scale: 1},
	1: datasheet.BNPCBase{Key: 1, Scale: 1},
	2: datasheet.BNPCBase{Key: 2, Scale: 1},
	3: datasheet.BNPCBase{Key: 3, Scale: 1.2},
}

// ExpectedBNPCNames derives from BNPCNameCSV
var ExpectedBNPCNames = map[uint32]datasheet.BNPCName{
	0: datasheet.BNPCName{Key: 0, Name: ""},
	1: datasheet.BNPCName{Key: 1, Name: ""},
	2: datasheet.BNPCName{Key: 2, Name: "ruins runner"},
	3: datasheet.BNPCName{Key: 3, Name: "antelope doe"},
}

// ExpectedModelCharas derives from ModelCharaCSV
var ExpectedModelCharas = map[uint32]datasheet.ModelChara{
	878: datasheet.ModelChara{Key: 878, Model: 8094},
	879: datasheet.ModelChara{Key: 879, Model: 8095},
	880: datasheet.ModelChara{Key: 880, Model: 8096},
	881: datasheet.ModelChara{Key: 881, Model: 8097},
	882: datasheet.ModelChara{Key: 882, Model: 8098},
	883: datasheet.ModelChara{Key: 883, Model: 8099},
}

// ExpectedModelSkeletons derives from ModelSkeletonCSV
var ExpectedModelSkeletons = map[uint32]datasheet.ModelSkeleton{
	8094: datasheet.ModelSkeleton{Key: 8094, Size: 0.2},
	8095: datasheet.ModelSkeleton{Key: 8095, Size: 0.2},
	8096: datasheet.ModelSkeleton{Key: 8096, Size: 0.2},
	8097: datasheet.ModelSkeleton{Key: 8097, Size: 0.2},
	8098: datasheet.ModelSkeleton{Key: 8098, Size: 0.2},
}

// ExpectedMapInfo derives from MapCSV
var ExpectedMapInfo = map[uint16]datasheet.MapInfo{
	0: datasheet.MapInfo{SizeFactor: 100},
	1: datasheet.MapInfo{Key: 1, ID: "default/00", SizeFactor: 100, PlaceName: 21,
		TerritoryType: 1,
	},
	2: datasheet.MapInfo{
		Key: 2, ID: "f1t1/00", SizeFactor: 200, PlaceName: 52, TerritoryType: 132,
	},
	3: datasheet.MapInfo{
		Key: 3, ID: "f1t2/00", SizeFactor: 200, PlaceName: 53, TerritoryType: 133,
	},
	14: datasheet.MapInfo{
		Key: 14, ID: "w1t2/01", SizeFactor: 200, PlaceName: 41, PlaceNameSub: 373,
		TerritoryType: 131,
	},
	73: datasheet.MapInfo{
		Key: 73, ID: "w1t2/02", SizeFactor: 200, PlaceName: 41, PlaceNameSub: 698,
		TerritoryType: 131,
	},
	178: datasheet.MapInfo{
		Key: 178, ID: "w1b4/00", SizeFactor: 200, OffsetX: -448, OffsetY: 0,
		PlaceName: 1409, TerritoryType: 196,
	},
	33: datasheet.MapInfo{
		Key: 33, ID: "s1fa/00", SizeFactor: 400, PlaceName: 359, TerritoryType: 206,
	},
	403: datasheet.MapInfo{
		Key: 403, ID: "s1fa/00", SizeFactor: 400, PlaceName: 359, PlaceNameSub: 19,
		TerritoryType: 293,
	},
}

// ExpectedPlaceNames derives from PlaceNameCSV
var ExpectedPlaceNames = map[uint16]datasheet.PlaceName{
	0:    datasheet.PlaceName{Key: 0, Name: ""},
	19:   datasheet.PlaceName{Key: 19, Name: ""},
	21:   datasheet.PlaceName{Key: 21, Name: "Eorzea"},
	41:   datasheet.PlaceName{Key: 41, Name: "Ul'dah - Steps of Thal"},
	52:   datasheet.PlaceName{Key: 52, Name: "New Gridania"},
	53:   datasheet.PlaceName{Key: 53, Name: "Old Gridania"},
	359:  datasheet.PlaceName{Key: 359, Name: "The Navel"},
	373:  datasheet.PlaceName{Key: 373, Name: "Merchant Strip"},
	698:  datasheet.PlaceName{Key: 698, Name: "Hustings Strip"},
	1409: datasheet.PlaceName{Key: 1409, Name: "The Burning Heart"},
}

// ExpectedTerritoryInfo derives from TerritoryTypeCSV
var ExpectedTerritoryInfo = map[uint16]datasheet.TerritoryInfo{
	1:   datasheet.TerritoryInfo{Key: 1, Name: "", Map: 0},
	128: datasheet.TerritoryInfo{Key: 128, Name: "s1t1", Map: 11},
	129: datasheet.TerritoryInfo{Key: 129, Name: "s1t2", Map: 12},
	130: datasheet.TerritoryInfo{Key: 130, Name: "w1t1", Map: 13},
	131: datasheet.TerritoryInfo{Key: 131, Name: "w1t2", Map: 14},
	132: datasheet.TerritoryInfo{Key: 132, Name: "f1t1", Map: 2},
	133: datasheet.TerritoryInfo{Key: 133, Name: "f1t2", Map: 3},
	196: datasheet.TerritoryInfo{Key: 196, Name: "w1b4", Map: 178},
	206: datasheet.TerritoryInfo{Key: 206, Name: "s1fa", Map: 33},
	293: datasheet.TerritoryInfo{Key: 293, Name: "s1fa_2", Map: 403},
	296: datasheet.TerritoryInfo{Key: 296, Name: "s1fa_3", Map: 403},
}

var ExpectedStatusData = map[uint32]datasheet.Status{
	0: datasheet.Status{Key: 0},
	1: datasheet.Status{Key: 1, Name: "Petrification", Description: "Stone-like rigidity is preventing the execution of actions."},
	2: datasheet.Status{Key: 2, Name: "Stun", Description: "Unable to execute actions."},
}

var ExpectedClassJobData = map[byte]datasheet.ClassJob{
	0: datasheet.ClassJob{Key: 0, Name: "adventurer", Abbreviation: "ADV"},
	1: datasheet.ClassJob{Key: 1, Name: "gladiator", Abbreviation: "GLA"},
	2: datasheet.ClassJob{Key: 2, Name: "pugilist", Abbreviation: "PGL"},
}

var ExpectedRecipeData = map[uint32]datasheet.Recipe{
	1: datasheet.Recipe{
		Key:              1,
		RecipeLevel:      1,
		ItemID:           5056,
		RecipeElement:    0,
		DifficultyFactor: 50,
		QualityFactor:    100,
		DurabilityFactor: 67,
		CanHQ:            true,
	},
	33067: datasheet.Recipe{
		Key:              33067,
		RecipeLevel:      320,
		ItemID:           23002,
		RecipeElement:    0,
		DifficultyFactor: 100,
		QualityFactor:    100,
		DurabilityFactor: 100,
		CanHQ:            true,
	},
	33068: datasheet.Recipe{
		Key:              33068,
		RecipeLevel:      320,
		ItemID:           23374,
		RecipeElement:    0,
		DifficultyFactor: 100,
		QualityFactor:    100,
		DurabilityFactor: 100,
		CanHQ:            true,
	},
	33073: datasheet.Recipe{
		Key:              33073,
		RecipeLevel:      380,
		ItemID:           23768,
		RecipeElement:    0,
		DifficultyFactor: 90,
		QualityFactor:    140,
		DurabilityFactor: 100,
		CanHQ:            true,
	},
	33074: datasheet.Recipe{
		Key:              33074,
		RecipeLevel:      380,
		ItemID:           23769,
		RecipeElement:    0,
		DifficultyFactor: 90,
		QualityFactor:    140,
		DurabilityFactor: 100,
		CanHQ:            true,
	},
}

var ExpectedRecipeLevelTableData = map[uint16]datasheet.RecipeLevel{
	320: datasheet.RecipeLevel{Key: 320, Difficulty: 3543, Quality: 15837, Durability: 70},
	380: datasheet.RecipeLevel{Key: 380, Difficulty: 4143, Quality: 21137, Durability: 70},
}

var ExpectedItemData = map[uint32]datasheet.Item{
	5056:  datasheet.Item{Key: 5056, Name: "Bronze Ingot"},
	23374: datasheet.Item{Key: 23374, Name: "Quaintrelle's Dress Shoes"},
	23768: datasheet.Item{Key: 23768, Name: "Rakshasa Blade"},
	23769: datasheet.Item{Key: 23769, Name: "Rakshasa Knuckles"},
}
