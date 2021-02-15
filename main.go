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

	/////// SYSTEMS //////////////

	keyboardListener := ecs.InputSystem{}

	ecsManager.AddSystem(&keyboardListener)

	/////// HANDLERS /////////////

	movementHandler := ecs.MoveHandler{}
	monsterHandler := ecs.MonsterHandler{}

	displayHandler := ecs.DisplayHandler{}
	inventoryHandler := ecs.InventoryHandler{}
	eventPrinter := ecs.EventPrinterHandler{}

	ecsManager.AddEventHandler(&movementHandler)
	ecsManager.AddEventHandler(&monsterHandler)
	ecsManager.AddEventHandler(&displayHandler)
	ecsManager.AddEventHandler(&inventoryHandler)
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
