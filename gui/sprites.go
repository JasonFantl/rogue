package gui

import (
	"image"
	_ "image/png"
	"log"

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
	imgFile, err := ebitenutil.OpenFile("data/first.png")
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
	TREE_TRUNK
)

// priorities :
// 0 - 10: floor
// 10 - 20: items
// 50 - 60: beings
func loadPremadeSprites() {
	preMadeSprites = make(map[int]Sprite, 0)

	baseOb := ebiten.DrawImageOptions{}
	leafOb := baseOb
	leafOb.ColorM.Scale(0.5, 0.5, 1, 0.5)

	preMadeSprites[GRASS_FLOOR] = Sprite{extractImage(0, 0), baseOb, 3}
	preMadeSprites[DIRT_FLOOR] = Sprite{extractImage(1, 0), baseOb, 2}
	preMadeSprites[STONE_FLOOR] = Sprite{extractImage(2, 0), baseOb, 4}
	preMadeSprites[STONE_WALL] = Sprite{extractImage(3, 0), baseOb, 98}
	preMadeSprites[BLOOD] = Sprite{extractImage(4, 0), baseOb, 9}
	preMadeSprites[TREE_TRUNK] = Sprite{extractImage(5, 0), baseOb, 91}

	preMadeSprites[PLAYER] = Sprite{extractImage(0, 1), baseOb, 59}
	preMadeSprites[MONSTER] = Sprite{extractImage(1, 1), baseOb, 58}

	preMadeSprites[TREASURE] = Sprite{extractImage(0, 2), baseOb, 11}
	preMadeSprites[POTION] = Sprite{extractImage(1, 2), baseOb, 12}
	preMadeSprites[WEAPON] = Sprite{extractImage(2, 2), baseOb, 13}

	preMadeSprites[LEAF] = Sprite{extractImage(0, 0), leafOb, 92}

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
