package main

import (
	"github.com/jasonfantl/rogue/ecs"
	"github.com/jasonfantl/rogue/gui"
)

func main() {

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

	// the order that these are added matters
	// they follow this order of execution

	ecsManager.AddEventHandler(&inputHandler)
	ecsManager.AddEventHandler(&monsterHandler)
	ecsManager.AddEventHandler(&attackHandler)
	ecsManager.AddEventHandler(&deathHandler)
	ecsManager.AddEventHandler(&movementHandler)
	ecsManager.AddEventHandler(&inventoryHandler)

	// display stuff, all happens on the same event, no queue
	// vision updates awarness
	// awareness updates memory
	ecsManager.AddEventHandler(&visionHandler)
	ecsManager.AddEventHandler(&memoryHandler)
	ecsManager.AddEventHandler(&displayHandler)
	ecsManager.AddEventHandler(&eventPrinter)

	//////// ENTITIES //////////////
	addPlayer(&ecsManager, 2, 2)

	generateRooms(&ecsManager, 20, 20)

	ecsManager.Start()

	for ecsManager.Running() {
		ecsManager.Run()
	}

}
