package ecs

type EventID int

type Event struct {
	ID     EventID
	data   interface{}
	entity Entity
}

const (
	TIMESTEP EventID = iota
	MOVED
	TRY_MOVE
	QUIT
	DISPLAY
	ERROR_EVENT
	TRY_PICK_UP
	PICKED_UP
)

type TimeStep struct{}

type Moved struct {
	x, y int
}

type TryMove struct {
	dx, dy int
}

type Quit struct {
}

type Display struct {
}

type TryPickUp struct {
	oneItem bool
	what    Entity
}

type PickedUp struct {
	byWho Entity
}

type ErrorEvent struct {
	err string
}
