package ecs

type VisionHandler struct{}

func (s *VisionHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// not sure how/when to update/handle this event, kinda tricky
	if event.ID == DISPLAY {

		entities, _ := m.getComponents(VISION)

		for entity := range entities {
			visionData, hasVision := m.getComponent(entity, VISION)
			awarnessData, hasAwarness := m.getComponent(entity, ENTITY_AWARENESS)
			positionData, hasPosition := m.getComponent(entity, POSITION)

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

				m.setComponent(entity, Component{ENTITY_AWARENESS, awarnessComponent})
			}
		}
	}

	return returnEvents
}
