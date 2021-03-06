package ecs

import (
	"time"

	"github.com/jasonfantl/rogue/gui"
)

type EventHandler interface {
	handleEvent(*Manager, Event) []Event
}

type Manager struct {
	lookupTable    map[ComponentID]map[Entity]interface{}
	positionLookup PositionLookup
	eventHandlers  []EventHandler
	entityCounter  Entity
	running        bool
	user           Entity
}

func New() Manager {
	newManager := Manager{}
	newManager.lookupTable = make(map[ComponentID]map[Entity]interface{})
	newManager.positionLookup = make(map[int]map[int][]Entity, 0)
	newManager.eventHandlers = make([]EventHandler, 0)
	newManager.entityCounter = 0
	newManager.running = false
	newManager.user = 0

	return newManager
}

func (m *Manager) Start() {
	m.running = true
	// make sure to display
	startingEvents := []Event{
		{DEBUG_EVENT, DebugEvent{"waking up handlers"}, m.user},
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

		// for position lookup
		if component.ID == POSITION {
			positionComponent := component.Data.(Position)
			m.positionLookup.add([]Entity{entity}, positionComponent.X, positionComponent.Y)
		}

		// manager special case
		if component.ID == PLAYER_CONTROLLER {
			m.user = entity
		}
		return true
	}
	return false
}

func (m *Manager) Run() {
	key, pressed := gui.GetKeyPress()

	if pressed {
		buttonEvent := []Event{{KEY_PRESSED, KeyPressed{key}, m.user}}
		m.sendEvents(buttonEvent)
	}
}

func (m *Manager) getComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	components, ok := m.lookupTable[componentID]
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
	_, ok := m.lookupTable[componentID]
	if ok {
		return m.lookupTable[componentID], true
	}
	return nil, false
}

func (m *Manager) setComponent(entity Entity, componentID ComponentID, data interface{}) {
	// for position lookup
	if componentID == POSITION {
		oldPositionData, hasPosition := m.getComponent(entity, POSITION)
		if hasPosition {
			oldPosition := oldPositionData.(Position)
			newPosition := data.(Position)
			m.positionLookup.move(entity, oldPosition.X, oldPosition.Y, newPosition.X, newPosition.Y)
		}
	}

	_, ok := m.lookupTable[componentID]
	if ok {
		m.lookupTable[componentID][entity] = data
	}
}

func (m *Manager) removeComponent(entity Entity, componentID ComponentID) {
	components, ok := m.lookupTable[componentID]
	if ok {
		// for position lookup
		if componentID == POSITION {
			positionData, hasPosition := components[entity]
			if hasPosition {
				positionComponent := positionData.(Position)
				m.positionLookup.remove(entity, positionComponent.X, positionComponent.Y)
			}
		}

		delete(m.lookupTable[componentID], entity)
	}
}

func (m *Manager) removeEntity(entity Entity) {
	for componentID := range m.lookupTable {
		m.removeComponent(entity, componentID)
	}
}

func (m *Manager) getEntitiesFromPos(x, y int) (entities []Entity) {
	col, ok := m.positionLookup[x]
	if ok {
		entities, ok := col[y]
		if ok {
			return entities
		}
	}
	return []Entity{}
}

func (m *Manager) sendEvents(events []Event) {

	/////// manager stuff //////
	timer := time.Now()
	blockFinished := Event{DEBUG_EVENT, DebugEvent{"-----------------------------"}, 999999999}
	events = append(events, blockFinished)
	sentDisplay := false
	userData, _ := m.getComponent(m.user, PLAYER_CONTROLLER)
	userComponent := userData.(PlayerController)
	//////////////////////////////

	// we re-display every time an independent event is fired, must clear the screen first
	gui.Clear()

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
		if sendingEvent.ID == DIED && sendingEvent.entity == userComponent.Controlling {
			monsters, isMonster := m.getComponents(MONSTER_CONTROLLER)
			if isMonster {
				for monster, _ := range monsters {
					userComponent.Controlling = monster
					break
				}
			}
			m.setComponent(m.user, PLAYER_CONTROLLER, userComponent)
		}

		// display stuff
		if !sentDisplay && len(events) == 0 {
			sentDisplay = true
			events = append(events, Event{DISPLAY, Display{}, userComponent.Controlling})
		}
	}
	gui.DrawCorner(time.Since(timer).String())

	gui.Show()

}
