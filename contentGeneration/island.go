package contentGeneration

import (
	"math/rand"
)

func generateIslandHeightmap(width, height, maxHeight int) [][]int {

	edgeRadius := 5

	tiles := make([][]int, width)

	// first fill with random noise, weighted at edges to sink
	for x := 0; x < width; x++ {
		tiles[x] = make([]int, height)
		for y := 0; y < height; y++ {
			if x < edgeRadius || x > width-edgeRadius || y < edgeRadius || y > height-edgeRadius {
				tiles[x][y] = rand.Intn(maxHeight / 4)
			} else {
				tiles[x][y] = rand.Intn(maxHeight)
			}
		}
	}

	// then smooth
	smoothRadius := 2
	smoothTimes := 8

	for i := 0; i < smoothTimes; i++ {
		newTiles := make([][]int, width)

		for x := 0; x < width; x++ {
			newTiles[x] = make([]int, height)
			for y := 0; y < height; y++ {
				average := 0
				for dx := -smoothRadius; dx < smoothRadius; dx++ {
					for dy := -smoothRadius; dy < smoothRadius; dy++ {
						rx := x + dx
						ry := y + dy
						if rx >= 0 && rx < width && ry >= 0 && ry < height {
							average += tiles[rx][ry]
						}
					}
				}
				average /= (smoothRadius * smoothRadius * 4)
				newTiles[x][y] = average
			}
		}
		tiles = newTiles
	}

	return tiles
}
