package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/jasonfantl/rogue/ecs"
	"github.com/jasonfantl/rogue/gui"
)

func main() {

	gui.Setup()
	defer gui.Quit()

	ecsManager := ecs.New()

	playerControlSystem := ecs.SystemControl{}
	ecsManager.AddSystem(&playerControlSystem)

	movementSystem := ecs.SystemMove{}
	ecsManager.AddSystem(&movementSystem)

	displaySystem := ecs.SystemDisplay{}
	displaySystem.SetSize(10)
	ecsManager.AddSystem(&displaySystem)

	player := []ecs.Component{
		{
			ID: ecs.DISPLAY,
			Data: ecs.Display{
				Character: '@',
				Priority:  999,
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
			ID: ecs.DESIRED_MOVE,
			Data: ecs.DesiredMove{
				X: 0,
				Y: 0,
			},
		},
		{
			ID: ecs.CONTROLLER,
			Data: ecs.Controller{
				Up:    tcell.KeyUp,
				Down:  tcell.KeyDown,
				Left:  tcell.KeyLeft,
				Right: tcell.KeyRight,
				Quit:  tcell.KeyEsc,
			},
		},
		{
			ID: ecs.QUIT_FLAG,
			Data: ecs.QuitFlag{
				HasQuit: false,
			},
		},
	}

	block := []ecs.Component{
		{
			ID: ecs.DISPLAY,
			Data: ecs.Display{
				Character: '#',
				Priority:  999,
			},
		},
		{
			ID: ecs.POSITION,
			Data: ecs.Position{
				X: 2,
				Y: 1,
			},
		},
		{
			ID:   ecs.BLOCKABLE_TAG,
			Data: ecs.BlockableTag{},
		},
	}

	notBlock := []ecs.Component{
		{
			ID: ecs.DISPLAY,
			Data: ecs.Display{
				Character: '.',
				Priority:  1,
			},
		},
		{
			ID: ecs.POSITION,
			Data: ecs.Position{
				X: 2,
				Y: 3,
			},
		},
	}

	playerID := ecsManager.AddEntity(player)
	ecsManager.AddEntity(block)
	ecsManager.AddEntity(notBlock)

	hasQuit := func() bool {
		data, ok := ecsManager.Lookup(playerID, ecs.QUIT_FLAG)
		if ok {
			return data.(ecs.QuitFlag).HasQuit
		}
		return true
	}

	fmt.Println("starting loop")

	for !hasQuit() {
		ecsManager.Run()
	}

}
