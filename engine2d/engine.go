package engine2d

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
}

const (
	ScreenHeight = 1080
	ScreenWidth  = 1960
)

func (g *Game) Update() error {
	return nil
}
func (g *Game) Draw(dst *ebiten.Image) {

}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
