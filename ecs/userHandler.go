package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type UserHandler struct {
}

func (h *UserHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// controll new entity when the old one dies

	if event.ID == DIED && event.entity == m.user.Controlling {
		brains, hasbrain := m.getComponents(BRAIN)
		if hasbrain {
			for brain := range brains {
				m.user.Controlling = brain
				break
			}
		}
	}

	if event.ID == KEY_PRESSED {
		keyPressedEvent := event.data.(KeyPressed)

		if m.user.Menu.active {
			returnEvents = append(returnEvents, h.handleMenu(m, keyPressedEvent.key)...)
		} else {
			returnEvents = append(returnEvents, h.handlePlaying(m, keyPressedEvent.key)...)
		}
	}

	return returnEvents
}

func (h *UserHandler) handlePlaying(m *Manager, key ebiten.Key) (returnEvents []Event) {
	timestep := false

	switch key {
	// remeber the screen has an inverted y, thats why we send these move values
	case m.user.DownKey:
		timestep = true
		returnEvents = append(returnEvents,
			Event{TRY_MOVE, TryMove{0, 1}, m.user.Controlling},
		)
	case m.user.UpKey:
		timestep = true
		returnEvents = append(returnEvents,
			Event{TRY_MOVE, TryMove{0, -1}, m.user.Controlling},
		)
	case m.user.LeftKey:
		timestep = true
		returnEvents = append(returnEvents,
			Event{TRY_MOVE, TryMove{-1, 0}, m.user.Controlling},
		)
	case m.user.RightKey:
		timestep = true
		returnEvents = append(returnEvents,
			Event{TRY_MOVE, TryMove{1, 0}, m.user.Controlling},
		)
	case m.user.PickupKey:
		timestep = true
		returnEvents = append(returnEvents,
			Event{PLAYER_TRY_PICK_UP, PlayerTryPickUp{}, m.user.Controlling},
		)
	case m.user.MenuKey:
		m.user.Menu.active = true
	case m.user.QuitKey:
		returnEvents = append(returnEvents,
			Event{QUIT, Quit{}, m.user.Controlling},
		)
	default:
		returnEvents = append(returnEvents,
			Event{DEBUG_EVENT, DebugEvent{"invalid key press"}, Entity(key)},
		)
	}

	if timestep {
		timeStepEvent := Event{TIMESTEP, TimeStep{}, m.user.Controlling}
		returnEvents = append([]Event{timeStepEvent}, returnEvents...)
	}

	return returnEvents
}

func (h *UserHandler) handleMenu(m *Manager, key ebiten.Key) (returnEvents []Event) {

	switch key {
	case m.user.DownKey:
		m.user.Menu.curserDown(m)
	case m.user.UpKey:
		m.user.Menu.curserUp(m)
	case m.user.PickupKey:
		returnEvents = append(returnEvents, m.user.Menu.curserSelect(m)...)
	case m.user.MenuKey:
		m.user.Menu.active = false
		m.user.Menu.curserReset()
	case m.user.QuitKey:
		returnEvents = append(returnEvents,
			Event{QUIT, Quit{}, m.user.Controlling},
		)
	default:
		returnEvents = append(returnEvents,
			Event{DEBUG_EVENT, DebugEvent{"invalid key press"}, Entity(key)},
		)
	}

	return returnEvents
}