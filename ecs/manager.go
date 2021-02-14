package ecs

import (
	"github.com/jasonfantl/rogue/gui"
)

// TODO:
// quick lookup function based on location
// maybe lookup based on multiple component tags

type Entity uint64

type EventHandler interface {
	handleEvent(*Manager, Event) (returnEvents []Event)
}

type Manager struct {
	lookupTable   map[ComponentID]map[Entity]interface{}
	eventHandlers []EventHandler
	systems       []System
	entityCounter Entity
	hasQuit       bool
}

func New() Manager {
	newManager := Manager{}
	newManager.lookupTable = make(map[ComponentID]map[Entity]interface{})
	newManager.eventHandlers = make([]EventHandler, 0)
	newManager.systems = make([]System, 0)
	newManager.entityCounter = 0
	newManager.hasQuit = false

	return newManager
}

func (m *Manager) Start() {
	// make sure to display
	startingEvents := []Event{
		Event{DISPLAY_EVENT, EventDisplayTrigger{}, 0},
	}
	m.sendEvents(startingEvents)
}

func (m *Manager) HasQuit() bool {
	return m.hasQuit
}

func (m *Manager) AddSystem(system System) {
	m.systems = append(m.systems, system)
}

func (m *Manager) AddEventHandler(eventHandler EventHandler) {
	m.eventHandlers = append(m.eventHandlers, eventHandler)
}

func (m *Manager) AddEntity(componenets []Component) Entity {
	entity := m.entityCounter
	m.entityCounter++

	for _, component := range componenets {
		m.AddComponenet(entity, component)
	}

	return entity
}

func (m *Manager) AddComponenet(entity Entity, component Component) bool {
	// check component map is initalized
	componentList, ok := m.lookupTable[component.ID]

	if !ok {
		m.lookupTable[component.ID] = make(map[Entity]interface{})
	}

	// check entity doesnt already have component
	if _, ok := componentList[entity]; !ok {
		m.lookupTable[component.ID][entity] = component.Data
		return true
	}
	return false
}

func (m *Manager) Run() {
	for _, system := range m.systems {
		system.run(m)
	}
}

func (m *Manager) getComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	entities, ok := m.lookupTable[componentID]
	if ok {
		data, ok := entities[entity]
		if ok {
			return data, true
		}
	}
	return nil, false
}

func (m *Manager) getEntitiesFromPos(x, y int) (entities []Entity) {
	for entity, positionData := range m.lookupTable[POSITION] {
		positionComponent := positionData.(Position)
		if positionComponent.X == x && positionComponent.Y == y {
			entities = append(entities, entity)
		}
	}
	return entities
}

func (m *Manager) sendEvent(event Event) {
	events := []Event{event}
	m.sendEvents(events)
}

func (m *Manager) sendEvents(events []Event) {

	// we re-display every time an independent event is fired, must clear the screen first
	gui.Clear()

	// queue style event handling
	for len(events) > 0 {
		sendingEvent := events[0] // pop
		events = events[1:]       // dequeue

		// special manager case
		if sendingEvent.ID == QUIT_EVENT {
			m.hasQuit = true
		}

		for _, eventHandler := range m.eventHandlers {

			respondingEvents := eventHandler.handleEvent(m, sendingEvent)

			events = append(events, respondingEvents...)
		}

		if len(events) > 100 {
			break
		}
	}

	// make sure to update the screen after events have all triggered
	for _, eventHandler := range m.eventHandlers {
		eventHandler.handleEvent(m, Event{DISPLAY_EVENT, EventDisplayTrigger{}, 0})
	}
}
