package gui

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

func Setup() {
	// Initialize screen
	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
	}

	termbox.SetOutputMode(termbox.OutputRGB)
}

func DrawText(x, y int, text string) {
	for i, r := range text {
		termbox.SetCell(x+i, y, r, termbox.RGBToAttribute(200, 200, 200), termbox.ColorDefault)
	}
}

func DrawFg(x, y int, r rune, c termbox.Attribute) {
	// _, h := s.Size()
	// invertedY := h - y
	termbox.SetFg(x, y, c)
	termbox.SetChar(x, y, r)
}

func DrawBg(x, y int, c termbox.Attribute) {
	// _, h := s.Size()
	// invertedY := h - y
	termbox.SetBg(x, y, c)
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
