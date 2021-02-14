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
		positionData, positionOk := m.getComponent(event.entity, POSITION)
		_, blockingOk := m.getComponent(event.entity, BLOCKABLE)

		if positionOk {
			positionComponent := positionData.(Position)

			// now check if new location is occupied
			newX := positionComponent.X + moveEvent.dx
			newY := positionComponent.Y + moveEvent.dy
			canMove := true
			// quick implementation, replace later
			for otherEntity, otherData := range m.lookupTable[POSITION] {
				otherPositionComponent := otherData.(Position)
				_, otherBlockableOk := m.lookupTable[BLOCKABLE][otherEntity]

				if otherPositionComponent.X == newX && otherPositionComponent.Y == newY {
					if otherBlockableOk || blockingOk {
						canMove = false
					}
				}
			}

			if canMove {
				positionComponent.X = newX
				positionComponent.Y = newY

				m.lookupTable[POSITION][event.entity] = positionComponent

				returnEvents = append(returnEvents, Event{MOVE_EVENT, EventMove{newX, newY}, event.entity})
			}
		}
	}
	return returnEvents
}
