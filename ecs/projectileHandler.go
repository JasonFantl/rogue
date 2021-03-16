package ecs

import (
	"math"

	"github.com/jasonfantl/rogue/gui"
)

// TODO: move in given direction, fix equipping handler to unequip when fired
type ProjectileHandler struct {
}

func (h *ProjectileHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == TRY_LAUNCH {
		tryLaunchEvent := event.data.(TryLaunch)

		projectileData, isProjectile := m.getComponent(tryLaunchEvent.what, PROJECTILE)

		if isProjectile {
			projectileComponent := projectileData.(Projectile)

			projectileComponent.MaxDistance = 2
			fighterData, isFighter := m.getComponent(event.entity, FIGHTER)
			if isFighter {
				fighterComponent := fighterData.(Fighter)
				projectileComponent.MaxDistance = fighterComponent.Strength
			}
			projectileComponent.goalDx, projectileComponent.goalDy = tryLaunchEvent.dx, tryLaunchEvent.dy
			projectileComponent.currentDx, projectileComponent.currentDy = 0, 0
			m.setComponent(tryLaunchEvent.what, PROJECTILE, projectileComponent)

			// how to let players walk on projectiles, but have them stop when fired?
			// bad solution for now
			returnEvents = append(returnEvents, h.tryMoveProjectile(m, tryLaunchEvent.what)...)
		}
	}

	// this listener is necessary since i dont have a better solution for stopping projectiles
	// try_move handled first by move handler, queing move event, then this handler removes the volume
	if event.ID == TRY_MOVE {
		_, isProjectile := m.getComponent(event.entity, PROJECTILE)
		if isProjectile {
			m.removeComponent(event.entity, VOLUME)
		}
	}

	if event.ID == MOVED {
		movedEvent := event.data.(Moved)

		projectileData, isProjectile := m.getComponent(event.entity, PROJECTILE)

		if isProjectile {
			projectileComponent := projectileData.(Projectile)
			projectileComponent.currentDx += movedEvent.toX - movedEvent.fromX
			projectileComponent.currentDy += movedEvent.toY - movedEvent.fromY
			m.setComponent(event.entity, PROJECTILE, projectileComponent)

			returnEvents = append(returnEvents, h.tryMoveProjectile(m, event.entity)...)

			// so we can see the patth, remove later
			path := []Component{
				{POSITION, Position{movedEvent.toX, movedEvent.toY}},
				{DISPLAYABLE, Displayable{gui.GetSprite(gui.BLOOD)}},
			}

			m.AddEntity(path)
		}
	}

	return returnEvents
}

func (s *ProjectileHandler) tryMoveProjectile(m *Manager, projectile Entity) (returnEvents []Event) {
	projectileData, isProjectile := m.getComponent(projectile, PROJECTILE)

	if isProjectile {
		projectileComponent := projectileData.(Projectile)

		// check were not at our goal
		if projectileComponent.goalDx != projectileComponent.currentDx || projectileComponent.goalDy != projectileComponent.currentDy {
			dx, dy := 0, 0

			if projectileComponent.goalDx > projectileComponent.currentDx {
				dx = 1
			} else if projectileComponent.goalDx < projectileComponent.currentDx {
				dx = -1
			}
			if projectileComponent.goalDy > projectileComponent.currentDy {
				dy = 1
			} else if projectileComponent.goalDy < projectileComponent.currentDy {
				dy = -1
			}

			if dx != 0 && dy != 0 {

				// how to best check next move?
				theta := math.Atan2(float64(projectileComponent.goalDx), float64(projectileComponent.goalDy))

				possibleX := float64(projectileComponent.currentDx + dx)
				possibleY := float64(projectileComponent.currentDy)
				rotatedXFrommovingX := math.Cos(theta)*possibleX - math.Sin(theta)*possibleY

				possibleX = float64(projectileComponent.currentDx)
				possibleY = float64(projectileComponent.currentDy + dy)
				rotatedXFrommovingY := math.Cos(theta)*possibleX - math.Sin(theta)*possibleY

				if math.Abs(rotatedXFrommovingX) < math.Abs(rotatedXFrommovingY) {
					dy = 0
				} else {
					dx = 0
				}
			}

			distance := projectileComponent.currentDx*projectileComponent.currentDx + projectileComponent.currentDy*projectileComponent.currentDy
			if distance < projectileComponent.MaxDistance*projectileComponent.MaxDistance {
				m.setComponent(projectile, VOLUME, Volume{})
				returnEvents = append(returnEvents, Event{TRY_MOVE, TryMove{dx, dy}, projectile})
			}
		}
	}
	return returnEvents
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
