package ecs

import (
	"sort"

	"github.com/jasonfantl/rogue/gui"
)

type MenuState uint

const (
	SHOWING_IVENTORY MenuState = iota
	SHOWING_INSPECTION
	SHOWING_PROJECTILE
	SHOWING_SETTINGS
)

type Menu struct {
	active           bool
	state            MenuState
	cursorX, cursorY int
	rememberedItem   Entity
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
}

func (menu *Menu) show(m *Manager) {
	switch menu.state {
	case SHOWING_IVENTORY:
		menu.showInventory(m)
	case SHOWING_PROJECTILE:
		menu.showProjectile(m)
	case SHOWING_SETTINGS:
		menu.showSettings(m)
	case SHOWING_INSPECTION:
		menu.showInspect(m)
	}
}

func (menu *Menu) moveCurser(m *Manager, dx, dy int) {
	menu.cursorX += dx
	menu.cursorY += dy
}

func (menu *Menu) selectAtCurser(m *Manager) (returnEvents []Event) {
	switch menu.state {
	case SHOWING_IVENTORY:
		returnEvents = append(returnEvents, menu.selectInventory(m)...)
	case SHOWING_INSPECTION:
		returnEvents = append(returnEvents, menu.selectInspect(m)...)
	case SHOWING_PROJECTILE:
		returnEvents = append(returnEvents, menu.selectProjectile(m)...)
	case SHOWING_SETTINGS:
		returnEvents = append(returnEvents, menu.selectSettings(m)...)
	}
	return returnEvents
}

func (menu *Menu) selectProjectile(m *Manager) (returnEvents []Event) {

	_, isProjectile := m.getComponent(menu.rememberedItem, PROJECTILE)
	if isProjectile {
		returnEvents = append(returnEvents, Event{TIMESTEP, TimeStep{}, m.user.Controlling})
		returnEvents = append(returnEvents, Event{TRY_LAUNCH, TryLaunch{menu.rememberedItem, menu.cursorX, menu.cursorY}, m.user.Controlling})
		menu.close(m)
	}
	return returnEvents
}

var menuSettings = map[string][]string{
	"Zoom": {"In", "Out"},
}

func (menu *Menu) selectSettings(m *Manager) (returnEvents []Event) {

	selectedSetting, selectedSettingValue := menu.getSelectedSetting()
	if selectedSetting == "menu switch" {
		menu.state = SHOWING_IVENTORY
	} else {
		switch selectedSetting {
		case "Zoom":
			switch selectedSettingValue {
			case "In":
				returnEvents = append(returnEvents, Event{SETTING_CHANGE, SettingChange{"zoom", -1}, m.user.Controlling})
			case "Out":
				returnEvents = append(returnEvents, Event{SETTING_CHANGE, SettingChange{"zoom", 1}, m.user.Controlling})
			}
		}
	}
	return returnEvents
}

func (menu *Menu) selectInventory(m *Manager) (returnEvents []Event) {
	selectedInventoryItem, selectedInventoryItemAction := menu.getSelectedInventoryItem(m)

	if selectedInventoryItem == 0 { // selected switch menu
		menu.state = SHOWING_SETTINGS
	} else {
		menu.rememberedItem = selectedInventoryItem

		switch selectedInventoryItemAction {
		case STASHABLE:
			returnEvents = append(returnEvents, Event{TRY_DROP, TryDrop{selectedInventoryItem}, m.user.Controlling})
		case CONSUMABLE:
			returnEvents = append(returnEvents, Event{TRY_CONSUME, TryConsume{selectedInventoryItem}, m.user.Controlling})
		case DAMAGE:
			returnEvents = append(returnEvents, Event{TRY_EQUIP_WEAPON, TryEquip{selectedInventoryItem}, m.user.Controlling})
		case DAMAGE_RESISTANCE:
			returnEvents = append(returnEvents, Event{TRY_EQUIP_ARMOR, TryEquip{selectedInventoryItem}, m.user.Controlling})
		case PROJECTILE:
			menu.cursorX, menu.cursorY = 0, 0
			menu.state = SHOWING_PROJECTILE
		case INFORMATION:
			menu.state = SHOWING_INSPECTION
		}
	}

	return returnEvents
}

func (menu *Menu) selectInspect(m *Manager) (returnEvents []Event) {
	menu.state = SHOWING_IVENTORY

	return returnEvents
}

func (menu *Menu) showProjectile(m *Manager) {
	gui.DisplaySprite(menu.cursorX, menu.cursorY, gui.GetSprite(gui.CURSER))
}

func (menu *Menu) showSettings(m *Manager) {
	menuText := "Settings: "

	selectedSetting, selectedSettingOption := menu.getSelectedSetting()
	if selectedSetting == "menu switch" {
		menuText += "switch to inventory"
	}

	for setting, _ := range menuSettings {
		settingText := setting
		if setting == selectedSetting {
			settingText = "\n- " + setting + " <- " + selectedSettingOption + " ->"
		} else {
			settingText = "\n  " + setting
		}
		menuText += settingText
	}

	gui.DrawTextUncentered(0, 0, menuText)
}

func (menu *Menu) getSelectedSetting() (string, string) {
	selectedLineY := menu.cursorY % (len(menuSettings) + 1)
	if selectedLineY < 0 {
		selectedLineY += len(menuSettings) + 1
	}

	if selectedLineY == 0 {
		return "menu switch", ""
	}
	selectedLineY--

	if len(menuSettings) > 0 {
		keys := make([]string, 0)
		for k := range menuSettings {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		setting := keys[selectedLineY]
		selectedLineX := menu.cursorX % len(menuSettings[setting])
		if selectedLineX < 0 {
			selectedLineX += len(menuSettings[setting])
		}
		option := menuSettings[setting][selectedLineX]

		return setting, option
	}
	return "menu switch", ""
}

func (menu *Menu) showInventory(m *Manager) {

	inventoryData, hasInventory := m.getComponent(m.user.Controlling, INVENTORY)

	selectedInventoryItem, selectedInventoryItemAction := menu.getSelectedInventoryItem(m)

	if hasInventory {
		inventoryComponent := inventoryData.(Inventory)

		keys := make([]int, 0)
		for k := range inventoryComponent.Items {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)

		inventoryText := "Inventory: "

		if selectedInventoryItem == 0 { //selecting inventory
			inventoryText += "switch to settings"
		}

		for _, key := range keys {
			item := Entity(key)

			informationString := "?"
			informationData, informationOk := m.getComponent(item, INFORMATION)
			if informationOk {
				informationComponent := informationData.(Information)
				informationString = informationComponent.Name
			}

			if item == selectedInventoryItem {
				informationString += ": <- "
				switch selectedInventoryItemAction {
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
				case INFORMATION:
					informationString += "inspect"
				}
				informationString += " ->"

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

func (menu *Menu) getSelectedInventoryItem(m *Manager) (Entity, ComponentID) {
	inventoryData, hasInventory := m.getComponent(m.user.Controlling, INVENTORY)

	if hasInventory {
		inventoryComponent := inventoryData.(Inventory)

		keys := make([]int, 0)
		for k := range inventoryComponent.Items {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)

		if len(keys) > 0 {
			selectedLine := menu.cursorY % (len(keys) + 1)
			if selectedLine < 0 {
				selectedLine += len(keys) + 1
			}

			if selectedLine != 0 {
				selectedLine--
				// select item
				selectedInventoryItem := Entity(keys[selectedLine])
				selectedInventoryItemAction := STASHABLE //defaults to dropping

				// select action
				actions := make([]ComponentID, 0)
				_, isPickupable := m.getComponent(selectedInventoryItem, STASHABLE)
				_, isWeapon := m.getComponent(selectedInventoryItem, DAMAGE)
				_, isArmor := m.getComponent(selectedInventoryItem, DAMAGE_RESISTANCE)
				_, isConsumable := m.getComponent(selectedInventoryItem, CONSUMABLE)
				_, isProjectile := m.getComponent(selectedInventoryItem, PROJECTILE)
				_, hasInformation := m.getComponent(selectedInventoryItem, INFORMATION)

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
				if hasInformation {
					actions = append(actions, INFORMATION)
				}

				if len(actions) > 0 {
					selectedLine = menu.cursorX % len(actions)
					if selectedLine < 0 {
						selectedLine += len(actions)
					}
					selectedInventoryItemAction = actions[selectedLine]
				}

				return selectedInventoryItem, selectedInventoryItemAction
			}
		}
	}
	return Entity(0), 0
}

func (menu *Menu) showInspect(m *Manager) {

	// main info
	name := "?"
	details := "no information"
	informationData, informationOk := m.getComponent(menu.rememberedItem, INFORMATION)
	if informationOk {
		informationComponent := informationData.(Information)
		name = informationComponent.Name
		details = informationComponent.Details
	}
	gui.DrawText(0, 0, name)
	gui.DrawText(0, 15, details)

	// what it can do

	_, isPickupable := m.getComponent(menu.rememberedItem, STASHABLE)
	_, isWeapon := m.getComponent(menu.rememberedItem, DAMAGE)
	_, isArmor := m.getComponent(menu.rememberedItem, DAMAGE_RESISTANCE)
	_, isConsumable := m.getComponent(menu.rememberedItem, CONSUMABLE)
	_, isProjectile := m.getComponent(menu.rememberedItem, PROJECTILE)

	lineY := 30
	gui.DrawText(0, lineY, "--------------")
	lineY += 10
	if isConsumable {
		gui.DrawText(0, lineY, "consumable")
		lineY += 10
	}
	if isWeapon {
		gui.DrawText(0, lineY, "is weapon")
		lineY += 10
	}
	if isArmor {
		gui.DrawText(0, lineY, "is armor")
		lineY += 10
	}
	if isPickupable {
		gui.DrawText(0, lineY, "pickupable")
		lineY += 10
	}
	if isProjectile {
		gui.DrawText(0, lineY, "throwable")
		lineY += 10
	}

	// an image
	displayData, isDisplayable := m.getComponent(menu.rememberedItem, DISPLAYABLE)
	if isDisplayable {
		displayComponent := displayData.(Displayable)
		gui.RawDisplaySprite(0, -60, 7.0, gui.GetSprite(gui.LEAF))
		gui.RawDisplaySprite(0, -60, 6.0, displayComponent.Sprite)
	}
}
