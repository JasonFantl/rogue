package ecs

type InventoryHandler struct{}

func (s *InventoryHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {
	if event.ID == TRY_PICK_UP_EVENT {
		tryPickUpEvent := event.data.(EventTryPickUp)

		// get entitys current position and make sure it has an inventory
		positionData, positionOk := m.getComponent(event.entity, POSITION)
		inventoryData, inventoryOk := m.getComponent(event.entity, INVENTORY)

		if positionOk && inventoryOk {
			positionComponent := positionData.(Position)
			inventoryComponent := inventoryData.(Inventory)

			// make sure the picker-uppers inventory is inited
			if cap(inventoryComponent.items) == 0 {
				inventoryComponent.items = make([]Entity, 0)
			}

			// handy functions
			// make sure it is pickupable and hasnt been stashed
			isTreasure := func(entity Entity) bool {
				_, pickupableOk := m.getComponent(entity, PICKUPABLE)
				_, stashedOk := m.getComponent(entity, STASHED)

				return pickupableOk && !stashedOk
			}

			// add stashed compoenet to item MAKE SURE TO REMOVE WHEN DROPPED!
			pickup := func(entity Entity) {
				stashedComponent := Component{STASHED, Stashed{event.entity}}
				m.AddComponenet(entity, stashedComponent)

				// then add it to our inventory
				inventoryComponent.items = append(inventoryComponent.items, entity)
				m.setComponent(event.entity, Component{INVENTORY, inventoryComponent})

				returnEvents = append(returnEvents, Event{PICKED_UP_EVENT, EventPickedUp{event.entity}, entity})
			}

			// check if were picking up one item or everything
			if tryPickUpEvent.oneItem {
				otherPositionData, otherPositionOk := m.getComponent(tryPickUpEvent.what, POSITION)
				// should you be able to pick up an item without position?
				if otherPositionOk {
					otherPositionComponent := otherPositionData.(Position)
					if otherPositionComponent.X == positionComponent.X && otherPositionComponent.Y == positionComponent.Y && isTreasure(tryPickUpEvent.what) {
						pickup(tryPickUpEvent.what)
					}
				}
			} else {
				// look for items below you to pickup
				for _, otherEntity := range m.getEntitiesFromPos(positionComponent.X, positionComponent.Y) {
					if isTreasure(otherEntity) {
						pickup(otherEntity)
					}
				}
			}
		}
	} else if event.ID == MOVE_EVENT {
		moveEvent := event.data.(EventMove)

		// check for all stashed
		components, ok := m.getComponents(STASHED)
		if ok {
			for stashedEntity, stashedData := range components {
				stashedComponent := stashedData.(Stashed)

				if stashedComponent.parent == event.entity {
					// only have to do anything if the thing has postition
					positionData, positionOk := m.getComponent(stashedEntity, POSITION)
					if positionOk {
						positionComponent := positionData.(Position)
						positionComponent.X = moveEvent.x
						positionComponent.Y = moveEvent.y
						m.setComponent(stashedEntity, Component{POSITION, positionComponent})

						returnEvents = append(returnEvents, Event{MOVE_EVENT, EventMove{positionComponent.X, positionComponent.Y}, stashedEntity})
					}
				}
			}
		}
	}

	return returnEvents
}
