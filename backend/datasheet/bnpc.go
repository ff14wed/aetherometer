package datasheet

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// BNPCInfo stores information about a game monster
type BNPCInfo struct {
	Name  string
	Size  float32
	Error byte
}

// BNPCStore stores information about all monsters
type BNPCStore struct {
	BNPCNames      map[uint32]BNPCName
	BNPCBases      map[uint32]BNPCBase
	ModelCharas    map[uint32]ModelChara
	ModelSkeletons map[uint32]ModelSkeleton
}

// BNPCName stores information about the name of the NPC
type BNPCName struct {
	ID   uint32 `json:"key"`
	Name string `json:"Singular"`
}

// BNPCBase stores information about the scale of the monster
type BNPCBase struct {
	ID    uint32  `json:"key"`
	Scale float32 `json:"Scale"`
}

// ModelChara stores information about the models used for the character
type ModelChara struct {
	ID    uint32 `json:"key"`
	Model uint32 `json:"Model"`
}

// ModelSkeleton stores information about the base size of the model
type ModelSkeleton struct {
	ID   uint32  `json:"key"`
	Size float32 `json:"0"`
}

// PopulateBNPCNames will populate the BNPCStore with BNPC data provided a
// path to the data sheet for BNPCNames
func (b *BNPCStore) PopulateBNPCNames(dataBytes io.Reader) error {
	b.BNPCNames = make(map[uint32]BNPCName)
	var rows []BNPCName
	d := json.NewDecoder(dataBytes)
	err := d.Decode(&rows)
	if err != nil {
		return fmt.Errorf("PopulateBNPCNames: %s", err)
	}
	for _, bNPCName := range rows {
		b.BNPCNames[bNPCName.ID] = bNPCName
	}
	return nil
}

// PopulateBNPCBases will populate the BNPCStore with BNPC data provided a
// path to the data sheet for BNPCBases
func (b *BNPCStore) PopulateBNPCBases(dataBytes io.Reader) error {
	b.BNPCBases = make(map[uint32]BNPCBase)
	var rows []BNPCBase
	d := json.NewDecoder(dataBytes)
	err := d.Decode(&rows)
	if err != nil {
		return fmt.Errorf("PopulateBNPCBases: %s", err)
	}
	for _, bNPCBase := range rows {
		b.BNPCBases[bNPCBase.ID] = bNPCBase
	}
	return nil
}

// PopulateModelCharas will populate the BNPCStore with BNPC data provided a
// path to the data sheet for ModelCharas
func (b *BNPCStore) PopulateModelCharas(dataBytes io.Reader) error {
	b.ModelCharas = make(map[uint32]ModelChara)
	var rows []ModelChara
	d := json.NewDecoder(dataBytes)
	err := d.Decode(&rows)
	if err != nil {
		return fmt.Errorf("PopulateModelCharas: %s", err)
	}
	for _, modelChara := range rows {
		b.ModelCharas[modelChara.ID] = modelChara
	}
	return nil
}

// PopulateModelSkeletons will populate the BNPCStore with BNPC data provided a
// path to the data sheet for ModelSkeletons
func (b *BNPCStore) PopulateModelSkeletons(dataBytes io.Reader) error {
	b.ModelSkeletons = make(map[uint32]ModelSkeleton)
	var rows []ModelSkeleton
	d := json.NewDecoder(dataBytes)
	err := d.Decode(&rows)
	if err != nil {
		return fmt.Errorf("PopulateModelSkeletons: %s", err)
	}
	for _, modelSkeleton := range rows {
		b.ModelSkeletons[modelSkeleton.ID] = modelSkeleton
	}
	return nil
}

// GetBNPCInfo returns a BNPCInfo object matching the provided parameters
// If the entry is not found, it returns nil
func (b *BNPCStore) GetBNPCInfo(bNPCNameID, bNPCBaseID, modelCharaID uint32) *BNPCInfo {
	var bNPCInfo BNPCInfo
	bNPCName, ok := b.BNPCNames[bNPCNameID]
	if !ok {
		return nil
	}
	bNPCInfo.Name = strings.Title(bNPCName.Name)

	bNPCBase, ok := b.BNPCBases[bNPCBaseID]
	var scale float32 = 1.0
	if ok {
		scale = bNPCBase.Scale
	}
	var baseSize float32 = 0.5
	modelChara, ok1 := b.ModelCharas[modelCharaID]
	modelSkeleton, ok2 := b.ModelSkeletons[modelChara.Model]
	if ok1 && ok2 {
		baseSize = modelSkeleton.Size
	}

	bNPCInfo.Size = scale * baseSize
	switch {
	case !ok:
		bNPCInfo.Error = 1
	case !ok1:
		bNPCInfo.Error = 2
	case !ok2:
		bNPCInfo.Error = 3
	}

	return &bNPCInfo
}
