package ecs

import (
	"sort"

	"github.com/jasonfantl/rogue/gui"
)

type MenuState uint

const (
	SHOWING_IVENTORY MenuState = iota
	SHOWING_PROJECTILE
)

type Menu struct {
	active                      bool
	state                       MenuState
	cursorX, cursorY            int
	selectedInventoryItem       Entity
	selectedInventoryItemAction ComponentID
}

func (menu *Menu) moveCurser(m *Manager, dx, dy int) {
	menu.cursorX += dx
	menu.cursorY += dy

	if menu.state == SHOWING_IVENTORY {
		menu.updateSelectedInventoryItem(m)
	}
}

func (menu *Menu) selectAtCurser(m *Manager) (returnEvents []Event) {
	switch menu.state {
	case SHOWING_IVENTORY:
		returnEvents = append(returnEvents, menu.selectInventory(m)...)
	case SHOWING_PROJECTILE:
		returnEvents = append(returnEvents, menu.selectProjectile(m)...)
	}

	return returnEvents
}

func (menu *Menu) close(m *Manager) {
	menu.active = false
	menu.reset(m)
}

func (menu *Menu) open(m *Manager) {
	menu.active = true
	menu.reset(m)
}

func (menu *Menu) reset(m *Manager) {
	menu.cursorX, menu.cursorY = 0, 0
	menu.state = SHOWING_IVENTORY
	menu.updateSelectedInventoryItem(m)
}

func (menu *Menu) selectProjectile(m *Manager) (returnEvents []Event) {

	_, isProjectile := m.getComponent(menu.selectedInventoryItem, PROJECTILE)
	if isProjectile {
		returnEvents = append(returnEvents, Event{TIMESTEP, TimeStep{}, m.user.Controlling})
		returnEvents = append(returnEvents, Event{TRY_LAUNCH, TryLaunch{menu.selectedInventoryItem, menu.cursorX, menu.cursorY}, m.user.Controlling})
		menu.close(m)
	}
	return returnEvents
}

func (menu *Menu) selectInventory(m *Manager) (returnEvents []Event) {
	menu.updateSelectedInventoryItem(m)

	switch menu.selectedInventoryItemAction {
	case STASHABLE:
		returnEvents = append(returnEvents, Event{TRY_DROP, TryDrop{menu.selectedInventoryItem}, m.user.Controlling})
	case CONSUMABLE:
		returnEvents = append(returnEvents, Event{TRY_CONSUME, TryConsume{menu.selectedInventoryItem}, m.user.Controlling})
	case DAMAGE:
		returnEvents = append(returnEvents, Event{TRY_EQUIP_WEAPON, TryEquip{menu.selectedInventoryItem}, m.user.Controlling})
	case DAMAGE_RESISTANCE:
		returnEvents = append(returnEvents, Event{TRY_EQUIP_ARMOR, TryEquip{menu.selectedInventoryItem}, m.user.Controlling})
	case PROJECTILE:
		menu.cursorX, menu.cursorY = 0, 0
		menu.state = SHOWING_PROJECTILE
	}

	return returnEvents
}

func (menu *Menu) show(m *Manager) {
	switch menu.state {
	case SHOWING_IVENTORY:
		menu.showInventory(m)
	case SHOWING_PROJECTILE:
		menu.showProjectile(m)
	}
}

func (menu *Menu) showProjectile(m *Manager) {
	gui.DisplaySprite(menu.cursorX, menu.cursorY, gui.GetSprite(gui.CURSER))
}

func (menu *Menu) showInventory(m *Manager) {
	menu.updateSelectedInventoryItem(m)

	inventoryData, hasInventory := m.getComponent(m.user.Controlling, INVENTORY)

	if hasInventory {
		inventoryComponent := inventoryData.(Inventory)

		keys := make([]int, 0)
		for k := range inventoryComponent.Items {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)

		inventoryText := "Inventory: "

		for _, key := range keys {
			item := Entity(key)

			informationString := "? : no information on item"
			informationData, informationOk := m.getComponent(item, INFORMATION)
			if informationOk {
				informationComponent := informationData.(Information)
				informationString = informationComponent.Name
			}

			if item == menu.selectedInventoryItem {
				informationString += " <- (e to "
				switch menu.selectedInventoryItemAction {
				case STASHABLE:
					informationString += "drop"
				case CONSUMABLE:
					informationString += "consume"
				case DAMAGE:
					informationString += "equip as weapon"
				case DAMAGE_RESISTANCE:
					informationString += "equip as armor"
				case PROJECTILE:
					informationString += "throw"
				}
				informationString += ") ->"

				// display info
				if informationOk {
					informationComponent := informationData.(Information)
					informationString += "\n    " + informationComponent.Details
				}
				informationString = "\n- " + informationString
			} else {
				informationString = "\n  " + informationString
			}
			inventoryText += informationString
		}

		if len(keys) > 0 {
			selectedLine := menu.cursorY % len(keys)
			if selectedLine < 0 {
				selectedLine += len(keys)
			}
		}
		gui.DrawTextUncentered(0, 0, inventoryText)
	}
}

func (menu *Menu) updateSelectedInventoryItem(m *Manager) {
	inventoryData, hasInventory := m.getComponent(m.user.Controlling, INVENTORY)

	menu.selectedInventoryItem = Entity(0)

	if hasInventory {
		inventoryComponent := inventoryData.(Inventory)

		keys := make([]int, 0)
		for k := range inventoryComponent.Items {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)

		if len(keys) > 0 {
			selectedLine := menu.cursorY % len(keys)
			if selectedLine < 0 {
				selectedLine += len(keys)
			}

			// select item
			menu.selectedInventoryItem = Entity(keys[selectedLine])

			// select action
			actions := make([]ComponentID, 0)
			_, isPickupable := m.getComponent(menu.selectedInventoryItem, STASHABLE)
			_, isWeapon := m.getComponent(menu.selectedInventoryItem, DAMAGE)
			_, isArmor := m.getComponent(menu.selectedInventoryItem, DAMAGE_RESISTANCE)
			_, isConsumable := m.getComponent(menu.selectedInventoryItem, CONSUMABLE)
			_, isProjectile := m.getComponent(menu.selectedInventoryItem, PROJECTILE)

			if isConsumable {
				actions = append(actions, CONSUMABLE)
			}
			if isWeapon {
				actions = append(actions, DAMAGE)
			}
			if isArmor {
				actions = append(actions, DAMAGE_RESISTANCE)
			}
			if isPickupable {
				actions = append(actions, STASHABLE)
			}
			if isProjectile {
				actions = append(actions, PROJECTILE)
			}

			if len(actions) > 0 {
				selectedLine = menu.cursorX % len(actions)
				if selectedLine < 0 {
					selectedLine += len(actions)
				}
				menu.selectedInventoryItemAction = actions[selectedLine]
			} else {
				// default is to drop
				menu.selectedInventoryItemAction = STASHABLE
			}
		}
	}
}
