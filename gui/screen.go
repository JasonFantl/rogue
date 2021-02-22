package gui

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

var defStyle = termbox.ColorDefault

func Setup() {
	// Initialize screen
	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
	}
}

func DrawText(x, y int, text string) {
	for i, r := range text {
		termbox.SetCell(x+i, y, r, defStyle, termbox.ColorMagenta)
	}
}

func DrawTile(x, y int, r rune) {
	// _, h := s.Size()
	// invertedY := h - y
	termbox.SetCell(x, y, r, defStyle, termbox.ColorMagenta)
}

var errorLine = 0

func UpdateErrors(toDisplay string) {
	w, _ := termbox.Size()
	// this assumes max error msg is 50 chars
	DrawText(w-50, errorLine, toDisplay)
	errorLine++
}

func Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	errorLine = 0
}

func Show() {
	termbox.Flush()
}

func Quit() {
	termbox.Close()
}
