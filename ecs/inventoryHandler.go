package ecs

type InventoryHandler struct{}

func (h *InventoryHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == PLAYER_TRY_PICK_UP {
		positionData, hasPosition := m.getComponent(event.entity, POSITION)
		if hasPosition {
			positionComponent := positionData.(Position)

			entities := m.getEntitiesFromPos(positionComponent.X, positionComponent.Y)
			for item := range entities {
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
		returnEvents = append(returnEvents, Event{TRY_DROP, TryDrop{event.entity}, consumedEvent.byWho})
	}

	if event.ID == TRY_LAUNCH {
		tryLaunchEvent := event.data.(TryLaunch)
		returnEvents = append(returnEvents, Event{TRY_DROP, TryDrop{tryLaunchEvent.what}, event.entity})
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
				returnEvents = append(returnEvents, Event{TRY_DROP, TryDrop{item}, event.entity})
			}
		}
	}

	if event.ID == TRY_DROP {
		tryDropEvent := event.data.(TryDrop)
		inventoryData, hasInventory := m.getComponent(event.entity, INVENTORY)
		// stashableData, hasStashable := m.getComponent(try, INVENTORY)

		if hasInventory {
			inventoryComponent := inventoryData.(Inventory)
			delete(inventoryComponent.Items, tryDropEvent.what)
		}

		stashableData, isStashable := m.getComponent(tryDropEvent.what, STASHABLE)
		if isStashable {
			stashableComponent := stashableData.(Stashable)
			stashableComponent.Stashed = false
			m.setComponent(tryDropEvent.what, STASHABLE, stashableComponent)
		}

		returnEvents = append(returnEvents, Event{DROPED, Dropped{event.entity}, tryDropEvent.what})
	}

	return returnEvents
}

func (h *InventoryHandler) isTreasure(m *Manager, entity Entity) bool {
	stashableData, hasStashable := m.getComponent(entity, STASHABLE)

	if hasStashable {
		stashableComponent := stashableData.(Stashable)
		return !stashableComponent.Stashed
	}
	return false
}

// add stashed component to item MAKE SURE TO REMOVE WHEN DROPPED!
func (h *InventoryHandler) pickup(m *Manager, item, parent Entity, inventoryComponent Inventory) Event {
	// update stashed flag
	stashableData, hasStashable := m.getComponent(item, STASHABLE)

	if hasStashable {
		stashableComponent := stashableData.(Stashable)
		stashableComponent.Stashed = true
		m.setComponent(item, STASHABLE, stashableComponent)

	}

	// make sure the inventory is inited
	if inventoryComponent.Items == nil {
		inventoryComponent.Items = map[Entity]bool{}
	}

	// then add it to our inventory
	inventoryComponent.Items[item] = true
	m.setComponent(parent, INVENTORY, inventoryComponent)

	return Event{PICKED_UP, PickedUp{parent}, item}
}
