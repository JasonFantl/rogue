package ecs

import "fmt"

type ChunkHandler struct {
}

func (h *ChunkHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == MOVED && event.entity == m.user.Controlling {

		movedEvent := event.data.(Moved)

		// moved chunks
		oldCx, oldCy := cordsToChunkCords(movedEvent.fromX, movedEvent.fromY)
		newCx, newCy := cordsToChunkCords(movedEvent.toX, movedEvent.toY)

		if newCx != oldCx || newCy != oldCy {
			fmt.Println("entered new chunk")
			// create sets of old and new chunks
			type cord struct{ x, y int }
			oldChunks := map[cord]bool{}
			newChunks := map[cord]bool{}

			for dcx := -RENDER_DISTANCE; dcx <= RENDER_DISTANCE; dcx++ {
				for dcy := -RENDER_DISTANCE; dcy <= RENDER_DISTANCE; dcy++ {
					oldChunks[cord{oldCx + dcx, oldCy + dcy}] = true
					newChunks[cord{newCx + dcx, newCy + dcy}] = true
				}
			}

			// then unload all old chunks not in new chunks
			for c := range oldChunks {
				if !newChunks[c] {
					m.entityManager.unloadChunk(c.x, c.y)
					delete(oldChunks, c)
				}
			}
			// then load all new chunks
			for c := range newChunks {
				if !oldChunks[c] {
					m.entityManager.loadChunk(c.x, c.y)
				}
			}
		}
	}

	return returnEvents
}
