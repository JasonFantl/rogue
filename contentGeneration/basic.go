package contentGeneration

import (
	"math/rand"

	"github.com/jasonfantl/rogue/ecs"
)

// we currently assume cave opening at width/2
func generateBasic(ecsManager *ecs.Manager, mask [][]bool, width, height int) {

	// generate path to cave
	basePathWidth := 3

	townHeight := height / 5

	pathX := width / 2
	pathY := height - 1

	for pathY >= 0 {

		pathWidth := basePathWidth
		if pathY > height-10 {
			pathWidth += pathY - (height - 10)
		}
		for dx := -pathWidth - rand.Intn(2); dx < pathWidth+rand.Intn(2); dx++ {
			x := pathX + dx
			y := pathY
			// randomly place stone in path, more the closer we are to the cave
			if rand.Intn((height-pathY)/4+1) == 0 {
				ecsManager.AddEntity(stoneFloor(x, y))
			} else {
				ecsManager.AddEntity(dirtFloor(x, y))
			}

			mask[x][y] = false
		}

		if pathY == height/2 {
			generateTown(ecsManager, mask, pathWidth, townHeight, pathX, pathY)
		}

		// lets not change dx in town
		if pathY < height/2 || pathY > height/2+townHeight {
			pathX += rand.Intn(3) - 1
		}
		pathY--
	}

	//tree likelyhood
	treeChance := 50

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			// check no path here
			if !mask[x][y] {
				continue
			}

			if rand.Intn(height-y) == 0 {
				ecsManager.AddEntity(stoneFloor(x, y))
			} else if rand.Intn(treeChance) == 0 {
				addTree(ecsManager, mask, x, y)
			} else {
				ecsManager.AddEntity(grassFloor(x, y))
			}

			mask[x][y] = false
		}
	}
}

// assumes path runs down along y axis
func generateTown(ecsManager *ecs.Manager, mask [][]bool, pathDistance, height, xOff, yOff int) {
	houseSpacing := 5

	for houseSide := -1; houseSide <= 1; houseSide += 2 {

		for houseY := 0; houseY < height; houseY += houseSpacing + rand.Intn(4) {
			houseWidth := rand.Intn(5) + 4
			houseHeight := rand.Intn(5) + 4

			houseX := (rand.Intn(3) + pathDistance + houseWidth) * houseSide

			generateHouse(ecsManager, mask, houseWidth, houseHeight, xOff+houseX, yOff+houseY, houseSide)
			houseY += houseHeight
		}
	}
}

// assumes path runs down along y axis
func generateHouse(ecsManager *ecs.Manager, mask [][]bool, width, height, xOff, yOff int, direction int) {
	for dx := -width / 2; dx <= width/2; dx++ {
		for dy := -height / 2; dy <= height/2; dy++ {
			x := dx + xOff
			y := dy + yOff

			mask[x][y] = false
			ecsManager.AddEntity(stoneFloor(x, y))

			if dx == -width/2 || dx == width/2 || dy == -height/2 || dy == height/2 {
				// leave space for door
				if dx == -direction*width/2 && dy == 0 {
					addDoor(ecsManager, mask, x, y)
				} else {
					ecsManager.AddEntity(treeTrunk(x, y))
				}
			}
		}
	}
	addTownsMember(ecsManager, xOff, yOff)
}
