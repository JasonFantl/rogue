package ecs

import (
	"github.com/jasonfantl/rogue/gui"
)

type EventHandler interface {
	handleEvent(*Manager, Event) (returnEvents []Event)
}

type Manager struct {
	lookupTable   map[ComponentID]map[Entity]interface{}
	eventHandlers []EventHandler
	entityCounter Entity
	running       bool
}

func New() Manager {
	newManager := Manager{}
	newManager.lookupTable = make(map[ComponentID]map[Entity]interface{})
	newManager.eventHandlers = make([]EventHandler, 0)
	newManager.entityCounter = 0
	newManager.running = false

	return newManager
}

func (m *Manager) Start() {
	m.running = true
	// make sure to display
	startingEvents := []Event{
		{DISPLAY, Display{}, 0},
	}
	m.sendEvents(startingEvents)
}

func (m *Manager) Running() bool {
	return m.running
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
	key, pressed := gui.GetKeyPress()

	if pressed {
		// send event as player so we know where to look for key mappings
		buttonEvent := []Event{{KEY_PRESSED, KeyPressed{key}, 0}}
		m.sendEvents(buttonEvent)
	}
}

func (m *Manager) getComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	components, ok := m.getComponents(componentID)
	if ok {
		data, ok := components[entity]
		if ok {
			return data, true
		}
	}
	return nil, false
}

func (m *Manager) getComponents(componentID ComponentID) (map[Entity]interface{}, bool) {
	_, ok := m.lookupTable[componentID]
	if ok {
		return m.lookupTable[componentID], true
	}
	return nil, false
}

func (m *Manager) setComponent(entity Entity, component Component) {
	_, ok := m.lookupTable[component.ID]
	if ok {
		m.lookupTable[component.ID][entity] = component.Data
	}
}

func (m *Manager) removeComponent(entity Entity, componentID ComponentID) {
	_, ok := m.lookupTable[componentID]
	if ok {
		delete(m.lookupTable[componentID], entity)
	}
}

func (m *Manager) removeEntity(entity Entity) {
	for componentID := range m.lookupTable {
		// special case so we dont remove the ability to quit the game
		if componentID != PLAYER_CONTROLLER {
			m.removeComponent(entity, componentID)
		}
	}
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

	blockFinished := Event{ERROR_EVENT, ErrorEvent{"-----------------------------"}, 999999999}
	sentDisplay := false

	events = append(events, blockFinished)
	// queue style event handling
	for len(events) > 0 {
		sendingEvent := events[0] // pop
		events = events[1:]       // dequeue

		// debug info
		if sendingEvent.entity == 999999999 && len(events) != 0 {
			events = append(events, blockFinished)
		}

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
			events = append(events, Event{DISPLAY, Display{}, 0})
		}
	}
	gui.Show()
}
