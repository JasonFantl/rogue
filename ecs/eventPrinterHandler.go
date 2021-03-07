package ecs

import (
	"fmt"
	"time"

	"github.com/jasonfantl/rogue/gui"
)

type EventPrinterHandler struct {
	lasteEventCall time.Time
}

func (h *EventPrinterHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	stringifiedEvent := ""
	//special case
	if event.ID == DEBUG_EVENT {
		stringifiedEvent = fmt.Sprintf("%T : %s : %v", event.data, event.data.(DebugEvent).err, event.entity)
	} else {
		duration := time.Since(h.lasteEventCall)
		stringifiedEvent = fmt.Sprintf("%-20T %-6v %s", event.data, event.entity, duration)
		h.lasteEventCall = time.Now()
	}

	// keep in mind this wont display any errors between display frames
	// would have to call gui.Show() on every event to make sure we see it, but that messes with the visuals
	gui.UpdateErrors(stringifiedEvent)

	return returnEvents
}
