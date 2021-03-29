package ecs

type PositionLookup map[Position]map[Entity]bool

func (pl PositionLookup) checkInited(p Position) {
	_, ok := pl[p]
	if !ok {
		pl[p] = make(map[Entity]bool)
	}
}

func (pl PositionLookup) addEntity(entity Entity, p Position) {
	pl.checkInited(p)
	alreadyHave := pl[p][entity]

	if !alreadyHave {
		pl[p][entity] = true
	}
}

func (pl PositionLookup) addEntities(entities map[Entity]bool, p Position) {
	for entity := range entities {
		pl.addEntity(entity, p)
	}
}

func (pl PositionLookup) removeEntity(entity Entity, p Position) {
	pl.checkInited(p)

	delete(pl[p], entity)
}

func (pl PositionLookup) getEntities(p Position) map[Entity]bool {
	_, ok := pl[p]
	if ok {
		return pl[p]
	}
	return make(map[Entity]bool)
}

func (pl PositionLookup) moveEntity(entity Entity, p, np Position) {
	pl.removeEntity(entity, p)
	pl.addEntity(entity, np)
}
