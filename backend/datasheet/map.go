package datasheet

import (
	"encoding/json"
	"fmt"
	"io"
)

// MapInfo stores the information for a game Map
type MapInfo struct {
	Key           uint16 `json:"key"`
	ID            string `json:"Id"`
	SizeFactor    uint16 `json:"SizeFactor"`
	OffsetX       int16  `json:"Offset{X}"`
	OffsetY       int16  `json:"Offset{Y}"`
	PlaceName     string `json:"PlaceName"`
	PlaceNameSub  string `json:"PlaceName{Sub}"`
	TerritoryType string `json:"TerritoryType"`
}

// TerritoryInfo stores the mapping between TerritoryID and MapID
type TerritoryInfo struct {
	ID    uint16 `json:"key"`
	Name  string `json:"Name"`
	MapID string `json:"Map"`
}

// MapStore stores information about all maps and territories
type MapStore struct {
	Maps        map[string]MapInfo
	Territories map[uint16]TerritoryInfo

	mapsForTerritories map[string][]MapInfo
}

// PopulateMaps will populate the MapStore with map data provided a
// path to the data sheet for Maps
func (m *MapStore) PopulateMaps(dataBytes io.Reader) error {
	m.Maps = make(map[string]MapInfo)
	var rows []MapInfo
	d := json.NewDecoder(dataBytes)
	err := d.Decode(&rows)
	if err != nil {
		return fmt.Errorf("PopulateMaps: %s", err)
	}
	for _, mapInfo := range rows {
		m.Maps[mapInfo.ID] = mapInfo
	}
	return nil
}

// PopulateTerritories will populate the MapStore with territory data provided a
// path to the data sheet
func (m *MapStore) PopulateTerritories(dataBytes io.Reader) error {
	m.Territories = make(map[uint16]TerritoryInfo)
	var rows []TerritoryInfo
	d := json.NewDecoder(dataBytes)
	err := d.Decode(&rows)
	if err != nil {
		return fmt.Errorf("PopulateTerritories: %s", err)
	}
	for _, territoryInfo := range rows {
		m.Territories[territoryInfo.ID] = territoryInfo
	}
	return nil
}

// GetMaps returns all Maps associated with the territory ID.
// If no entry is found, it returns nil
func (m *MapStore) GetMaps(ID uint16) []MapInfo {
	t, ok := m.Territories[ID]
	if !ok {
		return nil
	}

	if m.mapsForTerritories == nil {
		// Initialize the cache
		m.mapsForTerritories = make(map[string][]MapInfo)
	}

	_, ok = m.mapsForTerritories[t.Name]
	if !ok {
		// Regenerate cache
		if defaultMap, ok := m.Maps[t.MapID]; ok {
			m.mapsForTerritories[t.Name] = []MapInfo{defaultMap}
		}
		for mid, mv := range m.Maps {
			if mid == t.MapID {
				continue
			}
			m.mapsForTerritories[mv.TerritoryType] =
				append(m.mapsForTerritories[mv.TerritoryType], mv)
		}
	}

	return m.mapsForTerritories[t.Name]
}
