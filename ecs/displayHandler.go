package ecs

import (
	"sort"
	"strconv"

	"github.com/jasonfantl/rogue/gui"
	"github.com/nsf/termbox-go"
)

// Display is special, not like the other handlers
type DisplayHandler struct {
}

func (h *DisplayHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == DISPLAY {
		h.showEntity(m, event.entity)
	}

	return returnEvents
}

func (s *DisplayHandler) showEntity(m *Manager, entity Entity) {

	// get and check all the components necessary
	positionData, hasPosition := m.getComponent(entity, POSITION)
	memoryData, hasMemory := m.getComponent(entity, ENTITY_MEMORY)
	awarnessData, hasAwarness := m.getComponent(entity, ENTITY_AWARENESS)

	inventoryData, hasInventory := m.getComponent(entity, INVENTORY)

	informationData, hasInformation := m.getComponent(entity, INFORMATION)
	healthData, hasHealth := m.getComponent(entity, HEALTH)

	fighterData, isFighter := m.getComponent(entity, FIGHTER)

	/////////////// GRID /////////////////

	// need to keep track of priorities
	// maps 2d pos to unique int
	fgPriorities := make(map[int]map[int]int)
	bgPriorities := make(map[int]map[int]int)

	// main function
	display := func(x, y int, displayComponent Displayable) {
		priorityMap := bgPriorities
		if displayComponent.IsForeground {
			priorityMap = fgPriorities
		}

		_, ok := priorityMap[x]
		if !ok {
			priorityMap[x] = make(map[int]int)
		}
		currentPriority, ok := priorityMap[x][y]

		if !ok || displayComponent.Priority > currentPriority {
			if displayComponent.IsForeground {
				gui.DrawFg(x, y, displayComponent.Rune, displayComponent.Color)
			} else {
				gui.DrawBg(x, y, displayComponent.Color)
			}
			priorityMap[x][y] = displayComponent.Priority
		}
	}

	displayRadius := 20

	if hasPosition {
		positionComponent := positionData.(Position)

		visionData, hasVision := m.getComponent(entity, VISION)
		if hasVision {
			visionComponent := visionData.(Vision)
			displayRadius = visionComponent.Radius
		}

		// ---------- MAP ---------

		if hasAwarness && hasPosition {
			awarnessComponent := awarnessData.(EntityAwarness)
			positionComponent := positionData.(Position)

			displayCell := func(dx, dy int) {
				itemX := positionComponent.X + dx
				itemY := positionComponent.Y + dy

				// display items we are aware of
				wasAware := false
				items := awarnessComponent.AwareOf.get(itemX, itemY)

				for _, item := range items {
					itemDisplayData, itemHasDisplay := m.getComponent(item, DISPLAYABLE)
					if itemHasDisplay {
						itemDisplayComponent := itemDisplayData.(Displayable)

						wasAware = true
						display(dx, dy, itemDisplayComponent)
					}
				}

				// if no aware items, try memory
				if !wasAware && hasMemory {
					memoryComponent := memoryData.(EntityMemory)
					col, ok := memoryComponent.Memory[itemX]
					if ok {
						toDisplay, ok := col[itemY]
						if ok {
							// convert display to fadded memory
							r, g, b := termbox.AttributeToRGB(toDisplay.Color)
							fadeConstant := 4
							fadedColor := termbox.RGBToAttribute(r/uint8(fadeConstant), g/uint8(fadeConstant), b/uint8(fadeConstant))
							fadedMemory := Displayable{toDisplay.IsForeground, fadedColor, toDisplay.Rune, 1}
							display(dx, dy, fadedMemory)
						}
					}
				}
			}

			bounds := getOctantBounds(displayRadius)

			for row, col := range bounds {
				for dx := -col; dx < col; dx++ {
					displayCell(dx, row)
					displayCell(dx, -row)
				}
				if row == len(bounds)-1 || bounds[row+1] != col {
					for dx := -row; dx < row; dx++ {
						displayCell(dx, col)
						displayCell(dx, -col)
					}
				}
			}
		}

		// --------- BELOW YOU --------------

		items := make([]Entity, 0)
		belowYou := m.getEntitiesFromPos(positionComponent.X, positionComponent.Y)
		for _, item := range belowYou {
			_, hasPickupable := m.getComponent(item, PICKUPABLE)
			_, isStashed := m.getComponent(item, STASHED_FLAG)
			if hasPickupable && !isStashed {
				items = append(items, item)
			}
		}

		displayString := ""
		for _, item := range items {
			informationData, informationOk := m.getComponent(item, INFORMATION)
			if informationOk {
				informationComponent := informationData.(Information)
				displayString += informationComponent.Name + ", "
			} else {
				displayString += "?, "
			}
		}
		if displayString != "" {
			displayString = displayString[:len(displayString)-2]
			gui.DrawText(-len(displayString)/2, displayRadius, displayString)
		}
	}

	// ------------- PLAYER STATS ---------------

	if hasInformation {
		informationComponent := informationData.(Information)
		gui.DrawText(-len(informationComponent.Name)/2, -displayRadius-4, informationComponent.Name)
	}
	if hasHealth {
		healthComponent := healthData.(Health)
		healthString := "HP " + strconv.Itoa(healthComponent.Current) + "/" + strconv.Itoa(healthComponent.Max)
		gui.DrawText(-displayRadius-len(healthString)/2, -displayRadius, healthString)
	}
	if isFighter {
		fighterComponent := fighterData.(Fighter)
		strengthString := "STR " + strconv.Itoa(fighterComponent.Strength)
		gui.DrawText(displayRadius/2+len(strengthString)/2, -displayRadius, strengthString)
	}

	// -------------- INVENTORY ---------------

	if hasInventory {
		inventoryComponent := inventoryData.(Inventory)

		currentLineNum := -displayRadius + 4
		gui.DrawText(displayRadius+1, currentLineNum, "Inventory: ")
		currentLineNum++

		// then print each of its items (in sorted order)
		keys := make([]int, 0)
		for k := range inventoryComponent.Items {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)

		for _, key := range keys {
			item := Entity(key)
			informationData, informationOk := m.getComponent(item, INFORMATION)

			if informationOk {
				informationComponent := informationData.(Information)

				informationString := informationComponent.Name + " : " + informationComponent.Details
				gui.DrawText(displayRadius+5, currentLineNum, informationString)
				currentLineNum++
			} else {
				gui.DrawText(displayRadius+5, currentLineNum, "? : no information on item")
				currentLineNum++
			}
		}
	}

	// -------- EQUIPED ----------

	if isFighter {
		fighterComponent := fighterData.(Fighter)

		informationData, informationOk := m.getComponent(fighterComponent.Weapon, INFORMATION)

		if informationOk {
			informationComponent := informationData.(Information)

			gui.DrawText(-displayRadius-10, -1, "Weapon:")
			gui.DrawText(-displayRadius-5, 0, informationComponent.Name)

		}
	}

}
