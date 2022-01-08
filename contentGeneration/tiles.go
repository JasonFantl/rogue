package contentGeneration

import (
	"github.com/jasonfantl/rogue/ecs"
	"github.com/jasonfantl/rogue/gui"
)

func water(x, y int) map[ecs.ComponentID]interface{} {
	return floor(x, y, gui.GetSprite(gui.WATER))
}

func treeTrunk(x, y int) map[ecs.ComponentID]interface{} {
	return wall(x, y, gui.GetSprite(gui.TREE_TRUNK))
}

func sandFloor(x, y int) map[ecs.ComponentID]interface{} {
	return floor(x, y, gui.GetSprite(gui.SAND_FLOOR))
}

func grassFloor(x, y int) map[ecs.ComponentID]interface{} {
	return floor(x, y, gui.GetSprite(gui.GRASS_FLOOR))
}

func dirtFloor(x, y int) map[ecs.ComponentID]interface{} {
	return floor(x, y, gui.GetSprite(gui.DIRT_FLOOR))
}

func stoneFloor(x, y int) map[ecs.ComponentID]interface{} {
	return floor(x, y, gui.GetSprite(gui.STONE_FLOOR))
}

func stoneWall(x, y int) map[ecs.ComponentID]interface{} {
	return wall(x, y, gui.GetSprite(gui.STONE_WALL))
}

func leaf(x, y int) map[ecs.ComponentID]interface{} {
	return map[ecs.ComponentID]interface{}{
		ecs.POSITION:    ecs.Position{X: x, Y: y},
		ecs.DISPLAYABLE: ecs.Displayable{Sprite: gui.GetSprite(gui.LEAF)},
		ecs.MEMORABLE:   ecs.Memorable{},
	}
}

func floor(x, y int, sprite gui.Sprite) map[ecs.ComponentID]interface{} {
	return map[ecs.ComponentID]interface{}{
		ecs.POSITION:    ecs.Position{X: x, Y: y},
		ecs.DISPLAYABLE: ecs.Displayable{Sprite: sprite},
		ecs.MEMORABLE:   ecs.Memorable{},
	}
}

func wall(x, y int, sprite gui.Sprite) map[ecs.ComponentID]interface{} {
	return map[ecs.ComponentID]interface{}{
		ecs.POSITION:    ecs.Position{X: x, Y: y},
		ecs.DISPLAYABLE: ecs.Displayable{Sprite: sprite},
		ecs.MEMORABLE:   ecs.Memorable{},
		ecs.VOLUME:      ecs.Volume{},
		ecs.OPAQUE:      ecs.Opaque{},
	}
}
