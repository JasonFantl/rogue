package ecs

type EventID int

type Event struct {
	ID     EventID
	data   interface{}
	entity Entity
}

const (
	TIMESTEP EventID = iota
	MOVE_EVENT
	TRY_MOVE_EVENT
	QUIT_EVENT
	DISPLAY_EVENT
	ERROR_EVENT
	TRY_PICK_UP_EVENT
	PICKED_UP_EVENT
)

type EventTimeStep struct{}

type EventMove struct {
	x, y int
}

type EventTryMove struct {
	dx, dy int
}

type EventQuit struct {
}

type EventDisplayTrigger struct {
}

type EventTryPickUp struct {
	oneItem bool
	what    Entity
}

type EventPickedUp struct {
	byWho Entity
}

type EventError struct {
	err string
}
