package gui

import (
	"github.com/gdamore/tcell/v2"
)

func GetKeyPress() (tcell.Key, bool) {
	// Poll event
	ev := s.PollEvent()

	// Process event
	switch ev := ev.(type) {
	case *tcell.EventResize:
		s.Sync()
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyRune {
			return tcell.Key(ev.Rune()), true
		}
		return ev.Key(), true
	}

	// nothing happened
	return '~', false
}
