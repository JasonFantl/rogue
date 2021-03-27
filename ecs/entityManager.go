package ecs

type Entity uint64

type EntityManager struct {
	entityTable    EntityTable
	positionLookup ChunkPositionLookup
	entityCounter  Entity
}

func newEntityManager() EntityManager {
	newManager := EntityManager{}
	newManager.entityTable = make(EntityTable) // can we make chunk based?
	newManager.positionLookup = make(ChunkPositionLookup)
	newManager.entityCounter = 0

	return newManager
}

func (em *EntityManager) addEntity(m *Manager, componenets []Component) Entity {
	entity := em.entityCounter
	em.entityCounter++

	em.entityTable.addEntity(m, entity, componenets)

	positionData, hasPosition := em.entityTable.getComponent(entity, POSITION)
	if hasPosition {
		positionComponent := positionData.(Position)
		em.positionLookup.add(map[Entity]bool{entity: true}, positionComponent.X, positionComponent.Y)
	}

	return entity
}

func (m *EntityManager) getComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	return m.entityTable.getComponent(entity, componentID)
}

// can we reove this? promotes inefficient code
func (m *EntityManager) getComponents(componentID ComponentID) (map[Entity]interface{}, bool) {
	return m.entityTable.getComponents(componentID)
}

func (em *EntityManager) setComponent(m *Manager, entity Entity, componentID ComponentID, data interface{}) {

	// for position lookup
	if componentID == POSITION {
		newPosition := data.(Position)
		oldPositionData, hasPosition := em.getComponent(entity, POSITION)
		if hasPosition {
			oldPosition := oldPositionData.(Position)
			em.positionLookup.move(entity, oldPosition.X, oldPosition.Y, newPosition.X, newPosition.Y)
		} else {
			em.positionLookup.add(map[Entity]bool{entity: true}, newPosition.X, newPosition.Y)
		}
	}

	em.entityTable.setComponent(m, entity, componentID, data)
}

func (em *EntityManager) removeComponent(m *Manager, entity Entity, componentID ComponentID) {
	if componentID == POSITION {
		positionData, hasPosition := em.entityTable.getComponent(entity, POSITION)
		if hasPosition {
			positionComponent := positionData.(Position)
			em.positionLookup.remove(entity, positionComponent.X, positionComponent.Y)
		}
	}

	em.entityTable.removeComponent(m, entity, componentID)
}

func (em *EntityManager) removeEntity(m *Manager, entity Entity) {

	positionData, hasPosition := em.entityTable.getComponent(entity, POSITION)
	if hasPosition {
		positionComponent := positionData.(Position)
		em.positionLookup.remove(entity, positionComponent.X, positionComponent.Y)
	}

	em.entityTable.removeEntity(m, entity)
}

func (em *EntityManager) getEntitiesFromPos(x, y int) (entities map[Entity]bool) {
	return em.positionLookup.get(x, y)
}
