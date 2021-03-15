package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jasonfantl/rogue/gui"
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
	REACTIONS
	ENTITY_AWARENESS
	ENTITY_MEMORY
	VISION
	VOLUME
	FIGHTER
	DAMAGE
	DAMAGE_RESISTANCE
	OPAQUE
	DISPLAYABLE
	MEMORABLE
	STASHABLE
	USER
	BRAIN
	CONSUMABLE
	PROJECTILE
	LOCKABLE
)

type Position struct {
	X, Y int
}

type Displayable struct {
	Sprite gui.Sprite
}

type User struct {
	Controlling                                                    Entity
	UpKey, DownKey, LeftKey, RightKey, PickupKey, MenuKey, QuitKey ebiten.Key
	Menu                                                           Menu
}

type Brain struct {
	Desires []DesiredAction
}

type Vision struct {
	Radius int
}

type EntityAwarness struct {
	AwareOf PositionLookup
}

// perhaps make more general later. right now its only walls.
// and make it so mulitple things can be remembered on the same tile.
// this will be more complicated then it initially seems.
type EntityMemory struct {
	Memory map[int]map[int][]Displayable
}

type Memorable struct {
}

type Volume struct {
}

type Inventory struct {
	Items map[Entity]bool
}

type Information struct {
	Name, Details string
}

type Stashable struct {
	Stashed bool
}

type Consumable struct {
}

type Health struct {
	Max, Current int
}

type Fighter struct {
	Strength int
	Weapon   Entity
	Armor    Entity
}

type Damage struct {
	Amount int
}

type DamageResistance struct {
	Amount int
}

type Opaque struct {
}

type Reactions struct {
	Reactions []Reaction
}

type Lockable struct {
	Key                Entity
	Locked             bool
	LockedComponents   []Component
	UnlockedComponents []Component
}

type Projectile struct {
	MaxDistance          int
	currentDx, currentDy int
	goalDx, goalDy       int
}
