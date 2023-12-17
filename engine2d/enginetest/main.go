package main

import (
	"IJustWantToEscape/manager"
	"IJustWantToEscape/method"
	"IJustWantToEscape/objects"
	"IJustWantToEscape/utils"
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type Game struct {
	caller *sync.Once
	world  *resolv.Space
}

var testText = utils.Speech{
	&utils.Text{
		Size: utils.Normal,
		Text: "Normal",
	},
	&utils.Text{
		Size: utils.Big,
		Text: "Test",
	},
	&utils.Text{
		Size:  utils.Small,
		Text:  "Im SMall and red",
		Color: color.RGBA{255, 0, 0, 0},
		Endl:  true,
	},
	&utils.Text{
		Text:   "Another One!",
		Middle: true,
	},
	&utils.Text{
		Text:   "没中文的话我干嘛还自己导入字体",
		Middle: true,
	},
	&utils.Text{
		Text:          "Fonts by Yahei",
		StrikeThrough: true,
	},
	&utils.Text{
		Text:      "Underline",
		Underline: true,
		Size:      utils.Small,
	},
	&utils.Text{
		Text:          "BIGGGG",
		StrikeThrough: true,
		Underline:     true,
		Size:          utils.Big,
	},
}
var testObject = &objects.Structure{
	StatelessObject: method.StatelessObject{
		Sx: 1,
		Sy: 1,
		X:  100,
		Y:  100,
	},
}

// Draw implements ebiten.Game.
func (g *Game) Draw(screen *ebiten.Image) {
	if err := testObject.Draw(screen); err != nil {
		log.Fatal(err)
	}
	image, _ := utils.WriteString(testText)
	screen.DrawImage(image, nil)

}

// Layout implements ebiten.Game.
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return 1280, 960
}

// Update implements ebiten.Game.
func (g *Game) Update() error {
	return nil
}

func main() {
	manager.Init("./")
	var game = &Game{
		world:  resolv.NewSpace(640, 480, 64, 64),
		caller: &sync.Once{},
	}
	testObject.Init(game.world)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("(Interface) Test Engine")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
