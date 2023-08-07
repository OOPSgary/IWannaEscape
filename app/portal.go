package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type portal struct {
	Atlevel      int
	Tolevel      int
	Visibal      bool
	X, Y         float64
	SizeX, SizeY float64
	Fake         bool
	Object       *resolv.Object
}

var currentPortal *portal

func newCurrPortal(Tolevel int, X, Y, SizeX, SizeY float64, fake bool) *portal {
	currentPortal = &portal{
		Atlevel: Status,
		Tolevel: Tolevel,
		X:       X,
		Y:       Y,
		SizeX:   SizeX,
		SizeY:   SizeY,
		Fake:    fake,
		Object:  resolv.NewObject(X, Y, 42, 45, "portal"),
	}
	return currentPortal
}

func newPortal(Atlevel, Tolevel int, X, Y, SizeX, SizeY float64, fake bool) *portal {
	return &portal{
		Atlevel: Status,
		Tolevel: Tolevel,
		X:       X,
		Y:       Y,
		SizeX:   SizeX,
		SizeY:   SizeY,
		Fake:    fake,
		Object:  resolv.NewObject(X, Y, 42, 45, "portal"),
	}
}
func (p *portal) SetExist(Show bool) {
	if Show && !p.Visibal {
		p.Visibal = true
		World.Add(p.Object)
	} else if !Show && p.Visibal {
		p.Visibal = false
		World.Remove(p.Object)
	}
}
func (p *portal) SetCurrent() {
	currentPortal = p
}
func (p *portal) print(screen *ebiten.Image) {
	screen.DrawImage(portalImage, makeGeo(p.X, p.Y, p.SizeX, p.SizeY, 0, nil))
}
func (p *portal) Check(g *Game) bool {
	if p.Visibal && p.Object.Check(g.mainWorld.MainCharacter.SpeedX, g.mainWorld.MainCharacter.SpeedY, "character") != nil {
		Status = p.Tolevel
		p.SetExist(false)
		return true
	}
	return false
}
