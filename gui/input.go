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

	if ev.Ch == 0 {
		return ev.Key, true
	}
	return termbox.Key(ev.Ch), true
}
