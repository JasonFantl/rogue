package main

import (
	"math"
	"math/rand"

	"github.com/jasonfantl/rogue/ecs"
	"github.com/jasonfantl/rogue/gui"
)

func generateCaveMask(ecsManager *ecs.Manager, width, height int) [][]bool {

	deathLimit := 4
	birthLimit := 4
	iterationCount := 8

	// simple cellular automota implementation
	tiles := make([][]bool, width)

	// first: generate static noise with solid sides
	sideThickness := 2
	// openingSize := 30
	for x := 0; x < width; x++ {
		tiles[x] = make([]bool, height)
		for y := 0; y < height; y++ {
			if x < sideThickness-1 || x > width-sideThickness || y < sideThickness-1 || y > height-sideThickness {
				tiles[x][y] = false
				continue
			}
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

	return tiles
}

// we currently assume cave opening at width/2
func generateForest(ecsManager *ecs.Manager, width, height, xOff, yOff int) {

	placed := make(map[int]map[int]bool, 0)
	// generate path to cave
	basePathWidth := 3

	townHeight := 30

	pathX := width / 2
	pathY := height

	for pathY >= 0 {

		pathWidth := basePathWidth
		if pathY > height-10 {
			pathWidth += pathY - (height - 10)
		}
		for dx := -pathWidth - rand.Intn(2); dx < pathWidth+rand.Intn(2); dx++ {
			x := pathX + dx + xOff
			y := pathY + yOff
			// randomly place stone in path, more the closer we are to the cave
			if rand.Intn((height-pathY)/4+1) == 0 {
				ecsManager.AddEntity(stoneFloor(x, y))
			} else {
				ecsManager.AddEntity(dirtFloor(x, y))
			}

			_, inited := placed[x]
			if !inited {
				placed[x] = map[int]bool{}
			}

			placed[x][y] = true
		}

		// maybe generate house (town later)
		if pathY == height/2 {
			generateTown(ecsManager, placed, pathWidth, townHeight, xOff+pathX, yOff+pathY)
		}

		// lets not change dx in town
		if pathY < height/2 || pathY > height/2+townHeight {
			pathX += rand.Intn(3) - 1
		}
		pathY--
	}

	//tree likelyhood
	treeChance := 50

	for dx := 0; dx < width; dx++ {
		for dy := 0; dy < height; dy++ {
			x := xOff + dx
			y := yOff + dy

			// check no path here
			col, placed := placed[x]
			if placed {
				_, placed := col[y]
				if placed {
					continue
				}
			}

			if rand.Intn(height-dy) == 0 {
				ecsManager.AddEntity(stoneFloor(x, y))
				continue
			} else if rand.Intn(treeChance) == 0 {
				addTree(ecsManager, x, y)
				continue
			} else {
				ecsManager.AddEntity(grassFloor(x, y))
			}
		}
	}
}

// assumes path runs down along y axis
func generateTown(ecsManager *ecs.Manager, placed map[int]map[int]bool, pathDistance, height, xOff, yOff int) {
	houseSpacing := 5

	for houseSide := -1; houseSide <= 1; houseSide += 2 {

		for houseY := 0; houseY < height; houseY += houseSpacing + rand.Intn(4) {
			houseWidth := rand.Intn(5) + 4
			houseHeight := rand.Intn(5) + 4

			houseX := (rand.Intn(3) + pathDistance + houseWidth) * houseSide

			generateHouse(ecsManager, placed, houseWidth, houseHeight, xOff+houseX, yOff+houseY, houseSide)
			houseY += houseHeight
		}
	}
}

// assumes path runs down along y axis
func generateHouse(ecsManager *ecs.Manager, placed map[int]map[int]bool, width, height, xOff, yOff int, direction int) {
	for dx := -width / 2; dx <= width/2; dx++ {
		for dy := -height / 2; dy <= height/2; dy++ {
			x := dx + xOff
			y := dy + yOff

			_, inited := placed[x]
			if !inited {
				placed[x] = map[int]bool{}
			}
			placed[x][y] = true
			ecsManager.AddEntity(stoneFloor(x, y))

			if dx == -width/2 || dx == width/2 || dy == -height/2 || dy == height/2 {
				// leave space for door
				if !(dx == -direction*width/2 && dy == 0) {
					ecsManager.AddEntity(treeTrunk(x, y))
				}
			}
		}
	}
	addTownsMember(ecsManager, xOff, yOff)
}

func grassFloor(x, y int) []ecs.Component {
	return floor(x, y, gui.GetSprite(gui.GRASS_FLOOR))
}

func dirtFloor(x, y int) []ecs.Component {
	return floor(x, y, gui.GetSprite(gui.DIRT_FLOOR))
}

func addTree(ecsManager *ecs.Manager, x, y int) {

	ecsManager.AddEntity(treeTrunk(x, y))

	treeRadius := rand.Intn(5) + 2
	for dx := -treeRadius; dx <= treeRadius; dx++ {
		uy := treeRadius - int(math.Abs(float64(dx)))
		for dy := -uy; dy <= uy; dy++ {

			ecsManager.AddEntity(leaf(x+dx, y+dy))
		}
	}
}

func treeTrunk(x, y int) []ecs.Component {
	return wall(x, y, gui.GetSprite(gui.TREE_TRUNK))
}

func leaf(x, y int) []ecs.Component {
	return []ecs.Component{
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.LEAF)}},
		{ecs.MEMORABLE, ecs.Memorable{}},
	}
}

func stoneFloor(x, y int) []ecs.Component {
	return floor(x, y, gui.GetSprite(gui.STONE_FLOOR))
}
func stoneWall(x, y int) []ecs.Component {
	return wall(x, y, gui.GetSprite(gui.STONE_WALL))
}

func floor(x, y int, sprite gui.Sprite) []ecs.Component {
	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Displayable{sprite}
	memorableComponent := ecs.Memorable{}

	return []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAYABLE, displayComponent},
		{ecs.MEMORABLE, memorableComponent},
	}
}

func wall(x, y int, sprite gui.Sprite) []ecs.Component {
	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Displayable{sprite}
	memorableComponent := ecs.Memorable{}
	volumeTag := ecs.Volume{}
	opaqueTag := ecs.Opaque{}

	return []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAYABLE, displayComponent},
		{ecs.MEMORABLE, memorableComponent},
		{ecs.VOLUME, volumeTag},
		{ecs.OPAQUE, opaqueTag},
	}
}

// func vary(color termbox.Attribute, colorVariance int) termbox.Attribute {

// 	r, g, b := termbox.AttributeToRGB(color)
// 	dr := uint8(rand.Intn(colorVariance))
// 	dg := uint8(rand.Intn(colorVariance))
// 	db := uint8(rand.Intn(colorVariance))

// 	if r < 1<<7-dr {
// 		r += dr
// 	}
// 	if g < 1<<7-dg {
// 		g += dg
// 	}
// 	if b < 1<<7-db {
// 		b += db
// 	}
// 	return termbox.RGBToAttribute(r, g, b)
// }
