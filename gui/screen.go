package gui

import (
	"image/color"
	"sort"

	"github.com/hajimehoshi/bitmapfont/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	screenWidth  = 800
	screenHeight = 500
	screenScale  = 100.0
)

var (
	screen *ebiten.Image = ebiten.NewImage(Dimensions())
)

func Setup() {
	ebiten.SetFullscreen(false)
	loadSprites()
}

func DrawText(x, y int, inText string) {
	x, y = screenCords(x, y)
	// then center text
	textWidth := text.BoundString(bitmapfont.Face, inText).Dx()
	x -= textWidth / 2
	text.Draw(screen, inText, bitmapfont.Face, x, y, color.White)
}

// need to implement this properly
func DrawTextUncentered(x, y int, inText string) {
	x, y = screenCords(x, y)
	text.Draw(screen, inText, bitmapfont.Face, x, y, color.White)
}

func Debug(text string) {
	// ebitenutil.DebugPrint(screen, text)
}

func DisplaySprite(x, y int, sprite Sprite) {

	x, y = screenCords(x, y)

	sprite.Options.GeoM.Scale(screenScale/100.0, screenScale/100.0)
	sprite.Options.GeoM.Translate(float64(x), float64(y))

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

func screenCords(x, y int) (int, int) {
	return x*tileSize*screenScale/100 + screenWidth/2, y*tileSize*screenScale/100 + screenHeight/2
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
