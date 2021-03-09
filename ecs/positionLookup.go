package ecs

type PositionLookup map[int]map[int]map[Entity]bool

func (p PositionLookup) checkInited(x, y int) {
	_, ok := p[x]
	if !ok {
		p[x] = make(map[int]map[Entity]bool)
	}
	_, ok = p[x][y]
	if !ok {
		p[x][y] = map[Entity]bool{}
	}
}
func (p PositionLookup) add(entities map[Entity]bool, x, y int) {
	p.checkInited(x, y)
	for entity := range entities {
		alreadyHave := p[x][y][entity]

		if !alreadyHave {
			p[x][y][entity] = true
		}
	}
}

func (p PositionLookup) remove(entity Entity, x, y int) {
	p.checkInited(x, y)

	delete(p[x][y], entity)
}

func (p PositionLookup) get(x, y int) map[Entity]bool {
	_, ok := p[x]
	if ok {
		_, ok := p[x][y]
		if ok {
			return p[x][y]
		}
	}
	return make(map[Entity]bool, 0)
}

func (p PositionLookup) move(entity Entity, x, y, nx, ny int) {
	p.remove(entity, x, y)
	p.add(map[Entity]bool{entity: true}, nx, ny)
}
