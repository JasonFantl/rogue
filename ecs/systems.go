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

	timestep := false
	triggeredEvents := make([]Event, 0)

	playerControllerComponents, ok := m.getComponents(PLAYER_CONTROLLER)
	if ok {
		for entity, controllerData := range playerControllerComponents {

			controllerComponent := controllerData.(PlayerController)
			key, pressed := gui.GetKeyPress()

			if pressed {
				switch key {
				// remeber the screen has an inverted y, thats why we send these move values
				case controllerComponent.Down:
					timestep = true
					triggeredEvents = append(triggeredEvents,
						Event{TRY_MOVE, TryMove{0, 1}, entity},
					)
				case controllerComponent.Up:
					timestep = true
					triggeredEvents = append(triggeredEvents,
						Event{TRY_MOVE, TryMove{0, -1}, entity},
					)
				case controllerComponent.Left:
					timestep = true
					triggeredEvents = append(triggeredEvents,
						Event{TRY_MOVE, TryMove{-1, 0}, entity},
					)
				case controllerComponent.Right:
					timestep = true
					triggeredEvents = append(triggeredEvents,
						Event{TRY_MOVE, TryMove{1, 0}, entity},
					)
				case controllerComponent.Pickup:
					timestep = true
					triggeredEvents = append(triggeredEvents,
						Event{TRY_PICK_UP, TryPickUp{}, entity},
					)
				case controllerComponent.Quit:
					triggeredEvents = append(triggeredEvents,
						Event{QUIT, Quit{}, entity},
					)
				default:
					triggeredEvents = append(triggeredEvents,
						Event{ERROR_EVENT, ErrorEvent{"invalid key press"}, Entity(key)},
					)
				}
			}
		}

		if timestep {
			timeStepEvent := Event{TIMESTEP, TimeStep{}, 0}
			// shold we time step then move, or the other way around?
			triggeredEvents = append([]Event{timeStepEvent}, triggeredEvents...)
		}

		if len(triggeredEvents) > 0 {
			m.sendEvents(triggeredEvents)
		}
	}
}
