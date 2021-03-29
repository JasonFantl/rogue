package ecs

import (
	"github.com/jasonfantl/rogue/gui"
)

type UserHandler struct {
}

func (h *UserHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// controll new entity when the old one dies
	if event.ID == DIED && event.entity == m.user.Controlling {
		brains := m.getEntitiesWithComponent(BRAIN)
		for brain := range brains {
			m.user.Controlling = brain
			break
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

func (h *UserHandler) handlePlaying(m *Manager, key gui.Key) (returnEvents []Event) {
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
	case m.user.ActionKey:
		takenAction := false

		positionData, hasPosition := m.getComponent(m.user.Controlling, POSITION)
		if hasPosition {
			positionComponent := positionData.(Position)

			// if we are standing on anything, pick it up
			belowYou := m.getEntitiesAtPosition(positionComponent)
			for item := range belowYou {
				if isStashableTreasure(m, item) {
					returnEvents = append(returnEvents,
						Event{TRY_PICK_UP, TryPickUp{item}, m.user.Controlling},
					)
					takenAction = true
					timestep = true
					break
				}
			}

			// then try trading
			if !takenAction {
				for i := 0; i < 4; i++ {
					dx := (i / 2) * ((i%2)*2 - 1)
					dy := ((3 - i) / 2) * (((3-i)%2)*2 - 1)
					deltaPos := Position{positionComponent.X + dx, positionComponent.Y + dy}
					aroundYou := m.getEntitiesAtPosition(deltaPos)
					for entity := range aroundYou {
						_, hasInventory := m.getComponent(entity, INVENTORY)
						if hasInventory {
							m.user.Menu.open(m)
							m.user.Menu.state = SHOWING_TRADE
							m.user.Menu.rememberedEntity = entity
							break
						}
					}
				}
			}
		}

	case m.user.MenuKey:
		m.user.Menu.open(m)
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

func (h *UserHandler) handleMenu(m *Manager, key gui.Key) (returnEvents []Event) {

	switch key {
	case m.user.DownKey:
		m.user.Menu.moveCurser(m, 0, 1)
	case m.user.UpKey:
		m.user.Menu.moveCurser(m, 0, -1)
	case m.user.LeftKey:
		m.user.Menu.moveCurser(m, -1, 0)
	case m.user.RightKey:
		m.user.Menu.moveCurser(m, 1, 0)
	case m.user.ActionKey:
		returnEvents = append(returnEvents, m.user.Menu.selectAtCurser(m)...)
	case m.user.MenuKey:
		m.user.Menu.close(m)
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
