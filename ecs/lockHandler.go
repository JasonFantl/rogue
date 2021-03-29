package ecs

// ok, this isnt great, but entities need to add a component twice
type LockHandler struct{}

func (h *LockHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == WAKEUP_HANDLERS {

		entities := m.getEntities(LOCKABLE)
		for entity := range entities {
			lockableData, _ := m.getComponent(entity, LOCKABLE)
			lockableComponent := lockableData.(Lockable)

			if lockableComponent.Locked {
				for _, component := range lockableComponent.LockedComponents {
					m.setComponent(entity, component.ID, component.Data)
				}
			} else {
				for _, component := range lockableComponent.UnlockedComponents {
					m.setComponent(entity, component.ID, component.Data)
				}
			}
		}
	}

	if event.ID == TRY_MOVE {
		moveEvent := event.data.(TryMove)

		positionData, hasPosition := m.getComponent(event.entity, POSITION)

		if hasPosition {
			positionComponent := positionData.(Position)

			newPos := Position{positionComponent.X + moveEvent.dx, positionComponent.Y + moveEvent.dy}

			for lock := range m.getEntitiesAtPosition(newPos) {
				lockableData, hasLockable := m.getComponent(lock, LOCKABLE)
				if hasLockable {
					lockableComponent := lockableData.(Lockable)

					if lockableComponent.Locked {
						returnEvents = append(returnEvents, Event{TRY_UNLOCK, TryUnlock{lock}, event.entity})
					}
				}
			}
		}
	}

	if event.ID == TRY_UNLOCK {
		tryUnlockEvent := event.data.(TryUnlock)

		inventoryData, hasInventory := m.getComponent(event.entity, INVENTORY)
		lockableData, hasLockable := m.getComponent(tryUnlockEvent.what, LOCKABLE)

		if hasInventory && hasLockable {
			inventoryComponent := inventoryData.(Inventory)
			lockableComponent := lockableData.(Lockable)

			// should we check again that its not already unlocked? should only get event if its already unlocked
			for item := range inventoryComponent.Items {
				if item == lockableComponent.Key {
					lockableComponent.Locked = false
					m.setComponent(tryUnlockEvent.what, LOCKABLE, lockableComponent)
					h.unlockComponents(m, tryUnlockEvent.what)

					returnEvents = append(returnEvents, Event{UNLOCKED, Unlocked{}, tryUnlockEvent.what})
					break
				}
			}
		}
	}

	return returnEvents
}

func (h *LockHandler) unlockComponents(m *Manager, entity Entity) (returnEvents []Event) {
	lockableData, hasLockable := m.getComponent(entity, LOCKABLE)

	if hasLockable {
		lockableComponent := lockableData.(Lockable)

		h.swapComponents(m, entity, lockableComponent.UnlockedComponents, lockableComponent.LockedComponents)
		m.setComponent(entity, LOCKABLE, lockableComponent)
	}
	return returnEvents
}

func (h *LockHandler) swapComponents(m *Manager, entity Entity, inComponents, outComponents []Component) {

	// first remember and remove old components
	for i, component := range outComponents {
		componentData, hasComponent := m.getComponent(entity, component.ID)
		if hasComponent {
			// remember old components
			outComponents[i].Data = componentData
			// remove
			m.removeComponent(entity, component.ID)
		}
	}

	for _, component := range inComponents {
		m.setComponent(entity, component.ID, component.Data)
	}
}
