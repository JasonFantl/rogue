package ecs

// ok, this isnt great, but entities need to add a component twice
type LockHandler struct{}

func (h *LockHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == TRY_MOVE {
		moveEvent := event.data.(TryMove)

		positionData, hasPosition := m.getComponent(event.entity, POSITION)

		if hasPosition {
			positionComponent := positionData.(Position)

			newX := positionComponent.X + moveEvent.dx
			newY := positionComponent.Y + moveEvent.dy

			for lock := range m.getEntitiesFromPos(newX, newY) {
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

					if lockableComponent.Inversed {
						h.lockComponent(m, tryUnlockEvent.what)
					} else {
						h.unlockComponent(m, tryUnlockEvent.what)
					}
					returnEvents = append(returnEvents, Event{UNLOCKED, Unlocked{}, tryUnlockEvent.what})
					break
				}
			}
		}
	}

	return returnEvents
}

func (h *LockHandler) unlockComponent(m *Manager, entity Entity) (returnEvents []Event) {
	lockableData, hasLockable := m.getComponent(entity, LOCKABLE)

	if hasLockable {
		lockableComponent := lockableData.(Lockable)

		m.setComponent(entity, lockableComponent.Locking.ID, lockableComponent.Locking.Data)
	}
	return returnEvents
}

func (h *LockHandler) lockComponent(m *Manager, entity Entity) (returnEvents []Event) {
	lockableData, hasLockable := m.getComponent(entity, LOCKABLE)

	if hasLockable {
		lockableComponent := lockableData.(Lockable)

		updatedComponent, hasComponent := m.getComponent(entity, lockableComponent.Locking.ID)
		if hasComponent {
			// remember locked compoenent
			lockableComponent.Locking.Data = updatedComponent
			m.setComponent(entity, LOCKABLE, lockableComponent)

			// update manager
			m.removeComponent(entity, lockableComponent.Locking.ID)
		}
	}
	return returnEvents
}
