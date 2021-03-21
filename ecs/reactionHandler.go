package ecs

type Reaction struct {
	ReactionType EventID
	Reaction     EventHandler
}

func selectAffected(event Event) Entity {
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

type HealReaction struct {
	Amount int
}

func (reaction HealReaction) handleEvent(m *Manager, event Event) (returnEvents []Event) {
	affected := selectAffected(event)

	// should probably send heal event rather then actually heal, but fine for now
	healthData, hasHealth := m.getComponent(affected, HEALTH)
	if hasHealth {
		healthComponent := healthData.(Health)
		healthComponent.Current += reaction.Amount

		m.setComponent(affected, HEALTH, healthComponent)
	}

	return returnEvents
}

type VisionIncreaseReaction struct {
	Amount int
}

func (reaction VisionIncreaseReaction) handleEvent(m *Manager, event Event) (returnEvents []Event) {
	affected := selectAffected(event)

	visionData, hasVision := m.getComponent(affected, VISION)
	if hasVision {
		visionComponent := visionData.(Vision)
		visionComponent.Radius += reaction.Amount

		m.setComponent(affected, VISION, visionComponent)

		returnEvents = append(returnEvents, Event{SETTING_CHANGE, SettingChange{"zoom", 1}, m.user.Controlling})
	}

	return returnEvents
}

type StrengthIncreaseReaction struct {
	Amount int
}

func (reaction StrengthIncreaseReaction) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	affected := selectAffected(event)

	fighterData, isFighter := m.getComponent(affected, FIGHTER)
	if isFighter {
		fighterComponent := fighterData.(Fighter)
		fighterComponent.Strength += reaction.Amount

		m.setComponent(affected, FIGHTER, fighterComponent)
	}

	return returnEvents
}

type ReactionHandler struct{}

func (h *ReactionHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	reactionsData, hasReactions := m.getComponent(event.entity, REACTIONS)

	if hasReactions {
		reactionsComponent := reactionsData.(Reactions)

		for _, reaction := range reactionsComponent.Reactions {
			if reaction.ReactionType == event.ID {
				reactions := reaction.Reaction.handleEvent(m, event)
				returnEvents = append(returnEvents, reactions...)
			}
		}
	}

	return returnEvents
}
