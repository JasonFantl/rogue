package contentGeneration

import (
	"math/rand"

	"github.com/jasonfantl/rogue/ecs"
	"github.com/jasonfantl/rogue/gui"
)

func addPlayer(ecsManager *ecs.Manager, x, y int) {

	player := []ecs.Component{
		{ecs.BRAIN, ecs.Brain{
			[]ecs.DesiredAction{},
		}},
		{ecs.POSITION, ecs.Position{x, y}},
		{ecs.DISPLAYABLE, ecs.Displayable{gui.GetSprite(gui.PLAYER)}},
		{ecs.ENTITY_AWARENESS, ecs.EntityAwarness{}},
		{ecs.ENTITY_MEMORY, ecs.EntityMemory{}},
		{ecs.VISION, ecs.Vision{50}},
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
			UpKey:       gui.UP,
			DownKey:     gui.DOWN,
			LeftKey:     gui.LEFT,
			RightKey:    gui.RIGHT,
			ActionKey:   gui.ACTION,
			MenuKey:     gui.MENU,
			QuitKey:     gui.QUIT,
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
