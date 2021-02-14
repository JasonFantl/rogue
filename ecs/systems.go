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

	validPress := false
	timestep := false
	triggeredEvents := make([]Event, 0)

	for entity, controllerData := range m.lookupTable[PLAYER_CONTROLLER] {

		controllerComponent := controllerData.(PlayerController)
		key, pressed := gui.GetKeyPress()

		if pressed {
			validPress = true
			switch key {
			case controllerComponent.Down:
				timestep = true
				triggeredEvents = append(triggeredEvents,
					Event{TRY_MOVE_EVENT, EventTryMove{0, -1}, entity},
				)
			case controllerComponent.Up:
				timestep = true
				triggeredEvents = append(triggeredEvents,
					Event{TRY_MOVE_EVENT, EventTryMove{0, 1}, entity},
				)
			case controllerComponent.Left:
				timestep = true
				triggeredEvents = append(triggeredEvents,
					Event{TRY_MOVE_EVENT, EventTryMove{-1, 0}, entity},
				)
			case controllerComponent.Right:
				timestep = true
				triggeredEvents = append(triggeredEvents,
					Event{TRY_MOVE_EVENT, EventTryMove{1, 0}, entity},
				)
			case controllerComponent.Pickup:
				// should we time step when picking something up?
				triggeredEvents = append(triggeredEvents,
					Event{TRY_PICK_UP_EVENT, EventTryPickUp{}, entity},
				)
			case controllerComponent.Quit:
				triggeredEvents = append(triggeredEvents,
					Event{QUIT_EVENT, EventQuit{}, entity},
				)
			default:
				triggeredEvents = append(triggeredEvents,
					Event{ERROR_EVENT, EventError{"invalid key press"}, Entity(key)},
				)
			}
		}
	}

	if timestep {
		timeStepEvent := Event{TIMESTEP, EventTimeStep{}, 0}
		// shold we time step then move, or the other way around?
		triggeredEvents = append([]Event{timeStepEvent}, triggeredEvents...)
	}

	if validPress {
		m.sendEvents(triggeredEvents)
	}
}
