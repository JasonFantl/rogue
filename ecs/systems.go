package ecs

import (
	"github.com/jasonfantl/rogue/gui"
)

type System interface {
	run(*Manager)
}

type InputSystem struct {
}

func (s *InputSystem) run(m *Manager) {

	for entity, controllerData := range m.lookupTable[PLAYER_CONTROLLER] {

		controllerComponent := controllerData.(PlayerController)
		key, pressed := gui.GetKeyPress()

		if pressed {
			switch key {
			case controllerComponent.Down:
				m.sendEvent(Event{TRY_MOVE_EVENT, EventTryMove{0, -1}, entity})
			case controllerComponent.Up:
				m.sendEvent(Event{TRY_MOVE_EVENT, EventTryMove{0, 1}, entity})
			case controllerComponent.Left:
				m.sendEvent(Event{TRY_MOVE_EVENT, EventTryMove{-1, 0}, entity})
			case controllerComponent.Right:
				m.sendEvent(Event{TRY_MOVE_EVENT, EventTryMove{1, 0}, entity})
			case controllerComponent.Pickup:
				m.sendEvent(Event{TRY_PICK_UP_EVENT, EventTryPickUp{}, entity})
			case controllerComponent.Quit:
				m.sendEvent(Event{QUIT_EVENT, EventQuit{}, entity})
			default:
				m.sendEvent(Event{ERROR_EVENT, key, Entity(key)})
			}
		}
	}
}
