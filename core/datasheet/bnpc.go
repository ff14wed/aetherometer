package datasheet

import (
	"fmt"
	"io"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	Key  uint32 `datasheet:"key"`
	Name string `datasheet:"Singular"`
}

// BNPCBase stores information about the scale of the monster
type BNPCBase struct {
	Key   uint32  `datasheet:"key"`
	Scale float32 `datasheet:"Scale"`
}

// ModelChara stores information about the models used for the character
type ModelChara struct {
	Key   uint32 `datasheet:"key"`
	Model uint32 `datasheet:"Model"`
}

// ModelSkeleton stores information about the base size of the model
type ModelSkeleton struct {
	Key    uint32  `datasheet:"key"`
	Radius float32 `datasheet:"Radius"`
}

// PopulateBNPCNames will populate the BNPCStore with BNPC data provided a
// path to the data sheet for BNPCNames
func (b *BNPCStore) PopulateBNPCNames(dataReader io.Reader) error {
	b.BNPCNames = make(map[uint32]BNPCName)
	var rows []BNPCName
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateBNPCNames: %s", err)
	}
	for _, bNPCName := range rows {
		b.BNPCNames[bNPCName.Key] = bNPCName
	}
	return nil
}

// PopulateBNPCBases will populate the BNPCStore with BNPC data provided a
// path to the data sheet for BNPCBases
func (b *BNPCStore) PopulateBNPCBases(dataReader io.Reader) error {
	b.BNPCBases = make(map[uint32]BNPCBase)
	var rows []BNPCBase
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateBNPCBases: %s", err)
	}
	for _, bNPCBase := range rows {
		b.BNPCBases[bNPCBase.Key] = bNPCBase
	}
	return nil
}

// PopulateModelCharas will populate the BNPCStore with BNPC data provided a
// path to the data sheet for ModelCharas
func (b *BNPCStore) PopulateModelCharas(dataReader io.Reader) error {
	b.ModelCharas = make(map[uint32]ModelChara)
	var rows []ModelChara
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateModelCharas: %s", err)
	}
	for _, modelChara := range rows {
		b.ModelCharas[modelChara.Key] = modelChara
	}
	return nil
}

// PopulateModelSkeletons will populate the BNPCStore with BNPC data provided a
// path to the data sheet for ModelSkeletons
func (b *BNPCStore) PopulateModelSkeletons(dataReader io.Reader) error {
	b.ModelSkeletons = make(map[uint32]ModelSkeleton)
	var rows []ModelSkeleton
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateModelSkeletons: %s", err)
	}
	for _, modelSkeleton := range rows {
		b.ModelSkeletons[modelSkeleton.Key] = modelSkeleton
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

	bNPCInfo.Name = cases.Title(language.English).String(bNPCName.Name)

	bNPCBase, ok := b.BNPCBases[bNPCBaseID]
	var scale float32 = 1.0
	if ok {
		scale = bNPCBase.Scale
	}
	var baseSize float32 = 0.5
	modelChara, ok1 := b.ModelCharas[modelCharaID]
	modelSkeleton, ok2 := b.ModelSkeletons[modelChara.Model]
	if ok1 && ok2 {
		baseSize = modelSkeleton.Radius
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
