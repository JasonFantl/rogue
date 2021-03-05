package ecs

type Triggerable interface {
	trigger(*Manager, Event) []Event
}

type Effect struct {
	Trigger  EventID
	Reaction Triggerable
}

type HealEffect struct {
	Amount int
}

func (effect HealEffect) trigger(m *Manager, event Event) (returnEvents []Event) {

	affected := event.entity // probably not this entity

	switch event.ID {
	case CONSUMED:
		consumedEvent := event.data.(Consumed)
		affected = consumedEvent.byWho
	case PICKED_UP:
		pickedUpEvent := event.data.(PickedUp)
		affected = pickedUpEvent.byWho
	}

	// should probably send heal event rather then actually heal, but fine for now
	healthData, hasHealth := m.getComponent(affected, HEALTH)
	if hasHealth {
		healthComponent := healthData.(Health)
		healthComponent.Current += effect.Amount

		m.setComponent(affected, HEALTH, healthComponent)
	}

	return returnEvents
}

type EffectsHandler struct{}

func (h *EffectsHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	effectsData, hasEffects := m.getComponent(event.entity, EFFECTS)

	if hasEffects {
		effectsComponent := effectsData.(Effects)

		for _, effect := range effectsComponent.Effects {
			if effect.Trigger == event.ID {
				reactions := effect.Reaction.trigger(m, event)
				returnEvents = append(returnEvents, reactions...)
			}
		}
	}

	return returnEvents
}
