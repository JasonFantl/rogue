package gui

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

var s, sErr = tcell.NewScreen()
var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

func Setup() {
	// Initialize screen
	if sErr != nil {
		log.Fatalf("%+v", sErr)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default text style
	s.SetStyle(defStyle)
}

func drawText(x, y int, text string) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, defStyle)
	}
}

func Quit() {
	s.Fini()
}

func Show(toDisplay [][]rune) {
	s.Clear()
	drawText(0, 0, "Use arrow keys to move and esc to quit")
	for x := range toDisplay {
		for y := range toDisplay[x] {
			invertedY := cap(toDisplay) - y
			s.SetContent(x, invertedY, toDisplay[x][y], nil, defStyle)
		}
	}

	s.Show()
}
