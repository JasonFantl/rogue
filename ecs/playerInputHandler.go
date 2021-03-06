package ecs

type PlayerInputHandler struct {
}

func (h *PlayerInputHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == KEY_PRESSED {
		keyPressedEvent := event.data.(KeyPressed)

		controllerData, ok := m.getComponent(event.entity, PLAYER_CONTROLLER)

		timestep := false

		if ok {
			controllerComponent := controllerData.(PlayerController)

			switch keyPressedEvent.key {
			// remeber the screen has an inverted y, thats why we send these move values
			case controllerComponent.Down:
				timestep = true
				returnEvents = append(returnEvents,
					Event{TRY_MOVE, TryMove{0, 1}, controllerComponent.Controlling},
				)
			case controllerComponent.Up:
				timestep = true
				returnEvents = append(returnEvents,
					Event{TRY_MOVE, TryMove{0, -1}, controllerComponent.Controlling},
				)
			case controllerComponent.Left:
				timestep = true
				returnEvents = append(returnEvents,
					Event{TRY_MOVE, TryMove{-1, 0}, controllerComponent.Controlling},
				)
			case controllerComponent.Right:
				timestep = true
				returnEvents = append(returnEvents,
					Event{TRY_MOVE, TryMove{1, 0}, controllerComponent.Controlling},
				)
			case controllerComponent.Pickup:
				timestep = true
				returnEvents = append(returnEvents,
					Event{PLAYER_TRY_PICK_UP, PlayerTryPickUp{}, controllerComponent.Controlling},
				)
			case controllerComponent.Consume:
				timestep = true
				returnEvents = append(returnEvents,
					Event{TRY_CONSUME, TryConsume{}, controllerComponent.Controlling},
				)
			case controllerComponent.Quit:
				returnEvents = append(returnEvents,
					Event{QUIT, Quit{}, controllerComponent.Controlling},
				)
			default:
				returnEvents = append(returnEvents,
					Event{DEBUG_EVENT, DebugEvent{"invalid key press"}, Entity(keyPressedEvent.key)},
				)
			}
		}

		if timestep {
			timeStepEvent := Event{TIMESTEP, TimeStep{}, event.entity}
			returnEvents = append([]Event{timeStepEvent}, returnEvents...)
		}
	}

	return returnEvents
}
