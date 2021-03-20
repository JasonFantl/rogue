package ecs

import (
	"strconv"

	"github.com/jasonfantl/rogue/gui"
)

// Display is special, not like the other handlers
type DisplayHandler struct {
	displayRadius int // in pixles
	seeingRadius  int // in sprite blocks
}

func (h *DisplayHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == DISPLAY {
		gui.Clear()
		h.showEntity(m, m.user.Controlling)
		if m.user.Menu.active {
			h.showMenu(m)
		}
	}

	if event.ID == WAKEUP_HANDLERS {
		h.displayRadius = 100
		h.seeingRadius = 10
		gui.SpecialSetSpriteScale(h.displayRadius, h.seeingRadius)
	}

	if event.ID == SETTING_CHANGE {
		settingChange := event.data.(SettingChange)
		if settingChange.field == "zoom" {
			h.seeingRadius += settingChange.value
			if h.seeingRadius < 1 {
				h.seeingRadius = 1
			}

			gui.SpecialSetSpriteScale(h.displayRadius, h.seeingRadius)
		}
	}

	return returnEvents
}

func (h *DisplayHandler) showEntity(m *Manager, entity Entity) {

	h.showGrid(m, entity)
	h.showBelowYou(m, entity)
	h.showStats(m, entity)
	h.showEquiped(m, entity)

}

func (h *DisplayHandler) showGrid(m *Manager, entity Entity) {
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
					shouldDisplay := true
					stashableData, isStashable := m.getComponent(item, STASHABLE)
					if isStashable {
						stashableComponent := stashableData.(Stashable)
						shouldDisplay = !stashableComponent.Stashed
					}
					if shouldDisplay {
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

		bounds := getOctantBounds(h.seeingRadius)

		displayXY(0, 0)
		for octant := 0; octant < 8; octant++ {
			for row := 1; row < h.seeingRadius; row++ {
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

func (h *DisplayHandler) showBelowYou(m *Manager, entity Entity) {
	positionData, hasPosition := m.getComponent(entity, POSITION)
	_, hasAwarness := m.getComponent(entity, ENTITY_AWARENESS)

	if hasPosition && hasAwarness {
		positionComponent := positionData.(Position)

		items := make([]Entity, 0)
		belowYou := m.getEntitiesFromPos(positionComponent.X, positionComponent.Y)
		for item := range belowYou {
			if isStashableTreasure(m, item) {
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
			gui.DrawText(0, h.displayRadius*4/3, displayString)
		}
	}
}

func (h *DisplayHandler) showStats(m *Manager, entity Entity) {
	informationData, hasInformation := m.getComponent(entity, INFORMATION)
	healthData, hasHealth := m.getComponent(entity, HEALTH)
	fighterData, isFighter := m.getComponent(entity, FIGHTER)

	if hasInformation {
		informationComponent := informationData.(Information)
		gui.DrawText(0, -h.displayRadius*4/3, informationComponent.Name)
	}
	if hasHealth {
		healthComponent := healthData.(Health)
		healthString := "HP " + strconv.Itoa(healthComponent.Current) + "/" + strconv.Itoa(healthComponent.Max)
		gui.DrawText(-h.displayRadius*4/3, -h.displayRadius*4/3, healthString)

		bounds := getOctantBounds(h.seeingRadius)
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
		gui.DrawText(h.displayRadius*4/3, -h.displayRadius*4/3, strengthString)
	}
}

func (h *DisplayHandler) showEquiped(m *Manager, entity Entity) {
	fighterData, isFighter := m.getComponent(entity, FIGHTER)

	if isFighter {
		fighterComponent := fighterData.(Fighter)

		weaponInformationData, weaponInformationOk := m.getComponent(fighterComponent.Weapon, INFORMATION)
		armorInformationData, armorInformationOk := m.getComponent(fighterComponent.Armor, INFORMATION)

		if weaponInformationOk {
			informationComponent := weaponInformationData.(Information)

			gui.DrawText(-h.displayRadius*4/3, -20, "Weapon:")
			gui.DrawText(-h.displayRadius*4/3, -10, informationComponent.Name)
		}

		if armorInformationOk {
			informationComponent := armorInformationData.(Information)

			gui.DrawText(-h.displayRadius*4/3, 10, "Armor:")
			gui.DrawText(-h.displayRadius*4/3, 20, informationComponent.Name)
		}
	}
}

func (h *DisplayHandler) showMenu(m *Manager) {
	m.user.Menu.show(m)
}
