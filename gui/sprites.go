// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gui

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	tileSize = 16
)

var (
	tilesImage *ebiten.Image
)

func loadSprites() {
	imgFile, err := ebitenutil.OpenFile("data/first.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
	}

	tilesImage = ebiten.NewImageFromImage(img)

	loadPremadeSprites()
}

var preMadeSprites map[int]Sprite

const (
	PLAYER int = iota
	STONE_FLOOR
	STONE_WALL
	WEAPON
	MONSTER
	BLOOD
	GRASS_FLOOR
	DIRT_FLOOR
	POTION
	TREASURE
	LEAF
	FADED
)

// priorities :
// 0 - 10: floor
// 10 - 20: items
// 50 - 60: beings
func loadPremadeSprites() {
	preMadeSprites = make(map[int]Sprite, 0)

	preMadeSprites[GRASS_FLOOR] = Sprite{extractImage(0, 0), 1}
	preMadeSprites[DIRT_FLOOR] = Sprite{extractImage(1, 0), 1}
	preMadeSprites[STONE_FLOOR] = Sprite{extractImage(2, 0), 1}
	preMadeSprites[STONE_WALL] = Sprite{extractImage(3, 0), 99}

	preMadeSprites[PLAYER] = Sprite{extractImage(0, 1), 59}
	preMadeSprites[TREASURE] = Sprite{extractImage(1, 1), 11}
	preMadeSprites[POTION] = Sprite{extractImage(2, 1), 11}
	preMadeSprites[MONSTER] = Sprite{extractImage(3, 1), 58}

	preMadeSprites[LEAF] = Sprite{extractImage(1, 2), 91}
	preMadeSprites[BLOOD] = Sprite{extractImage(2, 2), 2}
	preMadeSprites[FADED] = Sprite{extractImage(3, 2), 99}

	preMadeSprites[WEAPON] = Sprite{extractImage(3, 1), 11}
}

func extractImage(x, y int) *ebiten.Image {
	x, y = x*tileSize, y*tileSize
	return tilesImage.SubImage(image.Rect(x, y, x+tileSize, y+tileSize)).(*ebiten.Image)
}

type Sprite struct {
	Image    *ebiten.Image
	Priority int
}

func GetSprite(id int) Sprite {
	return preMadeSprites[id]
}
