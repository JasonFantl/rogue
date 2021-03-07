package main

import (
	"math/rand"

	"github.com/jasonfantl/rogue/ecs"
	"github.com/nsf/termbox-go"
)

func generateGame(ecsManager *ecs.Manager, width, height int) {
	tiles := generateCaveMask(ecsManager, width, height)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if tiles[x][y] {
				ecsManager.AddEntity(stoneFloor(x, y))
			} else {
				ecsManager.AddEntity(stoneWall(x, y))
			}
		}
	}

	// then add cave entities
	itemCount := width
	for itemCount > 0 {
		x := rand.Intn(width)
		y := rand.Intn(height)
		if tiles[x][y] {
			r := rand.Intn(5)
			if r == 0 {
				addTreasure(ecsManager, x, y)
			} else if r == 1 {
				addMonster(ecsManager, x, y)
			} else if r == 2 {
				addPotion(ecsManager, x, y)
			} else if r == 3 {
				addWeapon(ecsManager, x, y)
			} else {
				addArmor(ecsManager, x, y)
			}
			itemCount--
		}
	}

	generateForest(ecsManager, width, height, 0, -height)

	addPlayer(ecsManager, width/2, -5)
}

func addPlayer(ecsManager *ecs.Manager, x, y int) {

	positionComponent := ecs.Position{x, y}
	displayComponent := ecs.Displayable{true, termbox.RGBToAttribute(200, 200, 250), '@', 199}
	visionComponent := ecs.Vision{20}
	awarnessComponent := ecs.EntityAwarness{}
	memoryComponent := ecs.EntityMemory{}
	inventoryComponent := ecs.Inventory{}
	informationComponent := ecs.Information{"Player", "the hero of our story"}
	volumeComponent := ecs.Volume{}
	fighterComponent := ecs.Fighter{10, 0, 0}
	damageComponent := ecs.Damage{1}
	healthComponent := ecs.Health{100, 80}

	player := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.DISPLAYABLE, displayComponent},
		{ecs.ENTITY_AWARENESS, awarnessComponent},
		{ecs.ENTITY_MEMORY, memoryComponent},
		{ecs.VISION, visionComponent},
		{ecs.INVENTORY, inventoryComponent},
		{ecs.INFORMATION, informationComponent},
		{ecs.VOLUME, volumeComponent},
		{ecs.FIGHTER, fighterComponent},
		{ecs.DAMAGE, damageComponent},
		{ecs.HEALTH, healthComponent},
	}

	playerID := ecsManager.AddEntity(player)

	user := []ecs.Component{
		{ecs.PLAYER_CONTROLLER, ecs.PlayerController{
			Controlling: playerID,
			Up:          termbox.Key('w'),
			Down:        termbox.Key('s'),
			Left:        termbox.Key('a'),
			Right:       termbox.Key('d'),
			Pickup:      termbox.Key('e'),
			Consume:     termbox.Key('r'),
			Quit:        termbox.KeyEsc,
		},
		}}

	ecsManager.AddEntity(user)
}

func addMonster(ecsManager *ecs.Manager, x, y int) {

	controllerComponent := ecs.MonsterController{
		[]ecs.MonsterAction{ecs.PICKUP, ecs.TREASURE_MOVE, ecs.RANDOM_MOVE, ecs.NOTHING},
	}
	positionComponent := ecs.Position{x, y}
	inventoryComponent := ecs.Inventory{}
	volumeComponent := ecs.Volume{}
	fighterComponent := ecs.Fighter{10, 0, 0}
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

func addArmor(ecsManager *ecs.Manager, x, y int) {

	armor := []ecs.Component{
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.PICKUPABLE, ecs.Pickupable{}},
	}

	armorInfos := [][]ecs.Component{
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.RGBToAttribute(150, 100, 10), 'Y', 102}},
			{ecs.INFORMATION, ecs.Information{"Leather armor", "sturdy and well worn"}},
			{ecs.DAMAGE_RESISTANCE, ecs.DamageResistance{5}},
		},
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.RGBToAttribute(212, 212, 212), 'Y', 102}},
			{ecs.INFORMATION, ecs.Information{"Metal plate", "shiny, dented"}},
			{ecs.DAMAGE_RESISTANCE, ecs.DamageResistance{10}},
		},
	}

	weaponInfo := armorInfos[rand.Intn(len(armorInfos))]

	armor = append(armor, weaponInfo...)

	ecsManager.AddEntity(armor)
}

func addPotion(ecsManager *ecs.Manager, x, y int) {

	potion := []ecs.Component{
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.PICKUPABLE, ecs.Pickupable{}},
	}

	potionInfos := [][]ecs.Component{
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.RGBToAttribute(140, 20, 40), 'o', 102}},
			{ecs.INFORMATION, ecs.Information{"Potion", "glowing red"}},
			{ecs.EFFECTS, ecs.Effects{[]ecs.Effect{
				ecs.Effect{
					ecs.PICKED_UP,
					ecs.HealEffect{10},
				},
			}}},
		},
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.RGBToAttribute(40, 50, 250), 'o', 102}},
			{ecs.INFORMATION, ecs.Information{"Potion", "dark blue, hard to see"}},
			{ecs.EFFECTS, ecs.Effects{[]ecs.Effect{
				ecs.Effect{
					ecs.PICKED_UP,
					ecs.VisionEffect{2},
				},
			}}},
		},
		{
			{ecs.DISPLAYABLE, ecs.Displayable{true, termbox.RGBToAttribute(40, 250, 50), 'o', 102}},
			{ecs.INFORMATION, ecs.Information{"Potion", "green, viscious"}},
			{ecs.EFFECTS, ecs.Effects{[]ecs.Effect{
				ecs.Effect{
					ecs.PICKED_UP,
					ecs.StrengthEffect{1},
				},
			}}},
		},
	}

	potionInfo := potionInfos[rand.Intn(len(potionInfos))]

	potion = append(potion, potionInfo...)

	ecsManager.AddEntity(potion)
}
