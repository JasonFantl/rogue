package ecs

import (
	"time"

	"github.com/jasonfantl/rogue/gui"
)

type EventHandler interface {
	handleEvent(*Manager, Event) []Event
}

type Manager struct {
	entityTable    map[ComponentID]map[Entity]interface{}
	positionLookup PositionLookup
	eventHandlers  []EventHandler
	entityCounter  Entity
	running        bool
	user           User
}

func New() Manager {
	newManager := Manager{}
	newManager.entityTable = make(map[ComponentID]map[Entity]interface{})
	newManager.positionLookup = PositionLookup{}
	newManager.eventHandlers = make([]EventHandler, 0)
	newManager.entityCounter = 0
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
	m.sendEvents(startingEvents)
}

func (m *Manager) Running() bool {
	return m.running
}

func (m *Manager) Run() {
	key, pressed := gui.GetKeyPress()

	if pressed {
		buttonEvent := []Event{{KEY_PRESSED, KeyPressed{key}, m.user.Controlling}}
		m.sendEvents(buttonEvent)
	}
}

func (m *Manager) AddEventHandler(eventHandler EventHandler) {
	m.eventHandlers = append(m.eventHandlers, eventHandler)
}

func (m *Manager) AddEntity(componenets []Component) Entity {
	entity := m.entityCounter
	m.entityCounter++

	for _, component := range componenets {
		m.setComponent(entity, component.ID, component.Data)
	}

	return entity
}

func (m *Manager) getComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	components, ok := m.entityTable[componentID]
	if ok {
		data, ok := components[entity]
		if ok {
			return data, true
		}
	}
	return nil, false
}

// can we reove this? promotes inefficient code
func (m *Manager) getComponents(componentID ComponentID) (map[Entity]interface{}, bool) {
	_, ok := m.entityTable[componentID]
	if ok {
		return m.entityTable[componentID], true
	}
	return nil, false
}

func (m *Manager) setComponent(entity Entity, componentID ComponentID, data interface{}) {
	// check component map is initalized
	_, ok := m.entityTable[componentID]
	if !ok {
		m.entityTable[componentID] = make(map[Entity]interface{})
	}

	// for position lookup
	if componentID == POSITION {
		newPosition := data.(Position)
		oldPositionData, hasPosition := m.getComponent(entity, POSITION)
		if hasPosition {
			oldPosition := oldPositionData.(Position)
			m.positionLookup.move(entity, oldPosition.X, oldPosition.Y, newPosition.X, newPosition.Y)
		} else {
			m.positionLookup.add(map[Entity]bool{entity: true}, newPosition.X, newPosition.Y)
		}
	}

	m.entityTable[componentID][entity] = data

	// manager special case
	if componentID == USER {
		m.user = data.(User)
	}
}

func (m *Manager) removeComponent(entity Entity, componentID ComponentID) {
	components, ok := m.entityTable[componentID]
	if ok {
		// for position lookup
		if componentID == POSITION {
			positionData, hasPosition := components[entity]
			if hasPosition {
				positionComponent := positionData.(Position)
				m.positionLookup.remove(entity, positionComponent.X, positionComponent.Y)
			}
		}

		delete(m.entityTable[componentID], entity)
	}
}

func (m *Manager) removeEntity(entity Entity) {
	for componentID := range m.entityTable {
		m.removeComponent(entity, componentID)
	}
}

func (m *Manager) getEntitiesFromPos(x, y int) (entities map[Entity]bool) {
	return m.positionLookup.get(x, y)
}

func (m *Manager) sendEvents(events []Event) {

	timer := time.Now()
	sentDisplay := false

	// queue style event handling
	for len(events) > 0 {
		sendingEvent := events[0] // pop
		events = events[1:]       // dequeue

		for _, eventHandler := range m.eventHandlers {
			respondingEvents := eventHandler.handleEvent(m, sendingEvent)
			events = append(events, respondingEvents...)
		}

		// special manager case
		if sendingEvent.ID == QUIT {
			m.running = false
		}

		// display stuff
		if !sentDisplay && len(events) == 0 {
			sentDisplay = true
			events = append(events, Event{DISPLAY, Display{}, m.user.Controlling})
		}
	}

	gui.Debug(time.Since(timer).String())
}
