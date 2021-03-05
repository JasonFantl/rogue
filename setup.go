package main

import (
	"math/rand"

	"github.com/jasonfantl/rogue/ecs"
	"github.com/nsf/termbox-go"
)

func generateGame(ecsManager *ecs.Manager, width, height int) {

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

	// then add people and items
	addedPlayer := false
	itemCount := 60
	for itemCount > 0 {
		x := rand.Intn(width-2) + 2
		y := rand.Intn(height-2) + 2
		if tiles[x-1][y-1] {
			if !addedPlayer {
				addedPlayer = true
				addPlayer(ecsManager, x, y)
				continue
			}

			r := rand.Intn(4)
			if r == 0 {
				addTreasure(ecsManager, x, y)
			} else if r == 1 {
				addMonster(ecsManager, x, y)
			} else if r == 2 {
				addPotion(ecsManager, x, y)
			} else {
				addWeapon(ecsManager, x, y)
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
	displayComponent := ecs.Displayable{false, termbox.RGBToAttribute(100, 100, 100), ' ', 101}
	memorableComponent := ecs.Memorable{}

	floor := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAYABLE, displayComponent},
		{ecs.MEMORABLE, memorableComponent},
	}

	ecsManager.AddEntity(floor)
}
func placeWall(ecsManager *ecs.Manager, x, y int) {
	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Displayable{false, termbox.RGBToAttribute(200, 200, 200), ' ', 199}
	memorableComponent := ecs.Memorable{}

	volumeTag := ecs.Volume{}
	opaqueTag := ecs.Opaque{}

	wall := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAYABLE, displayComponent},
		{ecs.MEMORABLE, memorableComponent},
		{ecs.VOLUME, volumeTag},
		{ecs.OPAQUE, opaqueTag},
	}

	ecsManager.AddEntity(wall)
}

func addPlayer(ecsManager *ecs.Manager, x, y int) {

	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Displayable{true, termbox.RGBToAttribute(200, 200, 250), '@', 199}
	visionComponent := ecs.Vision{10}
	awarnessComponent := ecs.EntityAwarness{}
	memoryComponent := ecs.EntityMemory{}
	controllerComponent := ecs.PlayerController{
		Up:      termbox.Key('w'),
		Down:    termbox.Key('s'),
		Left:    termbox.Key('a'),
		Right:   termbox.Key('d'),
		Pickup:  termbox.Key('e'),
		Consume: termbox.Key('r'),
		Quit:    termbox.KeyEsc,
	}
	inventoryComponent := ecs.Inventory{}
	informationComponent := ecs.Information{"Player", "the hero of our story"}
	volumeComponent := ecs.Volume{}
	fighterComponent := ecs.Fighter{10, 0}
	damageComponent := ecs.Damage{1}
	healthComponent := ecs.Health{100, 80}

	player := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAYABLE, displayComponent},
		{ecs.ENTITY_AWARENESS, awarnessComponent},
		{ecs.ENTITY_MEMORY, memoryComponent},
		{ecs.VISION, visionComponent},
		{ecs.PLAYER_CONTROLLER, controllerComponent},
		{ecs.INVENTORY, inventoryComponent},
		{ecs.INFORMATION, informationComponent},
		{ecs.VOLUME, volumeComponent},
		{ecs.FIGHTER, fighterComponent},
		{ecs.DAMAGE, damageComponent},
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
	fighterComponent := ecs.Fighter{10, 0}
	damageComponent := ecs.Damage{3}
	healthComponent := ecs.Health{50, 50}
	visionComponent := ecs.Vision{10}
	awarnessComponent := ecs.EntityAwarness{}

	monster := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.MONSTER_CONTROLLER, controllerComponent},
		{ecs.ENTITY_AWARENESS, awarnessComponent},
		{ecs.VISION, visionComponent},
		{ecs.INVENTORY, inventoryComponent},
		{ecs.VOLUME, volumeComponent},
		{ecs.FIGHTER, fighterComponent},
		{ecs.DAMAGE, damageComponent},
		{ecs.HEALTH, healthComponent},
	}

	monsterInfos := [][]ecs.Component{
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.RGBToAttribute(200, 20, 20), 'M', 199}},
			{ecs.INFORMATION, ecs.Information{"Monster", "generic"}},
		},
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.RGBToAttribute(20, 150, 80), 'O', 199}},
			{ecs.INFORMATION, ecs.Information{"Ogre", "Big and scary"}},
		},
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.RGBToAttribute(100, 200, 50), 'g', 199}},
			{ecs.INFORMATION, ecs.Information{"Goblin", "green and scrawny, sill scary though"}},
		},
	}

	monsterInfo := monsterInfos[rand.Intn(len(monsterInfos))]

	monster = append(monster, monsterInfo...)

	ecsManager.AddEntity(monster)
}

func addTreasure(ecsManager *ecs.Manager, x, y int) {
	treasure := []ecs.Component{
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.RGBToAttribute(200, 200, 20), '$', 102}},
		{ecs.PICKUPABLE, ecs.Pickupable{}},
	}

	treasureInfos := [][]ecs.Component{
		{
			{ecs.INFORMATION, ecs.Information{"gold coin", "scratched, but still usable"}}},
		{
			{ecs.INFORMATION, ecs.Information{"gem", "red and uncut"}}},
		{
			{ecs.INFORMATION, ecs.Information{"silver coin", "might buy you a mug"}}},
	}
	treasureInfo := treasureInfos[rand.Intn(len(treasureInfos))]

	treasure = append(treasure, treasureInfo...)

	ecsManager.AddEntity(treasure)
}

func addWeapon(ecsManager *ecs.Manager, x, y int) {

	weapon := []ecs.Component{
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.PICKUPABLE, ecs.Pickupable{}},
	}

	weaponInfos := [][]ecs.Component{
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.RGBToAttribute(150, 150, 150), '|', 102}},
			{ecs.INFORMATION, ecs.Information{"sword", "rusted"}},
			{ecs.DAMAGE, ecs.Damage{16}},
		},
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.RGBToAttribute(80, 40, 30), '|', 102}},
			{ecs.INFORMATION, ecs.Information{"stick", "primative, but better then nothing"}},
			{ecs.DAMAGE, ecs.Damage{8}},
		},
	}

	weaponInfo := weaponInfos[rand.Intn(len(weaponInfos))]

	weapon = append(weapon, weaponInfo...)

	ecsManager.AddEntity(weapon)
}

func addPotion(ecsManager *ecs.Manager, x, y int) {

	potion := []ecs.Component{
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.PICKUPABLE, ecs.Pickupable{}},
		{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.RGBToAttribute(240, 250, 0), 'o', 102}},
		{ecs.CONSUMABLE, ecs.Consumable{}},
	}

	potionInfos := [][]ecs.Component{
		{
			{ecs.EFFECTS, ecs.Effects{[]ecs.Effect{
				ecs.Effect{
					ecs.PICKED_UP,
					ecs.HealEffect{10},
				},
			}}},
		},
	}

	potionInfo := potionInfos[rand.Intn(len(potionInfos))]

	potion = append(potion, potionInfo...)

	ecsManager.AddEntity(potion)
}
