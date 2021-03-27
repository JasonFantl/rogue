package contentGeneration

import (
	"math/rand"

	"github.com/jasonfantl/rogue/ecs"
)

// we currently assume cave opening at width/2
func addForest(ecsManager *ecs.Manager, mask [][]bool, biomeMask [][]BiomeType) {

	forestBiomeMask := betBoolMaskFromBiomeMask(biomeMask, FOREST)

	//tree likelyhood
	treeChance := 50
	stoneChance := 100
	for x := 0; x < len(forestBiomeMask); x++ {
		for y := 0; y < len(forestBiomeMask[x]); y++ {
			if forestBiomeMask[x][y] {
				if rand.Intn(stoneChance) == 0 {
					addEntity(ecsManager, mask, false, true, x, y, stoneFloor(x, y))
					continue
				} else if rand.Intn(treeChance) == 0 {
					addTree(ecsManager, mask, x, y)
					continue
				} else {
					addEntity(ecsManager, mask, false, true, x, y, grassFloor(x, y))
				}
			}
		}
	}
}

func addTree(ecsManager *ecs.Manager, mask [][]bool, x, y int) {

	if !mask[x][y] {
		addEntity(ecsManager, mask, true, true, x, y, treeTrunk(x, y))

		treeRadius := rand.Intn(5) + 2
		for dx := -treeRadius; dx <= treeRadius; dx++ {
			uy := treeRadius - Abs(dx)
			for dy := -uy; dy <= uy; dy++ {
				addEntity(ecsManager, mask, true, false, x+dx, y+dy, leaf(x+dx, y+dy))
			}
		}
	}
}
