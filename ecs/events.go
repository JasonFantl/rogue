package ecs

import "github.com/jasonfantl/rogue/gui"

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
	PICKED_UP
	TRY_DROP
	DROPED
	TRY_EQUIP_WEAPON
	TRY_EQUIP_ARMOR
	EQUIPPED
	TRY_CONSUME
	TRY_UNLOCK
	UNLOCKED
	TRY_LAUNCH
	TRY_TRADE
	TRADED
	CONSUMED
	TRY_ATTACK
	DAMAGED
	DIED
	DEBUG_EVENT
	WAKEUP_HANDLERS
	SETTING_CHANGE
)

type TimeStep struct{}

type KeyPressed struct {
	key gui.Key
}

type Moved struct {
	fromX, fromY int
	toX, toY     int
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

type PickedUp struct {
	byWho Entity
}

type TryDrop struct {
	what Entity
}

type Dropped struct {
	byWho Entity
}

type TryEquip struct {
	what Entity
}

type Equipped struct {
	byWho Entity
}

type TryUnlock struct {
	what Entity
}

type Unlocked struct {
}

type TryAttack struct {
	who Entity
}

type TryLaunch struct {
	what   Entity
	dx, dy int
}

type TryTrade struct {
	who                  Entity
	offering, requesting map[Entity]bool
}

type Traded struct {
	taken, given []Entity
}

type Damaged struct {
}

type TryConsume struct {
	what Entity
}

type Consumed struct {
	byWho Entity
}

type Died struct {
}

type DebugEvent struct {
	err string
}

type WakeupHandlers struct {
}

type SettingChange struct {
	field string
	value int
}
