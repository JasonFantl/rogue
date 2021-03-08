package ecs

import (
	"sort"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jasonfantl/rogue/gui"
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

	visionData, hasVision := m.getComponent(entity, VISION)
	inventoryData, hasInventory := m.getComponent(entity, INVENTORY)

	informationData, hasInformation := m.getComponent(entity, INFORMATION)
	healthData, hasHealth := m.getComponent(entity, HEALTH)

	fighterData, isFighter := m.getComponent(entity, FIGHTER)

	// display is a circle
	displayRadius := 20

	if hasVision {
		visionComponent := visionData.(Vision)
		displayRadius = visionComponent.Radius
	}

	// ------- position related --------
	if hasPosition {
		positionComponent := positionData.(Position)

		// ---------- MAP ---------

		if hasAwarness {
			awarnessComponent := awarnessData.(EntityAwarness)

			displayXY := func(dx, dy int) {
				itemX := positionComponent.X + dx
				itemY := positionComponent.Y + dy

				// display items we are aware of
				items := awarnessComponent.AwareOf.get(itemX, itemY)

				displayables := make([]gui.Sprite, 0)
				for _, item := range items {
					itemDisplayData, itemHasDisplay := m.getComponent(item, DISPLAYABLE)
					if itemHasDisplay {
						itemDisplayComponent := itemDisplayData.(Displayable)
						displayables = append(displayables, itemDisplayComponent.Sprite)
					}
				}

				// display memory
				if hasMemory {
					memoryComponent := memoryData.(EntityMemory)
					col, ok := memoryComponent.Memory[itemX]
					if ok {
						toDisplay, ok := col[itemY]
						if ok {
							// add faded dispay to make memory different
							// displayables = append(displayables, toDisplay.Sprite)
							ob := ebiten.DrawImageOptions{}
							ob.ColorM.Scale(1, 1, 1, 0.5)
							gui.DisplayXY(dx, dy, toDisplay.Sprite.Image, &ob)
						}
					}
				}
				gui.DisplaySprites(dx, dy, displayables)
			}

			bounds := getOctantBounds(displayRadius)

			for row, col := range bounds {
				for dx := -col; dx < col; dx++ {
					displayXY(dx, row)
					displayXY(dx, -row)
				}
				if row == len(bounds)-1 || bounds[row+1] != col {
					for dx := -row; dx < row; dx++ {
						displayXY(dx, col)
						displayXY(dx, -col)
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
			gui.DrawText(-len(displayString)/2, displayRadius+5, displayString)
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

		bounds := getOctantBounds(displayRadius)
		healthRadians := healthComponent.Current * 8 * len(bounds) / healthComponent.Max
		healthDisplayed := 0
		for octant := 0; octant < 8; octant++ {
			for i := range bounds {
				if healthDisplayed <= healthRadians {
					k := i
					if octant%2 == 1 {
						k = len(bounds) - k - 1
					}
					dx, dy := transformOctant(k, bounds[k], octant)

					gui.DisplaySprites(dx, dy, []gui.Sprite{gui.GetSprite(gui.BLOOD)})
					healthDisplayed++
				}
			}
		}
	}
	if isFighter {
		fighterComponent := fighterData.(Fighter)
		strengthString := "STR " + strconv.Itoa(fighterComponent.Strength)
		gui.DrawText(displayRadius/2+len(strengthString)/2, -displayRadius, strengthString)
	}

	// -------------- INVENTORY ---------------

	if false && hasInventory {
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

		weaponInformationData, weaponInformationOk := m.getComponent(fighterComponent.Weapon, INFORMATION)
		armorInformationData, armorInformationOk := m.getComponent(fighterComponent.Armor, INFORMATION)

		if weaponInformationOk {
			informationComponent := weaponInformationData.(Information)

			nameOff := len(informationComponent.Name)
			gui.DrawText(-displayRadius-10-len("Weapon:")/2, -1, "Weapon:")
			gui.DrawText(-displayRadius-10-nameOff/2, 0, informationComponent.Name)
		}

		if armorInformationOk {
			informationComponent := armorInformationData.(Information)

			nameOff := len(informationComponent.Name)
			gui.DrawText(-displayRadius-10-len("Armor:")/2, 2, "Armor:")
			gui.DrawText(-displayRadius-10-nameOff/2, 3, informationComponent.Name)
		}
	}

}
