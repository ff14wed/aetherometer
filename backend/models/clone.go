package models

// Clone returns a deep copy of the Stream struct. Any changes made to this copy
// should not affect the original struct.
func (s Stream) Clone() Stream {
	s.Place = s.Place.Clone()
	s.Enmity = s.Enmity.Clone()

	if s.CraftingInfo != nil {
		craftInfoClone := *s.CraftingInfo
		s.CraftingInfo = &craftInfoClone
	}

	if len(s.EntitiesMap) > 0 {
		entitiesMap := make(map[uint64]*Entity)
		for id, ent := range s.EntitiesMap {
			if ent == nil {
				entitiesMap[id] = nil
				continue
			}
			entClone := ent.Clone()
			entitiesMap[id] = &entClone
		}
		s.EntitiesMap = entitiesMap
	}

	return s
}

// Clone returns a deep copy of the Place struct. Any changes made to this copy
// should not affect the original struct.
func (p Place) Clone() Place {
	if len(p.Maps) > 0 {
		maps := make([]MapInfo, len(p.Maps))
		copy(maps, p.Maps)
		p.Maps = maps
	}
	return p
}

// Clone returns a deep copy of the Enmity struct. Any changes made to this copy
// should not affect the original struct.
func (e Enmity) Clone() Enmity {
	if len(e.TargetHateRanking) > 0 {
		hateRankings := make([]HateRanking, len(e.TargetHateRanking))
		copy(hateRankings, e.TargetHateRanking)
		e.TargetHateRanking = hateRankings
	}

	if len(e.NearbyEnemyHate) > 0 {
		enemyHate := make([]HateEntry, len(e.NearbyEnemyHate))
		copy(enemyHate, e.NearbyEnemyHate)
		e.NearbyEnemyHate = enemyHate
	}
	return e
}

// Clone returns a deep copy of the Entity struct. Any changes made to this copy
// should not affect the original struct.
func (e Entity) Clone() Entity {
	if e.BNPCInfo != nil {
		bNPCInfoClone := e.BNPCInfo.Clone()
		e.BNPCInfo = &bNPCInfoClone
	}

	if e.LastAction != nil {
		lastActionClone := e.LastAction.Clone()
		e.LastAction = &lastActionClone
	}

	if e.CastingInfo != nil {
		castingInfoClone := *e.CastingInfo
		e.CastingInfo = &castingInfoClone
	}

	if len(e.Statuses) > 0 {
		statuses := make([]*Status, len(e.Statuses))
		for i, v := range e.Statuses {
			if v == nil {
				continue
			}
			vClone := *v
			statuses[i] = &vClone
		}
		e.Statuses = statuses

	}

	return e
}

// Clone returns a deep copy of the NPCInfo struct. Any changes made to this copy
// should not affect the original struct.
func (n NPCInfo) Clone() NPCInfo {
	if n.Name != nil {
		nameClone := *n.Name
		n.Name = &nameClone
	}

	if n.Size != nil {
		sizeClone := *n.Size
		n.Size = &sizeClone
	}

	return n
}

// Clone returns a deep copy of the Action struct. Any changes made to this copy
// should not affect the original struct.
func (a Action) Clone() Action {
	if len(a.Effects) > 0 {
		effects := make([]ActionEffect, len(a.Effects))
		copy(effects, a.Effects)
		a.Effects = effects
	}
	return a
}
