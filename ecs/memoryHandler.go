package ecs

type MemoryHandler struct{}

func (s *MemoryHandler) handleEvent(m *Manager, event Event) (returnEvents []Event) {

	if event.ID == DISPLAY {

		entities, _ := m.getComponents(ENTITY_MEMORY)

		for entity := range entities {
			memoryData, hasMemory := m.getComponent(entity, ENTITY_MEMORY)
			awarnessData, hasAwarness := m.getComponent(entity, ENTITY_AWARENESS)

			if hasMemory && hasAwarness {
				memoryComponent := memoryData.(EntityMemory)
				awarnessComponent := awarnessData.(EntityAwarness)

				for _, item := range awarnessComponent.AwareOf {
					_, itemIsMemorable := m.getComponent(item, MEMORABLE)
					itemDisplayData, itemHasDisplay := m.getComponent(item, DISPLAYABLE)
					itemPositionData, itemHasPosition := m.getComponent(item, POSITION)

					if itemIsMemorable && itemHasDisplay && itemHasPosition {
						itemDisplayComponent := itemDisplayData.(Displayable)
						itemPositionComponent := itemPositionData.(Position)

						x := itemPositionComponent.X
						y := itemPositionComponent.Y
						// make sure memory is inited
						if memoryComponent.Memory == nil {
							memoryComponent.Memory = make(map[int]map[int]Displayable)
						}
						if memoryComponent.Memory[x] == nil {
							memoryComponent.Memory[x] = make(map[int]Displayable)
						}

						memoryComponent.Memory[x][y] = itemDisplayComponent
					}
				}
				m.setComponent(entity, ENTITY_MEMORY, memoryComponent)
			}
		}
	}

	return returnEvents
}
