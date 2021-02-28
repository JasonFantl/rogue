package main

import (
	"math/rand"

	"github.com/jasonfantl/rogue/ecs"
	"github.com/nsf/termbox-go"
)

func generateRooms(ecsManager *ecs.Manager, width, height int) {

	addWallsAround(ecsManager, width+1, height+1)

	deathLimit := 4
	birthLimit := 4
	iterationCount := 8

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

	// finally add entities
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if tiles[x][y] {
				placeFloor(ecsManager, x+1, y+1)
			} else {
				placeWall(ecsManager, x+1, y+1)
			}
		}
	}

	// just for good touch, add treasure and monsters
	itemCount := 60
	for itemCount > 0 {
		x := rand.Intn(width-2) + 1
		y := rand.Intn(height-2) + 1
		if tiles[x][y] {
			if rand.Intn(8) > 1 {
				addTreasure(ecsManager, x+1, y+1)
			} else {
				addMonster(ecsManager, x+1, y+1)
			}
			itemCount--
		}
	}
}

func addWallsAround(ecsManager *ecs.Manager, width, height int) {
	for x := 1; x < width; x++ {
		placeWall(ecsManager, x, height)
		placeWall(ecsManager, x, 0)
	}
	for y := 1; y < height; y++ {
		placeWall(ecsManager, 0, y)
		placeWall(ecsManager, width, y)
	}
}

func placeFloor(ecsManager *ecs.Manager, x, y int) {
	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Displayable{false, termbox.ColorBlack, ' ', 1}

	floor := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAYABLE, displayComponent},
	}

	ecsManager.AddEntity(floor)
}
func placeWall(ecsManager *ecs.Manager, x, y int) {
	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Displayable{false, termbox.ColorCyan, ' ', 99}
	blockableTag := ecs.Volume{}

	wall := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAYABLE, displayComponent},
		{ecs.VOLUME, blockableTag},
	}

	ecsManager.AddEntity(wall)
}

func addPlayer(ecsManager *ecs.Manager, x, y int) {

	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Displayable{true, termbox.ColorCyan, '@', 99}
	visionComponent := ecs.Vision{10}
	awarnessComponent := ecs.EntityAwarness{}

	controllerComponent := ecs.PlayerController{
		Up:     termbox.Key('w'),
		Down:   termbox.Key('s'),
		Left:   termbox.Key('a'),
		Right:  termbox.Key('d'),
		Pickup: termbox.Key('e'),
		Quit:   termbox.KeyEsc,
	}
	inventoryComponent := ecs.Inventory{}
	informationComponent := ecs.Information{"Player", "the hero of our story"}
	volumeComponent := ecs.Volume{}
	violentComponent := ecs.Violent{4}
	healthComponent := ecs.Health{20, 20}

	player := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAYABLE, displayComponent},
		{ecs.ENTITY_AWARENESS, awarnessComponent},
		{ecs.VISION, visionComponent},
		{ecs.PLAYER_CONTROLLER, controllerComponent},
		{ecs.INVENTORY, inventoryComponent},
		{ecs.INFORMATION, informationComponent},
		{ecs.VOLUME, volumeComponent},
		{ecs.VIOLENT, violentComponent},
		{ecs.HEALTH, healthComponent},
	}

	ecsManager.AddEntity(player)
}

func addMonster(ecsManager *ecs.Manager, x, y int) {

	controllerComponent := ecs.MonsterController{
		[]ecs.MonsterAction{ecs.PICKUP, ecs.TREASURE_MOVE, ecs.RANDOM_MOVE, ecs.NOTHING},
	}
	positionComponent := ecs.Position{x, y}
	inventoryComponent := ecs.Inventory{}
	volumeComponent := ecs.Volume{}
	violentComponent := ecs.Violent{4}
	healthComponent := ecs.Health{20, 20}
	visionComponent := ecs.Vision{10}
	awarnessComponent := ecs.EntityAwarness{}

	monster := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.MONSTER_CONTROLLER, controllerComponent},
		{ecs.ENTITY_AWARENESS, awarnessComponent},
		{ecs.VISION, visionComponent},
		{ecs.INVENTORY, inventoryComponent},
		{ecs.VOLUME, volumeComponent},
		{ecs.VIOLENT, violentComponent},
		{ecs.HEALTH, healthComponent},
	}

	monsterInfos := [][]ecs.Component{
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.ColorRed, 'M', 99}},
			{ecs.INFORMATION, ecs.Information{"Monster", "generic"}},
		},
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.ColorGreen, 'O', 99}},
			{ecs.INFORMATION, ecs.Information{"Ogre", "Big and scary"}},
		},
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.ColorGreen, 'g', 99}},
			{ecs.INFORMATION, ecs.Information{"Goblin", "green and scrawny, sill scary though"}},
		},
	}

	monsterInfo := monsterInfos[rand.Intn(len(monsterInfos))]

	monster = append(monster, monsterInfo...)

	ecsManager.AddEntity(monster)
}

func addTreasure(ecsManager *ecs.Manager, x, y int) {

	treasureInfos := [][]string{{"gold coin", "scratched, but still usable"}, {"gem", "red and uncut"}, {"silver coin", "might buy you a mug"}}
	treasureInfo := treasureInfos[rand.Intn(len(treasureInfos))]

	treasure := []ecs.Component{
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.ColorYellow, '$', 2}},
		{ecs.PICKUPABLE, ecs.Pickupable{}},
		{ecs.DROPABLE, ecs.Dropable{}},
		{ecs.INFORMATION, ecs.Information{treasureInfo[0], treasureInfo[1]}},
	}

	ecsManager.AddEntity(treasure)
}
