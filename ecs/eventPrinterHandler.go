package ecs

import (
	"fmt"

	"github.com/jasonfantl/rogue/gui"
)

type EventPrinterHandler struct{}

func (s *EventPrinterHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	stringifiedEvent := fmt.Sprintf("%T : %v", event.data, event.entity)

	//special case
	if event.ID == ERROR_EVENT {
		stringifiedEvent = fmt.Sprintf("%T : %s : %v", event.data, event.data.(ErrorEvent).err, event.entity)
	}

	// keep in mind this wont display any errors between display frames
	// would have to call gui.Show() on every event to make sure we see it, but that messes with the visuals
	gui.UpdateErrors(stringifiedEvent)

	return returnEvents
}
