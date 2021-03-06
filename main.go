package main

import (
	"math/rand"
	"time"

	"github.com/jasonfantl/rogue/ecs"
	"github.com/jasonfantl/rogue/gui"
)

func main() {

	// disable to get the same results, good for testing
	rand.Seed(time.Now().UnixNano())

	gui.Setup()
	defer gui.Quit()

	ecsManager := ecs.New()

	/////// HANDLERS /////////////

	inputHandler := ecs.PlayerInputHandler{}
	visionHandler := ecs.VisionHandler{}
	movementHandler := ecs.MoveHandler{}
	monsterHandler := ecs.MonsterHandler{}
	attackHandler := ecs.AttackHandler{}
	deathHandler := ecs.DeathHandler{}
	displayHandler := ecs.DisplayHandler{}
	inventoryHandler := ecs.InventoryHandler{}
	eventPrinter := ecs.EventPrinterHandler{}
	memoryHandler := ecs.MemoryHandler{}
	effectHandler := ecs.EffectsHandler{}
	equippingHandler := ecs.EquippingHandler{}

	// the order that these are added matters
	// they follow this order of execution

	ecsManager.AddEventHandler(&inputHandler)
	ecsManager.AddEventHandler(&monsterHandler)
	ecsManager.AddEventHandler(&attackHandler)
	ecsManager.AddEventHandler(&movementHandler)
	ecsManager.AddEventHandler(&effectHandler)
	ecsManager.AddEventHandler(&equippingHandler)
	ecsManager.AddEventHandler(&inventoryHandler)
	// inventory before death, otherwise we cant drop all of its items
	ecsManager.AddEventHandler(&deathHandler)

	// display stuff, all happens on the same event, no queue
	// vision updates awarness
	// awareness updates memory
	ecsManager.AddEventHandler(&visionHandler)
	ecsManager.AddEventHandler(&memoryHandler)
	ecsManager.AddEventHandler(&displayHandler)
	ecsManager.AddEventHandler(&eventPrinter)

	generateGame(&ecsManager, 100, 100)

	ecsManager.Start()

	for ecsManager.Running() {
		ecsManager.Run()
	}

}
