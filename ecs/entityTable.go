package ecs

type EntityTable map[ComponentID]map[Entity]interface{}

func (entityTable EntityTable) addEntity(m *Manager, entity Entity, componenets []Component) {
	for _, component := range componenets {
		m.setComponent(entity, component.ID, component.Data)
	}
}

func (entityTable EntityTable) getComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	components, ok := entityTable[componentID]
	if ok {
		data, ok := components[entity]
		if ok {
			return data, true
		}
	}
	return nil, false
}

// can we reove this? promotes inefficient code
func (entityTable EntityTable) getComponents(componentID ComponentID) (map[Entity]interface{}, bool) {
	_, ok := entityTable[componentID]
	if ok {
		return entityTable[componentID], true
	}
	return nil, false
}

func (entityTable EntityTable) setComponent(m *Manager, entity Entity, componentID ComponentID, data interface{}) {
	// check component map is initalized
	_, ok := entityTable[componentID]
	if !ok {
		entityTable[componentID] = make(map[Entity]interface{})
	}

	entityTable[componentID][entity] = data
}

func (entityTable EntityTable) removeComponent(m *Manager, entity Entity, componentID ComponentID) {
	delete(entityTable[componentID], entity)
}

func (entityTable EntityTable) removeEntity(m *Manager, entity Entity) {
	for componentID := range entityTable {
		m.removeComponent(entity, componentID)
	}
}
