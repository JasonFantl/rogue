package ecs

type PositionLookup map[int]map[int][]Entity

func (p PositionLookup) checkInited(x, y int) {
	_, ok := p[x]
	if !ok {
		p[x] = map[int][]Entity{}
	}
	_, ok = p[x][y]
	if !ok {
		p[x][y] = make([]Entity, 0)
	}
}
func (p PositionLookup) add(entity Entity, x, y int) {
	p.checkInited(x, y)
	p[x][y] = append(p[x][y], entity)
}

func (p PositionLookup) remove(entity Entity, x, y int) {
	p.checkInited(x, y)

	index := -1
	entities := p[x][y]
	for i, e := range entities {
		if e == entity {
			index = i
			break
		}
	}

	if index != -1 {
		p[x][y][index] = p[x][y][len(p[x][y])-1]
		p[x][y] = p[x][y][:len(p[x][y])-1]
	}
}

func (p PositionLookup) move(entity Entity, x, y, nx, ny int) {
	p.remove(entity, x, y)
	p.checkInited(nx, ny)
	p[nx][ny] = append(p[nx][ny], entity)
}
