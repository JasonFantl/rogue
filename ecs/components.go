package ecs

import (
	"github.com/nsf/termbox-go"
)

type ComponentID uint64

type Component struct {
	ID   ComponentID
	Data interface{}
}

const (
	POSITION ComponentID = iota
	DESIRED_MOVE
	HEALTH
	INVENTORY
	INFORMATION
	VOLUME
	VIOLENT
	DISPLAYABLE
	PICKUPABLE
	DROPABLE
	STASHED_FLAG
	PLAYER_CONTROLLER
	MONSTER_CONTROLLER
)

type Position struct {
	X, Y int
}

type Displayable struct {
	Character rune
	// Style     tcell.Style
	Priority int
}

type PlayerController struct {
	Up, Down, Left, Right, Pickup, Quit termbox.Key
}

type MonsterController struct {
	// what data?
}

type Volume struct {
}

type Inventory struct {
	Items []Entity
}

type Information struct {
	Name, Details string
}

type Pickupable struct {
}
type Dropable struct {
}

type StashedFlag struct {
	Parent Entity
}

type Health struct {
	Max, Current int
}

type Violent struct {
	BaseAttackDmg int
}
