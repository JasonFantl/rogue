package contentGeneration

import (
	"math/rand"

	"github.com/jasonfantl/rogue/ecs"
)

// we currently assume cave opening at width/2
func addBeach(ecsManager *ecs.EntityManager, mask [][]bool, biomeMask [][]BiomeType) {

	beachBiomeMask := betBoolMaskFromBiomeMask(biomeMask, BEACH)

	//tree likelyhood
	treeChance := 100
	stoneChance := 100
	for x := 0; x < len(beachBiomeMask); x++ {
		for y := 0; y < len(beachBiomeMask[x]); y++ {
			if beachBiomeMask[x][y] {
				if rand.Intn(stoneChance) == 0 {
					addEntity(ecsManager, mask, true, true, x, y, stoneFloor(x, y))
				} else if rand.Intn(treeChance) == 0 {
					addPalmTree(ecsManager, mask, x, y)
				} else {
					addEntity(ecsManager, mask, true, true, x, y, sandFloor(x, y))
				}
			}
		}
	}
}

func addPalmTree(ecsManager *ecs.EntityManager, mask [][]bool, x, y int) {

	if !mask[x][y] {
		addEntity(ecsManager, mask, true, true, x, y, treeTrunk(x, y))

		treeRadius := rand.Intn(5) + 2
		for dx := -treeRadius; dx <= treeRadius; dx++ {
			addEntity(ecsManager, mask, true, false, x+dx, y, leaf(x+dx, y))
		}
		for dy := -treeRadius; dy <= treeRadius; dy++ {
			addEntity(ecsManager, mask, true, false, x, y+dy, leaf(x, y+dy))
		}
	}
}
