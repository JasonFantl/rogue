package ecs

import (
	"encoding/gob"
	"fmt"
)

type ChunkHandler struct {
}

func (h *ChunkHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == WAKEUP_HANDLERS {

		gob.Register(entityTable{})
		gob.Register(Position{})
		gob.Register(Displayable{})
		gob.Register(User{})
		gob.Register(Vision{})
		gob.Register(EntityAwarness{})
		gob.Register(EntityMemory{})
		gob.Register(Memorable{})
		gob.Register(Volume{})
		gob.Register(Inventory{})
		gob.Register(Information{})
		gob.Register(Stashable{})
		gob.Register(Consumable{})
		gob.Register(Health{})
		gob.Register(Fighter{})
		gob.Register(Brain{})
		gob.Register(Damage{})
		gob.Register(DamageResistance{})
		gob.Register(Opaque{})
		gob.Register(Reactions{})
		gob.Register(Lockable{})
		gob.Register(Projectile{})

		cx, cy := 0, 0
		positionData, hasPosition := m.getComponent(m.user.Controlling, POSITION)
		if hasPosition {
			positionComponent := positionData.(Position)
			cx, cy = positionComponent.X, positionComponent.Y
		}

		m.entityManager.entityTable.loadChunk(staticChunkPos)
		for dcx := -RENDER_DISTANCE; dcx <= RENDER_DISTANCE; dcx++ {
			for dcy := -RENDER_DISTANCE; dcy <= RENDER_DISTANCE; dcy++ {
				m.entityManager.entityTable.loadChunk(Position{cx + dcx, cy + dcy})
			}
		}
	}

	if event.ID == MOVED && event.entity == m.user.Controlling {

		movedEvent := event.data.(Moved)

		// moved chunks
		oldC := cordsToChunkCords(movedEvent.from)
		newC := cordsToChunkCords(movedEvent.to)

		if newC.X != oldC.X || newC.Y != oldC.Y {
			fmt.Println("entered new chunk")
			// create sets of old and new chunks
			oldChunks := map[Position]bool{}
			newChunks := map[Position]bool{}

			for dcx := -RENDER_DISTANCE; dcx <= RENDER_DISTANCE; dcx++ {
				for dcy := -RENDER_DISTANCE; dcy <= RENDER_DISTANCE; dcy++ {
					oldChunks[Position{oldC.X + dcx, oldC.Y + dcy}] = true
					newChunks[Position{newC.X + dcx, newC.Y + dcy}] = true
				}
			}

			// then unload all old chunks not in new chunks
			for c := range oldChunks {
				if !newChunks[c] {
					m.entityManager.entityTable.unloadChunk(c)
					delete(oldChunks, c)
				}
			}
			// then load all new chunks
			for c := range newChunks {
				if !oldChunks[c] {
					m.entityManager.entityTable.loadChunk(c)
				}
			}
		}
	}

	return returnEvents
}
