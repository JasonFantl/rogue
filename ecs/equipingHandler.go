package ecs

type EquippingHandler struct{}

func (h *EquippingHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// auto equipper
	if event.ID == PICKED_UP {
		pickedUpEvent := event.data.(PickedUp)

		FighterData, isFighter := m.getComponent(pickedUpEvent.byWho, FIGHTER)

		if isFighter {
			fighterComponent := FighterData.(Fighter)

			itemDamageData, itemIsWeapon := m.getComponent(event.entity, DAMAGE)
			if itemIsWeapon {
				itemDamageComponent := itemDamageData.(Damage)

				isBetter := true

				currentWeaponDamageData, currentDoesDamage := m.getComponent(fighterComponent.Weapon, DAMAGE)
				if currentDoesDamage {
					currentWeaponDamageComponent := currentWeaponDamageData.(Damage)
					if itemDamageComponent.Amount <= currentWeaponDamageComponent.Amount {
						isBetter = false
					}
				}

				if isBetter {
					returnEvents = append(returnEvents, Event{TRY_EQUIP_WEAPON, TryEquip{event.entity}, pickedUpEvent.byWho})
				}
			}

			itemResistanceData, itemIsArmor := m.getComponent(event.entity, DAMAGE_RESISTANCE)
			if itemIsArmor {
				itemResistanceComponent := itemResistanceData.(DamageResistance)

				isBetter := true

				currentArmorResistanceData, currentIsResistant := m.getComponent(fighterComponent.Armor, DAMAGE_RESISTANCE)
				if currentIsResistant {
					currentArmorResistanceComponent := currentArmorResistanceData.(DamageResistance)
					if itemResistanceComponent.Amount <= currentArmorResistanceComponent.Amount {
						isBetter = false
					}
				}
				if isBetter {
					returnEvents = append(returnEvents, Event{TRY_EQUIP_ARMOR, TryEquip{event.entity}, pickedUpEvent.byWho})
				}
			}
		}
	}

	if event.ID == TRY_EQUIP_WEAPON {
		tryEquipEvent := event.data.(TryEquip)
		FighterData, isFighter := m.getComponent(event.entity, FIGHTER)
		if isFighter {
			fighterComponent := FighterData.(Fighter)

			fighterComponent.Weapon = tryEquipEvent.what
			m.setComponent(event.entity, FIGHTER, fighterComponent)
		}
	}

	if event.ID == TRY_EQUIP_ARMOR {
		tryEquipEvent := event.data.(TryEquip)
		FighterData, isFighter := m.getComponent(event.entity, FIGHTER)
		if isFighter {
			fighterComponent := FighterData.(Fighter)

			fighterComponent.Armor = tryEquipEvent.what
			m.setComponent(event.entity, FIGHTER, fighterComponent)
		}
	}

	if event.ID == DROPED {
		droppedEvent := event.data.(Dropped)
		FighterData, isFighter := m.getComponent(droppedEvent.byWho, FIGHTER)

		if isFighter {
			fighterComponent := FighterData.(Fighter)

			if event.entity == fighterComponent.Weapon {
				fighterComponent.Weapon = Entity(0)
			}
			if event.entity == fighterComponent.Armor {
				fighterComponent.Armor = Entity(0)
			}

			m.setComponent(droppedEvent.byWho, FIGHTER, fighterComponent)

		}
	}

	return returnEvents
}
