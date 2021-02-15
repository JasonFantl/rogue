package gui

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

var s, sErr = tcell.NewScreen()
var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
var textStyle = defStyle.Foreground(tcell.ColorDarkGoldenrod)

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
	s.DisableMouse()
}

func DrawText(x, y int, text string) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, textStyle)
	}
}

func DrawTile(x, y int, r rune, style tcell.Style) {
	// _, h := s.Size()
	// invertedY := h - y
	s.SetContent(x, y, r, nil, style)
}

var errorLine = 0

func UpdateErrors(toDisplay string) {
	w, _ := s.Size()
	// this assumes max error msg is 50 chars
	DrawText(w-50, errorLine, toDisplay)
	errorLine++
}

func Clear() {
	s.Clear()
	errorLine = 0
}

func Show() {
	s.Show()
}

func Quit() {
	s.Fini()
}
