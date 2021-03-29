package ecs

type Entity uint64

type EntityManager struct {
	entityTable   chunkedEntityTable
	entityCounter Entity
}

func newEntityManager() EntityManager {
	newManager := EntityManager{}
	newManager.entityTable = chunkedEntityTable{}
	newManager.entityCounter = 0

	return newManager
}

func (em *EntityManager) AddEntity(componenets map[ComponentID]interface{}) Entity {
	entity := em.entityCounter
	em.entityCounter++
	em.entityTable.addEntity(entity, componenets)
	return entity
}

func (em *EntityManager) GetComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	return em.entityTable.getComponent(entity, componentID)
}

// can we reove this? promotes inefficient code
func (m *EntityManager) GetEntities(componentID ComponentID) map[Entity]bool {
	return m.entityTable.getEntities(componentID)
}

func (m *EntityManager) GetComponents(entity Entity) map[ComponentID]interface{} {
	return m.entityTable.getComponents(entity)
}

func (em *EntityManager) SetComponent(entity Entity, componentID ComponentID, data interface{}) {
	em.entityTable.setComponent(entity, componentID, data)
}

func (em *EntityManager) RemoveComponent(entity Entity, componentID ComponentID) {
	em.entityTable.removeComponent(entity, componentID)
}

func (em *EntityManager) RemoveEntity(entity Entity) {
	em.entityTable.removeEntity(entity)
}

func (em *EntityManager) GetEntitiesAtPosition(p Position) map[Entity]bool {
	return em.entityTable.getEntitiesAtPosition(p)
}

func (em *EntityManager) UnloadAllChunks() {
	em.entityTable.unloadAllChunks()
}
