package ecs

type InventoryHandler struct{}

func (s *InventoryHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// handy functions

	isTreasure := func(entity Entity) bool {
		_, pickupableOk := m.getComponent(entity, PICKUPABLE)
		_, stashedOk := m.getComponent(entity, STASHED_FLAG)

		return pickupableOk && !stashedOk
	}

	// add stashed component to item MAKE SURE TO REMOVE WHEN DROPPED!
	pickup := func(entity Entity, inventoryComponent Inventory) {
		stashedComponent := Component{STASHED_FLAG, StashedFlag{event.entity}}
		m.AddComponenet(entity, stashedComponent)

		// make sure the inventory is inited
		if inventoryComponent.Items == nil {
			inventoryComponent.Items = make(map[Entity]bool)
		}

		// then add it to our inventory
		inventoryComponent.Items[entity] = true
		m.setComponent(event.entity, INVENTORY, inventoryComponent)

		returnEvents = append(returnEvents, Event{PICKED_UP, PickedUp{event.entity}, entity})
	}

	if event.ID == PLAYER_TRY_PICK_UP {
		positionData, hasPosition := m.getComponent(event.entity, POSITION)
		if hasPosition {
			positionComponent := positionData.(Position)

			entities := m.getEntitiesFromPos(positionComponent.X, positionComponent.Y)
			for _, item := range entities {
				if isTreasure(item) {
					returnEvents = append(returnEvents, Event{TRY_PICK_UP, TryPickUp{item}, event.entity})
					break // dont need to check anymore
				}
			}
		}
	}

	if event.ID == TRY_PICK_UP {
		tryPickUpEvent := event.data.(TryPickUp)

		// get entitys current position and make sure it has an inventory
		positionData, hasPosition := m.getComponent(event.entity, POSITION)
		inventoryData, hasInventory := m.getComponent(event.entity, INVENTORY)

		if hasPosition && hasInventory {

			positionComponent := positionData.(Position)
			inventoryComponent := inventoryData.(Inventory)

			otherPositionData, otherHasPosition := m.getComponent(tryPickUpEvent.what, POSITION)
			if otherHasPosition {
				otherPositionComponent := otherPositionData.(Position)
				sameLocation := otherPositionComponent.X == positionComponent.X && otherPositionComponent.Y == positionComponent.Y
				if sameLocation && isTreasure(tryPickUpEvent.what) {
					pickup(tryPickUpEvent.what, inventoryComponent)
				}
			}
		}
	}

	if event.ID == PICKED_UP {
		pickedUpEvent := event.data.(PickedUp)

		FighterData, isFighter := m.getComponent(pickedUpEvent.byWho, FIGHTER)
		itemDamageData, itemDoesDamage := m.getComponent(event.entity, DAMAGE)

		if isFighter && itemDoesDamage {
			fighterComponent := FighterData.(Fighter)
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
	}

	if event.ID == CONSUMED {
		consumedEvent := event.data.(Consumed)
		inventoryData, hasInventory := m.getComponent(consumedEvent.byWho, INVENTORY)

		if hasInventory {
			inventoryComponent := inventoryData.(Inventory)
			delete(inventoryComponent.Items, event.entity)
		}
	}

	if event.ID == MOVED {
		moveEvent := event.data.(Moved)

		inventoryData, hasInventory := m.getComponent(event.entity, INVENTORY)

		if hasInventory {
			inventoryComponent := inventoryData.(Inventory)

			for item := range inventoryComponent.Items {
				// only have to do anything if the thing has postition
				positionData, hasPosition := m.getComponent(item, POSITION)
				if hasPosition {
					positionComponent := positionData.(Position)

					// should we announce each item we drag along? Silenced for now
					// returnEvents = append(returnEvents, Event{MOVED, Moved{positionComponent.X, positionComponent.Y, moveEvent.toX, moveEvent.toY}, item})

					positionComponent.X = moveEvent.toX
					positionComponent.Y = moveEvent.toY
					m.setComponent(item, POSITION, positionComponent)
				}
			}
		}
	}

	if event.ID == DIED {
		inventoryData, hasInventory := m.getComponent(event.entity, INVENTORY)

		if hasInventory {
			inventoryComponent := inventoryData.(Inventory)
			for item := range inventoryComponent.Items {
				m.removeComponent(item, STASHED_FLAG)
			}
		}
	}

	return returnEvents
}
