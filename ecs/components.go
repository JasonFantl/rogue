package ecs

import (
	"github.com/gdamore/tcell/v2"
)

type ComponentID uint64

type Component struct {
	ID   ComponentID
	Data interface{}
}

const (
	POSITION ComponentID = iota
	DESIRED_MOVE
	PLAYER_CONTROLLER
	MONSTER_CONTROLLER
	INVENTORY
	INFORMATION
	VOLUME
	DISPLAYABLE
	PICKUPABLE
	DROPABLE
	STASHED_FLAG
)

type Position struct {
	X, Y int
}

type Displayable struct {
	Character rune
	Style     tcell.Style
	Priority  int
}

type PlayerController struct {
	Up, Down, Left, Right, Pickup, Quit tcell.Key
}

type MonsterController struct {
	// what data?
}

type Volume struct {
}

type Inventory struct {
	items []Entity
}

type Information struct {
	Name, Details string
}

type Pickupable struct {
}
type Dropable struct {
}

type StashedFlag struct {
	parent Entity
}
