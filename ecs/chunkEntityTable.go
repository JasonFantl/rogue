package ecs

import (
	"encoding/gob"
	"fmt"
	"os"
)

const (
	CHUNK_SIZE      = 32
	RENDER_DISTANCE = 2
)

type chunkedEntityTable struct {
	dynamicChunks map[Position]*entityTable
	staticChunk   *entityTable
	chunkLookup   map[Entity]Position
}

func (et *chunkedEntityTable) addEntity(entity Entity, components map[ComponentID]interface{}) {

	// find the appropiate chunk
	chunk := et.staticChunk

	for componentID, componentData := range components {
		if componentID == POSITION {
			position := componentData.(Position)
			p := cordsToChunkCords(position)

			et.checkInited(p)
			chunk = et.dynamicChunks[p]

			et.chunkLookup[entity] = p

			break
		}
	}

	chunk.addEntity(entity, components)
}

func (et *chunkedEntityTable) removeEntity(entity Entity) {
	chunk := et.findChunk(entity)
	chunk.removeEntity(entity)
}

func (et *chunkedEntityTable) removeComponent(entity Entity, componentID ComponentID) {
	chunk := et.findChunk(entity)
	chunk.removeComponent(entity, componentID)
}

func (et *chunkedEntityTable) setComponent(entity Entity, componentID ComponentID, data interface{}) {

	chunk := et.findChunk(entity)

	positionData, hasPosition := et.getComponent(entity, componentID)

	if componentID == POSITION {
		newPosition := data.(Position)
		newChunkPos := cordsToChunkCords(newPosition)
		et.checkInited(newChunkPos)
		chunk = et.dynamicChunks[newChunkPos]
		et.chunkLookup[entity] = newChunkPos

		if hasPosition { // might be moving chunks
			oldPosition := positionData.(Position)
			oldChunkPos := cordsToChunkCords(oldPosition)

			if oldChunkPos.X != newChunkPos.X || oldChunkPos.Y != newChunkPos.Y {
				oldChunk := et.dynamicChunks[oldChunkPos]
				newChunk := et.dynamicChunks[newChunkPos]

				transferData := oldChunk.getComponents(entity)
				newChunk.addEntity(entity, transferData)
				oldChunk.removeEntity(entity)

			}
		}
	}

	chunk.setComponent(entity, componentID, data)
}

func (et *chunkedEntityTable) getComponent(entity Entity, componentID ComponentID) (interface{}, bool) {
	chunk := et.findChunk(entity)
	return chunk.getComponent(entity, componentID)
}

func (et *chunkedEntityTable) getComponents(entity Entity) map[ComponentID]interface{} {
	chunk := et.findChunk(entity)
	return chunk.getComponents(entity)
}

func (et *chunkedEntityTable) getEntities(componentID ComponentID) map[Entity]bool {
	entities := make(map[Entity]bool)

	for _, chunk := range et.dynamicChunks {
		chunkEntities := chunk.getEntitiesWithComponent(componentID)

		for e := range chunkEntities {
			entities[e] = true
		}
	}

	return entities
}

func (et *chunkedEntityTable) getEntitiesAtPosition(p Position) map[Entity]bool {
	chunkPos := cordsToChunkCords(p)
	et.checkInited(chunkPos)
	chunk := et.dynamicChunks[chunkPos]
	return chunk.getEntitiesAtPosition(p)
}

// how do we get chunk without gtting component?
func (et *chunkedEntityTable) findChunk(entity Entity) *entityTable {
	if et.staticChunk == nil {
		et.staticChunk = &entityTable{}
	}
	chunk := et.staticChunk

	chunkPos, hasPosition := et.chunkLookup[entity]
	if hasPosition {
		et.checkInited(chunkPos)
		chunk = et.dynamicChunks[chunkPos]
	}

	return chunk
}

func (et *chunkedEntityTable) checkInited(pos Position) {
	if et.chunkLookup == nil {
		et.chunkLookup = make(map[Entity]Position)
	}
	if et.dynamicChunks == nil {
		et.dynamicChunks = make(map[Position]*entityTable)
	}
	if et.dynamicChunks[pos] == nil {
		et.dynamicChunks[pos] = &entityTable{}
	}
}

func (et *chunkedEntityTable) unloadChunk(p Position) {
	data := et.dynamicChunks[p]

	if data == nil {
		return
	}

	dataFile, err := os.Create(cordsToFilename(p))

	if err != nil {
		fmt.Println(err)
		return
	}

	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(data)

	dataFile.Close()

	delete(et.dynamicChunks, p)
}

var staticChunkPos = Position{-999, -999}

func (et *chunkedEntityTable) unloadAllChunks() {
	for p := range et.dynamicChunks {
		et.unloadChunk(p)
	}
	et.unloadChunk(staticChunkPos)
}

func (et *chunkedEntityTable) loadChunk(p Position) {
	data := entityTable{}

	dataFile, err := os.Open(cordsToFilename(p))

	if err != nil {
		fmt.Println("error opening %d, %d", p.X, p.Y)
		fmt.Println(err)
		return
	}

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&data)

	if err != nil {
		fmt.Println("error decoding %d, %d", p.X, p.Y)
		fmt.Println(err)
		return
	}

	dataFile.Close()

	// do something with chunk data
	et.dynamicChunks[p] = &data
}

func cordsToFilename(p Position) string {
	return fmt.Sprintf("data/chunks/%d_%d.gob", p.X, p.Y)
}

func cordsToChunkCords(p Position) Position {
	cx, cy := p.X/CHUNK_SIZE, p.Y/CHUNK_SIZE
	if p.X < 0 {
		cx--
	}
	if p.Y < 0 {
		cy--
	}

	return Position{cx, cy}
}
