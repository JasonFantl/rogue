package ecs

import (
	"sort"
	"strconv"

	"github.com/jasonfantl/rogue/gui"
)

// Display is special, not like the other handlers
type DisplayHandler struct {
}

func (s *DisplayHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == DISPLAY {

		/////////////// GRID /////////////////

		// need to keep track of priorities
		// maps 2d pos to unique int
		priorities := make(map[int]int)
		maxX := 0

		// get new positions, the looping is currently as horrible as I can make it
		displayComponents, ok := m.getComponents(DISPLAYABLE)
		if ok {
			for entity, displayData := range displayComponents {
				positionData, positionOk := m.getComponent(entity, POSITION)

				if positionOk {
					positionComponent := positionData.(Position)
					displayComponent := displayData.(Displayable)

					x := positionComponent.X
					y := positionComponent.Y

					if x > maxX {
						maxX = x
					}

					uniqueID := x + (x+y)*(x+y+1)/2

					currentPriority, ok := priorities[uniqueID]

					if !ok || displayComponent.Priority > currentPriority {
						gui.DrawTile(x, y, displayComponent.Character)
						priorities[uniqueID] = displayComponent.Priority
					}
				}
			}
		}

		///////////// INVENTORY ///////////////////

		currentLineNum := 1
		inventoryComponents, ok := m.getComponents(INVENTORY)

		keys := make([]int, 0)
		for k, _ := range inventoryComponents {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)

		for _, key := range keys {
			entity := Entity(key)
			inventoryData := inventoryComponents[entity]

			inventoryComponent := inventoryData.(Inventory)

			// if we can, print entities information
			informationData, hasInformation := m.getComponent(entity, INFORMATION)
			healthData, hasHealth := m.getComponent(entity, HEALTH)

			if hasInformation {
				informationComponent := informationData.(Information)
				displayData := informationComponent.Name

				if hasHealth {
					healthComponent := healthData.(Health)
					displayData += " : " + strconv.Itoa(healthComponent.Current) + "/" + strconv.Itoa(healthComponent.Max)
				}
				gui.DrawText(maxX+3, currentLineNum, displayData)
				currentLineNum++
			}

			// then print each of its items
			for _, entity := range inventoryComponent.Items {
				informationData, informationOk := m.getComponent(entity, INFORMATION)

				if informationOk {
					informationComponent := informationData.(Information)

					informationString := informationComponent.Name + " : " + informationComponent.Details
					gui.DrawText(maxX+5, currentLineNum, informationString)
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
