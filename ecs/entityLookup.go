package ecs

type EntityLookup map[Entity]map[ComponentID]interface{}

func (el EntityLookup) addEntity(entity Entity, components map[ComponentID]interface{}) {
	el[entity] = components
}

func (el EntityLookup) getComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	components, ok := el[entity]
	if ok {
		data, ok := components[componentID]
		if ok {
			return data, true
		}
	}
	return nil, false
}

// can we reove this? promotes inefficient code
func (el EntityLookup) getComponents(entity Entity) map[ComponentID]interface{} {
	_, ok := el[entity]
	if ok {
		return el[entity]
	}
	return make(map[ComponentID]interface{})
}

func (el EntityLookup) setComponent(entity Entity, componentID ComponentID, data interface{}) {
	// check component map is initalized
	_, ok := el[entity]
	if !ok {
		el[entity] = make(map[ComponentID]interface{})
	}

	el[entity][componentID] = data
}

func (entityTable EntityLookup) removeComponent(entity Entity, componentID ComponentID) {
	delete(entityTable[entity], componentID)
}

func (el EntityLookup) removeEntity(entity Entity) {
	delete(el, entity)
}
