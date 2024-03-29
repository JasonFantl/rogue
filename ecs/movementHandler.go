package ecs

type MoveHandler struct {
}

func (s *MoveHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {
	// somwhere in here the components are copied, so we cant edit them
	// that why we set the components at the end
	// there must be a way to maintain the pointer to the component

	if event.ID == TRY_MOVE {
		// unpack event data
		moveEvent := event.data.(TryMove)

		// get entitys current position and if its blocking
		positionData, hasPosition := m.getComponent(event.entity, POSITION)
		_, hasVolume := m.getComponent(event.entity, VOLUME)

		if hasPosition {
			positionComponent := positionData.(Position)

			// now check if new location is occupied
			newPos := Position{positionComponent.X + moveEvent.dx, positionComponent.Y + moveEvent.dy}
			canMove := true

			for otherEntity := range m.getEntitiesAtPosition(newPos) {
				// since we use getEntitiesFromPos, it must have the same position
				_, otherHasVolume := m.getComponent(otherEntity, VOLUME)

				if otherHasVolume && hasVolume {
					canMove = false
				}
			}

			if canMove {
				returnEvents = append(returnEvents, Event{MOVED, Moved{positionComponent, newPos}, event.entity})

				positionComponent = newPos

				m.setComponent(event.entity, POSITION, positionComponent)

			}
		}
	}
	return returnEvents
}
