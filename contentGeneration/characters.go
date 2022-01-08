package contentGeneration

import (
	"github.com/jasonfantl/rogue/ecs"
	"github.com/jasonfantl/rogue/gui"
)

func addPlayer(ecsManager *ecs.Manager, m *ecs.Manager, x, y int) {

	player := map[ecs.ComponentID]interface{}{
		ecs.BRAIN:            ecs.Brain{Desires: []ecs.DesiredAction{}},
		ecs.POSITION:         ecs.Position{X: x, Y: y},
		ecs.DISPLAYABLE:      ecs.Displayable{Sprite: gui.GetSprite(gui.PLAYER)},
		ecs.ENTITY_AWARENESS: ecs.EntityAwareness{},
		ecs.ENTITY_MEMORY:    ecs.EntityMemory{},
		ecs.VISION:           ecs.Vision{Radius: 12},
		ecs.INVENTORY:        ecs.Inventory{},
		ecs.INFORMATION:      ecs.Information{Name: "Player", Details: "the hero of our story"},
		ecs.VOLUME:           ecs.Volume{},
		ecs.FIGHTER:          ecs.Fighter{Strength: 10, Weapon: 0, Armor: 0},
		ecs.DAMAGE:           ecs.Damage{Amount: 1},
		ecs.HEALTH:           ecs.Health{Max: 100, Current: 90},
	}

	playerID := ecsManager.AddEntity(player)

	user := ecs.User{
		Controlling: playerID,
		UpKey:       gui.UP,
		DownKey:     gui.DOWN,
		LeftKey:     gui.LEFT,
		RightKey:    gui.RIGHT,
		ActionKey:   gui.ACTION,
		MenuKey:     gui.MENU,
		QuitKey:     gui.QUIT,
		Menu:        ecs.Menu{},
	}

	ecsManager.SetUser(user)
}

func addMonster(ecsManager *ecs.Manager, x, y int) {

	monster := map[ecs.ComponentID]interface{}{
		ecs.POSITION: ecs.Position{X: x, Y: y},
		ecs.BRAIN: ecs.Brain{
			Desires: []ecs.DesiredAction{ecs.PICKUP, ecs.TREASURE_MOVE, ecs.RANDOM_MOVE, ecs.DO_NOTHING},
		},
		ecs.ENTITY_AWARENESS: ecs.EntityAwareness{},
		ecs.VISION:           ecs.Vision{Radius: 10},
		ecs.VOLUME:           ecs.Volume{},
		ecs.FIGHTER:          ecs.Fighter{Strength: 10, Weapon: 0, Armor: 0},
		ecs.DAMAGE:           ecs.Damage{Amount: 3},
		ecs.HEALTH:           ecs.Health{Max: 50, Current: 50},
	}

	monsterInfos := []map[ecs.ComponentID]interface{}{
		{
			ecs.DISPLAYABLE: ecs.Displayable{Sprite: gui.GetSprite(gui.MONSTER3)},
			ecs.INFORMATION: ecs.Information{Name: "Monster", Details: "generic"},
		},
		{
			ecs.DISPLAYABLE: ecs.Displayable{Sprite: gui.GetSprite(gui.MONSTER2)},
			ecs.INVENTORY:   ecs.Inventory{},
			ecs.INFORMATION: ecs.Information{Name: "Ogre", Details: "Big and scary"},
		},
		{
			ecs.DISPLAYABLE: ecs.Displayable{Sprite: gui.GetSprite(gui.MONSTER1)},
			ecs.INVENTORY:   ecs.Inventory{},
			ecs.INFORMATION: ecs.Information{Name: "Goblin", Details: "green and scrawny, sill scary though"},
		},
	}

	addRandom(monster, monsterInfos)

	ecsManager.AddEntity(monster)
}

func addTownsMember(ecsManager *ecs.Manager, x, y int) {

	controllerComponent := ecs.Brain{
		Desires: []ecs.DesiredAction{ecs.PICKUP, ecs.RANDOM_MOVE, ecs.DO_NOTHING},
	}

	positionComponent := ecs.Position{X: x, Y: y}
	inventoryComponent := ecs.Inventory{}
	volumeComponent := ecs.Volume{}
	damageComponent := ecs.Damage{Amount: 3}
	healthComponent := ecs.Health{Max: 50, Current: 50}
	visionComponent := ecs.Vision{Radius: 10}
	awarenessComponent := ecs.EntityAwareness{}

	townsMember := map[ecs.ComponentID]interface{}{
		ecs.POSITION:         positionComponent,
		ecs.BRAIN:            controllerComponent,
		ecs.ENTITY_AWARENESS: awarenessComponent,
		ecs.VISION:           visionComponent,
		ecs.INVENTORY:        inventoryComponent,
		ecs.VOLUME:           volumeComponent,
		ecs.DAMAGE:           damageComponent,
		ecs.HEALTH:           healthComponent,
		ecs.DISPLAYABLE:      ecs.Displayable{Sprite: gui.GetSprite(gui.PLAYER)},
		ecs.INFORMATION:      ecs.Information{Name: "Human", Details: "lives in town"},
	}

	ecsManager.AddEntity(townsMember)
}
