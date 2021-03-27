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

func (m *Manager) AddEntity(components []Component) Entity {
	return m.entityManager.addEntity(m, components)
}

func (m *Manager) getComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	return m.entityManager.getComponent(entity, componentID)
}

// can we reove this? promotes inefficient code
func (m *Manager) getComponents(componentID ComponentID) (map[Entity]interface{}, bool) {
	return m.entityManager.getComponents(componentID)
}

func (m *Manager) setComponent(entity Entity, componentID ComponentID, data interface{}) {
	// manager special case
	if componentID == USER {
		m.user = data.(User)
	}

	m.entityManager.setComponent(m, entity, componentID, data)
}

func (m *Manager) removeComponent(entity Entity, componentID ComponentID) {
	m.entityManager.removeComponent(m, entity, componentID)
}

func (m *Manager) removeEntity(entity Entity) {
	m.entityManager.removeEntity(m, entity)
}

func (m *Manager) getEntitiesFromPos(x, y int) (entities map[Entity]bool) {
	return m.entityManager.getEntitiesFromPos(x, y)
}
