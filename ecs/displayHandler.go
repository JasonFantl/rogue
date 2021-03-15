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
		gui.Clear()
		h.showEntity(m, m.user.Controlling)
		if m.user.Menu.active {
			h.showMenu(m)
		}
	}

	return returnEvents
}

func (h *DisplayHandler) showEntity(m *Manager, entity Entity) {

	// display is a circle
	displayRadius := 10

	visionData, hasVision := m.getComponent(entity, VISION)
	if hasVision {
		visionComponent := visionData.(Vision)
		displayRadius = visionComponent.Radius
	}

	h.showGrid(m, entity, displayRadius)
	h.showBelowYou(m, entity, displayRadius)
	h.showStats(m, entity, displayRadius)
	// h.showInventory(m, entity, displayRadius)
	h.showEquiped(m, entity, displayRadius)

}

func (s *DisplayHandler) showGrid(m *Manager, entity Entity, displayRadius int) {
	positionData, hasPosition := m.getComponent(entity, POSITION)
	memoryData, hasMemory := m.getComponent(entity, ENTITY_MEMORY)
	awarnessData, hasAwarness := m.getComponent(entity, ENTITY_AWARENESS)

	if hasPosition && hasAwarness {
		positionComponent := positionData.(Position)
		awarnessComponent := awarnessData.(EntityAwarness)

		displayXY := func(dx, dy int) {
			itemX := positionComponent.X + dx
			itemY := positionComponent.Y + dy

			// display items we are aware of
			items := awarnessComponent.AwareOf.get(itemX, itemY)

			displayables := make([]gui.Sprite, 0)
			for item := range items {
				itemDisplayData, itemHasDisplay := m.getComponent(item, DISPLAYABLE)
				if itemHasDisplay {
					// we want to ignore stashed items
					_, itemIsStashed := m.getComponent(item, STASHED_FLAG)
					if !itemIsStashed {
						itemDisplayComponent := itemDisplayData.(Displayable)
						displayables = append(displayables, itemDisplayComponent.Sprite)
					}
				}
			}

			// display memory
			if hasMemory {
				memoryComponent := memoryData.(EntityMemory)
				col, ok := memoryComponent.Memory[itemX]
				if ok {
					items, ok := col[itemY]
					if ok {
						for _, item := range items {
							// add faded dispay to make memory look different
							displayables = append(displayables, gui.Fade(item.Sprite))
						}
					}
				}
			}
			gui.DisplaySprites(dx, dy, displayables)
		}

		bounds := getOctantBounds(displayRadius)

		displayXY(0, 0)
		for octant := 0; octant < 8; octant++ {
			for row := 1; row < displayRadius; row++ {
				for col := 0; col <= row; col++ {
					if bounds[col] == row {
						break
					}
					if octant%2 == 0 {
						if col == 0 || col == row {
							continue
						}
					}
					// in bounds, continue on
					dx, dy := transformOctant(row, col, octant)
					displayXY(dx, dy)
				}
			}
		}
	}
}

func (s *DisplayHandler) showBelowYou(m *Manager, entity Entity, displayRadius int) {
	positionData, hasPosition := m.getComponent(entity, POSITION)
	_, hasAwarness := m.getComponent(entity, ENTITY_AWARENESS)

	if hasPosition && hasAwarness {
		positionComponent := positionData.(Position)

		items := make([]Entity, 0)
		belowYou := m.getEntitiesFromPos(positionComponent.X, positionComponent.Y)
		for item := range belowYou {
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
			gui.DrawText(0, displayRadius+5, displayString)
		}
	}
}

func (s *DisplayHandler) showStats(m *Manager, entity Entity, displayRadius int) {
	informationData, hasInformation := m.getComponent(entity, INFORMATION)
	healthData, hasHealth := m.getComponent(entity, HEALTH)
	fighterData, isFighter := m.getComponent(entity, FIGHTER)

	if hasInformation {
		informationComponent := informationData.(Information)
		gui.DrawText(0, -displayRadius-4, informationComponent.Name)
	}
	if hasHealth {
		healthComponent := healthData.(Health)
		healthString := "HP " + strconv.Itoa(healthComponent.Current) + "/" + strconv.Itoa(healthComponent.Max)
		gui.DrawText(-displayRadius, -displayRadius, healthString)

		bounds := getOctantBounds(displayRadius)
		healthPercent := healthComponent.Current * 8 * len(bounds) / healthComponent.Max
		healthDisplayed := 0
		for octant := 2; octant < 6; octant++ {
			for i := range bounds {
				if healthDisplayed <= healthPercent {
					k := i
					if octant%2 == 1 {
						k = len(bounds) - k - 1
					}
					dx, dy := transformOctant(k, bounds[k], octant)

					gui.DisplaySprites(dx, -dy, []gui.Sprite{gui.GetSprite(gui.BLOOD)})
					gui.DisplaySprites(-dx, -dy, []gui.Sprite{gui.GetSprite(gui.BLOOD)})
					healthDisplayed += 2
				}
			}
		}
	}
	if isFighter {
		fighterComponent := fighterData.(Fighter)
		strengthString := "STR " + strconv.Itoa(fighterComponent.Strength)
		gui.DrawText(displayRadius, -displayRadius, strengthString)
	}
}

func (s *DisplayHandler) showInventory(m *Manager, entity Entity, displayRadius int) {
	inventoryData, hasInventory := m.getComponent(entity, INVENTORY)

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
}

func (s *DisplayHandler) showEquiped(m *Manager, entity Entity, displayRadius int) {
	fighterData, isFighter := m.getComponent(entity, FIGHTER)

	if isFighter {
		fighterComponent := fighterData.(Fighter)

		weaponInformationData, weaponInformationOk := m.getComponent(fighterComponent.Weapon, INFORMATION)
		armorInformationData, armorInformationOk := m.getComponent(fighterComponent.Armor, INFORMATION)

		if weaponInformationOk {
			informationComponent := weaponInformationData.(Information)

			gui.DrawText(-displayRadius-10, -1, "Weapon:")
			gui.DrawText(-displayRadius-10, 0, informationComponent.Name)
		}

		if armorInformationOk {
			informationComponent := armorInformationData.(Information)

			gui.DrawText(-displayRadius-10, 2, "Armor:")
			gui.DrawText(-displayRadius-10, 3, informationComponent.Name)
		}
	}
}

func (h *DisplayHandler) showMenu(m *Manager) {
	m.user.Menu.show(m)
}
