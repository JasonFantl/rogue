package ecs

type VisionHandler struct{}

func (s *VisionHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	// not sure how/when to update/handle this event, kinda tricky
	if event.ID == DISPLAY {

		entities, _ := m.getComponents(VISION)

		for entity := range entities {
			visionData, hasVision := m.getComponent(entity, VISION)
			awarnessData, hasAwarness := m.getComponent(entity, ENTITY_AWARENESS)
			positionData, hasPosition := m.getComponent(entity, POSITION)

			if hasVision && hasAwarness && hasPosition {
				visionComponent := visionData.(Vision)
				awarnessComponent := awarnessData.(EntityAwarness)
				positionComponent := positionData.(Position)

				// clear old awarness first
				awarnessComponent.AwareOf = make([]Entity, 0)

				updateAwareOf(m, positionComponent, visionComponent, &awarnessComponent.AwareOf)

				m.setComponent(entity, Component{ENTITY_AWARENESS, awarnessComponent})
			}
		}
	}

	return returnEvents
}

func updateAwareOf(m *Manager, position Position, vision Vision, awareOf *[]Entity) {
	// add where we are
	entities := m.getEntitiesFromPos(position.X, position.Y)
	*awareOf = append(*awareOf, entities...)

	for octant := 0; octant < 8; octant++ {
		updateOctant(m, position, vision, awareOf, octant)
	}
}

func updateOctant(m *Manager, position Position, vision Vision, awareOf *[]Entity, octant int) {
	line := ShadowLine{}

	for row := 1; row < vision.Radius; row++ {
		for col := 0; col <= row; col++ {
			delta := transformOctant(row, col, octant)
			x := position.X + delta.X
			y := position.Y + delta.Y

			// Set the visibility of this tile.
			bottomVisible := !line.isInShadow(float64(row), float64(col)-0.5)
			leftVisible := !line.isInShadow(float64(row)-0.5, float64(col))

			visible := bottomVisible || leftVisible

			if visible {
				entities := m.getEntitiesFromPos(x, y)
				*awareOf = append(*awareOf, entities...)

				// Add any opaque tiles to the shadow map.
				isOpaque := false
				for _, entity := range entities {
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

func transformOctant(row, col, octant int) Position {
	switch octant {
	case 0:
		return Position{col, -row}
	case 1:
		return Position{row, -col}
	case 2:
		return Position{row, col}
	case 3:
		return Position{col, row}
	case 4:
		return Position{-col, row}
	case 5:
		return Position{-row, col}
	case 6:
		return Position{-row, -col}
	case 7:
		return Position{-col, -row}
	}
	return Position{0, 0}
}

// ---- shadow line -----
type ShadowLine struct {
	Shadows []Shadow
}

func (s ShadowLine) isInShadow(r, c float64) bool {
	for _, shadow := range s.Shadows {
		if shadow.contains(r, c) {
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

func (s Shadow) contains(r, c float64) bool {
	slope := c / r

	return s.Start <= slope && slope <= s.End
}

func projectTile(r, c int) Shadow {
	row := float64(r)
	col := float64(c)
	start := (col - 0.5) / (row + 0.5)
	end := (col + 0.5) / (row - 0.5)
	return Shadow{start, end}
}
