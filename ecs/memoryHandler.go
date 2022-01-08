package ecs

type MemoryHandler struct{}

func (s *MemoryHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == DISPLAY {

		entities := m.getEntitiesWithComponent(ENTITY_MEMORY)

		for entity := range entities {
			memoryData, hasMemory := m.getComponent(entity, ENTITY_MEMORY)
			awarenessData, hasAwarness := m.getComponent(entity, ENTITY_AWARENESS)

			if hasMemory && hasAwarness {
				memoryComponent := memoryData.(EntityMemory)
				awarenessComponent := awarenessData.(EntityAwareness)

				// make sure memory is initalized
				if memoryComponent.Memory == nil {
					memoryComponent.Memory = make(map[Position][]Displayable)
				}

				for pos, items := range awarenessComponent.AwareOf {
					// make sure memory is initalized
					if memoryComponent.Memory[pos] == nil {
						memoryComponent.Memory[pos] = make([]Displayable, 0)
					}
					updatedEntities := make([]Displayable, 0)
					for item := range items {
						_, itemIsMemorable := m.getComponent(item, MEMORABLE)
						itemDisplayData, itemHasDisplay := m.getComponent(item, DISPLAYABLE)

						if itemIsMemorable && itemHasDisplay {
							itemDisplayComponent := itemDisplayData.(Displayable)

							updatedEntities = append(updatedEntities, itemDisplayComponent)
						}
					}
					memoryComponent.Memory[pos] = updatedEntities
				}

				m.setComponent(entity, ENTITY_MEMORY, memoryComponent)
			}
		}
	}

	return returnEvents
}
