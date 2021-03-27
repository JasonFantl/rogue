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

// for manager
type ChunkPositionLookup map[int]map[int]PositionLookup

func (p ChunkPositionLookup) checkInited(x, y int) {
	_, ok := p[x]
	if !ok {
		p[x] = make(map[int]PositionLookup)
	}
	_, ok = p[x][y]
	if !ok {
		p[x][y] = PositionLookup{}
	}
}

func (p ChunkPositionLookup) add(entities map[Entity]bool, x, y int) {
	cx, cy := cordsToChunkCords(x, y)
	p.checkInited(cx, cy)
	p[cx][cy].add(entities, x, y)
}

func (p ChunkPositionLookup) remove(entity Entity, x, y int) {
	cx, cy := cordsToChunkCords(x, y)
	p.checkInited(cx, cy)
	p[cx][cy].remove(entity, x, y)
}

func (p ChunkPositionLookup) get(x, y int) map[Entity]bool {
	cx, cy := cordsToChunkCords(x, y)
	p.checkInited(cx, cy)
	return p[cx][cy].get(x, y)
}

func (p ChunkPositionLookup) move(entity Entity, x, y, nx, ny int) {
	cx, cy := cordsToChunkCords(x, y)
	ncx, ncy := cordsToChunkCords(nx, ny)

	p.checkInited(cx, cy)
	if cx != ncx || cy != ncy {
		p.checkInited(ncx, ncy)
		p[cx][cy].remove(entity, x, y)
		p[ncx][ncy].add(map[Entity]bool{entity: true}, nx, ny)
	} else {
		p[cx][cy].move(entity, x, y, nx, ny)
	}
}

func (p ChunkPositionLookup) setChunk(cx, cy int, positionLookup PositionLookup) {
	p.checkInited(cx, cy)

	p[cx][cy] = positionLookup
}

func (p ChunkPositionLookup) popChunk(cx, cy int) PositionLookup {
	p.checkInited(cx, cy)

	positionLookup := p[cx][cy]
	delete(p[cx], cy)
	return positionLookup
}
