package gui

import (
	"image/color"
	"log"
	"sort"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 1000
	screenHeight = 500
	screenScale  = 50.0
)

var (
	screen *ebiten.Image
)

var mplusNormalFont font.Face

func Setup() {
	ebiten.SetFullscreen(false)
	loadSprites()
	loadFont()
}

func loadFont() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 80
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// need to implement this properly
func DrawText(x, y int, inText string) {
	x, y = screenCords(x, y)
	text.Draw(screen, inText, mplusNormalFont, x, y, color.White)
	// ebitenutil.DebugPrint(Screen, text)
}

var debugString = ""

func Debug(text string) {
	debugString += text + "\n"
	// ebitenutil.DebugPrint(screen, debugString)
}

func DisplayXY(x, y int, image *ebiten.Image, op *ebiten.DrawImageOptions) {

	x, y = screenCords(x, y)

	op.GeoM.Scale(screenScale/100.0, screenScale/100.0)
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(image, op)
}

func DisplaySprites(x, y int, sprites []Sprite) {
	sort.Slice(sprites, func(i, j int) bool {
		return sprites[i].Priority < sprites[j].Priority
	})

	for _, sprite := range sprites {
		op := &ebiten.DrawImageOptions{}
		DisplayXY(x, y, sprite.Image, op)
	}
}

func screenCords(x, y int) (int, int) {
	return x*tileSize*screenScale/100 + screenWidth/2, y*tileSize*screenScale/100 + screenHeight/2
}

func Clear() {
	screen = ebiten.NewImage(Dimensions())
	debugString = strconv.Itoa(int(ebiten.Key('e'))) + "\n"
}

func GetImage() *ebiten.Image {
	return screen
}

func Dimensions() (int, int) {
	return screenWidth, screenHeight
}
