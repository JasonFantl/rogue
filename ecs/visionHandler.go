package ecs

type VisionHandler struct{}

func (s *VisionHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// not sure how/when to update/handle this event, kinda tricky
	if event.ID == DISPLAY {

		entities := m.getEntitiesWithComponent(VISION)

		for entity := range entities {
			visionData, hasVision := m.getComponent(entity, VISION)
			awarenessData, hasAwarness := m.getComponent(entity, ENTITY_AWARENESS)
			positionData, hasPosition := m.getComponent(entity, POSITION)

			if hasVision && hasAwarness && hasPosition {
				visionComponent := visionData.(Vision)
				awarenessComponent := awarenessData.(EntityAwareness)
				positionComponent := positionData.(Position)

				// clear old awareness first
				awarenessComponent.AwareOf = PositionLookup{}

				updateAwareOf(m, positionComponent, visionComponent, awarenessComponent.AwareOf)

				m.setComponent(entity, ENTITY_AWARENESS, awarenessComponent)
			}
		}
	}

	return returnEvents
}

func updateAwareOf(m *Manager, position Position, vision Vision, awareOf PositionLookup) {
	// add where we are
	entities := m.getEntitiesAtPosition(position)
	awareOf.addEntities(entities, position)

	for octant := 0; octant < 8; octant++ {
		updateOctant(m, position, vision, awareOf, octant)
	}
}

func getOctantBounds(radius int) []int {
	//---from display
	// // get bounds, just do quarter circle
	// type bound struct{ row, col int }
	// bounds := make([]bound, 0)

	// circleX := displayRadius
	// circleY := 0

	// // Initializing the value of P
	// P := 1 - displayRadius
	// for circleX > circleY {
	// 	//circle math
	// 	circleY++
	// 	// Mid-point is inside or on the perimeter
	// 	if P <= 0 {
	// 		P = P + 2*circleY + 1
	// 	} else { // Mid-point is outside the perimeter
	// 		circleX--
	// 		P = P + 2*circleY - 2*circleX + 1
	// 	}
	// 	// All the perimeter points have already been displayed
	// 	if circleX < circleY {
	// 		break
	// 	}

	// 	bounds = append(bounds, bound{circleY, circleX})
	// 	if circleX != circleY && P > 0 {
	// 		bounds = append(bounds, bound{circleX, circleY})
	// 	}
	// }
	bounds := make([]int, 0)

	circleX := radius
	circleY := 0
	// Initializing the value of P
	P := 1 - radius

	bounds = append(bounds, circleX)
	for circleX > circleY {
		//circle math
		circleY++
		// Mid-point is inside or on the perimeter
		if P < 0 {
			P = P + 2*circleY + 1
		} else { // Mid-point is outside the perimeter
			circleX--
			P = P + 2*circleY - 2*circleX + 1
		}
		// All the perimeter points have already been displayed
		if circleX < circleY {
			break
		}
		bounds = append(bounds, circleX)
	}

	return bounds
}

func updateOctant(m *Manager, position Position, vision Vision, awareOf PositionLookup, octant int) {
	// first create bounds of octant
	bounds := getOctantBounds(vision.Radius)

	line := ShadowLine{}

	for row := 1; row < vision.Radius; row++ {
		for col := 0; col <= row; col++ {
			// check if out of bounds
			if bounds[col] == row {
				break
			}
			// in bounds, continue on
			dx, dy := transformOctant(row, col, octant)
			x := position.X + dx
			y := position.Y + dy

			// Set the visibility of this tile.
			visible := !line.isInShadow(projectTile(row, col))

			if visible {
				pos := Position{x, y}
				entities := m.getEntitiesAtPosition(pos)
				awareOf.addEntities(entities, pos)

				// Add any opaque tiles to the shadow map.
				isOpaque := false
				for entity := range entities {
					_, hasOpaque := m.getComponent(entity, OPAQUE)
					if hasOpaque {
						isOpaque = true
						break
					}
				}

				if isOpaque {
					line.add(projectTile(row, col))
				}
			}
		}
	}
}

func transformOctant(row, col, octant int) (int, int) {
	switch octant {
	case 0:
		return col, -row
	case 1:
		return row, -col
	case 2:
		return -row, -col
	case 3:
		return -col, -row
	case 4:
		return -col, row
	case 5:
		return -row, col
	case 6:
		return row, col
	case 7:
		return col, row
	}
	return 0, 0
}

// ---- shadow line -----
type ShadowLine struct {
	Shadows []Shadow
}

func (s ShadowLine) isInShadow(inS Shadow) bool {
	for _, shadow := range s.Shadows {
		if shadow.contains(inS) {
			return true
		}
	}
	return false
}

func (s *ShadowLine) add(shadow Shadow) {
	// Figure out where to slot the new shadow in the list.
	index := 0
	for ; index < len(s.Shadows); index++ {
		// Stop when we hit the insertion point.
		if s.Shadows[index].Start >= shadow.Start {
			break
		}
	}

	// The new shadow is going here. See if it overlaps the
	// previous or next.
	overlappingPrevious := false
	if index > 0 && s.Shadows[index-1].End >= shadow.Start {
		overlappingPrevious = true
	}

	overlappingNext := false
	if index < len(s.Shadows) && s.Shadows[index].Start <= shadow.End {
		overlappingNext = true
	}

	// Insert and unify with overlapping shadows.
	if overlappingNext {
		if overlappingPrevious {
			// Overlaps both, so unify one and delete the other.
			s.Shadows[index-1].End = s.Shadows[index].End
			s.Shadows = append(s.Shadows[:index], s.Shadows[index+1:]...)
		} else {
			// Overlaps the next one, so unify it with that.
			s.Shadows[index].Start = shadow.Start
		}
	} else {
		if overlappingPrevious {
			// Overlaps the previous one, so unify it with that.
			s.Shadows[index-1].End = shadow.End
		} else {
			// Does not overlap anything, so insert.
			s.Shadows = append(s.Shadows, Shadow{})
			copy(s.Shadows[index+1:], s.Shadows[index:])
			s.Shadows[index] = shadow
		}
	}
}

// ---- shadow ----
type Shadow struct {
	Start, End float64
}

func (s Shadow) contains(inS Shadow) bool {
	return s.Start <= inS.Start && inS.End <= s.End
}

func projectTile(r, c int) Shadow {
	row := float64(r)
	col := float64(c)
	cellRadius := 0.5
	start := (col - cellRadius) / (row + cellRadius)
	end := (col + cellRadius) / (row - cellRadius)
	return Shadow{start, end}
}
