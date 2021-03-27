package ecs

import (
	"fmt"
	"time"

	"github.com/jasonfantl/rogue/gui"
)

type EventHandler interface {
	handleEvent(*Manager, Event) []Event
}

type EventManager struct {
	eventHandlers []EventHandler
}

func NewEventHandler() EventManager {
	newManager := EventManager{}
	newManager.eventHandlers = make([]EventHandler, 0)

	return newManager
}

func (e *EventManager) addEventHandler(eventHandler EventHandler) {
	e.eventHandlers = append(e.eventHandlers, eventHandler)
}

func (e *EventManager) sendEvents(m *Manager, events []Event) {

	timer := time.Now()
	sentDisplay := false

	// queue style event handling
	for len(events) > 0 {
		sendingEvent := events[0] // pop
		events = events[1:]       // dequeue

		for _, eventHandler := range e.eventHandlers {
			respondingEvents := eventHandler.handleEvent(m, sendingEvent)
			events = append(events, respondingEvents...)
		}

		// display stuff
		if !sentDisplay && len(events) == 0 {
			sentDisplay = true
			events = append(events, Event{DISPLAY, Display{}, m.user.Controlling})
		}
		if sendingEvent.ID == QUIT {
			fmt.Printf("Event quitting")

			m.running = false
		}
	}

	gui.Debug(time.Since(timer).String())
}
