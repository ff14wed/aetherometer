package datasheet

// Collection encapsulates a collection of datasheets
type Collection struct {
	MapData    MapStore
	BNPCData   BNPCStore
	ActionData ActionStore
	StatusData StatusStore
}
