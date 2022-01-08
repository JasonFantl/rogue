package ecs

import (
	"math/rand"
)

type DesiredAction uint

const (
	DO_NOTHING DesiredAction = iota
	PICKUP
	RANDOM_MOVE
	TREASURE_MOVE
)

type BrainHandler struct{}

func (h *BrainHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == TIMESTEP {

		brains := m.getEntitiesWithComponent(BRAIN)

		for brainEntity := range brains {
			brainData, _ := m.getComponent(brainEntity, BRAIN)
			brainComponent := brainData.(Brain)

			actionPossibilities := []DesiredAction{DO_NOTHING}

			// unpack all the components we will need
			positionData, hasPosition := m.getComponent(brainEntity, POSITION)
			awarenessData, hasAwarness := m.getComponent(brainEntity, ENTITY_AWARENESS)
			_, hasInventory := m.getComponent(brainEntity, INVENTORY)

			// ------------ RANDOM WALK ------------------
			if hasPosition {
				actionPossibilities = append(actionPossibilities, RANDOM_MOVE)
			}

			// ------------ MOVE TO TREASURE --------------
			if hasPosition && hasAwarness {
				awarenessComponent := awarenessData.(EntityAwareness)
				positionComponent := positionData.(Position)
				seeTreasure := func() bool {
					for itemPos, items := range awarenessComponent.AwareOf {
						dx := positionComponent.X - itemPos.X
						dy := positionComponent.Y - itemPos.Y
						if dx != 0 || dy != 0 {
							for item := range items {
								if isTreasure(m, item) {
									return true
								}
							}
						}
					}
					return false
				}
				if seeTreasure() {
					actionPossibilities = append(actionPossibilities, TREASURE_MOVE)
				}
			}

			// ------------ PICK UP ----------------------
			if hasPosition && hasInventory {
				positionComponent := positionData.(Position)

				entities := m.getEntitiesAtPosition(positionComponent)
				for item := range entities {
					if isStashableTreasure(m, item) {
						actionPossibilities = append(actionPossibilities, PICKUP)
						break // don't need to check anymore
					}
				}
			}

			// ---------------------------------
			// second, choose an action
			// ----------------------------------

			decidedAction := actionPossibilities[0]
			// by looping from begginning, we make the highest priority start at 0th index
			found := false
			for _, desiredAction := range brainComponent.Desires {
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
			// third, execute action
			// ----------------------------------

			switch decidedAction {
			case RANDOM_MOVE:
				returnEvents = append(returnEvents, h.moveRandom(m, brainEntity)...)
			case TREASURE_MOVE:
				// make it 50% chance we go after treasure
				if rand.Intn(2) == 0 {
					returnEvents = append(returnEvents, h.moveToTreasure(m, brainEntity)...)
				} else {
					returnEvents = append(returnEvents, h.moveRandom(m, brainEntity)...)
				}
			case PICKUP:
				returnEvents = append(returnEvents, h.pickup(m, brainEntity)...)
			}
		}
	}

	return returnEvents
}

func (h *BrainHandler) moveRandom(m *Manager, brain Entity) (returnEvents []Event) {
	moveX := 0
	moveY := 0

	diceRoll := rand.Intn(3)
	if diceRoll == 0 {
		moveY = rand.Intn(3) - 1
	} else if diceRoll == 1 {
		moveX = rand.Intn(3) - 1
	}

	if moveX != 0 || moveY != 0 {
		returnEvents = append(returnEvents, Event{TRY_MOVE, TryMove{moveX, moveY}, brain})
	}

	return returnEvents
}

func (h *BrainHandler) moveToTreasure(m *Manager, brain Entity) (returnEvents []Event) {

	awarenessData, _ := m.getComponent(brain, ENTITY_AWARENESS)
	positionData, _ := m.getComponent(brain, POSITION)

	awarenessComponent := awarenessData.(EntityAwareness)
	positionComponent := positionData.(Position)

	// location of nearest treasure
	dx := 999
	dy := 999

	for itemPos, items := range awarenessComponent.AwareOf {
		newDx := itemPos.X - positionComponent.X
		newDy := itemPos.Y - positionComponent.Y
		if newDx != 0 || newDy != 0 {
			newDis := newDx*newDx + newDy*newDy
			for item := range items {
				if isTreasure(m, item) {

					oldDis := dx*dx + dy*dy

					if newDis < oldDis {
						dx = newDx
						dy = newDy
					}
				}
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
		returnEvents = append(returnEvents, Event{TRY_MOVE, TryMove{moveX, moveY}, brain})
	}

	return returnEvents
}

// Needs to control for picking up only one item; right now just takes first treasure from the pile
func (s *BrainHandler) pickup(m *Manager, brain Entity) (returnEvents []Event) {
	positionData, _ := m.getComponent(brain, POSITION)
	positionComponent := positionData.(Position)

	entities := m.getEntitiesAtPosition(positionComponent)
	for entity := range entities {
		if isStashableTreasure(m, entity) {
			returnEvents = append(returnEvents, Event{TRY_PICK_UP, TryPickUp{entity}, brain})
			break // don't need to check anymore
		}
	}

	return returnEvents
}

func isStashableTreasure(m *Manager, item Entity) bool {
	stashableData, isStashable := m.getComponent(item, STASHABLE)
	if isStashable {
		stashableComponent := stashableData.(Stashable)
		return !stashableComponent.Stashed
	}
	return false
}

func isTreasure(m *Manager, item Entity) bool {
	_, isStashable := m.getComponent(item, STASHABLE)
	return isStashable
}
