package app

import (
	"image/color"
	"log"
	"time"

	"github.com/adrg/sysfont"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

type strike struct {
	Pos          movement
	SizeX, SizeY float64
	Angle        float64
	Shape        *resolv.ConvexPolygon
	Trigger      chan TrapTrigger
	KillSingal   chan interface{}
	Online       bool //Not using this again
}
type TrapTrigger struct {
	Movement []trapmovement
}

const (
	TrapAppear = iota + 1
	TrapDisappear
	TrapMove
	TrapSleep
)

type trapmovement struct {
	Mode int
	//10MillCounter
	Time     float64
	Movement movementPlus
	//Any Function in here will be process after the Sleep/Movement
	//So it's better to just use TrapSleep to excuse that
	Func func()
}
type movementPlus struct {
	Angle        float64
	X, Y         float64
	SetPos       bool
	SizeX, SizeY float64
}

func NewStrike(Pos movement, SizeX, SizeY float64) *strike {

	shape := resolv.NewConvexPolygon(
		0, 0,
		32/2, 0,
		0, 32,
		32, 32,
	)
	shape.SetPosition(Pos.x, Pos.y)
	shape.SetScale(SizeX, SizeY)
	return &strike{
		Pos:        Pos,
		SizeX:      SizeX,
		SizeY:      SizeY,
		Angle:      0,
		Shape:      shape,
		Trigger:    make(chan TrapTrigger),
		KillSingal: make(chan interface{}),
		Online:     false,
	}

}
func renderStrikes(s map[any]*strike, screen *ebiten.Image) {
	for _, singleS := range s {
		if singleS != nil {
			singleS.Render(screen)
		}
	}
}
func cleanStrikes(s map[any]*strike) {
	for _, s := range strikeList {
		if !isChanClosed(s.KillSingal) {
			s.KillSingal <- 1
		}
	}
	strikeList = make(map[any]*strike)
}
func (s *strike) Render(screen *ebiten.Image) {
	if s.Online {
		geo := &ebiten.DrawImageOptions{}
		geo.GeoM.Scale(s.SizeX, s.SizeY)
		geo.GeoM.Rotate(getRadian(s.Angle))
		geo.GeoM.Translate(s.Pos.x, s.Pos.y)
		screen.DrawImage(StrikePhoto, geo)

		// Try to underSTAND the rotate between Resolv and Ebitengine
		shape := s.Shape
		verts := shape.Transformed()
		for i := 0; i < len(verts); i++ {
			vert := verts[i]
			next := verts[0]

			if i < len(verts)-1 {
				next = verts[i+1]
			}
			vector.StrokeLine(screen, float32(vert.X()), float32(vert.Y()), float32(next.X()), float32(next.Y()), 2, color.RGBA{255, 0, 0, 1}, true)

			// ebitenutil.DrawLine(screen, vert.X(), vert.Y(), next.X(), next.Y(), color.RGBA{255, 0, 0, 1})

		}
	}

}
func (s *strike) Send(t TrapTrigger) error {
	select {
	case <-s.KillSingal:
		return errSendOnClosedChannel
	default:
		s.Trigger <- t
	}
	return nil
}
func (s *strike) Load() {
	go s.load()
}
func (s *strike) load() {
	for {
		if s.process() {
			waitKeepProcessing.Add(1)
			close(s.Trigger)
			if !isChanClosed(s.KillSingal) {
				close(s.KillSingal)
			}
			s.Online = false
			waitKeepProcessing.Done()
			break
		}
	}
}
func (s *strike) process() (stop bool) {
	var a TrapTrigger
	select {
	case a = <-s.Trigger:
	case <-s.KillSingal:
		return true
	}
	sysfont.NewFinder(nil).List()
	for _, action := range a.Movement {
		log.Println("Excuse movement ", action)
		go newSoundPlayer(strikeSound).Play()
		var stop bool
		switch action.Mode {
		case 1:
			stop = s.handlerAppear(action)
		case 2:
			stop = s.handlerDisAppear(action)
		case 3:
			stop = s.handlerMovement(action)
		case 4:
			stop = s.handlerWaiting(action)
		}
		if stop {
			return true
		}
	}
	return false
}
func (s *strike) handlerAppear(action trapmovement) (stop bool) {
	select {
	case <-time.After(time.Millisecond * 10 * time.Duration(action.Time)):
		s.Online = true
		if action.Func != nil {
			action.Func()
		}
		return false
	case <-s.KillSingal:
		return true
	}
}
func (s *strike) handlerDisAppear(action trapmovement) (stop bool) {
	select {
	case <-time.After(time.Millisecond * 10 * time.Duration(action.Time)):
		s.Online = false
		if action.Func != nil {
			action.Func()
		}
		return false
	case <-s.KillSingal:
		return true
	}
}
func (s *strike) handlerMovement(action trapmovement) (stop bool) {
	PreferData := struct {
		sizeX, sizeY float64
		PosX, PosY   float64
		Angle        float64
	}{
		sizeX: s.SizeX,
		sizeY: s.SizeY,
		PosX:  s.Pos.x,
		PosY:  s.Pos.y,
		Angle: s.Angle,
	}
	delayProcess := func() {
		if action.Movement.SetPos {
			dx := (action.Movement.X - PreferData.PosX)
			dy := (action.Movement.Y - PreferData.PosY)
			s.Pos.x += dx / float64(action.Time)
			s.Pos.y += dy / float64(action.Time)
		} else {
			s.Pos.x += action.Movement.X / float64(action.Time)
			s.Pos.y += action.Movement.Y / float64(action.Time)
		}

		s.SizeX += (ifPositiveNum(s.SizeX, action.Movement.SizeX) - PreferData.sizeX) / float64(action.Time)
		s.SizeY += (ifPositiveNum(s.SizeY, action.Movement.SizeY) - PreferData.sizeY) / float64(action.Time)
		s.Angle += action.Movement.Angle / action.Time
		s.Shape.SetScale(s.SizeX, s.SizeY)
		s.Shape.SetPosition(s.Pos.x, s.Pos.y)
		s.Shape.SetRotation(-getRadian(s.Angle))
	}
	if action.Time <= 0 {
		if action.Movement.SetPos {
			s.Pos.x = action.Movement.X
			s.Pos.y = action.Movement.Y
		} else {
			s.Pos.x += action.Movement.X
			s.Pos.y += action.Movement.Y
		}
		s.SizeX = ifPositiveNum(s.SizeX, action.Movement.SizeX)
		s.SizeY = ifPositiveNum(s.SizeY, action.Movement.SizeY)
		s.Angle += action.Movement.Angle
		s.Shape.ScaleW = s.SizeX
		s.Shape.ScaleH = s.SizeY
		s.Shape.SetPosition(s.Pos.x, s.Pos.y)
		s.Shape.SetRotation(-getRadian(s.Angle))
		if action.Func != nil {
			action.Func()
		}
	} else {
		for i := float64(0); i <= action.Time; i++ {
			select {
			case <-time.After(time.Millisecond * 10):
				delayProcess()
			case <-s.KillSingal:
				return true
			}
		}
		s.Angle = PreferData.Angle + action.Movement.Angle
		s.Pos.x = PreferData.PosX + action.Movement.X
		s.Pos.y = PreferData.PosY + action.Movement.Y
		s.Shape.SetRotation(-getRadian(s.Angle))
		s.Shape.SetPosition(s.Pos.x, s.Pos.y)
		if action.Func != nil {
			action.Func()
		}
	}
	return false
}
func (s *strike) handlerWaiting(action trapmovement) (stop bool) {
	select {
	case <-s.KillSingal:
		return true
	case <-time.After(time.Millisecond * 10 * time.Duration(action.Time)):
		if action.Func != nil {
			action.Func()
		}
		return false
	}
}
func ifPositiveNum(value, replaceValue float64) float64 {
	if replaceValue > 0 {
		return replaceValue
	}
	return value
}
