package ecs

import (
	"sort"
	"strconv"

	"github.com/jasonfantl/rogue/gui"
)

// Display is special, not like the other handlers
type DisplayHandler struct {
}

func (h *DisplayHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == DISPLAY {
		_, playerExist := m.getComponent(event.entity, INVENTORY)
		if playerExist {
			h.showEntity(m, event.entity)
		} else {
			h.showAll(m)
		}
	}

	return returnEvents
}

func (s *DisplayHandler) showEntity(m *Manager, entity Entity) {

	/////////////// GRID /////////////////

	// need to keep track of priorities
	// maps 2d pos to unique int
	fgPriorities := make(map[int]int)
	bgPriorities := make(map[int]int)

	maxX := 0
	displayOffset := 10

	awarnessData, hasAwarness := m.getComponent(entity, ENTITY_AWARENESS)
	positionData, hasPosition := m.getComponent(entity, POSITION)

	if hasAwarness && hasPosition {
		awarnessComponent := awarnessData.(EntityAwarness)
		positionComponent := positionData.(Position)

		for _, item := range awarnessComponent.AwareOf {
			displayData, hasDisplay := m.getComponent(item, DISPLAYABLE)
			positionData, hasPosition := m.getComponent(item, POSITION)

			if hasDisplay && hasPosition {
				seenDisplayComponent := displayData.(Displayable)
				seenPositionComponent := positionData.(Position)

				x := seenPositionComponent.X - positionComponent.X + displayOffset
				y := seenPositionComponent.Y - positionComponent.Y + displayOffset

				if x > maxX {
					maxX = x
				}

				uniqueID := x + (x+y)*(x+y+1)/2

				if seenDisplayComponent.IsForeground {
					currentPriority, ok := fgPriorities[uniqueID]
					if !ok || seenDisplayComponent.Priority > currentPriority {
						gui.DrawFg(x, y, seenDisplayComponent.Rune, seenDisplayComponent.Color)
						fgPriorities[uniqueID] = seenDisplayComponent.Priority
					}
				} else {
					currentPriority, ok := bgPriorities[uniqueID]
					if !ok || seenDisplayComponent.Priority > currentPriority {
						gui.DrawBg(x, y, seenDisplayComponent.Color)
						bgPriorities[uniqueID] = seenDisplayComponent.Priority
					}
				}
			}
		}
	}

	///////////// INVENTORY ///////////////////

	currentLineNum := 1
	inventoryData, hasInventory := m.getComponent(entity, INVENTORY)

	if hasInventory {
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
		for _, item := range inventoryComponent.Items {
			informationData, informationOk := m.getComponent(item, INFORMATION)

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
}

func (s *DisplayHandler) showAll(m *Manager) {

	/////////////// GRID /////////////////

	// need to keep track of priorities
	// maps 2d pos to unique int
	fgPriorities := make(map[int]int)
	bgPriorities := make(map[int]int)

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

				if displayComponent.IsForeground {
					currentPriority, ok := fgPriorities[uniqueID]
					if !ok || displayComponent.Priority > currentPriority {
						gui.DrawFg(x, y, displayComponent.Rune, displayComponent.Color)
						fgPriorities[uniqueID] = displayComponent.Priority
					}
				} else {
					currentPriority, ok := bgPriorities[uniqueID]
					if !ok || displayComponent.Priority > currentPriority {
						gui.DrawBg(x, y, displayComponent.Color)
						bgPriorities[uniqueID] = displayComponent.Priority
					}
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
}
