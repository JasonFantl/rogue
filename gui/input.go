package gui

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var previouslyPressed []ebiten.Key

var lastPressTime time.Time
var keyDelay int64 = 100     // in Ms
var quickTimeThresh int = 60 // in frames

func GetKeyPress() (ebiten.Key, bool) {

	// currently get most recent key

	pressed := ebiten.Key(0)
	durration := -1

	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			dur := inpututil.KeyPressDuration(k)
			if durration == -1 || dur < durration {

				pressed = k
				durration = dur
			}
		}
	}

	if durration >= 0 { // got a key
		dt := time.Since(lastPressTime).Milliseconds()

		if inpututil.IsKeyJustPressed(pressed) || durration > quickTimeThresh || dt > keyDelay {
			lastPressTime = time.Now()
			return pressed, true
		}
	}

	return 0, false
}
