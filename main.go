package main

import (
	"fmt"

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

	// the order that these are added matters
	// they follow this order of execution
	ecsManager.AddEventHandler(&inputHandler)
	ecsManager.AddEventHandler(&monsterHandler)
	ecsManager.AddEventHandler(&attackHandler)
	ecsManager.AddEventHandler(&deathHandler)
	ecsManager.AddEventHandler(&movementHandler)
	ecsManager.AddEventHandler(&inventoryHandler)
	ecsManager.AddEventHandler(&visionHandler)
	ecsManager.AddEventHandler(&displayHandler)
	ecsManager.AddEventHandler(&eventPrinter)

	//////// ENTITIES //////////////
	addPlayer(&ecsManager, 2, 2)
	addMonster(&ecsManager, 50, 30)

	generateRooms(&ecsManager, 90, 40)

	fmt.Println("starting loop")

	ecsManager.Start()

	for ecsManager.Running() {
		ecsManager.Run()
	}

}
