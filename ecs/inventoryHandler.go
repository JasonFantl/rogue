package ecs

type InventoryHandler struct{}

func (h *InventoryHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == PLAYER_TRY_PICK_UP {
		positionData, hasPosition := m.getComponent(event.entity, POSITION)
		if hasPosition {
			positionComponent := positionData.(Position)

			entities := m.getEntitiesFromPos(positionComponent.X, positionComponent.Y)
			for _, item := range entities {
				if h.isTreasure(m, item) {
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
				if sameLocation && h.isTreasure(m, tryPickUpEvent.what) {
					response := h.pickup(m, tryPickUpEvent.what, event.entity, inventoryComponent)
					returnEvents = append(returnEvents, response)
				}
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

func (h *InventoryHandler) isTreasure(m *Manager, entity Entity) bool {
	_, pickupableOk := m.getComponent(entity, PICKUPABLE)
	_, stashedOk := m.getComponent(entity, STASHED_FLAG)

	return pickupableOk && !stashedOk
}

// add stashed component to item MAKE SURE TO REMOVE WHEN DROPPED!
func (h *InventoryHandler) pickup(m *Manager, entity, parent Entity, inventoryComponent Inventory) Event {
	stashedComponent := Component{STASHED_FLAG, StashedFlag{parent}}
	m.AddComponenet(entity, stashedComponent)

	// make sure the inventory is inited
	if inventoryComponent.Items == nil {
		inventoryComponent.Items = make(map[Entity]bool)
	}

	// then add it to our inventory
	inventoryComponent.Items[entity] = true
	m.setComponent(parent, INVENTORY, inventoryComponent)

	return Event{PICKED_UP, PickedUp{parent}, entity}
}
