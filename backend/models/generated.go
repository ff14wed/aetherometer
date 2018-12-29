// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"time"
)

type Action struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	TargetID uint64    `json:"targetID"`
	UseTime  time.Time `json:"useTime"`
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

type CraftingInfo struct {
	LastCraftAction   int `json:"lastCraftAction"`
	StepNum           int `json:"stepNum"`
	TotalProgress     int `json:"totalProgress"`
	ProgressDelta     int `json:"progressDelta"`
	TotalQuality      int `json:"totalQuality"`
	QualityDelta      int `json:"qualityDelta"`
	HqChance          int `json:"hqChance"`
	Durability        int `json:"durability"`
	DurabilityDelta   int `json:"durabilityDelta"`
	CurrentCondition  int `json:"currentCondition"`
	PreviousCondition int `json:"previousCondition"`
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
	Class            int          `json:"class"`
	IsNPC            bool         `json:"isNPC"`
	IsEnemy          bool         `json:"isEnemy"`
	IsPet            bool         `json:"isPet"`
	BNPCInfo         *NPCInfo     `json:"bNPCInfo"`
	Resources        Resources    `json:"resources"`
	Location         Location     `json:"location"`
	LastAction       *Action      `json:"lastAction"`
	Statuses         []*Status    `json:"statuses"`
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

type Status struct {
	ID          int       `json:"id"`
	Extra       int       `json:"extra"`
	Name        string    `json:"name"`
	StartedTime time.Time `json:"startedTime"`
	Duration    time.Time `json:"duration"`
	ActorID     uint64    `json:"actorID"`
	LastTick    time.Time `json:"lastTick"`
	BaseDamage  int       `json:"baseDamage"`
	CritRate    int       `json:"critRate"`
}

type StreamEvent struct {
	StreamID int             `json:"streamID"`
	Type     StreamEventType `json:"type"`
}

type StreamEventType interface {
	IsStreamEventType()
}

type UpdateCastingInfo struct {
	CastingInfo CastingInfo `json:"castingInfo"`
}

func (UpdateCastingInfo) IsEntityEventType() {}

type UpdateClass struct {
	Class int `json:"class"`
}

func (UpdateClass) IsEntityEventType() {}

type UpdateCraftingInfo struct {
	CraftingInfo CraftingInfo `json:"craftingInfo"`
}

func (UpdateCraftingInfo) IsStreamEventType() {}

type UpdateEnmity struct {
	Enmity Enmity `json:"enmity"`
}

func (UpdateEnmity) IsStreamEventType() {}

type UpdateLastAction struct {
	Action Action `json:"action"`
}

func (UpdateLastAction) IsEntityEventType() {}

type UpdateLocation struct {
	Location Location `json:"location"`
}

func (UpdateLocation) IsEntityEventType() {}

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
