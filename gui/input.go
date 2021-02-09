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
		return ev.Key(), true
		// if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
		// 	// quit()
		// 	return 'q', true
		// } else if ev.Key() == tcell.KeyCtrlL {
		// 	s.Sync()
		// } else {
		// 	return ev.Rune(), true
		// }
	}

	// nothing happened
	return '~', false
}
