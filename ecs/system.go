package ecs

import (
	"github.com/jasonfantl/rogue/gui"
)

type System interface {
	run(*manager)
}

type SystemDisplay struct {
	ground     [][]rune
	priorities [][]int
}

func (s *SystemDisplay) SetSize(size int) {
	s.ground = make([][]rune, size)
	s.priorities = make([][]int, size)

	for i := range s.ground {
		s.ground[i] = make([]rune, size)
		s.priorities[i] = make([]int, size)
	}
}

func (s *SystemDisplay) run(m *manager) {
	// first clear old map
	for i := range s.ground {
		for j := range s.ground[i] {
			s.ground[i][j] = ' '
			s.priorities[i][j] = -1
		}
	}

	// get new positions, the looping is currently as horrible as I can make it
	for _, componenets := range m.lookupTable {
		positionData, positionOk := componenets[POSITION]
		displayData, displayOk := componenets[DISPLAY]

		if positionOk && displayOk {
			positionComponent := positionData.(Position)
			displayComponent := displayData.(Display)

			x := positionComponent.X
			y := positionComponent.Y

			if x >= 0 && x < cap(s.ground) {
				if y >= 0 && y < cap(s.ground[x]) {
					if displayComponent.Priority > s.priorities[x][y] {
						s.ground[x][y] = displayComponent.Character
						s.priorities[x][y] = displayComponent.Priority
					}
				}
			}
		}
	}

	gui.Show(s.ground)
}

type SystemControl struct {
}

func (s *SystemControl) run(m *manager) {

	// somwhere in here the components are copied, so we cant edit them
	// that why we set the components at the end
	// there must be a way to maintain the pointer to the component

	for _, componenets := range m.lookupTable {
		controllerData, controllerOk := componenets[CONTROLLER]
		quitFlagData, quitFlagOk := componenets[QUIT_FLAG]
		desiredMoveData, desiredMoveOk := componenets[DESIRED_MOVE]

		if controllerOk && quitFlagOk && desiredMoveOk {
			// is it here at the type conversion that we lose refrence?
			controllerComponent := controllerData.(Controller)
			quitFlagComponent := quitFlagData.(QuitFlag)
			desiredMoveComponent := desiredMoveData.(DesiredMove)

			key, pressed := gui.GetKeyPress()

			if pressed {
				switch key {
				case controllerComponent.Down:
					desiredMoveComponent.X = 0
					desiredMoveComponent.Y = -1
				case controllerComponent.Up:
					desiredMoveComponent.X = 0
					desiredMoveComponent.Y = 1
				case controllerComponent.Left:
					desiredMoveComponent.X = -1
					desiredMoveComponent.Y = 0
				case controllerComponent.Right:
					desiredMoveComponent.X = 1
					desiredMoveComponent.Y = 0
				case controllerComponent.Quit:
					quitFlagComponent.HasQuit = true
				}
			}

			componenets[QUIT_FLAG] = quitFlagComponent
			componenets[DESIRED_MOVE] = desiredMoveComponent

		}
	}
}

// suuuuuper basic, slow and bad
type SystemMove struct {
}

func (s *SystemMove) run(m *manager) {

	for _, componenets := range m.lookupTable {
		positionData, positionOk := componenets[POSITION]
		desiredMoveData, desiredMoveOk := componenets[DESIRED_MOVE]
		_, blockingOk := componenets[BLOCKABLE_TAG]

		if positionOk && desiredMoveOk {

			positionComponent := positionData.(Position)
			desiredMoveComponent := desiredMoveData.(DesiredMove)

			// now check if new location is occupied
			newX := positionComponent.X + desiredMoveComponent.X
			newY := positionComponent.Y + desiredMoveComponent.Y

			canMove := true
			// quick implementation, replace later
			for _, othersComponenets := range m.lookupTable {
				otherPositionData, otherPositionOk := othersComponenets[POSITION]
				_, otherBlockableOk := othersComponenets[BLOCKABLE_TAG]

				if otherPositionOk {
					othersPositionComponent := otherPositionData.(Position)
					if othersPositionComponent.X == newX && othersPositionComponent.Y == newY {
						if otherBlockableOk || blockingOk {
							canMove = false
						}
					}
				}
			}

			if canMove {
				positionComponent.X = newX
				positionComponent.Y = newY
			}

			componenets[POSITION] = positionComponent
		}
	}
}
