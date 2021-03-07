package main

import (
	"math/rand"

	"github.com/jasonfantl/rogue/ecs"
	"github.com/nsf/termbox-go"
)

func generateCaveMask(ecsManager *ecs.Manager, width, height int) [][]bool {

	deathLimit := 4
	birthLimit := 4
	iterationCount := 8

	// simple cellular automota implementation
	tiles := make([][]bool, width)

	// first: generate static noise with solid sides
	sideThickness := 2
	openingSize := 20
	for x := 0; x < width; x++ {
		tiles[x] = make([]bool, height)
		for y := 0; y < height; y++ {
			if (x < sideThickness || x > width-sideThickness || y < sideThickness || y > height-sideThickness) &&
				!(y < sideThickness && x > width/2-openingSize/2 && x < width/2+openingSize/2) {
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
	pathWidth := 3

	pathX := width / 2
	pathY := height

	for pathY >= 0 {

		for dx := -pathWidth - rand.Intn(2); dx < pathWidth+rand.Intn(2); dx++ {
			x := pathX + dx + xOff
			y := pathY + yOff
			// randomly place stone in path
			if rand.Intn(5) == 0 {
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

		pathX += rand.Intn(3) - 1

		pathY--
	}
	//tree likelyhood
	treeChance := 50
	stoneChance := 50

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

			if rand.Intn(stoneChance) == 0 {
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

func grassFloor(x, y int) []ecs.Component {
	return floor(x, y, vary(termbox.RGBToAttribute(15, 110, 0), 6))
}

func dirtFloor(x, y int) []ecs.Component {
	return floor(x, y, vary(termbox.RGBToAttribute(70, 50, 0), 10))
}

func addTree(ecsManager *ecs.Manager, x, y int) {

	ecsManager.AddEntity(wall(x, y, termbox.RGBToAttribute(60, 40, 0)))

	treeRadius := rand.Intn(5) + 2
	for dx := -treeRadius; dx <= treeRadius; dx++ {
		for dy := -treeRadius; dy <= treeRadius; dy++ {
			thickness := 1
			hasLeaf := (rand.Intn(treeRadius-thickness)+thickness)*treeRadius > (dx*dx + dy*dy)
			if (dx == 0 && dy == 0) || !hasLeaf {
				continue
			}

			colorTint := uint8(treeRadius*treeRadius - (dx*dx+dy*dy)/2)
			leaf := []ecs.Component{
				{ecs.POSITION, ecs.Position{x + dx, y + dy}},
				{ecs.DISPLAYABLE, ecs.Displayable{false, vary(termbox.RGBToAttribute(40, 110-colorTint, 20), 4), ' ', 103}},
				{ecs.MEMORABLE, ecs.Memorable{}},
			}

			ecsManager.AddEntity(leaf)
		}
	}
}

func stoneFloor(x, y int) []ecs.Component {
	return floor(x, y, vary(termbox.RGBToAttribute(100, 100, 100), 6))
}
func stoneWall(x, y int) []ecs.Component {
	return wall(x, y, termbox.RGBToAttribute(250, 250, 250))
}

func floor(x, y int, color termbox.Attribute) []ecs.Component {
	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Displayable{false, color, ' ', 101}
	memorableComponent := ecs.Memorable{}

	return []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAYABLE, displayComponent},
		{ecs.MEMORABLE, memorableComponent},
	}
}

func wall(x, y int, color termbox.Attribute) []ecs.Component {
	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Displayable{false, color, ' ', 199}
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

func vary(color termbox.Attribute, colorVariance int) termbox.Attribute {

	r, g, b := termbox.AttributeToRGB(color)
	dr := uint8(rand.Intn(colorVariance))
	dg := uint8(rand.Intn(colorVariance))
	db := uint8(rand.Intn(colorVariance))

	if r < 1<<7-dr {
		r += dr
	}
	if g < 1<<7-dg {
		g += dg
	}
	if b < 1<<7-db {
		b += db
	}
	return termbox.RGBToAttribute(r, g, b)
}
