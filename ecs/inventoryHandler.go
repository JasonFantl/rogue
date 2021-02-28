package ecs

type InventoryHandler struct{}

func (s *InventoryHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {
	if event.ID == PLAYER_TRY_PICK_UP {
		positionData, hasPosition := m.getComponent(event.entity, POSITION)
		if hasPosition {
			positionComponent := positionData.(Position)

			entities := m.getEntitiesFromPos(positionComponent.X, positionComponent.Y)
			for _, entity := range entities {
				_, hasPickUpAble := m.getComponent(entity, PICKUPABLE)
				_, hasStashed := m.getComponent(entity, STASHED_FLAG)
				isTreasure := hasPickUpAble && !hasStashed
				if isTreasure {
					returnEvents = append(returnEvents, Event{TRY_PICK_UP, TryPickUp{entity}, event.entity})
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
			if event.entity == 0 {
				returnEvents = append(returnEvents, Event{ERROR_EVENT, ErrorEvent{"player picking up"}, 0})
			}

			otherPositionData, otherHasPosition := m.getComponent(tryPickUpEvent.what, POSITION)
			// should you be able to pick up an item without position?
			if otherHasPosition {
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
					positionData, hasPosition := m.getComponent(stashedEntity, POSITION)
					if hasPosition {
						positionComponent := positionData.(Position)

						returnEvents = append(returnEvents, Event{MOVED, Moved{positionComponent.X, positionComponent.Y, moveEvent.toX, moveEvent.toY}, stashedEntity})

						positionComponent.X = moveEvent.toX
						positionComponent.Y = moveEvent.toY
						m.setComponent(stashedEntity, Component{POSITION, positionComponent})
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
