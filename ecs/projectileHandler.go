package ecs

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
			projectileComponent.Traveled = 0
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
		projectileData, isProjectile := m.getComponent(event.entity, PROJECTILE)

		if isProjectile {
			projectileComponent := projectileData.(Projectile)
			projectileComponent.Traveled += 1
			m.setComponent(event.entity, PROJECTILE, projectileComponent)

			returnEvents = append(returnEvents, h.tryMoveProjectile(m, event.entity)...)
		}
	}

	return returnEvents
}

func (s *ProjectileHandler) tryMoveProjectile(m *Manager, projectile Entity) (returnEvents []Event) {
	dx, dy := 0, 1 // how to calculate? line algorithm?

	projectileData, isProjectile := m.getComponent(projectile, PROJECTILE)

	if isProjectile {
		projectileComponent := projectileData.(Projectile)

		if projectileComponent.Traveled < projectileComponent.MaxDistance {
			m.AddComponenet(projectile, Component{VOLUME, Volume{}})
			returnEvents = append(returnEvents, Event{TRY_MOVE, TryMove{dx, dy}, projectile})
		}
	}
	return returnEvents
}
