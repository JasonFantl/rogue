package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	// disable to get the same results, good for testing
	rand.Seed(time.Now().UnixNano())

	game := &Game{}

	game.Start()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
