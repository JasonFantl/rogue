package ecs

type DeathHandler struct {
}

func (s *DeathHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == DAMAGED {
		// get entitys health
		healthData, hasHealth := m.getComponent(event.entity, HEALTH)

		if hasHealth {
			healthComponent := healthData.(Health)
			if healthComponent.Current <= 0 {
				returnEvents = append(returnEvents, Event{DIED, Died{}, event.entity})
			}
		}
	}

	// might be nice if there left a dead body
	// right now we just delete the entities
	if event.ID == DIED {
		m.removeEntity(event.entity)
	}
	return returnEvents
}
