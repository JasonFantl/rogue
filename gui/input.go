package gui

import (
	"github.com/nsf/termbox-go"
)

var ev termbox.Event
var newEvent = false

func GetKeyPress() (termbox.Key, bool) {

	if newEvent {
		newEvent = false
		if ev.Type == termbox.EventError {
			panic(ev.Err)
		}

		if ev.Ch == 0 {
			return ev.Key, true
		}
		return termbox.Key(ev.Ch), true
	}
	return termbox.Key(0), false
}

func eventListener() {
	for true {
		ev = termbox.PollEvent()
		newEvent = true
	}
}
