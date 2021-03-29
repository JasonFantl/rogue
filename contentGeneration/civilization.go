package contentGeneration

import (
	"math/rand"

	"github.com/jasonfantl/rogue/ecs"
)

func addTowns(ecsManager *ecs.EntityManager, mask [][]bool, heightmap [][]int, biomeMask [][]BiomeType) {

	width := len(mask)
	height := len(mask[0])

	townCount := rand.Intn(5) + 5
	for townCount > 0 {
		x, y := rand.Intn(width), rand.Intn(height)
		size := rand.Intn(10) + 1

		// check were on forest
		if biomeMask[x][y] == FOREST {
			addTown(ecsManager, mask, x, y, size)
			townCount--
		}
	}
}

// not yet implemented
func addTown(ecsManager *ecs.EntityManager, mask [][]bool, x, y, size int) {

	addHouse(ecsManager, mask, x, y)

}

// not yet implemented
// probably want variaty of houses (proc gen?)
func addHouse(ecsManager *ecs.EntityManager, mask [][]bool, x, y int) {
	width, height := rand.Intn(5)+3, rand.Intn(5)+3

	for dx := -width / 2; dx <= width/2; dx++ {
		for dy := -height / 2; dy <= height/2; dy++ {
			nx := dx + x
			ny := dy + y

			addEntity(ecsManager, mask, true, true, nx, ny, stoneFloor(nx, ny))

			if dx == -width/2 || dx == width/2 || dy == -height/2 || dy == height/2 {
				// leave space for door
				if dx == width/2 && dy == 0 {
					addDoor(ecsManager, mask, nx, ny)
				} else {
					addEntity(ecsManager, mask, true, true, nx, ny, treeTrunk(nx, ny))
				}
			}
		}
	}
	addTownsMember(ecsManager, x, y)
}

// type mapLocation struct {
// 	x, y int
// 	size int
// }

// // not yet implemented
// func generateRoads(ecsManager *ecs.EntityManager, mask [][]bool, locations []mapLocation) {

// 	// tmp
// 	for i := range locations {
// 		for j := i + 1; j < len(locations); j++ {
// 			generateRoad(ecsManager, mask, locations[i], locations[j])
// 		}
// 	}
// }

// // not yet implemented
// func generateRoad(ecsManager *ecs.EntityManager, mask [][]bool, locationA, locationB mapLocation) {

// 	// tmp
// 	if locationA.x < locationB.x {
// 		locationA, locationB = locationB, locationA
// 	}
// 	for x := locationA.x; x < locationB.x; x++ {
// 		y := locationA.y + (locationB.x-x)*(x-locationA.x)/(locationB.x-locationA.x)
// 		addEntity(ecsManager, mask, true, true, x, y, dirtFloor(x, y))
// 	}
// }
