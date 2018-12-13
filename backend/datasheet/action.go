package datasheet

import (
	"encoding/json"
	"fmt"
	"io"
)

// ActionStore stores all of the Action data
type ActionStore map[uint32]Action

// Action stores the data for a game Action
type Action struct {
	ID             uint32 `json:"key"`
	Name           string `json:"Name"`
	ActionCategory string `json:"ActionCategory"`
	ClassJob       string `json:"ClassJob"`
	Range          int8   `json:"Range"`
	TargetArea     bool   `json:"TargetArea"`
	CastType       byte   `json:"CastType"`
	EffectRange    byte   `json:"EffectRange"`
	XAxisModifier  byte   `json:"XAxisModifier"`
	CostType       byte   `json:"Cost{Type}"`
	Cost           uint16 `json:"Cost"`
	ComboAction    string `json:"Action{Combo}"`
	Cast           uint16 `json:"Cast\u003c100ms\u003e"`
	Recast         uint16 `json:"Recast\u003c100ms\u003e"`
	GainedStatus   string `json:"Status{GainSelf}"`
	Omen           string `json:"Omen"`
}

// PopulateActions will populate the ActionStore with Action data provided a
// path to the data sheet for Actions
func (a *ActionStore) PopulateActions(dataBytes io.Reader) error {
	*a = make(map[uint32]Action)
	var rows []Action
	d := json.NewDecoder(dataBytes)
	err := d.Decode(&rows)
	if err != nil {
		return fmt.Errorf("PopulateActions: %s", err)
	}
	for _, action := range rows {
		(*a)[action.ID] = action
	}
	return nil
}
