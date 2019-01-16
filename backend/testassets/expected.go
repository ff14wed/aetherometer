package testassets

import "github.com/ff14wed/sibyl/backend/datasheet"

// ExpectedMapInfo derives from MapJSON
var ExpectedMapInfo = map[uint16]datasheet.MapInfo{
	0: datasheet.MapInfo{SizeFactor: 100},
	1: datasheet.MapInfo{Key: 1, ID: "default/00", SizeFactor: 100, PlaceName: "Eorzea"},
	2: datasheet.MapInfo{
		Key: 2, ID: "f1t1/00", SizeFactor: 200, PlaceName: "New Gridania", TerritoryType: "f1t1",
	},
	3: datasheet.MapInfo{
		Key: 3, ID: "f1t2/00", SizeFactor: 200,
		PlaceName: "Old Gridania", TerritoryType: "f1t2",
	},
	14: datasheet.MapInfo{
		Key: 14, ID: "w1t2/01", SizeFactor: 200, PlaceName: "Ul'dah - Steps of Thal",
		PlaceNameSub: "Merchant Strip", TerritoryType: "w1t2",
	},
	73: datasheet.MapInfo{
		Key: 73, ID: "w1t2/02", SizeFactor: 200, PlaceName: "Ul'dah - Steps of Thal",
		PlaceNameSub: "Hustings Strip", TerritoryType: "w1t2",
	},
	178: datasheet.MapInfo{
		Key: 178, ID: "w1b4/00", SizeFactor: 200, OffsetX: -448, OffsetY: 0,
		PlaceName: "The Burning Heart", TerritoryType: "w1b4",
	},
	33: datasheet.MapInfo{
		Key: 33, ID: "s1fa/00", SizeFactor: 400, PlaceName: "The Navel", TerritoryType: "s1fa",
	},
	403: datasheet.MapInfo{
		Key: 403, ID: "s1fa/00", SizeFactor: 400, PlaceName: "The Navel", TerritoryType: "s1fa_2",
	},
}

var ExpectedDefaultMapsForMapIDs = map[string]uint16{
	"":           0,
	"default/00": 1,
	"f1t1/00":    2,
	"f1t2/00":    3,
	"w1t2/01":    14,
	"w1t2/02":    73,
	"w1b4/00":    178,
	"s1fa/00":    33,
}

// ExpectedTerritoryInfo derives from TerritoryTypeJSON
var ExpectedTerritoryInfo = map[uint16]datasheet.TerritoryInfo{
	1:   datasheet.TerritoryInfo{ID: 1, Name: "", MapID: ""},
	128: datasheet.TerritoryInfo{ID: 128, Name: "s1t1", MapID: "s1t1/01"},
	129: datasheet.TerritoryInfo{ID: 129, Name: "s1t2", MapID: "s1t2/00"},
	130: datasheet.TerritoryInfo{ID: 130, Name: "w1t1", MapID: "w1t1/01"},
	131: datasheet.TerritoryInfo{ID: 131, Name: "w1t2", MapID: "w1t2/01"},
	132: datasheet.TerritoryInfo{ID: 132, Name: "f1t1", MapID: "f1t1/00"},
	133: datasheet.TerritoryInfo{ID: 133, Name: "f1t2", MapID: "f1t2/00"},
	196: datasheet.TerritoryInfo{ID: 196, Name: "w1b4", MapID: "w1b4/00"},
	206: datasheet.TerritoryInfo{ID: 206, Name: "s1fa", MapID: "s1fa/00"},
	293: datasheet.TerritoryInfo{ID: 293, Name: "s1fa_2", MapID: "s1fa/00"},
	296: datasheet.TerritoryInfo{ID: 296, Name: "s1fa_3", MapID: "s1fa/00"},
}

// ExpectedActionData derives from ActionJSON
var ExpectedActionData = datasheet.ActionStore{
	0: datasheet.Action{ID: 0},
	2: datasheet.Action{
		ID: 2, Name: "Interaction", ActionCategory: "Event", ClassJob: "adventurer",
		CastType: 1, Range: 3, Cast: 50,
	},
	3: datasheet.Action{
		ID: 3, Name: "Sprint", ActionCategory: "System", ClassJob: "adventurer",
		CastType: 1, Recast: 600,
	},
	4: datasheet.Action{
		ID: 4, Name: "Mount", ActionCategory: "Item", ClassJob: "adventurer",
		CastType: 1, Cast: 10,
	},
	5: datasheet.Action{
		ID: 5, Name: "Teleport", ActionCategory: "System", ClassJob: "adventurer",
		CastType: 1, Cast: 50,
	},
	7: datasheet.Action{
		ID: 7, Name: "Attack", ActionCategory: "Auto-attack",
		ClassJob: "adventurer", Range: -1, CastType: 1,
	},
	9: datasheet.Action{
		ID: 9, Name: "Fast Blade", ActionCategory: "Weaponskill",
		ClassJob: "gladiator", Range: -1, CastType: 1, CostType: 5, Cost: 60,
		Recast: 25,
	},
	11: datasheet.Action{
		ID: 11, Name: "Savage Blade", ActionCategory: "Weaponskill",
		ClassJob: "gladiator", Range: -1, CastType: 1, CostType: 5, Cost: 60,
		ComboAction: "Fast Blade", Recast: 25,
	},
	26: datasheet.Action{
		ID: 26, Name: "Sword Oath", ActionCategory: "Spell", ClassJob: "paladin",
		CastType: 1, CostType: 3, Cost: 10, Recast: 25, GainedStatus: "Sword Oath",
	},
	50: datasheet.Action{
		ID: 50, Name: "Unchained", ActionCategory: "Ability", ClassJob: "warrior",
		CastType: 1, CostType: 0, Cost: 0, Recast: 900,
	},
	102: datasheet.Action{
		ID: 102, Name: "Flaming Arrow", ActionCategory: "Ability", ClassJob: "",
		Range: -1, TargetArea: true, CastType: 7, EffectRange: 5, Cast: 0,
		Recast: 600,
	},
	203: datasheet.Action{
		ID: 203, Name: "Skyshard", ActionCategory: "Limit Break", ClassJob: "adventurer",
		Range: 25, TargetArea: true, CastType: 2, EffectRange: 8, CostType: 11,
		Cost: 0, Cast: 20, Omen: "general_1bf",
	},
	4238: datasheet.Action{
		ID: 4238, Name: "Big Shot", ActionCategory: "Limit Break", ClassJob: "adventurer",
		Range: 30, CastType: 4, EffectRange: 30, XAxisModifier: 4, CostType: 11,
		Cast: 20, Omen: "general02f",
	},
}

// ExpectedBNPCBases derives from BNPCBaseJSON
var ExpectedBNPCBases = map[uint32]datasheet.BNPCBase{
	0: datasheet.BNPCBase{ID: 0, Scale: 1},
	1: datasheet.BNPCBase{ID: 1, Scale: 1},
	2: datasheet.BNPCBase{ID: 2, Scale: 1},
	3: datasheet.BNPCBase{ID: 3, Scale: 1.2},
}

// ExpectedBNPCNames derives from BNPCNameJSON
var ExpectedBNPCNames = map[uint32]datasheet.BNPCName{
	0: datasheet.BNPCName{ID: 0, Name: ""},
	1: datasheet.BNPCName{ID: 1, Name: ""},
	2: datasheet.BNPCName{ID: 2, Name: "ruins runner"},
	3: datasheet.BNPCName{ID: 3, Name: "antelope doe"},
}

// ExpectedModelCharas derives from ModelCharaJSON
var ExpectedModelCharas = map[uint32]datasheet.ModelChara{
	878: datasheet.ModelChara{ID: 878, Model: 8094},
	879: datasheet.ModelChara{ID: 879, Model: 8095},
	880: datasheet.ModelChara{ID: 880, Model: 8096},
	881: datasheet.ModelChara{ID: 881, Model: 8097},
	882: datasheet.ModelChara{ID: 882, Model: 8098},
	883: datasheet.ModelChara{ID: 883, Model: 8099},
}

// ExpectedModelSkeletons derives from ModelSkeletonJSON
var ExpectedModelSkeletons = map[uint32]datasheet.ModelSkeleton{
	8094: datasheet.ModelSkeleton{ID: 8094, Size: 0.2},
	8095: datasheet.ModelSkeleton{ID: 8095, Size: 0.2},
	8096: datasheet.ModelSkeleton{ID: 8096, Size: 0.2},
	8097: datasheet.ModelSkeleton{ID: 8097, Size: 0.2},
	8098: datasheet.ModelSkeleton{ID: 8098, Size: 0.2},
}
