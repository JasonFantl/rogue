package ecs

type manager struct {
	lookupTable   map[Entity]map[ComponentID]interface{}
	systems       []System
	entityCounter Entity
	hasQuit       bool
}

func New() manager {
	newManager := manager{}
	newManager.lookupTable = make(map[Entity]map[ComponentID]interface{})
	newManager.systems = make([]System, 0)
	newManager.entityCounter = 0

	newManager.hasQuit = false

	return newManager
}

func (m *manager) AddEntity(componenets []Component) Entity {
	entity := m.entityCounter
	m.entityCounter++

	// init component lookup table
	m.lookupTable[entity] = make(map[ComponentID]interface{})

	for _, component := range componenets {
		m.AddComponenet(entity, component)
	}

	return entity
}

func (m *manager) AddComponenet(entity Entity, component Component) {
	// check entity exist
	componentList, ok := m.lookupTable[entity]
	if ok {
		// check componenet doesnt already exist
		if _, ok := componentList[component.ID]; !ok {
			m.lookupTable[entity][component.ID] = component.Data
		}
	}
}

func (m *manager) AddSystem(system System) {
	m.systems = append(m.systems, system)
}

func (m *manager) Run() {
	for _, system := range m.systems {
		system.run(m)
	}
}

func (m *manager) Lookup(entity Entity, componentID ComponentID) (interface{}, bool) {
	components, ok := m.lookupTable[entity]
	if ok {
		data, ok := components[componentID]
		if ok {
			return data, true
		}
	}
	return nil, false
}

func (m *manager) Quit() bool {
	return m.hasQuit
}
