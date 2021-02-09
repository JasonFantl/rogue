package ecs

import (
	"github.com/gdamore/tcell/v2"
)

type ComponentID uint64

const (
	POSITION ComponentID = iota
	DESIRED_MOVE
	DISPLAY
	CONTROLLER
	QUIT_FLAG
	BLOCKABLE_TAG
)

type Component struct {
	ID   ComponentID
	Data interface{}
}

type Position struct {
	X, Y int
}

type DesiredMove struct {
	X, Y int
}

type Display struct {
	Character rune
	Priority  int
}

type Controller struct {
	Up, Down, Left, Right, Quit tcell.Key
}

type QuitFlag struct {
	HasQuit bool
}

type BlockableTag struct {
}
