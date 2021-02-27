package ecs

type VisionHandler struct{}

func (s *VisionHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// not sure how/when to update/handle this event, kinda tricky
	if event.ID == DISPLAY {

		visionData, hasVision := m.getComponent(event.entity, VISION)
		awarnessData, hasAwarness := m.getComponent(event.entity, ENTITY_AWARENESS)
		positionData, hasPosition := m.getComponent(event.entity, POSITION)

		if hasVision && hasAwarness && hasPosition {
			visionComponent := visionData.(Vision)
			awarnessComponent := awarnessData.(EntityAwarness)
			positionComponent := positionData.(Position)

			// clear old awarness first
			awarnessComponent.AwareOf = make([]Entity, 0)

			// later imlement FOV, for now just display everything withen the raduis
			for dx := -visionComponent.Reach; dx <= visionComponent.Reach; dx++ {
				for dy := -visionComponent.Reach; dy <= visionComponent.Reach; dy++ {
					x := positionComponent.X + dx
					y := positionComponent.Y + dy

					awarnessComponent.AwareOf = append(awarnessComponent.AwareOf, m.getEntitiesFromPos(x, y)...)
				}
			}

			m.setComponent(event.entity, Component{ENTITY_AWARENESS, awarnessComponent})
		}
	}

	return returnEvents
}
