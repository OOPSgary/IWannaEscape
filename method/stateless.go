package method

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type StatelessObject struct {
	X, Y, Sx, Sy float64
	Object       *resolv.Object
}

// Only use for interface
func (o *StatelessObject) Update() error { return nil }
func (o *StatelessObject) Init(space *resolv.Space) error {
	{
		if o.Object != nil {
			space.Add(o.Object)
		}
		return nil
	}
}

func (o *StatelessObject) Position() (x, y float64) { return o.X, o.Y }

func (o *StatelessObject) Draw(dst *ebiten.Image) error {
	panic("ShouldNotBeReached")
}
func (o *StatelessObject) DrawPictrue(src, dst *ebiten.Image) error {
	option := &ebiten.DrawImageOptions{}
	option.GeoM.Scale(o.Sx, o.Sy)
	option.GeoM.Translate(o.Position())
	dst.DrawImage(src, option)
	return nil
}
func (o *StatelessObject) Quit() error { return nil }
