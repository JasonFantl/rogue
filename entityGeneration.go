package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jasonfantl/rogue/ecs"
	"github.com/jasonfantl/rogue/gui"
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
	itemCount := width + height
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
	addWeapon(ecsManager, width/2+1, -5)

}

func addPlayer(ecsManager *ecs.Manager, x, y int) {

	player := []ecs.Component{
		{ecs.BRAIN, ecs.Brain{
			[]ecs.DesiredAction{},
		}},
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.PLAYER)}},
		{ecs.ENTITY_AWARENESS, ecs.EntityAwarness{}},
		{ecs.ENTITY_MEMORY, ecs.EntityMemory{}},
		{ecs.VISION, ecs.Vision{10}},
		{ecs.INVENTORY, ecs.Inventory{}},
		{ecs.INFORMATION, ecs.Information{"Player", "the hero of our story"}},
		{ecs.VOLUME, ecs.Volume{}},
		{ecs.FIGHTER, ecs.Fighter{10, 0, 0}},
		{ecs.DAMAGE, ecs.Damage{1}},
		{ecs.HEALTH, ecs.Health{100, 90}},
	}

	playerID := ecsManager.AddEntity(player)

	user := []ecs.Component{
		{ecs.USER, ecs.User{
			Controlling: playerID,
			UpKey:       ebiten.KeyW,
			DownKey:     ebiten.KeyS,
			LeftKey:     ebiten.KeyA,
			RightKey:    ebiten.KeyD,
			ActionKey:   ebiten.KeyE,
			MenuKey:     ebiten.KeyQ,
			QuitKey:     ebiten.KeyEscape,
			Menu:        ecs.Menu{},
		},
		}}

	ecsManager.AddEntity(user)
}

func addMonster(ecsManager *ecs.Manager, x, y int) {

	monster := []ecs.Component{
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.BRAIN, ecs.Brain{
			[]ecs.DesiredAction{ecs.PICKUP, ecs.TREASURE_MOVE, ecs.RANDOM_MOVE, ecs.DO_NOTHING},
		}},
		{ecs.ENTITY_AWARENESS, ecs.EntityAwarness{}},
		{ecs.VISION, ecs.Vision{10}},
		{ecs.VOLUME, ecs.Volume{}},
		{ecs.FIGHTER, ecs.Fighter{10, 0, 0}},
		{ecs.DAMAGE, ecs.Damage{3}},
		{ecs.HEALTH, ecs.Health{50, 50}},
	}

	monsterInfos := [][]ecs.Component{
		{
			{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.MONSTER3)}},
			{ecs.INFORMATION, ecs.Information{"Monster", "generic"}},
		},
		{
			{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.MONSTER2)}},
			{ecs.INVENTORY, ecs.Inventory{}},
			{ecs.INFORMATION, ecs.Information{"Ogre", "Big and scary"}},
		},
		{
			{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.MONSTER1)}},
			{ecs.INVENTORY, ecs.Inventory{}},
			{ecs.INFORMATION, ecs.Information{"Goblin", "green and scrawny, sill scary though"}},
		},
	}

	monsterInfo := monsterInfos[rand.Intn(len(monsterInfos))]

	monster = append(monster, monsterInfo...)

	ecsManager.AddEntity(monster)
}

func addTownsMember(ecsManager *ecs.Manager, x, y int) {

	controllerComponent := ecs.Brain{
		[]ecs.DesiredAction{ecs.PICKUP, ecs.RANDOM_MOVE, ecs.DO_NOTHING},
	}

	positionComponent := ecs.Position{x, y}
	inventoryComponent := ecs.Inventory{}
	volumeComponent := ecs.Volume{}
	damageComponent := ecs.Damage{3}
	healthComponent := ecs.Health{50, 50}
	visionComponent := ecs.Vision{10}
	awarnessComponent := ecs.EntityAwarness{}

	townsMember := []ecs.Component{
		{ecs.POSITION, positionComponent},
		{ecs.BRAIN, controllerComponent},
		{ecs.ENTITY_AWARENESS, awarnessComponent},
		{ecs.VISION, visionComponent},
		{ecs.INVENTORY, inventoryComponent},
		{ecs.VOLUME, volumeComponent},
		{ecs.DAMAGE, damageComponent},
		{ecs.HEALTH, healthComponent},
	}

	townsMemberInfos := [][]ecs.Component{
		{
			{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.PLAYER)}},
			{ecs.INFORMATION, ecs.Information{"Human", "lives in town"}},
		},
	}

	townsMemberInfo := townsMemberInfos[rand.Intn(len(townsMemberInfos))]

	townsMember = append(townsMember, townsMemberInfo...)

	ecsManager.AddEntity(townsMember)
}

func addTreasure(ecsManager *ecs.Manager, x, y int) {
	treasure := []ecs.Component{
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.STASHABLE, ecs.Stashable{}},
	}

	treasureInfos := [][]ecs.Component{
		{
			{ecs.INFORMATION, ecs.Information{"gold coin", "scratched, but still usable"}},
			{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.GOLD_COIN)}},
		},
		{
			{ecs.INFORMATION, ecs.Information{"gem", "red and uncut"}},
			{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.GEM)}},
		},
		{
			{ecs.INFORMATION, ecs.Information{"silver coin", "might buy you a mug"}},
			{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.SILVER_COIN)}},
		},
	}
	treasureInfo := treasureInfos[rand.Intn(len(treasureInfos))]

	treasure = append(treasure, treasureInfo...)

	ecsManager.AddEntity(treasure)
}

func addWeapon(ecsManager *ecs.Manager, x, y int) {

	weapon := []ecs.Component{
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.STASHABLE, ecs.Stashable{}},
		{ecs.PROJECTILE, ecs.Projectile{}},
	}

	weaponInfos := [][]ecs.Component{
		{
			{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.SWORD)}},
			{ecs.INFORMATION, ecs.Information{"sword", "rusted"}},
			{ecs.DAMAGE, ecs.Damage{16}},
		},
		{
			{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.STICK)}},
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
		{ecs.STASHABLE, ecs.Stashable{}},
	}

	armorInfos := [][]ecs.Component{
		{
			{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.LEATHER_ARMOR)}},
			{ecs.INFORMATION, ecs.Information{"Leather armor", "sturdy and well worn"}},
			{ecs.DAMAGE_RESISTANCE, ecs.DamageResistance{5}},
		},
		{
			{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.METAL_ARMOR)}},
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
		{ecs.STASHABLE, ecs.Stashable{}},
		{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.POTION)}},
		{ecs.CONSUMABLE, ecs.Consumable{}},
	}

	potionInfos := [][]ecs.Component{
		{
			{ecs.INFORMATION, ecs.Information{"Potion", "glowing red"}},
			{ecs.REACTIONS, ecs.Reactions{[]ecs.Reaction{
				ecs.Reaction{
					ecs.CONSUMED,
					ecs.HealReaction{10},
				},
			}}},
		},
		{
			{ecs.INFORMATION, ecs.Information{"Potion", "dark blue, hard to see"}},
			{ecs.REACTIONS, ecs.Reactions{[]ecs.Reaction{
				ecs.Reaction{
					ecs.CONSUMED,
					&ecs.VisionIncreaseReaction{2},
				},
			}}},
		},
		{
			{ecs.INFORMATION, ecs.Information{"Potion", "green, viscious"}},
			{ecs.REACTIONS, ecs.Reactions{[]ecs.Reaction{
				ecs.Reaction{
					ecs.CONSUMED,
					ecs.StrengthIncreaseReaction{1},
				},
			}}},
		},
	}

	potionInfo := potionInfos[rand.Intn(len(potionInfos))]

	potion = append(potion, potionInfo...)

	ecsManager.AddEntity(potion)
}

// how to seperate lock and key pair?
func addDoor(ecsManager *ecs.Manager, x, y int) {

	key := []ecs.Component{
		{ecs.POSITION, ecs.Position{x + 1, y}},
		{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.KEY)}},
		{ecs.STASHABLE, ecs.Stashable{}},
	}

	keyEntity := ecsManager.AddEntity(key)

	// locked compoentn isnt great, have to add compoennt twce if inversed
	door := []ecs.Component{
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.LOCKABLE, ecs.Lockable{
			keyEntity,
			true,
			[]ecs.Component{
				{ecs.VOLUME, ecs.Volume{}},
				{ecs.OPAQUE, ecs.Opaque{}},
				{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.CLOSED_DOOR)}}},
			[]ecs.Component{
				{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.OPEN_DOOR)}}},
		}},
	}

	ecsManager.AddEntity(door)
}
