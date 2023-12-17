package objects

import (
	"IJustWantToEscape/manager"
	"IJustWantToEscape/method"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

// 作为Stateless代表
type Structure struct{ method.StatelessObject }

func init() {
	manager.Manager.Sign("structure.png", "structure", manager.ImageFile)
}
func (s *Structure) Draw(dst *ebiten.Image) error {
	src, _, _ := manager.Manager.GetImage("structure")
	return s.StatelessObject.DrawPictrue(src, dst)
}
func (s *Structure) Init(world *resolv.Space) error {
	s.Object = resolv.NewObject(s.X, s.Y, s.Sx*64, s.Sy*64, "stateless")
	return s.StatelessObject.Init(world)
}
