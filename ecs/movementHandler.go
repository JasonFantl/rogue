package ecs

type MoveHandler struct {
}

func (s *MoveHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {
	// somwhere in here the components are copied, so we cant edit them
	// that why we set the components at the end
	// there must be a way to maintain the pointer to the component

	if event.ID == TRY_MOVE_EVENT {
		// unpack event data
		moveEvent := event.data.(EventTryMove)

		// get entitys current position and if its blocking
		positionData, hasPosition := m.getComponent(event.entity, POSITION)
		_, hasBlockable := m.getComponent(event.entity, BLOCKABLE)

		if hasPosition {
			positionComponent := positionData.(Position)

			// now check if new location is occupied
			newX := positionComponent.X + moveEvent.dx
			newY := positionComponent.Y + moveEvent.dy
			canMove := true

			for _, otherEntity := range m.getEntitiesFromPos(newX, newY) {
				// since we use getEntitiesFromPos, it must have the same position
				_, otherHasBlockable := m.getComponent(otherEntity, BLOCKABLE)

				if otherHasBlockable || hasBlockable {
					canMove = false
				}
			}

			if canMove {
				positionComponent.X = newX
				positionComponent.Y = newY

				m.setComponent(event.entity, Component{POSITION, positionComponent})

				returnEvents = append(returnEvents, Event{MOVE_EVENT, EventMove{newX, newY}, event.entity})
			}
		}
	}
	return returnEvents
}
