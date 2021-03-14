package ecs

import (
	"fmt"
	"time"

	"github.com/jasonfantl/rogue/gui"
)

type EventPrinterHandler struct {
	debugString   string
	lastEventCall time.Time
}

func (h *EventPrinterHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {
	duration := time.Since(h.lastEventCall)
	h.lastEventCall = time.Now()

	stringifiedEvent := "pressed"
	//special case
	if event.ID == DEBUG_EVENT {
		stringifiedEvent = fmt.Sprintf("%-20s : %-6v %s", event.data.(DebugEvent).err, event.entity, duration)
	} else {
		stringifiedEvent = fmt.Sprintf("%-20T : %-6v %s", event.data, event.entity, duration)
	}

	h.debugString += stringifiedEvent + "\n"

	if event.ID == DISPLAY {
		gui.Debug(h.debugString)

		// we start on a new line to leave space for frame rate in screen
		h.debugString = "\n"
	}

	return returnEvents
}
