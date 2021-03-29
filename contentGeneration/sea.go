package contentGeneration

import (
	"github.com/jasonfantl/rogue/ecs"
)

// we currently assume cave opening at width/2
func addSea(ecsManager *ecs.EntityManager, mask [][]bool, biomeMask [][]BiomeType) {

	waterBiomeMask := betBoolMaskFromBiomeMask(biomeMask, WATER)

	for x := 0; x < len(waterBiomeMask); x++ {
		for y := 0; y < len(waterBiomeMask[x]); y++ {
			if waterBiomeMask[x][y] {
				addEntity(ecsManager, mask, true, true, x, y, water(x, y))
			}
		}
	}
}
