package ecs

import (
	"github.com/nsf/termbox-go"
)

type Entity uint64

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
	ENTITY_AWARENESS
	ENTITY_MEMORY
	VISION
	VOLUME
	VIOLENT
	OPAQUE
	DISPLAYABLE
	MEMORABLE
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
	IsForeground bool
	Color        termbox.Attribute
	Rune         rune
	Priority     int // reserve 0-99 for memories, 100-199 for displays
}

type PlayerController struct {
	Up, Down, Left, Right, Pickup, Quit termbox.Key
}

type MonsterController struct {
	ActionPriorities []MonsterAction
}

type Vision struct {
	Radius int
}

type EntityAwarness struct {
	AwareOf []Entity
}

// perhaps make more general later. right now its only walls.
// and make it so mulitple things can be remembered on the same tile.
// this will be more complicated then it initially seems.
type EntityMemory struct {
	Memory map[int]map[int]Displayable
}

type Memorable struct {
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

type Opaque struct {
}
