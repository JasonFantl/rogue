package gui

import (
	"github.com/nsf/termbox-go"
)

func GetKeyPress() (termbox.Key, bool) {
	// Poll event
	ev := termbox.PollEvent()

	// Process event
	if ev.Type == termbox.EventError {
		panic(ev.Err)
	}

	return termbox.Key(ev.Ch), true
}
