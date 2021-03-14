package ecs

type Triggerable interface {
	trigger(*Manager, Event) []Event
}

type Effect struct {
	Trigger  EventID
	Reaction Triggerable
}

func selectEffected(event Event) Entity {
	switch event.ID {
	case CONSUMED:
		consumedEvent := event.data.(Consumed)
		return consumedEvent.byWho
	case PICKED_UP:
		pickedUpEvent := event.data.(PickedUp)
		return pickedUpEvent.byWho
	}
	return event.entity
}

type HealEffect struct {
	Amount int
}

func (effect HealEffect) trigger(m *Manager, event Event) (returnEvents []Event) {
	affected := selectEffected(event)

	// should probably send heal event rather then actually heal, but fine for now
	healthData, hasHealth := m.getComponent(affected, HEALTH)
	if hasHealth {
		healthComponent := healthData.(Health)
		healthComponent.Current += effect.Amount

		m.setComponent(affected, HEALTH, healthComponent)
	}

	return returnEvents
}

type VisionEffect struct {
	Amount int
}

func (effect VisionEffect) trigger(m *Manager, event Event) (returnEvents []Event) {
	affected := selectEffected(event)

	visionData, hasVision := m.getComponent(affected, VISION)
	if hasVision {
		visionComponent := visionData.(Vision)
		visionComponent.Radius += effect.Amount

		m.setComponent(affected, VISION, visionComponent)
	}

	return returnEvents
}

type StrengthEffect struct {
	Amount int
}

func (effect StrengthEffect) trigger(m *Manager, event Event) (returnEvents []Event) {

	affected := selectEffected(event)

	fighterData, isFighter := m.getComponent(affected, FIGHTER)
	if isFighter {
		fighterComponent := fighterData.(Fighter)
		fighterComponent.Strength += effect.Amount

		m.setComponent(affected, FIGHTER, fighterComponent)
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
