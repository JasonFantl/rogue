package contentGeneration

import (
	"math/rand"

	"github.com/jasonfantl/rogue/ecs"
)

func addCaves(ecsManager *ecs.Manager, mask [][]bool) {

	caveMask := generateCaveMask(mask)

	width := len(mask)
	height := len(mask[0])
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if mask[x][y] {
				if !caveMask[x][y] {
					addEntity(ecsManager, mask, false, true, x, y, stoneWall(x, y))
				} else {
					addEntity(ecsManager, mask, false, true, x, y, stoneFloor(x, y))
				}
			}
		}
	}

	// then add cave entities
	itemCount := width + height
	for itemCount > 0 {
		x := rand.Intn(width)
		y := rand.Intn(height)
		if caveMask[x][y] {
			r := rand.Intn(5)
			if r == 0 {
				addTreasure(ecsManager, x, y)
			} else if r == 1 {
				addMonster(ecsManager, x, y)
			} else if r == 2 {
				addPotion(ecsManager, x, y)
			} else if r == 3 {
				addWeapon(ecsManager, x, y)
			} else {
				addArmor(ecsManager, x, y)
			}
			itemCount--
		}
	}
}

func generateCaveMask(mask [][]bool) [][]bool {

	deathLimit := 4
	birthLimit := 4
	iterationCount := 8

	width := len(mask)
	height := len(mask[0])

	// simple cellular automota implementation
	tiles := make([][]bool, width)

	for x := 0; x < width; x++ {
		tiles[x] = make([]bool, height)
		for y := 0; y < height; y++ {
			tiles[x][y] = (rand.Intn(2) == 1)
		}
	}

	// second: run CA a couple of times
	for step := 0; step < iterationCount; step++ {
		newTiles := make([][]bool, width)
		for x := 0; x < width; x++ {
			newTiles[x] = make([]bool, height)
			for y := 0; y < height; y++ {
				// count neighbors
				nCount := 0
				for dx := -1; dx < 2; dx++ {
					for dy := -1; dy < 2; dy++ {
						testX := x + dx
						testY := y + dy
						if testX < 0 || testX >= width || testY < 0 || testY >= height {
							nCount++
						} else if testX != x || testY != 0 {
							if tiles[testX][testY] {
								nCount++
							}
						}
					}
				}

				// run rules
				if tiles[x][y] && nCount > deathLimit {
					newTiles[x][y] = true
				} else if !tiles[x][y] && nCount > birthLimit {
					newTiles[x][y] = true
				}
			}
		}
		tiles = newTiles
	}

	return andMasks(mask, tiles)
}
