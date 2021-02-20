package ecs

import "math/rand"

type MonsterHandler struct{}

func (s *MonsterHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// decision making is done at the beginning of the turn
	if event.ID == TIMESTEP {

		// well just have each monster move one at a time
		monsters, monstersOk := m.getComponents(MONSTER_CONTROLLER)
		if monstersOk {
			for monster, _ := range monsters {
				// randomly walks, weighted towards treasure

				// unpack all the components we will need
				positionData, hasPosition := m.getComponent(monster, POSITION)

				if hasPosition {
					positionComponent := positionData.(Position)

					// handy func
					isTreasure := func(entity Entity) bool {
						_, hasPickUpAble := m.getComponent(entity, PICKUPABLE)

						_, hasStashed := m.getComponent(entity, STASHED_FLAG)
						// funny side effect if we enable, they chase people and monsters with treasure
						hasStashed = false

						return hasPickUpAble && !hasStashed
					}

					// first check if were on top of any treasure
					entites := m.getEntitiesFromPos(positionComponent.X, positionComponent.Y)
					for _, entity := range entites {
						if isTreasure(entity) {
							// if so, try to pick it up
							returnEvents = append(returnEvents, Event{TRY_PICK_UP, TryPickUp{true, entity}, monster})
						}
					}
					// should this end its turn? no for now

					// randomly move, weighted towards treasure
					moveX := rand.Intn(3) - 1
					moveY := 0

					if rand.Intn(2) == 0 {
						moveX = 0
						moveY = rand.Intn(3) - 1
					}

					treasures, _ := m.getComponents(PICKUPABLE)
					for e := range treasures {
						_, hasStashed := m.getComponent(e, STASHED_FLAG)
						if !hasStashed {
							itemPositionData, hasPosition := m.getComponent(e, POSITION)
							if hasPosition {
								itemPositionComponent := itemPositionData.(Position)
								dx := itemPositionComponent.X - positionComponent.X
								dy := itemPositionComponent.Y - positionComponent.Y

								if dx*dx+dy*dy < rand.Intn(80) {
									if dx*dx > dy*dy {
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
