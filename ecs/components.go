package ecs

import (
	"github.com/gdamore/tcell/v2"
)

type ComponentID uint64

const (
	POSITION ComponentID = iota
	DESIRED_MOVE
	DISPLAY
	PLAYER_CONTROLLER
	BLOCKABLE
	INVENTORY
	INFORMATION
	PICKUPABLE
	DROPABLE
	STASHED
)

type Component struct {
	ID   ComponentID
	Data interface{}
}

type Position struct {
	X, Y int
}

type Display struct {
	Character rune
	Style     tcell.Style
	Priority  int
}

type PlayerController struct {
	Up, Down, Left, Right, Pickup, Quit tcell.Key
}

type Blockable struct {
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

type Stashed struct {
	parent Entity
}
