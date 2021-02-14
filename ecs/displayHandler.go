package ecs

import (
	"github.com/jasonfantl/rogue/gui"
)

// Display is special, not like the other handlers
type DisplayHandler struct {
}

func (s *DisplayHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == DISPLAY_EVENT {

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
	}

	return returnEvents
}
