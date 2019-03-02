package datasheet

import (
	"fmt"
	"io"

	"github.com/ff14wed/sibyl/backend/models"
)

// MapInfo stores the information for a game Map
type MapInfo struct {
	Key           uint16 `datasheet:"key"`
	ID            string `datasheet:"Id"`
	SizeFactor    uint16 `datasheet:"SizeFactor"`
	OffsetX       int16  `datasheet:"Offset{X}"`
	OffsetY       int16  `datasheet:"Offset{Y}"`
	PlaceName     uint16 `datasheet:"PlaceName"`
	PlaceNameSub  uint16 `datasheet:"PlaceName{Sub}"`
	TerritoryType uint16 `datasheet:"TerritoryType"`
}

type PlaceName struct {
	Key  uint16 `datasheet:"key"`
	Name string `datasheet:"Name"`
}

// TerritoryInfo stores information about the territory
type TerritoryInfo struct {
	Key  uint16 `datasheet:"key"`
	Name string `datasheet:"Name"`
	Map  uint16 `datasheet:"Map"`
}

// MapStore stores information about all maps and territories
type MapStore struct {
	Maps        map[uint16]MapInfo
	Territories map[uint16]TerritoryInfo
	PlaceNames  map[uint16]PlaceName

	mapsForTerritories map[uint16][]uint16
}

// PopulateMaps will populate the MapStore with map data provided a
// path to the data sheet for Maps
func (m *MapStore) PopulateMaps(dataReader io.Reader) error {
	m.Maps = make(map[uint16]MapInfo)

	var rows []MapInfo
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateMaps: %s", err)
	}
	for _, mapInfo := range rows {
		m.Maps[mapInfo.Key] = mapInfo
	}
	return nil
}

// PopulatePlaceNames will populate the MapStore with place name data provided
// a path to the data sheet for PlaceNames
func (m *MapStore) PopulatePlaceNames(dataReader io.Reader) error {
	m.PlaceNames = make(map[uint16]PlaceName)

	var rows []PlaceName
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateMaps: %s", err)
	}
	for _, placeName := range rows {
		m.PlaceNames[placeName.Key] = placeName
	}
	return nil
}

// PopulateTerritories will populate the MapStore with territory data provided a
// path to the data sheet
func (m *MapStore) PopulateTerritories(dataReader io.Reader) error {
	m.Territories = make(map[uint16]TerritoryInfo)
	var rows []TerritoryInfo
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateTerritories: %s", err)
	}
	for _, territoryInfo := range rows {
		m.Territories[territoryInfo.Key] = territoryInfo
	}
	return nil
}

// GetMaps returns all Maps associated with the territory ID.
// If no entry is found, it returns nil
func (m *MapStore) GetMaps(territoryID uint16) []models.MapInfo {
	t, found := m.Territories[territoryID]
	if !found {
		return nil
	}

	if m.mapsForTerritories == nil {
		m.mapsForTerritories = make(map[uint16][]uint16)

		for _, mv := range m.Maps {
			m.mapsForTerritories[mv.TerritoryType] =
				append(m.mapsForTerritories[mv.TerritoryType], mv.Key)
		}
	}

	mapKeys := m.mapsForTerritories[t.Key]
	if len(mapKeys) == 0 {
		foundMap := m.GetMap(t.Map)
		if foundMap.Key != 0 {
			return []models.MapInfo{foundMap}
		}
		return []models.MapInfo{}
	}

	var mapInfos []models.MapInfo
	for _, mapKey := range mapKeys {
		foundMap := m.GetMap(mapKey)
		if foundMap.Key != 0 {
			mapInfos = append(mapInfos, foundMap)
		}
	}
	return mapInfos
}

// GetMap returns the models.MapInfo associated with the Map key.
// It returns an empty models.MapInfo if no entry is found.
func (m *MapStore) GetMap(key uint16) models.MapInfo {
	mapInfo, found := m.Maps[key]
	if !found {
		return models.MapInfo{}
	}

	var (
		placeName     string
		placeNameSub  string
		territoryType string
	)

	if n, found := m.PlaceNames[mapInfo.PlaceName]; found {
		placeName = n.Name
	}

	if n, found := m.PlaceNames[mapInfo.PlaceNameSub]; found {
		placeNameSub = n.Name
	}

	if n, found := m.Territories[mapInfo.TerritoryType]; found {
		territoryType = n.Name
	}

	return models.MapInfo{
		Key:           int(mapInfo.Key),
		ID:            mapInfo.ID,
		SizeFactor:    int(mapInfo.SizeFactor),
		OffsetX:       int(mapInfo.OffsetX),
		OffsetY:       int(mapInfo.OffsetY),
		PlaceName:     placeName,
		PlaceNameSub:  placeNameSub,
		TerritoryType: territoryType,
	}
}
