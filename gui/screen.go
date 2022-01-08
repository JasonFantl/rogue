package gui

import (
	"image/color"
	"sort"

	"github.com/hajimehoshi/bitmapfont/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	screenWidth  = 800
	screenHeight = 500
)

var (
	screen      *ebiten.Image = ebiten.NewImage(Dimensions())
	spriteScale               = 1.0
)

func Setup() {
	ebiten.SetWindowResizable(true)
	loadSprites()
}

func DrawText(x, y int, inText string) {
	nx, ny := screenCords(float64(x), float64(y))
	// then center text
	textWidth := text.BoundString(bitmapfont.Face, inText).Dx()
	nx -= float64(textWidth) / 2
	text.Draw(screen, inText, bitmapfont.Face, int(nx), int(ny), color.White)
}

// need to implement this properly
func DrawTextUncentered(x, y int, inText string) {
	nx, ny := screenCords(float64(x), float64(y))
	text.Draw(screen, inText, bitmapfont.Face, int(nx), int(ny), color.White)
}

func Debug(text string) {
	// return
	ebitenutil.DebugPrint(screen, text)
}

func DisplaySprite(x, y int, sprite Sprite) {

	scaledX := float64(x*tileSize) * spriteScale
	scaledY := float64(y*tileSize) * spriteScale

	scaledX, scaledY = screenCords(scaledX, scaledY)

	// offset by sprite size since drawn from corner
	scaledX -= tileSize
	scaledY -= tileSize

	sprite.Options.GeoM.Scale(spriteScale, spriteScale)
	sprite.Options.GeoM.Translate(scaledX, scaledY)

	screen.DrawImage(sprite.Image, &sprite.Options)
}

func RawDisplaySprite(x, y int, scale float64, sprite Sprite) {

	nx, ny := screenCords(float64(x), float64(y))
	nx -= scale * tileSize / 2
	ny -= scale * tileSize / 2

	sprite.Options.GeoM.Scale(scale, scale)
	sprite.Options.GeoM.Translate(nx, ny)

	screen.DrawImage(sprite.Image, &sprite.Options)
}

func DisplaySprites(x, y int, sprites []Sprite) {
	sort.Slice(sprites, func(i, j int) bool {
		return sprites[i].Priority < sprites[j].Priority
	})

	for _, sprite := range sprites {
		DisplaySprite(x, y, sprite)
	}
}

func SpecialSetSpriteScale(total, new int) {
	spriteScale = float64(total) / float64(new*tileSize)
}

func screenCords(x, y float64) (float64, float64) {
	return x + screenWidth/2, y + screenHeight/2
}

func Clear() {
	screen = ebiten.NewImage(Dimensions())
}

func GetImage() *ebiten.Image {
	return screen
}

func Dimensions() (int, int) {
	return screenWidth, screenHeight
}
