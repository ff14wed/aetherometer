package datasheet

import (
	"fmt"
	"io"
)

// ActionStore stores all of the Action data. Note that querying actions
// directly on the map will result in an empty Omen field. You should use
// GetAction in order to have an action with a correctly populated field.
type ActionStore struct {
	Actions map[uint32]Action
	Omens   map[uint16]Omen
}

// Action stores the data for a game Action
type Action struct {
	Key           uint32 `datasheet:"key"`
	Name          string `datasheet:"Name"`
	Range         int8   `datasheet:"Range"`
	TargetArea    bool   `datasheet:"TargetArea"`
	CastType      byte   `datasheet:"CastType"`
	EffectRange   byte   `datasheet:"EffectRange"`
	XAxisModifier byte   `datasheet:"XAxisModifier"`
	OmenID        uint16 `datasheet:"Omen"`
	Omen          string
}

// Omen stores the data for a game action Omen
type Omen struct {
	Key  uint16 `datasheet:"key"`
	Name string `datasheet:"FileName"`
}

// PopulateActions will populate the ActionStore with Action data provided a
// path to the data sheet for Actions.
func (a *ActionStore) PopulateActions(dataReader io.Reader) error {
	a.Actions = make(map[uint32]Action)

	var rows []Action
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateActions: %s", err)
	}
	for _, action := range rows {
		a.Actions[action.Key] = action
	}
	return nil
}

// PopulateOmens will populate the ActionStore with Omen data provided a
// path to the data sheet for Omens
func (a *ActionStore) PopulateOmens(dataReader io.Reader) error {
	a.Omens = make(map[uint16]Omen)

	var rows []Omen
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateOmens: %s", err)
	}
	for _, omen := range rows {
		a.Omens[omen.Key] = omen
	}
	return nil
}

// GetAction returns the Action associated with the action key. It will
// also populate the Omen field on the action.
func (a *ActionStore) GetAction(key uint32) Action {
	action, found := a.Actions[key]
	if !found {
		return Action{}
	}
	if omen, found := a.Omens[action.OmenID]; found {
		action.Omen = omen.Name
	}
	return action
}
