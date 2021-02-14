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
	displayHandler := ecs.DisplayHandler{}
	inventoryHandler := ecs.InventoryHandler{}
	eventPrinter := ecs.EventPrinterHandler{}

	ecsManager.AddEventHandler(&movementHandler)
	ecsManager.AddEventHandler(&displayHandler)
	ecsManager.AddEventHandler(&inventoryHandler)
	ecsManager.AddEventHandler(&eventPrinter)

	//////// ENTITIES //////////////
	addPlayer(&ecsManager, 2, 2)
	generateRooms(&ecsManager, 50, 20)

	fmt.Println("starting loop")

	ecsManager.Start()

	for !ecsManager.HasQuit() {
		ecsManager.Run()
	}

}
