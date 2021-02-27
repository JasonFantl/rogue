package ecs

import "github.com/nsf/termbox-go"

type EventID int

type Event struct {
	ID     EventID
	data   interface{}
	entity Entity
}

const (
	TIMESTEP EventID = iota
	KEY_PRESSED
	MOVED
	TRY_MOVE
	QUIT
	DISPLAY
	TRY_PICK_UP
	PLAYER_TRY_PICK_UP // need to get rid of later, not necessary
	PICKED_UP
	TRY_ATTACK
	DAMAGED
	DIED
	ERROR_EVENT
)

type TimeStep struct{}

type KeyPressed struct {
	key termbox.Key
}

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
	what Entity
}
type PlayerTryPickUp struct {
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
