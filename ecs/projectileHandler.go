package ecs

// TODO: move in given direction, fix equipping handler to unequip when fired
type ProjectileHandler struct {
}

func (h *ProjectileHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == TIMESTEP {
		projectiles := m.getEntities(PROJECTILE)
		for projectile := range projectiles {
			projectileData, _ := m.getComponent(projectile, PROJECTILE)
			projectileComponent := projectileData.(Projectile)
			projectileComponent.inFlight = false

			m.setComponent(projectile, PROJECTILE, projectileComponent)
		}
	}

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
			projectileComponent.inFlight = true
			m.setComponent(tryLaunchEvent.what, PROJECTILE, projectileComponent)

			// how to let players walk on projectiles, but have them stop when fired?
			// bad solution for now
			returnEvents = append(returnEvents, h.tryMoveProjectile(m, tryLaunchEvent.what)...)
		}
	}

	// this listener is necessary since i dont have a better solution for stopping projectiles
	// try_move handled first by move handler, queing move event, then this handler removes the volume
	if event.ID == TRY_MOVE {
		projectileData, isProjectile := m.getComponent(event.entity, PROJECTILE)
		if isProjectile {
			projectileComponent := projectileData.(Projectile)

			if projectileComponent.inFlight {
				m.removeComponent(event.entity, VOLUME)
			}
		}
	}

	if event.ID == MOVED {
		movedEvent := event.data.(Moved)

		projectileData, isProjectile := m.getComponent(event.entity, PROJECTILE)

		if isProjectile {
			projectileComponent := projectileData.(Projectile)
			projectileComponent.currentDx += movedEvent.to.X - movedEvent.from.X
			projectileComponent.currentDy += movedEvent.to.Y - movedEvent.from.Y
			m.setComponent(event.entity, PROJECTILE, projectileComponent)

			if projectileComponent.inFlight {
				returnEvents = append(returnEvents, h.tryMoveProjectile(m, event.entity)...)
			}
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

				x2 := projectileComponent.goalDx
				y2 := projectileComponent.goalDy

				y3 := projectileComponent.currentDy
				x3 := projectileComponent.currentDx + dx
				rxToY := x2*y3 - y2*x3

				y3 = projectileComponent.currentDy + dy
				x3 = projectileComponent.currentDx
				ryToY := x2*y3 - y2*x3

				if Abs(rxToY) < Abs(ryToY) {
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
