package ecs

import "math/rand"

type MonsterHandler struct{}

func (s *MonsterHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// basid decision making at beginning of turn
	if event.ID == TIMESTEP {

		// well just have each monster move one at a time
		monsters, monstersOk := m.getComponents(MONSTER_CONTROLLER)
		if monstersOk {
			for monster, _ := range monsters {
				// moves to nearest treasure and pick up

				// unpack all the components we will need
				positionData, hasPosition := m.getComponent(monster, POSITION)

				if hasPosition {
					positionComponent := positionData.(Position)

					// handy func, also, maybe more specs for what a monster considers treasure?
					isTreasure := func(entity Entity) bool {
						_, hasPickUpAble := m.getComponent(entity, PICKUPABLE)
						_, hasStashed := m.getComponent(entity, STASHED)
						return hasPickUpAble && !hasStashed
					}

					// first check if were on top of any treasure
					entites := m.getEntitiesFromPos(positionComponent.X, positionComponent.Y)
					for _, entity := range entites {
						if isTreasure(entity) {
							// if so, try to pick it up
							returnEvents = append(returnEvents, Event{TRY_PICK_UP_EVENT, EventTryPickUp{true, entity}, monster})
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
						_, hasStashed := m.getComponent(e, STASHED)
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
						returnEvents = append(returnEvents, Event{TRY_MOVE_EVENT, EventTryMove{moveX, moveY}, monster})
					}
				}
			}
		}
	}

	return returnEvents
}
