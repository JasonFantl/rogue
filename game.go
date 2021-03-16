package main

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jasonfantl/rogue/ecs"
	"github.com/jasonfantl/rogue/gui"
)

type Game struct {
	ecsManager ecs.Manager
}

func (g *Game) Start() {
	gui.Setup()

	g.ecsManager = ecs.New()

	screenWidth, screenHeight := gui.Dimensions()
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Rogue")

	// the order that these are added matters
	// they follow this order of execution

	g.ecsManager.AddEventHandler(&ecs.UserHandler{})
	g.ecsManager.AddEventHandler(&ecs.BrainHandler{})
	g.ecsManager.AddEventHandler(&ecs.AttackHandler{})
	g.ecsManager.AddEventHandler(&ecs.MoveHandler{})
	g.ecsManager.AddEventHandler(&ecs.ReactionHandler{})
	g.ecsManager.AddEventHandler(&ecs.LockHandler{})
	g.ecsManager.AddEventHandler(&ecs.ProjectileHandler{})
	g.ecsManager.AddEventHandler(&ecs.EquippingHandler{})
	g.ecsManager.AddEventHandler(&ecs.InventoryHandler{})
	// inventory before death, otherwise we cant drop all of its items
	// should place handler that deletes items very carefully, proabably here at the end
	g.ecsManager.AddEventHandler(&ecs.ConsumptionHandler{})
	g.ecsManager.AddEventHandler(&ecs.DeathHandler{})

	// display stuff, all happens on the same event, no queue
	// vision updates awarness
	// awareness updates memory
	g.ecsManager.AddEventHandler(&ecs.VisionHandler{})
	g.ecsManager.AddEventHandler(&ecs.MemoryHandler{})
	g.ecsManager.AddEventHandler(&ecs.DisplayHandler{})
	g.ecsManager.AddEventHandler(&ecs.EventPrinterHandler{})

	generateGame(&g.ecsManager, 200, 200)

	g.ecsManager.Start()

}

func (g *Game) Update() error {
	if !g.ecsManager.Running() {
		return errors.New("player hit quit")
	}

	g.ecsManager.Run()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	screen.DrawImage(gui.GetImage(), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return gui.Dimensions()
}
