package ecs

import (
	"math/rand"
)

type MonsterAction uint

const (
	NOTHING MonsterAction = iota
	PICKUP
	RANDOM_MOVE
	TREASURE_MOVE
)

type MonsterHandler struct{}

func (h *MonsterHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// decision making is done at the beginning of the turn
	if event.ID == TIMESTEP {

		// well just have each monster move one at a time
		monsters, _ := m.getComponents(MONSTER_CONTROLLER)

		for monster, monsterData := range monsters {
			monsterComponent := monsterData.(MonsterController)
			// ---------------------------------
			// first find all possible actions
			// ----------------------------------
			actionPossibilities := []MonsterAction{NOTHING}

			// unpack all the components we will need
			positionData, hasPosition := m.getComponent(monster, POSITION)
			awarnessData, hasAwarness := m.getComponent(monster, ENTITY_AWARENESS)

			// ------------ RANDOM WALK ------------------
			if hasPosition {
				actionPossibilities = append(actionPossibilities, RANDOM_MOVE)
			}

			// ------------ MOVE TO TREASURE --------------
			if hasPosition && hasAwarness {
				awarnessComponent := awarnessData.(EntityAwarness)

				for _, item := range awarnessComponent.AwareOf {
					if h.isTreasure(m, monster, item) {
						actionPossibilities = append(actionPossibilities, TREASURE_MOVE)

						break
					}
				}
			}

			// ------------ PICK UP ----------------------
			if hasPosition {
				positionComponent := positionData.(Position)

				entites := m.getEntitiesFromPos(positionComponent.X, positionComponent.Y)
				for _, entity := range entites {
					_, hasPickUpAble := m.getComponent(entity, PICKUPABLE)
					_, hasStashed := m.getComponent(entity, STASHED_FLAG)
					isTreasure := hasPickUpAble && !hasStashed
					if isTreasure {
						actionPossibilities = append(actionPossibilities, PICKUP)
						break // dont need to check anymore
					}
				}
			}

			// ---------------------------------
			// second, choose an action
			// ----------------------------------

			decidedAction := actionPossibilities[0]
			// by looping from begginning, we make the highest priiority start at 0th index
			found := false
			for _, desiredAction := range monsterComponent.ActionPriorities {
				for _, possibleAction := range actionPossibilities {
					if desiredAction == possibleAction {
						decidedAction = desiredAction
						found = true
						break
					}
				}
				if found {
					break
				}
			}

			// ---------------------------------
			// third, execute ation
			// ----------------------------------

			switch decidedAction {
			case RANDOM_MOVE:
				returnEvents = append(returnEvents, h.moveRandom(m, monster)...)
			case TREASURE_MOVE:
				// make it 50% chance we go after treasure
				if rand.Intn(2) == 0 {
					returnEvents = append(returnEvents, h.moveToTreasure(m, monster)...)
				} else {
					returnEvents = append(returnEvents, h.moveRandom(m, monster)...)
				}
			case PICKUP:
				returnEvents = append(returnEvents, h.pickup(m, monster)...)
			}
		}
	}

	return returnEvents
}

func (h *MonsterHandler) moveRandom(m *Manager, monster Entity) (returnEvents []Event) {
	moveX := 0
	moveY := 0

	diceRoll := rand.Intn(3)
	if diceRoll == 0 {
		moveY = rand.Intn(3) - 1
	} else if diceRoll == 1 {
		moveX = rand.Intn(3) - 1
	}

	if moveX != 0 || moveY != 0 {
		returnEvents = append(returnEvents, Event{TRY_MOVE, TryMove{moveX, moveY}, monster})
	}

	return returnEvents
}

func (h *MonsterHandler) moveToTreasure(m *Manager, monster Entity) (returnEvents []Event) {

	awarnessData, _ := m.getComponent(monster, ENTITY_AWARENESS)
	positionData, _ := m.getComponent(monster, POSITION)

	awarnessComponent := awarnessData.(EntityAwarness)
	positionComponent := positionData.(Position)

	// location of nearest treasure
	dx := 999
	dy := 999

	for _, item := range awarnessComponent.AwareOf {
		if h.isTreasure(m, monster, item) {

			itemPositionData, _ := m.getComponent(item, POSITION)
			itemPositionComponent := itemPositionData.(Position)

			newDx := itemPositionComponent.X - positionComponent.X
			newDy := itemPositionComponent.Y - positionComponent.Y

			oldDis := dx*dx + dy*dy
			newDis := newDx*newDx + newDy*newDy

			if newDis < oldDis {
				dx = newDx
				dy = newDy
			}
		}
	}

	moveX := 0
	moveY := 0

	tiebreaker := rand.Intn(2)
	if dx*dx+tiebreaker > dy*dy {
		moveY = 0
		if dx > 0 {
			moveX = 1
		} else {
			moveX = -1
		}
	} else {
		moveX = 0
		if dy > 0 {
			moveY = 1
		} else {
			moveY = -1
		}
	}

	if moveX != 0 || moveY != 0 {
		returnEvents = append(returnEvents, Event{TRY_MOVE, TryMove{moveX, moveY}, monster})
	}

	return returnEvents
}

// Needs to control for picking up only one item; right now just takes first treasure from the pile
func (s *MonsterHandler) pickup(m *Manager, monster Entity) (returnEvents []Event) {
	positionData, _ := m.getComponent(monster, POSITION)
	positionComponent := positionData.(Position)

	entities := m.getEntitiesFromPos(positionComponent.X, positionComponent.Y)
	for _, entity := range entities {
		_, hasPickUpAble := m.getComponent(entity, PICKUPABLE)
		_, hasStashed := m.getComponent(entity, STASHED_FLAG)
		isTreasure := hasPickUpAble && !hasStashed
		if isTreasure {
			returnEvents = append(returnEvents, Event{TRY_PICK_UP, TryPickUp{entity}, monster})
			break // dont need to check anymore
		}
	}

	return returnEvents
}

// Needs to control for picking up only one item; right now just takes first treasure from the pile
func (s *MonsterHandler) isTreasure(m *Manager, monster, item Entity) bool {
	_, hasPickupable := m.getComponent(item, PICKUPABLE)
	stashedData, hasStashed := m.getComponent(item, STASHED_FLAG)

	if hasPickupable && hasStashed {
		stashedComponent := stashedData.(StashedFlag)
		if stashedComponent.Parent != monster {
			return true
		}
		return false
	}

	return hasPickupable
}
