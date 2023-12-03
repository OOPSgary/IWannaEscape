package Object

import (
	"IJustWantToEscape/manager"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type StatelessObject struct {
	X, Y, Sx, Sy float64
	TargetImage  string
	Object       *resolv.Object
}

// Only use for interface
func (o *StatelessObject) Update() error { return nil }
func (o *StatelessObject) Run(space *resolv.Space) error {
	{
		if o.Object != nil {
			space.Add(o.Object)
		}
		return nil
	}
}

func (o *StatelessObject) Position() (x, y float64) { return o.X, o.Y }

func (o *StatelessObject) Draw(dst *ebiten.Image) error {
	if o.TargetImage == "" {
		return nil
	}
	option := &ebiten.DrawImageOptions{}
	option.GeoM.Translate(o.Position())
	option.GeoM.Scale(o.Sx, o.Sy)
	i, _, err := manager.Manager.GetImage(o.TargetImage)
	if err != nil {
		return err
	}
	i.DrawImage(dst, option)
	return nil
}
func (o *StatelessObject) Quit() error
