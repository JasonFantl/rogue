package contentGeneration

import (
	"math/rand"

	"github.com/jasonfantl/rogue/ecs"
	"github.com/jasonfantl/rogue/gui"
)

func GenerateGame(ecsManager *ecs.Manager, width, height int) {

	mask := make([][]bool, width)
	for x := 0; x < width; x++ {
		mask[x] = make([]bool, height)
		for y := 0; y < height; y++ {
			mask[x][y] = true
		}
	}

	generateBasic(ecsManager, mask, width, height/2)
	addCaves(ecsManager, mask)

	addPlayer(ecsManager, ecsManager, width/2, height/2)
	addWeapon(ecsManager, width/2+1, height/2)
}

func addTreasure(ecsManager *ecs.Manager, x, y int) {
	treasure := map[ecs.ComponentID]interface{}{
		ecs.POSITION:  ecs.Position{X: x, Y: y},
		ecs.STASHABLE: ecs.Stashable{},
	}

	treasureInfos := []map[ecs.ComponentID]interface{}{
		{
			ecs.INFORMATION: ecs.Information{Name: "gold coin", Details: "scratched, but still usable"},
			ecs.DISPLAYABLE: ecs.Displayable{Sprite: gui.GetSprite(gui.GOLD_COIN)},
		},
		{
			ecs.INFORMATION: ecs.Information{Name: "gem", Details: "red and uncut"},
			ecs.DISPLAYABLE: ecs.Displayable{Sprite: gui.GetSprite(gui.GEM)},
		},
		{
			ecs.INFORMATION: ecs.Information{Name: "silver coin", Details: "might buy you a mug"},
			ecs.DISPLAYABLE: ecs.Displayable{Sprite: gui.GetSprite(gui.SILVER_COIN)},
		},
	}

	addRandom(treasure, treasureInfos)

	ecsManager.AddEntity(treasure)
}

func addWeapon(ecsManager *ecs.Manager, x, y int) {

	weapon := map[ecs.ComponentID]interface{}{
		ecs.POSITION:   ecs.Position{X: x, Y: y},
		ecs.STASHABLE:  ecs.Stashable{},
		ecs.PROJECTILE: ecs.Projectile{},
	}

	weaponInfos := []map[ecs.ComponentID]interface{}{
		{
			ecs.DISPLAYABLE: ecs.Displayable{Sprite: gui.GetSprite(gui.SWORD)},
			ecs.INFORMATION: ecs.Information{Name: "sword", Details: "rusted"},
			ecs.DAMAGE:      ecs.Damage{Amount: 16},
		},
		{
			ecs.DISPLAYABLE: ecs.Displayable{Sprite: gui.GetSprite(gui.STICK)},
			ecs.INFORMATION: ecs.Information{Name: "stick", Details: "primitive, but better then nothing"},
			ecs.DAMAGE:      ecs.Damage{Amount: 8},
		},
	}

	addRandom(weapon, weaponInfos)

	ecsManager.AddEntity(weapon)
}

func addArmor(ecsManager *ecs.Manager, x, y int) {

	armor := map[ecs.ComponentID]interface{}{
		ecs.POSITION:  ecs.Position{X: x, Y: y},
		ecs.STASHABLE: ecs.Stashable{},
	}

	armorInfos := []map[ecs.ComponentID]interface{}{
		{
			ecs.DISPLAYABLE:       ecs.Displayable{Sprite: gui.GetSprite(gui.LEATHER_ARMOR)},
			ecs.INFORMATION:       ecs.Information{Name: "Leather armor", Details: "sturdy and well worn"},
			ecs.DAMAGE_RESISTANCE: ecs.DamageResistance{Amount: 5},
		},
		{
			ecs.DISPLAYABLE:       ecs.Displayable{Sprite: gui.GetSprite(gui.METAL_ARMOR)},
			ecs.INFORMATION:       ecs.Information{Name: "Metal plate", Details: "shiny, dented"},
			ecs.DAMAGE_RESISTANCE: ecs.DamageResistance{Amount: 10},
		},
	}

	addRandom(armor, armorInfos)

	ecsManager.AddEntity(armor)
}

func addPotion(ecsManager *ecs.Manager, x, y int) {

	potion := map[ecs.ComponentID]interface{}{
		ecs.POSITION:    ecs.Position{X: x, Y: y},
		ecs.STASHABLE:   ecs.Stashable{},
		ecs.DISPLAYABLE: ecs.Displayable{Sprite: gui.GetSprite(gui.POTION)},
		ecs.CONSUMABLE:  ecs.Consumable{},
	}

	potionInfos := []map[ecs.ComponentID]interface{}{
		{
			ecs.INFORMATION: ecs.Information{Name: "Potion", Details: "glowing red"},
			ecs.REACTIONS: ecs.Reactions{Reactions: []ecs.Reaction{
				{
					ReactionType: ecs.CONSUMED,
					Reaction:     ecs.HealReaction{Amount: 10},
				},
			}},
		},
		{
			ecs.INFORMATION: ecs.Information{Name: "Potion", Details: "dark blue, hard to see"},
			ecs.REACTIONS: ecs.Reactions{Reactions: []ecs.Reaction{
				{
					ReactionType: ecs.CONSUMED,
					Reaction:     &ecs.VisionIncreaseReaction{Amount: 2},
				},
			}},
		},
		{
			ecs.INFORMATION: ecs.Information{Name: "Potion", Details: "green, viscous"},
			ecs.REACTIONS: ecs.Reactions{Reactions: []ecs.Reaction{
				{
					ReactionType: ecs.CONSUMED,
					Reaction:     ecs.StrengthIncreaseReaction{Amount: 1},
				},
			}},
		},
	}

	addRandom(potion, potionInfos)

	ecsManager.AddEntity(potion)
}

// how to separate lock and key pair?
func addDoor(ecsManager *ecs.Manager, mask [][]bool, x, y int) {

	key := map[ecs.ComponentID]interface{}{
		ecs.POSITION:    ecs.Position{X: x + 1, Y: y},
		ecs.DISPLAYABLE: ecs.Displayable{Sprite: gui.GetSprite(gui.KEY)},
		ecs.STASHABLE:   ecs.Stashable{},
	}

	keyEntity := addEntity(ecsManager, mask, true, false, x, y, key)

	// locked component isn't great, have to add component twice if inverted
	door := map[ecs.ComponentID]interface{}{
		ecs.POSITION: ecs.Position{X: x, Y: y},
		ecs.LOCKABLE: ecs.Lockable{
			Key:    keyEntity,
			Locked: true,
			LockedComponents: []ecs.Component{
				{ID: ecs.VOLUME, Data: ecs.Volume{}},
				{ID: ecs.OPAQUE, Data: ecs.Opaque{}},
				{ID: ecs.DISPLAYABLE, Data: ecs.Displayable{Sprite: gui.GetSprite(gui.CLOSED_DOOR)}}},
			UnlockedComponents: []ecs.Component{
				{ID: ecs.DISPLAYABLE, Data: ecs.Displayable{Sprite: gui.GetSprite(gui.OPEN_DOOR)}}},
		},
	}

	addEntity(ecsManager, mask, true, true, x, y, door)

}

func addEntity(ecsManager *ecs.Manager, mask [][]bool, ignoreMask, effectMask bool, x, y int, entity map[ecs.ComponentID]interface{}) ecs.Entity {
	if x >= 0 && x < len(mask) {
		if y >= 0 && y < len(mask[x]) {
			if ignoreMask || mask[x][y] {
				mask[x][y] = mask[x][y] && !effectMask
				return ecsManager.AddEntity(entity)
			}
		}
	}
	return 0
}

func addRandom(entity map[ecs.ComponentID]interface{}, componentsList []map[ecs.ComponentID]interface{}) {
	components := componentsList[rand.Intn(len(componentsList))]
	for key, val := range components {
		entity[key] = val
	}
}
