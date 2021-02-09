package ecs

import (
	"github.com/gdamore/tcell/v2"
)

type ComponentID uint64

const (
	POSITION ComponentID = iota
	DISPLAY
	CONTROLLER
)

type Component struct {
	ID   ComponentID
	Data interface{}
}

type Position struct {
	X, Y int
}

type Display struct {
	Character string
}

type Controller struct {
	Up, Down, Left, Right, Quit tcell.Key
	HasQuit                     bool
}
