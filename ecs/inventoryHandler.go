package ecs

type InventoryHandler struct{}

func (s *InventoryHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {
	if event.ID == TRY_PICK_UP {
		tryPickUpEvent := event.data.(TryPickUp)

		// get entitys current position and make sure it has an inventory
		positionData, positionOk := m.getComponent(event.entity, POSITION)
		inventoryData, inventoryOk := m.getComponent(event.entity, INVENTORY)

		if positionOk && inventoryOk {
			positionComponent := positionData.(Position)
			inventoryComponent := inventoryData.(Inventory)

			// make sure the picker-uppers inventory is inited
			if cap(inventoryComponent.Items) == 0 {
				inventoryComponent.Items = make([]Entity, 0)
			}

			// handy functions
			// make sure it is pickupable and hasnt been stashed
			isTreasure := func(entity Entity) bool {
				_, pickupableOk := m.getComponent(entity, PICKUPABLE)
				_, stashedOk := m.getComponent(entity, STASHED_FLAG)

				return pickupableOk && !stashedOk
			}

			// add stashed compoenet to item MAKE SURE TO REMOVE WHEN DROPPED!
			pickup := func(entity Entity) {
				stashedComponent := Component{STASHED_FLAG, StashedFlag{event.entity}}
				m.AddComponenet(entity, stashedComponent)

				// then add it to our inventory
				inventoryComponent.Items = append(inventoryComponent.Items, entity)
				m.setComponent(event.entity, Component{INVENTORY, inventoryComponent})

				returnEvents = append(returnEvents, Event{PICKED_UP, PickedUp{event.entity}, entity})
			}

			// check if were picking up one item or everything

			otherPositionData, otherPositionOk := m.getComponent(tryPickUpEvent.what, POSITION)
			// should you be able to pick up an item without position?
			if otherPositionOk {
				otherPositionComponent := otherPositionData.(Position)
				if otherPositionComponent.X == positionComponent.X && otherPositionComponent.Y == positionComponent.Y && isTreasure(tryPickUpEvent.what) {
					pickup(tryPickUpEvent.what)
				}
			}

		}
	}

	if event.ID == MOVED {
		moveEvent := event.data.(Moved)

		// check for all stashed
		stashedComponents, ok := m.getComponents(STASHED_FLAG)
		if ok {
			for stashedEntity, stashedData := range stashedComponents {
				stashedComponent := stashedData.(StashedFlag)

				if stashedComponent.Parent == event.entity {
					// only have to do anything if the thing has postition
					positionData, positionOk := m.getComponent(stashedEntity, POSITION)
					if positionOk {
						positionComponent := positionData.(Position)
						positionComponent.X = moveEvent.x
						positionComponent.Y = moveEvent.y
						m.setComponent(stashedEntity, Component{POSITION, positionComponent})

						returnEvents = append(returnEvents, Event{MOVED, Moved{positionComponent.X, positionComponent.Y}, stashedEntity})
					}
				}
			}
		}
	}

	if event.ID == DIED {
		// loop through all stashed
		components, ok := m.getComponents(STASHED_FLAG)
		if ok {
			for stashedEntity, stashedData := range components {
				stashedComponent := stashedData.(StashedFlag)

				if stashedComponent.Parent == event.entity {
					m.removeComponent(stashedEntity, STASHED_FLAG)
				}
			}
		}
	}

	return returnEvents
}
