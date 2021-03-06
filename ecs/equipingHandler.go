package ecs

type EquippingHandler struct{}

func (s *EquippingHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == PICKED_UP {
		pickedUpEvent := event.data.(PickedUp)

		FighterData, isFighter := m.getComponent(pickedUpEvent.byWho, FIGHTER)

		if isFighter {
			fighterComponent := FighterData.(Fighter)

			itemDamageData, itemIsWeapon := m.getComponent(event.entity, DAMAGE)
			itemResistanceData, itemIsArmor := m.getComponent(event.entity, DAMAGE_RESISTANCE)

			if itemIsWeapon {
				itemDamageComponent := itemDamageData.(Damage)

				isBetterWeapon := true

				currentWeaponDamageData, currentDoesDamage := m.getComponent(fighterComponent.Weapon, DAMAGE)
				if currentDoesDamage {
					currentWeaponDamageComponent := currentWeaponDamageData.(Damage)
					if itemDamageComponent.Amount < currentWeaponDamageComponent.Amount {
						isBetterWeapon = false
					}
				}

				if isBetterWeapon {
					fighterComponent.Weapon = event.entity
					m.setComponent(pickedUpEvent.byWho, FIGHTER, fighterComponent)
				}
			}

			if itemIsArmor {
				itemResistanceComponent := itemResistanceData.(DamageResistance)

				isBetterArmor := true

				currentArmorResistanceData, currentIsResistant := m.getComponent(fighterComponent.Armor, DAMAGE_RESISTANCE)
				if currentIsResistant {
					currentArmorResistanceComponent := currentArmorResistanceData.(DamageResistance)
					if itemResistanceComponent.Amount < currentArmorResistanceComponent.Amount {
						isBetterArmor = false
					}
				}

				if isBetterArmor {
					fighterComponent.Armor = event.entity
				}
			}

			m.setComponent(pickedUpEvent.byWho, FIGHTER, fighterComponent)
		}
	}

	return returnEvents
}
