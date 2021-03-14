package ecs

import "sort"

type Menu struct {
	active  bool
	cursorY int
}

func (menu *Menu) curserUp(m *Manager) {
	menu.cursorY--
}

func (menu *Menu) curserDown(m *Manager) {
	menu.cursorY++
}

func (menu *Menu) getSelected(m *Manager) Entity {
	inventoryData, hasInventory := m.getComponent(m.user.Controlling, INVENTORY)

	if hasInventory {
		inventoryComponent := inventoryData.(Inventory)

		keys := make([]int, 0)
		for k := range inventoryComponent.Items {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)

		selectedLine := menu.cursorY % len(keys)
		if selectedLine < 0 {
			selectedLine += len(keys)
		}

		if len(keys) > 0 {
			return Entity(keys[selectedLine])
		}
	}
	return Entity(0)
}

func (menu *Menu) curserSelect(m *Manager) (returnEvents []Event) {
	selected := menu.getSelected(m)
	_, isConsumable := m.getComponent(selected, CONSUMABLE)
	if isConsumable {
		returnEvents = append(returnEvents, Event{TRY_CONSUME, TryConsume{selected}, m.user.Controlling})
	}

	return returnEvents
}

func (menu *Menu) curserReset() {
	menu.cursorY = 0
}
