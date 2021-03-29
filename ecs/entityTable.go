package ecs

type entityTable struct {
	EL EntityLookup    // entity to components
	CL ComponentLookup // componentID to entities
	PL PositionLookup  // position to entities
}

// updating functions

func (et *entityTable) addEntity(entity Entity, components map[ComponentID]interface{}) {
	et.checkInited()

	et.EL.addEntity(entity, components)
	et.CL.addEntity(entity, components)

	positionData, hasPosition := et.EL.getComponent(entity, POSITION)
	if hasPosition {
		positionComponent := positionData.(Position)
		et.PL.addEntity(entity, positionComponent)
	}
}

func (et *entityTable) removeEntity(entity Entity) {
	et.checkInited()
	positionData, hasPosition := et.EL.getComponent(entity, POSITION)
	if hasPosition {
		positionComponent := positionData.(Position)
		et.PL.removeEntity(entity, positionComponent)
	}

	et.EL.removeEntity(entity)
	et.CL.removeEntity(entity)
}

func (et *entityTable) setComponent(entity Entity, componentID ComponentID, data interface{}) {
	et.checkInited()

	if componentID == POSITION {
		positionData, hasPosition := et.EL.getComponent(entity, POSITION)
		newPosition := data.(Position)
		if hasPosition {
			oldPosition := positionData.(Position)
			et.PL.moveEntity(entity, oldPosition, newPosition)
		} else {
			et.PL.addEntity(entity, newPosition)
		}
	}

	et.EL.setComponent(entity, componentID, data)
	et.CL.setComponent(entity, componentID, data)
}

func (et *entityTable) removeComponent(entity Entity, componentID ComponentID) {
	et.checkInited()
	if componentID == POSITION {
		positionData, hasPosition := et.EL.getComponent(entity, POSITION)
		if hasPosition {
			positionComponent := positionData.(Position)
			et.PL.removeEntity(entity, positionComponent)
		}
	}

	et.EL.removeComponent(entity, componentID)
	et.CL.removeComponent(entity, componentID)
}

// getter functions

func (et *entityTable) getComponents(entity Entity) map[ComponentID]interface{} {
	et.checkInited()
	return et.EL.getComponents(entity)
}

func (et *entityTable) getComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	// could use either entityLookup or componentLookup, doesnt really matter
	et.checkInited()
	return et.EL.getComponent(entity, componentID)
}

func (et *entityTable) getEntitiesWithComponent(componentID ComponentID) map[Entity]bool {
	et.checkInited()
	return et.CL.getEntities(componentID)
}

func (et *entityTable) getEntitiesAtPosition(p Position) map[Entity]bool {
	et.checkInited()
	return et.PL.getEntities(p)
}

func (et *entityTable) checkInited() {
	if et.CL == nil {
		et.CL = make(ComponentLookup)
	}
	if et.EL == nil {
		et.EL = make(EntityLookup)
	}
	if et.PL == nil {
		et.PL = make(PositionLookup)
	}
}
