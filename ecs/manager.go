package ecs

import (
	"github.com/jasonfantl/rogue/gui"
)

type Manager struct {
	entityManager EntityManager
	eventManager  EventManager
	running       bool
	user          User
}

func New() Manager {
	newManager := Manager{}
	newManager.entityManager = newEntityManager()
	newManager.eventManager = NewEventHandler()
	newManager.running = false
	newManager.user = User{}

	return newManager
}

func (m *Manager) Start() {
	m.running = true
	// make sure to display
	startingEvents := []Event{
		{WAKEUP_HANDLERS, WakeupHandlers{}, m.user.Controlling},
	}
	m.eventManager.sendEvents(m, startingEvents)
}

func (m *Manager) Running() bool {
	return m.running
}

func (m *Manager) Run() {
	key, pressed := gui.GetKeyPress()

	if pressed {
		buttonEvent := []Event{{KEY_PRESSED, KeyPressed{key}, m.user.Controlling}}
		m.eventManager.sendEvents(m, buttonEvent)
	}
}

func (m *Manager) AddEventHandler(eventHandler EventHandler) {
	m.eventManager.addEventHandler(eventHandler)
}

func (m *Manager) AddEntity(components map[ComponentID]interface{}) Entity {
	return m.entityManager.addEntity(components)
}

func (m *Manager) getComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	return m.entityManager.getComponent(entity, componentID)
}

// can we reove this? promotes inefficient code
func (m *Manager) getEntitiesWithComponent(componentID ComponentID) map[Entity]bool {
	return m.entityManager.getEntitiesWithComponent(componentID)
}

func (m *Manager) setComponent(entity Entity, componentID ComponentID, data interface{}) {
	// manager special case
	if componentID == USER {
		m.user = data.(User)
	}

	m.entityManager.setComponent(entity, componentID, data)
}

func (m *Manager) removeComponent(entity Entity, componentID ComponentID) {
	m.entityManager.removeComponent(entity, componentID)
}

func (m *Manager) removeEntity(entity Entity) {
	m.entityManager.removeEntity(entity)
}

func (m *Manager) getEntitiesAtPosition(pos Position) (entities map[Entity]bool) {
	return m.entityManager.getEntitiesAtPosition(pos)
}

func (m *Manager) SetUser(user User) {
	m.user = user
}
