package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type portal struct {
	AtLevel      int
	ToLevel      int
	Active       bool
	X, Y         float64
	SizeX, SizeY float64
	Fake         bool
	Object       *resolv.Object
}

var currentPortal []*portal

func newCurrPortal(To int, X, Y, SizeX, SizeY float64, fake bool) []*portal {
	currentPortal = []*portal{{
		AtLevel: Status,
		ToLevel: To,
		X:       X,
		Y:       Y,
		SizeX:   SizeX,
		SizeY:   SizeY,
		Fake:    fake,
		Object:  resolv.NewObject(X, Y, 42, 45, "portal"),
	},
	}
	return currentPortal
}

func newPortal(AtLevel, ToLevel int, X, Y, SizeX, SizeY float64, fake bool) *portal {
	return &portal{
		AtLevel: AtLevel,
		ToLevel: ToLevel,
		X:       X,
		Y:       Y,
		SizeX:   SizeX,
		SizeY:   SizeY,
		Fake:    fake,
		Object:  resolv.NewObject(X, Y, 42, 45, "portal"),
	}
}
func (g *Game) checkCurrentPortals() bool {
	for _, portals := range currentPortal {
		if portals.Check(g) {
			Status = portals.ToLevel
			goto Collision
		}
	}
	return false
Collision:
	currentPortal = currentPortal[:0]
	return true
}
func (p *portal) SetExist(Show bool) {
	if Show && !p.Active {
		p.Active = true
		World.Add(p.Object)
	} else if !Show && p.Active {
		p.Active = false
		World.Remove(p.Object)
	}
}
func (p *portal) SetCurrent() {
	currentPortal = []*portal{p}
}
func (p *portal) AddCurrent() {
	currentPortal = append(currentPortal, p)
}
func (p *portal) Draw(screen *ebiten.Image) {
	screen.DrawImage(portalImage, makeGeo(p.X, p.Y, p.SizeX, p.SizeY, 0, nil))
}
func (p *portal) Check(g *Game) bool {
	if p.Active && p.Object.Check(g.mainWorld.MainCharacter.SpeedX, g.mainWorld.MainCharacter.SpeedY, "character") != nil {
		Status = p.ToLevel
		p.SetExist(false)
		return true
	}
	return false
}
