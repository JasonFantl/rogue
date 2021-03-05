package ecs

type ConsumableHandler struct{}

func (h *ConsumableHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == TRY_CONSUME {
		tryConsumeEvent := event.data.(TryConsume)

		_, canConsume := m.getComponent(tryConsumeEvent.what, CONSUMABLE)

		if canConsume {
			returnEvents = append(returnEvents, Event{CONSUMED, Consumed{event.entity}, tryConsumeEvent.what})
		}
	}

	return returnEvents
}
