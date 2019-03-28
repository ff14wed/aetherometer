package datasheet

import (
	"fmt"
	"io"
)

// ActionStore stores all of the Action data. Note that querying actions
// directly on the map will result in an empty Omen field. You should use
// GetAction in order to have an action with a correctly populated field.
type ActionStore struct {
	Actions      map[uint32]Action
	Omens        map[uint16]Omen
	CraftActions map[uint32]CraftAction
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

type CraftAction struct {
	Key  uint32 `datasheet:"key"`
	Name string `datasheet:"Name"`
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

// PopulateCraftActions will populate the ActionStore with CraftAction data
// provided a path to the data sheet for CraftActions
func (a *ActionStore) PopulateCraftActions(dataReader io.Reader) error {
	a.CraftActions = make(map[uint32]CraftAction)

	var rows []CraftAction
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateCraftActions: %s", err)
	}
	for _, craftAction := range rows {
		a.CraftActions[craftAction.Key] = craftAction
	}
	return nil
}

// GetAction returns the Action associated with the action key. It will
// also populate the Omen field on the action.
// If the action is not found in the standard Action store, it will attempt
// to return an action from the CraftAction store.
func (a *ActionStore) GetAction(key uint32) Action {
	if action, found := a.Actions[key]; found {
		if omen, found := a.Omens[action.OmenID]; found {
			action.Omen = omen.Name
		}
		return action
	}
	if craftAction, found := a.CraftActions[key]; found {
		return Action{
			Key:  craftAction.Key,
			Name: craftAction.Name,
		}
	}
	return Action{}
}
