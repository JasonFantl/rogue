package ecs

type EventID int

type Event struct {
	ID     EventID
	data   interface{}
	entity Entity
}

const (
	MOVE_EVENT EventID = iota
	TRY_MOVE_EVENT
	QUIT_EVENT
	DISPLAY_EVENT
	ERROR_EVENT
	TRY_PICK_UP_EVENT
	PICKED_UP_EVENT
)

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
}

type EventPickedUp struct {
	byWho Entity
}

type EventError struct {
	err string
}
