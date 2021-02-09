package ecs

import (
	"github.com/jasonfantl/rogue/gui"
)

type System interface {
	run(*manager)
}

type SystemDisplay struct {
	ground [][]string
}

func (s *SystemDisplay) SetSize(size int) {
	s.ground = make([][]string, size)
	for i := range s.ground {
		s.ground[i] = make([]string, size)
	}
}

func (s *SystemDisplay) run(m *manager) {
	// first clear old map
	for i := range s.ground {
		for j := range s.ground[i] {
			s.ground[i][j] = " "
		}
	}

	// get new positions, the looping is currently as horrible as I can make it
	for _, componenets := range m.lookupTable {
		positionData, positionOk := componenets[POSITION]
		displayData, displayOk := componenets[DISPLAY]

		if positionOk && displayOk {
			positionComponent := positionData.(Position)
			displayComponent := displayData.(Display)

			inXbounds := positionComponent.X >= 0 && positionComponent.X < cap(s.ground)
			if inXbounds {
				inYbounds := positionComponent.Y >= 0 && positionComponent.Y < cap(s.ground[positionComponent.X])
				if inYbounds {
					s.ground[positionComponent.X][positionComponent.Y] = displayComponent.Character
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
		positionData, positionOk := componenets[POSITION]

		if controllerOk && positionOk {
			// is it here at the type conversion that we lose refrence?
			controllerComponent := controllerData.(Controller)
			positionComponent := positionData.(Position)

			key, pressed := gui.GetKeyPress()

			if pressed {
				switch key {
				case controllerComponent.Down:
					positionComponent.Y--
				case controllerComponent.Up:
					positionComponent.Y++
				case controllerComponent.Left:
					positionComponent.X--
				case controllerComponent.Right:
					positionComponent.X++
				case controllerComponent.Quit:
					controllerComponent.HasQuit = true
				}
			}

			componenets[CONTROLLER] = controllerComponent
			componenets[POSITION] = positionComponent
		}
	}
}
