// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"time"
)

type Action struct {
	TargetID          uint64         `json:"targetID"`
	Name              string         `json:"name"`
	GlobalCounter     int            `json:"globalCounter"`
	AnimationLockTime float64        `json:"animationLockTime"`
	HiddenAnimation   int            `json:"hiddenAnimation"`
	Location          Location       `json:"location"`
	ID                int            `json:"id"`
	Variation         int            `json:"variation"`
	EffectDisplayType int            `json:"effectDisplayType"`
	IsAoE             bool           `json:"isAoE"`
	Effects           []ActionEffect `json:"effects"`
	EffectFlags       int            `json:"effectFlags"`
	UseTime           time.Time      `json:"useTime"`
}

type ActionEffect struct {
	TargetID        uint64 `json:"targetID"`
	Type            int    `json:"type"`
	HitSeverity     int    `json:"hitSeverity"`
	Param           int    `json:"param"`
	BonusPercent    int    `json:"bonusPercent"`
	ValueMultiplier int    `json:"valueMultiplier"`
	Flags           int    `json:"flags"`
	Value           int    `json:"value"`
}

type AddEntity struct {
	Entity Entity `json:"entity"`
}

func (AddEntity) IsEntityEventType() {}

type AddStream struct {
	Stream Stream `json:"stream"`
}

func (AddStream) IsStreamEventType() {}

type CastingInfo struct {
	ActionID      int       `json:"actionID"`
	ActionName    string    `json:"actionName"`
	StartTime     time.Time `json:"startTime"`
	CastTime      time.Time `json:"castTime"`
	TargetID      uint64    `json:"targetID"`
	Location      Location  `json:"location"`
	CastType      int       `json:"castType"`
	EffectRange   int       `json:"effectRange"`
	XAxisModifier int       `json:"xAxisModifier"`
	Omen          string    `json:"omen"`
}

type ClassJob struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

type CraftingInfo struct {
	Recipe              RecipeInfo `json:"recipe"`
	LastCraftActionID   int        `json:"lastCraftActionID"`
	LastCraftActionName string     `json:"lastCraftActionName"`
	StepNum             int        `json:"stepNum"`
	Progress            int        `json:"progress"`
	ProgressDelta       int        `json:"progressDelta"`
	Quality             int        `json:"quality"`
	QualityDelta        int        `json:"qualityDelta"`
	HqChance            int        `json:"hqChance"`
	Durability          int        `json:"durability"`
	DurabilityDelta     int        `json:"durabilityDelta"`
	CurrentCondition    int        `json:"currentCondition"`
	PreviousCondition   int        `json:"previousCondition"`
	ReuseProc           bool       `json:"reuseProc"`
}

type Enmity struct {
	TargetHateRanking []HateRanking `json:"targetHateRanking"`
	NearbyEnemyHate   []HateEntry   `json:"nearbyEnemyHate"`
}

type Entity struct {
	ID               uint64       `json:"id"`
	Index            int          `json:"index"`
	Name             string       `json:"name"`
	TargetID         uint64       `json:"targetID"`
	OwnerID          uint64       `json:"ownerID"`
	Level            int          `json:"level"`
	ClassJob         ClassJob     `json:"classJob"`
	IsNPC            bool         `json:"isNPC"`
	IsEnemy          bool         `json:"isEnemy"`
	IsPet            bool         `json:"isPet"`
	BNPCInfo         *NPCInfo     `json:"bNPCInfo"`
	Resources        Resources    `json:"resources"`
	Location         Location     `json:"location"`
	LastAction       *Action      `json:"lastAction"`
	Statuses         []*Status    `json:"statuses"`
	LockonMarker     int          `json:"lockonMarker"`
	CastingInfo      *CastingInfo `json:"castingInfo"`
	RawSpawnJSONData string       `json:"rawSpawnJSONData"`
}

type EntityEvent struct {
	StreamID int             `json:"streamID"`
	EntityID uint64          `json:"entityID"`
	Type     EntityEventType `json:"type"`
}

type EntityEventType interface {
	IsEntityEventType()
}

type HateEntry struct {
	EnemyID     uint64 `json:"enemyID"`
	HatePercent int    `json:"hatePercent"`
}

type HateRanking struct {
	ActorID uint64 `json:"actorID"`
	Hate    int    `json:"hate"`
}

type Location struct {
	X           float64   `json:"x"`
	Y           float64   `json:"y"`
	Z           float64   `json:"z"`
	Orientation float64   `json:"orientation"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type MapInfo struct {
	Key           int    `json:"key"`
	ID            string `json:"id"`
	SizeFactor    int    `json:"SizeFactor"`
	OffsetX       int    `json:"OffsetX"`
	OffsetY       int    `json:"OffsetY"`
	PlaceName     string `json:"PlaceName"`
	PlaceNameSub  string `json:"PlaceNameSub"`
	TerritoryType string `json:"TerritoryType"`
}

type NPCInfo struct {
	NameID  int      `json:"nameID"`
	BaseID  int      `json:"baseID"`
	ModelID int      `json:"modelID"`
	Name    *string  `json:"name"`
	Size    *float64 `json:"size"`
	Error   int      `json:"error"`
}

type Place struct {
	MapID       int       `json:"mapID"`
	TerritoryID int       `json:"territoryID"`
	Maps        []MapInfo `json:"maps"`
}

type RecipeInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	RecipeLevel int    `json:"recipeLevel"`
	Element     int    `json:"element"`
	CanHQ       bool   `json:"canHQ"`
	Difficulty  int    `json:"difficulty"`
	Quality     int    `json:"quality"`
	Durability  int    `json:"durability"`
}

type RemoveEntity struct {
	ID uint64 `json:"id"`
}

func (RemoveEntity) IsEntityEventType() {}

type RemoveStatus struct {
	Index int `json:"index"`
}

func (RemoveStatus) IsEntityEventType() {}

type RemoveStream struct {
	ID int `json:"id"`
}

func (RemoveStream) IsStreamEventType() {}

type Resources struct {
	Hp       int       `json:"hp"`
	Mp       int       `json:"mp"`
	Tp       int       `json:"tp"`
	MaxHP    int       `json:"maxHP"`
	MaxMP    int       `json:"maxMP"`
	LastTick time.Time `json:"lastTick"`
}

type SetEntities struct {
	Entities []Entity `json:"entities"`
}

func (SetEntities) IsEntityEventType() {}

type Status struct {
	ID          int       `json:"id"`
	Param       int       `json:"param"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartedTime time.Time `json:"startedTime"`
	Duration    time.Time `json:"duration"`
	ActorID     uint64    `json:"actorID"`
	LastTick    time.Time `json:"lastTick"`
}

type StreamEvent struct {
	StreamID int             `json:"streamID"`
	Type     StreamEventType `json:"type"`
}

type StreamEventType interface {
	IsStreamEventType()
}

type StreamRequest struct {
	StreamID int    `json:"streamID"`
	Data     string `json:"data"`
}

type UpdateCastingInfo struct {
	CastingInfo *CastingInfo `json:"castingInfo"`
}

func (UpdateCastingInfo) IsEntityEventType() {}

type UpdateClass struct {
	ClassJob ClassJob `json:"classJob"`
	Level    int      `json:"level"`
}

func (UpdateClass) IsEntityEventType() {}

type UpdateCraftingInfo struct {
	CraftingInfo *CraftingInfo `json:"craftingInfo"`
}

func (UpdateCraftingInfo) IsStreamEventType() {}

type UpdateEnmity struct {
	Enmity Enmity `json:"enmity"`
}

func (UpdateEnmity) IsStreamEventType() {}

type UpdateIDs struct {
	ServerID    int    `json:"serverID"`
	CharacterID uint64 `json:"characterID"`
	InstanceNum int    `json:"instanceNum"`
}

func (UpdateIDs) IsStreamEventType() {}

type UpdateLastAction struct {
	Action Action `json:"action"`
}

func (UpdateLastAction) IsEntityEventType() {}

type UpdateLocation struct {
	Location Location `json:"location"`
}

func (UpdateLocation) IsEntityEventType() {}

type UpdateLockonMarker struct {
	LockonMarker int `json:"lockonMarker"`
}

func (UpdateLockonMarker) IsEntityEventType() {}

type UpdateMap struct {
	Place Place `json:"place"`
}

func (UpdateMap) IsStreamEventType() {}

type UpdateResources struct {
	Resources Resources `json:"resources"`
}

func (UpdateResources) IsEntityEventType() {}

type UpdateTarget struct {
	TargetID uint64 `json:"targetID"`
}

func (UpdateTarget) IsEntityEventType() {}

type UpsertStatus struct {
	Index  int    `json:"index"`
	Status Status `json:"status"`
}

func (UpsertStatus) IsEntityEventType() {}
