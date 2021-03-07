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

	go eventListener()
}

func DrawText(x, y int, text string) {
	x, y = offset(x, y)
	for i, r := range text {
		termbox.SetCell(x+i, y, r, termbox.RGBToAttribute(200, 200, 200), termbox.ColorDefault)
	}
}

func DrawCorner(text string) {
	for i, r := range text {
		termbox.SetCell(i, 0, r, termbox.RGBToAttribute(250, 250, 250), termbox.ColorDefault)
	}
}

func DrawFg(x, y int, r rune, c termbox.Attribute) {
	x, y = offset(x, y)
	termbox.SetFg(x, y, c)
	termbox.SetChar(x, y, r)
}

func DrawBg(x, y int, c termbox.Attribute) {
	x, y = offset(x, y)
	termbox.SetBg(x, y, c)
}

var errorLine = 0

func UpdateErrors(text string) {
	// w, _ := termbox.Size()
	// // this assumes max error msg is 50 chars
	// for i, r := range text {
	// 	termbox.SetCell(w-50+i, errorLine, r, termbox.RGBToAttribute(150, 150, 150), termbox.ColorDefault)
	// }
	// errorLine++
}

func offset(x, y int) (int, int) {
	width, height := termbox.Size()
	return x + width/2, y + height/2
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
