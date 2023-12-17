package method

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type Object interface {
	Position() (x, y float64)
	Init(world *resolv.Space) error
	Draw(dst *ebiten.Image)
	Update() error
}

/*
type Object interface {
	Movement
	Render
	Collision
}
type Movement interface {
	Position(X, Y float64) error
	SetPosition(X, Y float64) error
	Angle(degree float64) error
	SetAngle(degree float64) error
	CenterAngle(degree float64) error
	SetCenterAngle(degree float64) error
	Scale(xt, yt float64) error
	ReScale(xt, yt float64) error
	Size(x, y float64) error
	ReSize(x, y float64) error
	X() (X float64)
	Y() (Y float64)
	H() (H float64)
	W() (W float64)
	Sx() (Sx float64)
	Sy() (Sy float64)
}

type Collision interface {
	Join() error
	Leave() error
	SetTag(tag ...string) error
	Tag() (tag []string)
	AddTag(tag ...string) error
	Check(tag ...string) *resolv.Collision
	CheckWithSpeed(dx, dy float64, tag ...string) *resolv.Collision
	Object() *resolv.Object
}
type Render interface {
	Draw(I *ebiten.Image) error
	DrawRPG(I *ebiten.Image) error
	DrawWithOpt(I *ebiten.Image, opt *ebiten.DrawImageOptions) error
	DrawRPGWithOpt(I *ebiten.Image, opt *ebiten.DrawImageOptions) error
	Invis(status bool)
	RPGFlightMode(h float64)
	Status() (RenderImage string)
	RPGStatus() (RenderImage string)
	DiscuzPhoto(SignName string) *ebiten.Image
	DefaultDiscuzPhoto(SignName string) *ebiten.Image
}
*/
