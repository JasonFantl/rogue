package ecs

import "github.com/nsf/termbox-go"

type AttackHandler struct {
}

func (s *AttackHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// trying to attack
	if event.ID == TRY_MOVE {
		// unpack event data
		moveEvent := event.data.(TryMove)

		// get entitys current position and if it can attack
		positionData, hasPosition := m.getComponent(event.entity, POSITION)
		violentData, hasViolent := m.getComponent(event.entity, VIOLENT)

		if hasPosition && hasViolent {
			positionComponent := positionData.(Position)
			violentComponent := violentData.(Violent)

			// now check if new location is occupied by something with health
			newX := positionComponent.X + moveEvent.dx
			newY := positionComponent.Y + moveEvent.dy

			// should you be able to attack everthing in the tile at once? yes for now
			// attacks each individually
			for _, otherEntity := range m.getEntitiesFromPos(newX, newY) {
				_, otherHasHealth := m.getComponent(otherEntity, HEALTH)

				if otherHasHealth {
					attackAmount := violentComponent.BaseAttackDmg
					returnEvents = append(returnEvents, Event{TRY_ATTACK, TryAttack{otherEntity, attackAmount}, event.entity})
				}
			}
		}
	}

	// do dmg
	if event.ID == TRY_ATTACK {
		// unpack event data
		tryAttackEvent := event.data.(TryAttack)

		// get attacked entities health
		healthData, hasHealth := m.getComponent(tryAttackEvent.who, HEALTH)

		if hasHealth {
			healthComponent := healthData.(Health)
			healthComponent.Current -= tryAttackEvent.dmg
			m.setComponent(tryAttackEvent.who, Component{HEALTH, healthComponent})

			returnEvents = append(returnEvents, Event{DAMAGED, Damaged{}, tryAttackEvent.who})
		}
	}

	// blood handler
	if event.ID == DAMAGED {

		// get relevant components
		healthData, hasHealth := m.getComponent(event.entity, HEALTH)
		positionData, hasPosition := m.getComponent(event.entity, POSITION)

		if hasHealth && hasPosition {
			healthComponent := healthData.(Health)
			positionComponent := positionData.(Position)

			if healthComponent.Current < healthComponent.Max/2 {

				bloodInfo := "hard to tell whos or whats blood this is"
				informationData, hasInformation := m.getComponent(event.entity, INFORMATION)
				if hasInformation {
					informationComponent := informationData.(Information)
					bloodInfo = "the blood of " + informationComponent.Name
				}

				blood := []Component{
					{POSITION, Position{positionComponent.X, positionComponent.Y}},
					{DISPLAYABLE, Displayable{false, termbox.ColorRed, ' ', 2}},
					{INFORMATION, Information{"Blood", bloodInfo}},
				}

				m.AddEntity(blood)
			}
		}
	}
	return returnEvents
}
