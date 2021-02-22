package ecs

import "math/rand"

type MonsterAction uint

const (
	NOTHING MonsterAction = iota
	PICKUP
	MOVE
)

type MonsterHandler struct{}

func (mh *MonsterHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// decision making is done at the beginning of the turn
	if event.ID == TIMESTEP {

		// well just have each monster move one at a time
		monsters, monstersExist := m.getComponents(MONSTER_CONTROLLER)
		if monstersExist {

			actionPossibilities := []MonsterAction{NOTHING}

			for monster, monsterData := range monsters {
				monsterComponent := monsterData.(MonsterController)
				// ---------------------------------
				// first find all possible actions
				// ----------------------------------

				// unpack all the components we will need
				positionData, hasPosition := m.getComponent(monster, POSITION)

				if hasPosition {
					positionComponent := positionData.(Position)

					// ------------ MOVE ------------------
					actionPossibilities = append(actionPossibilities, MOVE)

					// ------------ PICK UP ----------------------
					// first check if we are on top of any treasure
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
				case MOVE:
					returnEvents = append(returnEvents, mh.move(m, monster, positionData.(Position))...)
				case PICKUP:
					returnEvents = append(returnEvents, mh.pickup(m, monster, positionData.(Position))...)
				}
			}
		}
	}

	return returnEvents
}

func (mh *MonsterHandler) move(m *Manager, monster Entity, positionComponent Position) (returnEvents []Event) {

	moveX := 0
	moveY := 0

	if rand.Intn(2) == 0 {
		moveY = rand.Intn(3) - 1
	} else {
		moveX = rand.Intn(3) - 1
	}

	treasures, _ := m.getComponents(PICKUPABLE)
	for e := range treasures {
		itemPositionData, hasPosition := m.getComponent(e, POSITION)
		_, hasStashed := m.getComponent(e, STASHED_FLAG)
		// enable this to chase entities carrying items
		hasStashed = false

		if hasPosition && !hasStashed {
			itemPositionComponent := itemPositionData.(Position)
			dx := itemPositionComponent.X - positionComponent.X
			dy := itemPositionComponent.Y - positionComponent.Y

			if dx != 0 || dy != 0 {
				if rand.Intn(2) == 0 && dx*dx+dy*dy < rand.Intn(100) {
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
				}
			}
		}
	}

	if moveX != 0 || moveY != 0 {
		returnEvents = append(returnEvents, Event{TRY_MOVE, TryMove{moveX, moveY}, monster})
	}

	return returnEvents
}

// Needs to control for picking up only one item; right now just takes first treasure from the pile
func (s *MonsterHandler) pickup(m *Manager, monster Entity, positionComponent Position) (returnEvents []Event) {
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
