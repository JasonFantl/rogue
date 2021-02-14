package ecs

type InventoryHandler struct{}

func (s *InventoryHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {
	if event.ID == TRY_PICK_UP_EVENT {

		// get entitys current position and make sure it has an inventory
		positionData, positionOk := m.getComponent(event.entity, POSITION)
		inventoryData, inventoryOk := m.getComponent(event.entity, INVENTORY)

		if positionOk && inventoryOk {
			positionComponent := positionData.(Position)
			inventoryComponent := inventoryData.(Inventory)

			// look for items below you to pickup
			for _, otherEntity := range m.getEntitiesFromPos(positionComponent.X, positionComponent.Y) {
				// make sure it is pickupable and hasnt been stashed
				_, pickupableOk := m.getComponent(otherEntity, PICKUPABLE)
				_, stashedOk := m.getComponent(otherEntity, STASHED)

				if pickupableOk && !stashedOk {

					// make sure the picker-uppers inventory is inited
					if cap(inventoryComponent.items) == 0 {
						inventoryComponent.items = make([]Entity, 0)
					}
					// how to pickup?

					// add stashed compoenet to item MAKE SURE TO REMOVE WHEN DROPPED!
					stashedComponent := Component{STASHED, Stashed{event.entity}}
					m.AddComponenet(otherEntity, stashedComponent)

					// then add it to our inventory
					inventoryComponent.items = append(inventoryComponent.items, otherEntity)
					m.lookupTable[INVENTORY][event.entity] = inventoryComponent
					returnEvents = append(returnEvents, Event{PICKED_UP_EVENT, EventPickedUp{event.entity}, otherEntity})
				}
			}
		}
	}

	// handle moving items in inventory with its parent
	if event.ID == MOVE_EVENT {
		moveEvent := event.data.(EventMove)

		// check for all stashed
		for stashedEntity, stashedData := range m.lookupTable[STASHED] {
			stashedComponent := stashedData.(Stashed)

			if stashedComponent.parent == event.entity {
				// they better have position
				positionData, positionOk := m.lookupTable[POSITION][stashedEntity]
				if positionOk {
					positionComponent := positionData.(Position)
					positionComponent.X = moveEvent.x
					positionComponent.Y = moveEvent.y
					m.lookupTable[POSITION][stashedEntity] = positionComponent
					// should we return an event? it would stop other stashed items from being moved
				}
			}
		}
	}

	return returnEvents
}
