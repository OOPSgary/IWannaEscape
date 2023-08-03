package main

import (
	"IJustWantToEscape/app"
	"image/color"
	"log"
	"sync"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	buttonImage, _ := loadButtonImage()
	f := app.GetFont(4)
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Start the game", f, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   15,
			Right:  15,
			Top:    2,
			Bottom: 2,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			app.Dead = true
		}),
	)
	rootContainer.AddChild(button)
	log.Println("IJustWantTOEscape")
	ebiten.SetWindowTitle("IJustwantToEscape")
	ebiten.SetWindowSize(640, 480)
	ebiten.SetTPS(120)
	ebiten.SetVsyncEnabled(false)
	go func() {
		if ebiten.IsFocused() {

		} else {

		}
	}()
	var NewGame *app.Game = &app.Game{
		Wait: &sync.WaitGroup{},
	}
	NewGame.HomePage = &ebitenui.UI{
		Container: rootContainer,
	}
	if err := ebiten.RunGame(NewGame); err != nil {
		log.Fatal(err)
	}

}
func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}
