package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/jasonfantl/rogue/ecs"
	"github.com/jasonfantl/rogue/gui"
)

func main() {

	gui.Setup()

	ecsManager := ecs.New()

	displaySystem := ecs.SystemDisplay{}
	displaySystem.SetSize(10)
	ecsManager.AddSystem(&displaySystem)

	playerControlSystem := ecs.SystemControl{}
	ecsManager.AddSystem(&playerControlSystem)

	player := []ecs.Component{
		{
			ID: ecs.DISPLAY,
			Data: ecs.Display{
				Character: "@",
			},
		},
		{
			ID: ecs.POSITION,
			Data: ecs.Position{
				X: 1,
				Y: 1,
			},
		},
		{
			ID: ecs.CONTROLLER,
			Data: ecs.Controller{
				Up:      tcell.KeyUp,
				Down:    tcell.KeyDown,
				Left:    tcell.KeyLeft,
				Right:   tcell.KeyRight,
				Quit:    tcell.KeyEsc,
				HasQuit: false,
			},
		},
	}

	playerID := ecsManager.AddEntity(player)

	hasQuit := func() bool {
		data, ok := ecsManager.Lookup(playerID, ecs.CONTROLLER)
		if ok {
			return data.(ecs.Controller).HasQuit
		}
		return true
	}

	fmt.Println("starting loop")

	for !hasQuit() {
		ecsManager.Run()
	}
	gui.Quit()
}
