package contentGeneration

import (
	"github.com/jasonfantl/rogue/ecs"
	"github.com/jasonfantl/rogue/gui"
)

func GenerateGame(ecsManager *ecs.Manager, width, height int) {

	m := new(ecs.EntityManager)
	mask := make([][]bool, width)
	for x := 0; x < width; x++ {
		mask[x] = make([]bool, height)
		for y := 0; y < height; y++ {
			mask[x][y] = false
		}
	}

	island := generateIslandHeightmap(width, height, 100)
	biomes := biomesMask(island, []BiomeType{FOREST})

	addTowns(m, mask, island, biomes)

	// for x := 0; x < width; x++ {
	// 	for y := 0; y < height; y++ {
	// 		symbol := " "
	// 		if mask[x][y] {
	// 			symbol = "@"
	// 		} else if island[x][y] > 44 {
	// 			symbol = "."
	// 		}
	// 		fmt.Printf(symbol)
	// 	}
	// 	fmt.Printf("\n")
	// }

	addSea(m, mask, biomes)
	addBeach(m, mask, biomes)

	addCaves(m, mask, biomes)
	addForest(m, mask, biomes)

	// // then add cave entities
	// itemCount := width + height
	// for itemCount > 0 {
	// 	x := rand.Intn(width)
	// 	y := rand.Intn(height)
	// 	if tiles[x][y] {
	// 		r := rand.Intn(5)
	// 		if r == 0 {
	// 			addTreasure(ecsManager, x, y)
	// 		} else if r == 1 {
	// 			addMonster(ecsManager, x, y)
	// 		} else if r == 2 {
	// 			addPotion(ecsManager, x, y)
	// 		} else if r == 3 {
	// 			addWeapon(ecsManager, x, y)
	// 		} else {
	// 			addArmor(ecsManager, x, y)
	// 		}
	// 		itemCount--
	// 	}
	// }

	addPlayer(m, ecsManager, width/2, height/2)
	addWeapon(m, width/2+1, height/2)

	m.UnloadAllChunks()
}

func addTreasure(ecsManager *ecs.EntityManager, x, y int) {
	treasure := map[ecs.ComponentID]interface{}{
		ecs.POSITION:  ecs.Position{x, y},
		ecs.STASHABLE: ecs.Stashable{},
	}

	// treasureInfos := []map[ecs.ComponentID]interface{}{
	// 	{
	// 		{ecs.INFORMATION, ecs.Information{"gold coin", "scratched, but still usable"}},
	// 		{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.GOLD_COIN)}},
	// 	},
	// 	{
	// 		{ecs.INFORMATION, ecs.Information{"gem", "red and uncut"}},
	// 		{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.GEM)}},
	// 	},
	// 	{
	// 		{ecs.INFORMATION, ecs.Information{"silver coin", "might buy you a mug"}},
	// 		{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.SILVER_COIN)}},
	// 	},
	// }
	// treasureInfo := treasureInfos[rand.Intn(len(treasureInfos))]

	// treasure = append(treasure, treasureInfo...)

	ecsManager.AddEntity(treasure)
}

func addWeapon(ecsManager *ecs.EntityManager, x, y int) {

	weapon := map[ecs.ComponentID]interface{}{
		ecs.POSITION:   ecs.Position{x, y},
		ecs.STASHABLE:  ecs.Stashable{},
		ecs.PROJECTILE: ecs.Projectile{},
	}

	// weaponInfos := []map[ecs.ComponentID]interface{}{
	// 	{
	// 		{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.SWORD)}},
	// 		{ecs.INFORMATION, ecs.Information{"sword", "rusted"}},
	// 		{ecs.DAMAGE, ecs.Damage{16}},
	// 	},
	// 	{
	// 		{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.STICK)}},
	// 		{ecs.INFORMATION, ecs.Information{"stick", "primative, but better then nothing"}},
	// 		{ecs.DAMAGE, ecs.Damage{8}},
	// 	},
	// }

	// weaponInfo := weaponInfos[rand.Intn(len(weaponInfos))]

	// weapon = append(weapon, weaponInfo...)

	ecsManager.AddEntity(weapon)
}

func addArmor(ecsManager *ecs.EntityManager, x, y int) {

	armor := map[ecs.ComponentID]interface{}{
		ecs.POSITION:  ecs.Position{x, y},
		ecs.STASHABLE: ecs.Stashable{},
	}

	// armorInfos := []map[ecs.ComponentID]interface{}{
	// 	{
	// 		{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.LEATHER_ARMOR)}},
	// 		{ecs.INFORMATION, ecs.Information{"Leather armor", "sturdy and well worn"}},
	// 		{ecs.DAMAGE_RESISTANCE, ecs.DamageResistance{5}},
	// 	},
	// 	{
	// 		{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.METAL_ARMOR)}},
	// 		{ecs.INFORMATION, ecs.Information{"Metal plate", "shiny, dented"}},
	// 		{ecs.DAMAGE_RESISTANCE, ecs.DamageResistance{10}},
	// 	},
	// }

	// weaponInfo := armorInfos[rand.Intn(len(armorInfos))]

	// armor = append(armor, weaponInfo...)

	ecsManager.AddEntity(armor)
}

func addPotion(ecsManager *ecs.EntityManager, x, y int) {

	potion := map[ecs.ComponentID]interface{}{
		ecs.POSITION:    ecs.Position{x, y},
		ecs.STASHABLE:   ecs.Stashable{},
		ecs.DISPLAYABLE: ecs.Displayable{gui.GetSprite(gui.POTION)},
		ecs.CONSUMABLE:  ecs.Consumable{},
	}

	// potionInfos := []map[ecs.ComponentID]interface{}{
	// 	{
	// 		ecs.INFORMATION: ecs.Information{"Potion", "glowing red"},
	// 		ecs.REACTIONS: ecs.Reactions{[]ecs.Reaction{
	// 			ecs.Reaction{
	// 				ecs.CONSUMED,
	// 				ecs.HealReaction{10},
	// 			},
	// 		}},
	// 	},
	// 	{
	// 		ecs.INFORMATION: ecs.Information{"Potion", "dark blue, hard to see"},
	// 		ecs.REACTIONS: ecs.Reactions{[]ecs.Reaction{
	// 			ecs.Reaction{
	// 				ecs.CONSUMED,
	// 				&ecs.VisionIncreaseReaction{2},
	// 			},
	// 		}},
	// 	},
	// 	{
	// 		ecs.INFORMATION: ecs.Information{"Potion", "green, viscious"},
	// 		ecs.REACTIONS: ecs.Reactions{[]ecs.Reaction{
	// 			ecs.Reaction{
	// 				ecs.CONSUMED,
	// 				ecs.StrengthIncreaseReaction{1},
	// 			},
	// 		}},
	// 	},
	// }

	// potionInfo := potionInfos[rand.Intn(len(potionInfos))]

	// potion = append(potion, potionInfo...)

	ecsManager.AddEntity(potion)
}

// how to seperate lock and key pair?
func addDoor(ecsManager *ecs.EntityManager, mask [][]bool, x, y int) {

	key := map[ecs.ComponentID]interface{}{
		ecs.POSITION:    ecs.Position{x + 1, y},
		ecs.DISPLAYABLE: ecs.Displayable{gui.GetSprite(gui.KEY)},
		ecs.STASHABLE:   ecs.Stashable{},
	}

	keyEntity := addEntity(ecsManager, mask, true, false, x, y, key)

	// locked compoentn isnt great, have to add compoennt twce if inversed
	door := map[ecs.ComponentID]interface{}{
		ecs.POSITION: ecs.Position{x, y},
		ecs.LOCKABLE: ecs.Lockable{
			keyEntity,
			true,
			[]ecs.Component{
				{ecs.VOLUME, ecs.Volume{}},
				{ecs.OPAQUE, ecs.Opaque{}},
				{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.CLOSED_DOOR)}}},
			[]ecs.Component{
				{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.OPEN_DOOR)}}},
		},
	}

	addEntity(ecsManager, mask, true, true, x, y, door)

}

func addEntity(ecsManager *ecs.EntityManager, mask [][]bool, ignoreMask, effectMask bool, x, y int, entity map[ecs.ComponentID]interface{}) ecs.Entity {
	if x >= 0 && x < len(mask) {
		if y >= 0 && y < len(mask[x]) {
			if ignoreMask || !mask[x][y] {
				mask[x][y] = mask[x][y] || effectMask
				return ecsManager.AddEntity(entity)
			}
		}
	}
	return 0
}
