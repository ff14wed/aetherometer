package datasheet

import (
	"fmt"
	"io"

	"github.com/ff14wed/aetherometer/core/models"
)

// WorldStore stores all of the World data.
type WorldStore map[uint32]World

// World stores the data for a game World
type World struct {
	Key  uint32 `datasheet:"key"`
	Name string `datasheet:"Name"`
}

// PopulateWorlds will populate the WorldStore with World data provided a
// path to the data sheet for Worlds.
func (w *WorldStore) PopulateWorlds(dataReader io.Reader) error {
	*w = make(map[uint32]World)

	var rows []World
	err := UnmarshalReader(dataReader, &rows)
	if err != nil {
		return fmt.Errorf("PopulateWorlds: %s", err)
	}
	for _, world := range rows {
		(*w)[world.Key] = world
	}
	return nil
}

func (w *WorldStore) Lookup(worldID int) models.World {
	if resolved, ok := (*w)[uint32(worldID)]; ok {
		return models.World{
			ID:   worldID,
			Name: resolved.Name,
		}
	}

	return models.World{
		ID:   worldID,
		Name: fmt.Sprintf("Unknown_%d", worldID),
	}
}
