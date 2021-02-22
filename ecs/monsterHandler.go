package ecs

import "math/rand"

type MonsterHandler struct{}

func (s *MonsterHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// decision making is done at the beginning of the turn
	if event.ID == TIMESTEP {

		// well just have each monster move one at a time
		monsters, monstersExist := m.getComponents(MONSTER_CONTROLLER)
		if monstersExist {
			for monster, _ := range monsters {
				// randomly walks, weighted towards treasure

				// unpack all the components we will need
				positionData, hasPosition := m.getComponent(monster, POSITION)

				if hasPosition {
					positionComponent := positionData.(Position)

					// first check if we are on top of any treasure
					entites := m.getEntitiesFromPos(positionComponent.X, positionComponent.Y)
					for _, entity := range entites {
						_, hasPickUpAble := m.getComponent(entity, PICKUPABLE)
						_, hasStashed := m.getComponent(entity, STASHED_FLAG)
						isTreasure := hasPickUpAble && !hasStashed
						if isTreasure {
							// if so, try to pick it up
							returnEvents = append(returnEvents, Event{TRY_PICK_UP, TryPickUp{true, entity}, monster})
						}
					}
					// should this end its turn? no for now

					// randomly move, weighted towards treasure
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
				}
			}
		}
	}

	return returnEvents
}
