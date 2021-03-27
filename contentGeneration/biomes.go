package contentGeneration

import (
	"math/rand"
)

type BiomeType int

const (
	BEACH BiomeType = iota
	MOUNTAIN
	FOREST
	WATER
)

// note:
func biomesMask(heightmap [][]int, biomeTypes []BiomeType) [][]BiomeType {

	// first we generate some biome points
	type biomePoint struct {
		x, y int
		t    BiomeType
	}

	width := len(heightmap)
	height := len(heightmap[0])
	biomePointCount := width * height / 100

	biomePoints := make([]biomePoint, biomePointCount)

	for i := 0; i < biomePointCount; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		t := biomeTypes[rand.Intn(len(biomeTypes))]

		biomePoints = append(biomePoints, biomePoint{x, y, t})
	}

	tiles := make([][]BiomeType, width)

	waterHeight := 45
	beachHeight := 46
	mountainHeight := 50
	// then we populate based on those points
	for x := 0; x < width; x++ {
		tiles[x] = make([]BiomeType, height)
		for y := 0; y < height; y++ {
			if heightmap[x][y] < waterHeight {
				tiles[x][y] = WATER
			} else if heightmap[x][y] < beachHeight {
				tiles[x][y] = BEACH
			} else if heightmap[x][y] > mountainHeight {
				tiles[x][y] = MOUNTAIN
			} else {
				// find closest point
				nearest := biomePoint{-1, -1, '-'}
				for _, point := range biomePoints {
					dx := nearest.x - x
					dy := nearest.y - y
					oldDist := dx*dx + dy*dy
					dx = point.x - x
					dy = point.y - y
					newDist := dx*dx + dy*dy

					if newDist < oldDist || nearest.x == -1 {
						nearest = point
					}
				}

				tiles[x][y] = nearest.t
			}
		}
	}

	return tiles
}

func betBoolMaskFromBiomeMask(biomeMask [][]BiomeType, biome BiomeType) [][]bool {
	boolBiomeMask := make([][]bool, len(biomeMask))
	for x := range boolBiomeMask {
		boolBiomeMask[x] = make([]bool, len(biomeMask[x]))
		for y := range boolBiomeMask[x] {
			boolBiomeMask[x][y] = biomeMask[x][y] == biome
		}
	}

	return boolBiomeMask
}

func andMasks(mask1, mask2 [][]bool) [][]bool {
	andMask := make([][]bool, len(mask1))
	for x := range andMask {
		andMask[x] = make([]bool, len(mask1[x]))
		for y := range andMask[x] {
			andMask[x][y] = mask1[x][y]
			if x < len(mask2) {
				if y < len(mask2[x]) {
					andMask[x][y] = mask1[x][y] && mask2[x][y]
				}
			}
		}
	}
	return andMask
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
