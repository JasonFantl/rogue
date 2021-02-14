package ecs

import (
	"fmt"

	"github.com/jasonfantl/rogue/gui"
)

/*
How Handlers work:

They get passed all events, they can then check if its something they are intereted in

	if event.ID == RELEVANT_EVENT

then we type convert the event

		relevantEvent := event.data.(RelevantEvent)

then the relevant components are found for the entity that triggered the event

		someData, someOk := m.getComponent(event.entity, SOME)
		if someOk {
			someComponent := someData.(Some)
			... code using component ...
		}

and sometimes we want to update the component, but that requires we refrence the actual compoenent

			m.lookupTable[SOME][event.entity] = someComponent
*/

type EventHandler interface {
	handleEvent(*Manager, Event) (returnEvents []Event)
}

// Display is special, not like the other handlers
type DisplayHandler struct {
}

func (s *DisplayHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	/////////////// GRID /////////////////
	// need to keep track of priorities
	// maps 2d pos to unique int
	priorities := make(map[int]int)
	maxX := 0

	// get new positions, the looping is currently as horrible as I can make it
	for entity, displayData := range m.lookupTable[DISPLAY] {
		positionData, positionOk := m.lookupTable[POSITION][entity]

		if positionOk {
			positionComponent := positionData.(Position)
			displayComponent := displayData.(Display)

			x := positionComponent.X
			y := positionComponent.Y

			if x > maxX {
				maxX = x
			}

			uniqueID := x + (x+y)*(x+y+1)/2

			currentPriority, ok := priorities[uniqueID]

			if !ok || displayComponent.Priority > currentPriority {
				gui.DrawTile(x, y, displayComponent.Character, displayComponent.Style)
				priorities[uniqueID] = displayComponent.Priority
			}
		}
	}

	///////////// INVENTORY ///////////////////

	currentLineNum := 1
	for entity, inventoryData := range m.lookupTable[INVENTORY] {

		inventoryComponent := inventoryData.(Inventory)

		// if we can, print information
		informationData, informationOk := m.lookupTable[INFORMATION][entity]
		if informationOk {
			informationComponent := informationData.(Information)
			gui.DrawText(maxX+3, currentLineNum, informationComponent.Name)
			currentLineNum++
		}

		for _, entity := range inventoryComponent.items {
			informationData := m.lookupTable[INFORMATION][entity]
			informationComponent, informationOk := informationData.(Information)

			if informationOk {
				itemData := informationComponent.Name + " : " + informationComponent.Details
				gui.DrawText(maxX+5, currentLineNum, itemData)
				currentLineNum++
			} else {
				gui.DrawText(maxX+5, currentLineNum, "? : no information on item")
				currentLineNum++
			}
		}
	}

	gui.Show()

	return returnEvents
}

type MoveHandler struct {
}

func (s *MoveHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {
	// somwhere in here the components are copied, so we cant edit them
	// that why we set the components at the end
	// there must be a way to maintain the pointer to the component

	if event.ID == TRY_MOVE_EVENT {
		// unpack event data
		moveEvent := event.data.(EventTryMove)

		// get entitys current position and if its blocking
		positionData, positionOk := m.getComponent(event.entity, POSITION)
		_, blockingOk := m.getComponent(event.entity, BLOCKABLE)

		if positionOk {
			positionComponent := positionData.(Position)

			// now check if new location is occupied
			newX := positionComponent.X + moveEvent.dx
			newY := positionComponent.Y + moveEvent.dy
			canMove := true
			// quick implementation, replace later
			for otherEntity, otherData := range m.lookupTable[POSITION] {
				otherPositionComponent := otherData.(Position)
				_, otherBlockableOk := m.lookupTable[BLOCKABLE][otherEntity]

				if otherPositionComponent.X == newX && otherPositionComponent.Y == newY {
					if otherBlockableOk || blockingOk {
						canMove = false
					}
				}
			}

			if canMove {
				positionComponent.X = newX
				positionComponent.Y = newY

				m.lookupTable[POSITION][event.entity] = positionComponent

				returnEvents = append(returnEvents, Event{MOVE_EVENT, EventMove{newX, newY}, event.entity})
			}
		}
	}
	return returnEvents
}

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

type EventPrinterHandler struct{}

func (s *EventPrinterHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	stringifiedEvent := fmt.Sprintf("%T : %v", event.data, event.entity)

	//special case
	if event.ID == ERROR_EVENT {
		stringifiedEvent = fmt.Sprintf("%T : %s : %v", event.data, event.data.(EventError).err, event.entity)
	}

	gui.UpdateErrors(stringifiedEvent)
	gui.Show()

	return returnEvents
}
