// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"time"
)

type EntityEventType interface {
	IsEntityEventType()
}

type StreamEventType interface {
	IsStreamEventType()
}

type Action struct {
	TargetID          uint64         `json:"targetID"`
	Name              string         `json:"name"`
	GlobalCounter     int            `json:"globalCounter"`
	AnimationLockTime float64        `json:"animationLockTime"`
	HiddenAnimation   int            `json:"hiddenAnimation"`
	Location          *Location      `json:"location" validate:"nil=false"`
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
	Entity *Entity `json:"entity" validate:"nil=false"`
}

func (AddEntity) IsEntityEventType() {}

type AddStream struct {
	Stream *Stream `json:"stream" validate:"nil=false"`
}

func (AddStream) IsStreamEventType() {}

type CastingInfo struct {
	ActionID      int       `json:"actionID"`
	ActionName    string    `json:"actionName"`
	StartTime     time.Time `json:"startTime"`
	CastTime      time.Time `json:"castTime"`
	TargetID      uint64    `json:"targetID"`
	Location      *Location `json:"location" validate:"nil=false"`
	CastType      int       `json:"castType"`
	EffectRange   int       `json:"effectRange"`
	XAxisModifier int       `json:"xAxisModifier"`
	Omen          string    `json:"omen"`
}

type ChatEvent struct {
	ChannelID    uint64 `json:"channelID"`
	ChannelWorld *World `json:"channelWorld" validate:"nil=false"`
	ChannelType  string `json:"channelType"`
	ContentID    uint64 `json:"contentID"`
	EntityID     uint64 `json:"entityID"`
	World        *World `json:"world" validate:"nil=false"`
	Name         string `json:"name"`
	Message      string `json:"message"`
}

func (ChatEvent) IsStreamEventType() {}

type ClassJob struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

type CraftingInfo struct {
	Recipe              *RecipeInfo `json:"recipe" validate:"nil=false"`
	LastCraftActionID   int         `json:"lastCraftActionID"`
	LastCraftActionName string      `json:"lastCraftActionName"`
	StepNum             int         `json:"stepNum"`
	Progress            int         `json:"progress"`
	ProgressDelta       int         `json:"progressDelta"`
	Quality             int         `json:"quality"`
	QualityDelta        int         `json:"qualityDelta"`
	HqChance            int         `json:"hqChance"`
	Durability          int         `json:"durability"`
	DurabilityDelta     int         `json:"durabilityDelta"`
	CurrentCondition    int         `json:"currentCondition"`
	PreviousCondition   int         `json:"previousCondition"`
	Completed           bool        `json:"completed"`
	Failed              bool        `json:"failed"`
	ReuseProc           bool        `json:"reuseProc"`
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
	ClassJob         *ClassJob    `json:"classJob" validate:"nil=false"`
	IsNpc            bool         `json:"isNPC"`
	IsEnemy          bool         `json:"isEnemy"`
	IsPet            bool         `json:"isPet"`
	BNPCInfo         *NPCInfo     `json:"bNPCInfo"`
	Resources        *Resources   `json:"resources" validate:"nil=false"`
	Location         *Location    `json:"location" validate:"nil=false"`
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
	ItemID      int    `json:"itemID"`
	Element     int    `json:"element"`
	CanHq       bool   `json:"canHQ"`
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
	MaxHp    int       `json:"maxHP"`
	MaxMp    int       `json:"maxMP"`
	LastTick time.Time `json:"lastTick"`
}

type SetEntities struct {
	Entities []Entity `json:"entities"`
}

func (SetEntities) IsEntityEventType() {}

type Stats struct {
	Strength            int `json:"strength"`
	Dexterity           int `json:"dexterity"`
	Vitality            int `json:"vitality"`
	Intelligence        int `json:"intelligence"`
	Mind                int `json:"mind"`
	Piety               int `json:"piety"`
	Hp                  int `json:"hp"`
	Mp                  int `json:"mp"`
	Tp                  int `json:"tp"`
	Gp                  int `json:"gp"`
	Cp                  int `json:"cp"`
	Delay               int `json:"delay"`
	Tenacity            int `json:"tenacity"`
	AttackPower         int `json:"attackPower"`
	Defense             int `json:"defense"`
	DirectHitRate       int `json:"directHitRate"`
	Evasion             int `json:"evasion"`
	MagicDefense        int `json:"magicDefense"`
	CriticalHit         int `json:"criticalHit"`
	AttackMagicPotency  int `json:"attackMagicPotency"`
	HealingMagicPotency int `json:"healingMagicPotency"`
	ElementalBonus      int `json:"elementalBonus"`
	Determination       int `json:"determination"`
	SkillSpeed          int `json:"skillSpeed"`
	SpellSpeed          int `json:"spellSpeed"`
	Haste               int `json:"haste"`
	Craftsmanship       int `json:"craftsmanship"`
	Control             int `json:"control"`
	Gathering           int `json:"gathering"`
	Perception          int `json:"perception"`
}

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

type StreamRequest struct {
	StreamID int    `json:"streamID"`
	Data     string `json:"data"`
}

type UpdateCastingInfo struct {
	CastingInfo *CastingInfo `json:"castingInfo"`
}

func (UpdateCastingInfo) IsEntityEventType() {}

type UpdateClass struct {
	ClassJob *ClassJob `json:"classJob" validate:"nil=false"`
	Level    int       `json:"level"`
}

func (UpdateClass) IsEntityEventType() {}

type UpdateCraftingInfo struct {
	CraftingInfo *CraftingInfo `json:"craftingInfo"`
}

func (UpdateCraftingInfo) IsStreamEventType() {}

type UpdateEnmity struct {
	Enmity *Enmity `json:"enmity" validate:"nil=false"`
}

func (UpdateEnmity) IsStreamEventType() {}

type UpdateIDs struct {
	ServerID     int    `json:"serverID"`
	InstanceNum  int    `json:"instanceNum"`
	CharacterID  uint64 `json:"characterID"`
	HomeWorld    *World `json:"homeWorld" validate:"nil=false"`
	CurrentWorld *World `json:"currentWorld" validate:"nil=false"`
}

func (UpdateIDs) IsStreamEventType() {}

type UpdateLastAction struct {
	Action *Action `json:"action" validate:"nil=false"`
}

func (UpdateLastAction) IsEntityEventType() {}

type UpdateLocation struct {
	Location *Location `json:"location" validate:"nil=false"`
}

func (UpdateLocation) IsEntityEventType() {}

type UpdateLockonMarker struct {
	LockonMarker int `json:"lockonMarker"`
}

func (UpdateLockonMarker) IsEntityEventType() {}

type UpdateMap struct {
	Place *Place `json:"place" validate:"nil=false"`
}

func (UpdateMap) IsStreamEventType() {}

type UpdateResources struct {
	Resources *Resources `json:"resources" validate:"nil=false"`
}

func (UpdateResources) IsEntityEventType() {}

type UpdateStats struct {
	Stats *Stats `json:"stats" validate:"nil=false"`
}

func (UpdateStats) IsStreamEventType() {}

type UpdateTarget struct {
	TargetID uint64 `json:"targetID"`
}

func (UpdateTarget) IsEntityEventType() {}

type UpsertStatus struct {
	Index  int     `json:"index"`
	Status *Status `json:"status" validate:"nil=false"`
}

func (UpsertStatus) IsEntityEventType() {}

type World struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
