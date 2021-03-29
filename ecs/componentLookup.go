package ecs

type ComponentLookup map[ComponentID]map[Entity]bool

func (cl ComponentLookup) addEntity(entity Entity, componenets map[ComponentID]interface{}) {
	for componentID, data := range componenets {
		cl.setComponent(entity, componentID, data)
	}
}

// can we remove this? promotes inefficient code
func (cl ComponentLookup) getEntities(componentID ComponentID) map[Entity]bool {
	_, ok := cl[componentID]
	if ok {
		return cl[componentID]
	}
	return map[Entity]bool{}
}

func (cl ComponentLookup) setComponent(entity Entity, componentID ComponentID, data interface{}) {
	// check component map is initalized
	_, ok := cl[componentID]
	if !ok {
		cl[componentID] = make(map[Entity]bool)
	}

	cl[componentID][entity] = true
}

func (cl ComponentLookup) removeComponent(entity Entity, componentID ComponentID) {
	delete(cl[componentID], entity)
}

func (cl ComponentLookup) removeEntity(entity Entity) {
	for componentID := range cl {
		cl.removeComponent(entity, componentID)
	}
}
