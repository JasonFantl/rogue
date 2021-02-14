package main

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/jasonfantl/rogue/ecs"
)

func generateRooms(ecsManager *ecs.Manager, width, height int) {

	deathLimit := 3
	birthLimit := 4
	// simple cellular automota implementation
	tiles := make([][]bool, width)

	// first: generate static noise
	for x := 0; x < width; x++ {
		tiles[x] = make([]bool, height)
		for y := 0; y < height; y++ {
			tiles[x][y] = (rand.Intn(2) == 1)
		}
	}

	// second: run CA a couple of times
	for step := 0; step < 3; step++ {
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

	// finally add entities
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if tiles[x][y] {
				placeFloor(ecsManager, x, y)
			} else {
				placeWall(ecsManager, x, y)
			}
		}
	}

	// just for good touch, add treasure
	tp := 10
	for tp > 0 {
		x := rand.Intn(width)
		y := rand.Intn(height)
		if tiles[x][y] {
			addTreasure(ecsManager, x, y)
			tp--
		}
	}
}

var floorStyle = tcell.StyleDefault.Background(tcell.ColorDimGray)
var wallStyle = tcell.StyleDefault.Background(tcell.ColorDarkGrey)

func placeFloor(ecsManager *ecs.Manager, x, y int) {
	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Display{' ', floorStyle, 1}

	floor := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAY, displayComponent},
	}

	ecsManager.AddEntity(floor)
}
func placeWall(ecsManager *ecs.Manager, x, y int) {
	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Display{' ', wallStyle, 99}
	blockableTag := ecs.Blockable{}

	wall := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAY, displayComponent},
		{ecs.BLOCKABLE, blockableTag},
	}

	ecsManager.AddEntity(wall)
}

func addPlayer(ecsManager *ecs.Manager, x, y int) {

	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Display{'@', floorStyle.Foreground(tcell.ColorDarkOrange), 99}
	controllerComponent := ecs.PlayerController{
		Up:     tcell.Key('w'),
		Down:   tcell.Key('s'),
		Left:   tcell.Key('a'),
		Right:  tcell.Key('d'),
		Pickup: tcell.Key(' '),
		Quit:   tcell.KeyEsc,
	}
	inventoryComponent := ecs.Inventory{}
	informationComponent := ecs.Information{"Player", "the hero of our story"}

	player := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAY, displayComponent},
		{ecs.PLAYER_CONTROLLER, controllerComponent},
		{ecs.INVENTORY, inventoryComponent},
		{ecs.INFORMATION, informationComponent},
	}

	ecsManager.AddEntity(player)
}

func addTreasure(ecsManager *ecs.Manager, x, y int) {

	treasure := []ecs.Component{
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.DISPLAY, ecs.Display{'$', floorStyle.Foreground(tcell.ColorYellow), 2}},
		{ecs.PICKUPABLE, ecs.Pickupable{}},
		{ecs.DROPABLE, ecs.Dropable{}},
		{ecs.INFORMATION, ecs.Information{"gold coin", "scratched, but still usable"}},
	}

	ecsManager.AddEntity(treasure)
}
