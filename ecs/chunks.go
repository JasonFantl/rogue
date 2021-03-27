package ecs

import (
	"encoding/gob"
	"fmt"
	"os"
)

const (
	CHUNK_SIZE      = 8
	RENDER_DISTANCE = 2
)

type chunkData PositionLookup

func (em *EntityManager) unloadChunk(cx, cy int) {
	data := em.positionLookup.popChunk(cx, cy)

	dataFile, err := os.Create(cordsToFilename(cx, cy))

	if err != nil {
		fmt.Println(err)
		return
	}

	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(data)

	dataFile.Close()
}

func (em *EntityManager) loadChunk(cx, cy int) {
	data := chunkData{}

	dataFile, err := os.Open(cordsToFilename(cx, cy))

	if err != nil {
		fmt.Println(err)
		return
	}

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&data)

	if err != nil {
		fmt.Println(err)
		return
	}

	dataFile.Close()

	// do something with chunk data
	em.positionLookup.setChunk(cx, cy, PositionLookup(data))
}

func cordsToFilename(x, y int) string {
	return fmt.Sprintf("data/chunks/%d_%d.gob", x, y)
}

func cordsToChunkCords(x, y int) (int, int) {
	cx, cy := x/CHUNK_SIZE, y/CHUNK_SIZE
	if x < 0 {
		cx--
	}
	if y < 0 {
		cy--
	}

	return cx, cy
}
