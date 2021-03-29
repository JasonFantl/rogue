package ecs

type Entity uint64

type EntityManager struct {
	entityTable    entityTable
	positionLookup PositionLookup
	entityCounter  Entity
}

func newEntityManager() EntityManager {
	newManager := EntityManager{}
	newManager.entityTable = newEntityTable()
	newManager.entityCounter = 0

	return newManager
}

func (em *EntityManager) addEntity(componenets map[ComponentID]interface{}) Entity {
	entity := em.entityCounter
	em.entityCounter++

	em.entityTable.addEntity(entity, componenets)

	return entity
}

func (m *EntityManager) getComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	return m.entityTable.getComponent(entity, componentID)
}

// can we reove this? promotes inefficient code
func (m *EntityManager) getEntitiesWithComponent(componentID ComponentID) map[Entity]bool {
	return m.entityTable.getEntitiesWithComponent(componentID)
}

func (em *EntityManager) setComponent(entity Entity, componentID ComponentID, data interface{}) {
	em.entityTable.setComponent(entity, componentID, data)
}

func (em *EntityManager) removeComponent(entity Entity, componentID ComponentID) {
	em.entityTable.removeComponent(entity, componentID)
}

func (em *EntityManager) removeEntity(entity Entity) {
	em.entityTable.removeEntity(entity)
}

func (em *EntityManager) getEntitiesAtPosition(pos Position) (entities map[Entity]bool) {
	return em.entityTable.getEntitiesAtPosition(pos)
}
