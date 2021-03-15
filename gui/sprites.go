package gui

import (
	"image"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	tileSize = 8
)

var (
	tileSheet *ebiten.Image
)

func loadSprites() {
	imgFile, err := ebitenutil.OpenFile("data/spritesheet.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
	}

	tileSheet = ebiten.NewImageFromImage(img)

	loadPremadeSprites()
}

var preMadeSprites map[int]Sprite

const (
	GRASS_FLOOR int = iota
	DIRT_FLOOR
	STONE_FLOOR
	STONE_WALL
	BLOOD
	TREE_TRUNK
	LEAF
	PLAYER
	MONSTER1
	MONSTER2
	MONSTER3
	STICK
	SWORD
	METAL_ARMOR
	LEATHER_ARMOR
	GOLD_COIN
	SILVER_COIN
	GEM
	POTION
	KEY
	CLOSED_DOOR
	OPEN_DOOR
	CURSER
)

// priorities :
// 0 - 10: floor
// 10 - 20: items
// 50 - 60: beings
// 100 > : special sprites (menu)
func loadPremadeSprites() {
	preMadeSprites = make(map[int]Sprite, 0)

	baseOb := ebiten.DrawImageOptions{}
	leafOb := baseOb
	leafOb.ColorM.Scale(0.5, 0.5, 1, 0.5)
	openDoorOb := baseOb
	openDoorOb.GeoM.Rotate(math.Pi / 2)
	openDoorOb.GeoM.Translate(tileSize, 3)

	preMadeSprites[GRASS_FLOOR] = Sprite{extractImage(0, 0), baseOb, 3}
	preMadeSprites[LEAF] = Sprite{extractImage(0, 0), leafOb, 92}
	preMadeSprites[DIRT_FLOOR] = Sprite{extractImage(1, 0), baseOb, 2}
	preMadeSprites[STONE_FLOOR] = Sprite{extractImage(2, 0), baseOb, 4}
	preMadeSprites[STONE_WALL] = Sprite{extractImage(3, 0), baseOb, 98}
	preMadeSprites[BLOOD] = Sprite{extractImage(4, 0), baseOb, 9}
	preMadeSprites[TREE_TRUNK] = Sprite{extractImage(5, 0), baseOb, 91}
	preMadeSprites[CLOSED_DOOR] = Sprite{extractImage(6, 0), baseOb, 91}
	preMadeSprites[OPEN_DOOR] = Sprite{extractImage(6, 0), openDoorOb, 91}

	preMadeSprites[PLAYER] = Sprite{extractImage(0, 1), baseOb, 59}
	preMadeSprites[MONSTER1] = Sprite{extractImage(1, 1), baseOb, 58}
	preMadeSprites[MONSTER2] = Sprite{extractImage(2, 1), baseOb, 58}
	preMadeSprites[MONSTER3] = Sprite{extractImage(3, 1), baseOb, 58}

	preMadeSprites[STICK] = Sprite{extractImage(0, 2), baseOb, 11}
	preMadeSprites[SWORD] = Sprite{extractImage(1, 2), baseOb, 11}
	preMadeSprites[METAL_ARMOR] = Sprite{extractImage(2, 2), baseOb, 11}
	preMadeSprites[LEATHER_ARMOR] = Sprite{extractImage(3, 2), baseOb, 11}

	preMadeSprites[GOLD_COIN] = Sprite{extractImage(0, 3), baseOb, 11}
	preMadeSprites[SILVER_COIN] = Sprite{extractImage(1, 3), baseOb, 11}
	preMadeSprites[GEM] = Sprite{extractImage(2, 3), baseOb, 11}
	preMadeSprites[POTION] = Sprite{extractImage(3, 3), baseOb, 11}
	preMadeSprites[KEY] = Sprite{extractImage(4, 3), baseOb, 11}

	preMadeSprites[CURSER] = Sprite{extractImage(0, 4), baseOb, 100}

}

func extractImage(x, y int) *ebiten.Image {
	x, y = x*tileSize, y*tileSize
	return tileSheet.SubImage(image.Rect(x, y, x+tileSize, y+tileSize)).(*ebiten.Image)
}

type Sprite struct {
	Image    *ebiten.Image
	Options  ebiten.DrawImageOptions
	Priority int
}

func GetSprite(id int) Sprite {
	return preMadeSprites[id]
}

// handleing alpha is weird, need to keep priorities, but not overturn aware tiles
func Fade(sprite Sprite) Sprite {
	sprite.Options.ColorM.Scale(0.5, 0.5, 0.5, 1)
	sprite.Priority -= 100
	return sprite
}
