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
	TRY_PICK_UP
	PICKED_UP
	TRY_ATTACK
	DAMAGED
	DIED
	ERROR_EVENT
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

type TryAttack struct {
	who Entity
	dmg int
}

type Damaged struct {
}

type Died struct {
}

type ErrorEvent struct {
	err string
}
