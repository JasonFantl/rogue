package ecs

type ConsumptionHandler struct {
}

func (h *ConsumptionHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == TRY_CONSUME {
		tryConsumeEvent := event.data.(TryConsume)

		_, isConsumable := m.getComponent(tryConsumeEvent.what, CONSUMABLE)

		if isConsumable {
			returnEvents = append(returnEvents, Event{CONSUMED, Consumed{event.entity}, tryConsumeEvent.what})
		}
	}

	return returnEvents
}
