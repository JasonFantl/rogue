package main

import (
	"github.com/jasonfantl/rogue/gui"
)

var x, y = 0, 1

func main() {

	gui.Setup()
	quit := false

	for !quit {
		gui.Show(x, y)
		key, pressed := gui.GetKeyPress()

		if pressed {
			switch key {
			case 'q':
				quit = true
				gui.Quit()
			case 'w':
				y = y - 1
			case 'a':
				x = x - 1
			case 's':
				y = y + 1
			case 'd':
				x = x + 1
			}
		}
	}
}
