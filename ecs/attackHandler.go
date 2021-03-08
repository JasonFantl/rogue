package ecs

import (
	"math/rand"

	"github.com/jasonfantl/rogue/gui"
)

type AttackHandler struct {
}

func (s *AttackHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// trying to attack
	if event.ID == TRY_MOVE {
		// unpack event data
		moveEvent := event.data.(TryMove)

		// get entitys current position and if it can attack
		positionData, hasPosition := m.getComponent(event.entity, POSITION)
		_, isfighter := m.getComponent(event.entity, FIGHTER)

		if hasPosition && isfighter {
			positionComponent := positionData.(Position)

			// now check if new location is occupied by something with health
			newX := positionComponent.X + moveEvent.dx
			newY := positionComponent.Y + moveEvent.dy

			// should you be able to attack everthing in the tile at once? yes for now
			// attacks each individually
			for _, otherEntity := range m.getEntitiesFromPos(newX, newY) {
				_, otherHasHealth := m.getComponent(otherEntity, HEALTH)

				if otherHasHealth {
					returnEvents = append(returnEvents, Event{TRY_ATTACK, TryAttack{otherEntity}, event.entity})
				}
			}
		}
	}

	// do dmg
	if event.ID == TRY_ATTACK {
		// unpack event data
		tryAttackEvent := event.data.(TryAttack)

		// get components
		attackedHealthData, attackedHasHealth := m.getComponent(tryAttackEvent.who, HEALTH)
		fighterData, isFighter := m.getComponent(event.entity, FIGHTER)

		if attackedHasHealth && isFighter {
			attackedHealthComponent := attackedHealthData.(Health)
			fighterComponent := fighterData.(Fighter)

			// base dmage
			damage := fighterComponent.Strength

			// weapon damage
			// check if we have a weapon, otherwise use self as weapon
			weapon := fighterComponent.Weapon
			if weapon == 0 {
				weapon = event.entity
			}
			damageData, doesDamage := m.getComponent(weapon, DAMAGE)
			if doesDamage {
				damageComponent := damageData.(Damage)
				damage += damageComponent.Amount
			}

			// armor protection
			attackedFighterData, attackedIsFighter := m.getComponent(tryAttackEvent.who, FIGHTER)
			if attackedIsFighter {
				attackedFighterComponent := attackedFighterData.(Fighter)

				// check if we have a weapon, otherwise use self as weapon
				armor := attackedFighterComponent.Armor
				if armor == 0 {
					armor = tryAttackEvent.who
				}

				armorData, attackedHasArmor := m.getComponent(armor, DAMAGE_RESISTANCE)
				if attackedHasArmor {
					armorComponent := armorData.(DamageResistance)
					// how to use AC?
					if rand.Intn(damage+1) < armorComponent.Amount {
						damage = 0
					}
				}
			}

			attackedHealthComponent.Current -= damage
			m.setComponent(tryAttackEvent.who, HEALTH, attackedHealthComponent)

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
					{DISPLAYABLE, Displayable{gui.GetSprite(gui.BLOOD)}},
					{INFORMATION, Information{"Blood", bloodInfo}},
				}

				m.AddEntity(blood)
			}
		}
	}
	return returnEvents
}
