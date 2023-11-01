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

type Strike struct {
	Pos          movement
	SizeX, SizeY float64
	Angle        float64
	Shape        *resolv.ConvexPolygon
	Trigger      chan TrapTrigger
	KillSignal   chan interface{}
	Online       bool //Not using this again
}
type TrapTrigger struct {
	Movement []trapMovement
}

const (
	AppearTrap = iota + 1
	DisappearTrap
	MoveTrap
	SleepTrap
)

type trapMovement struct {
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

func NewStrike(Pos movement, SizeX, SizeY float64) *Strike {

	shape := resolv.NewConvexPolygon(
		0, 0,
		32/2, 0,
		0, 32,
		32, 32,
	)
	shape.SetPosition(Pos.x, Pos.y)
	shape.SetScale(SizeX, SizeY)
	return &Strike{
		Pos:        Pos,
		SizeX:      SizeX,
		SizeY:      SizeY,
		Angle:      0,
		Shape:      shape,
		Trigger:    make(chan TrapTrigger),
		KillSignal: make(chan interface{}),
		Online:     false,
	}

}
func renderStrikes(s map[any]*Strike, screen *ebiten.Image) {
	for _, singleS := range s {
		if singleS != nil {
			singleS.Render(screen)
		}
	}
}
func cleanStrikes(s map[any]*Strike) {
	for _, s := range s {
		if !isChanClosed(s.KillSignal) {
			s.KillSignal <- 1
		}
	}
	strikeList = make(map[any]*Strike)
}
func (s *Strike) Render(screen *ebiten.Image) {
	if s.Online {
		geo := &ebiten.DrawImageOptions{}
		geo.GeoM.Scale(s.SizeX, s.SizeY)
		geo.GeoM.Rotate(getRadian(s.Angle))
		geo.GeoM.Translate(s.Pos.x, s.Pos.y)
		screen.DrawImage(StrikePhoto, geo)

		shape := s.Shape
		Vert := shape.Transformed()
		for i := 0; i < len(Vert); i++ {
			vert := Vert[i]
			next := Vert[0]

			if i < len(Vert)-1 {
				next = Vert[i+1]
			}
			vector.StrokeLine(screen, float32(vert.X()), float32(vert.Y()), float32(next.X()), float32(next.Y()), 2, color.RGBA{R: 255, A: 1}, true)

			// ebitenutil.DrawLine(screen, vert.X(), vert.Y(), next.X(), next.Y(), color.RGBA{255, 0, 0, 1})

		}
	}

}
func (s *Strike) Send(t TrapTrigger) error {
	select {
	case <-s.KillSignal:
		return errSendOnClosedChannel
	case s.Trigger <- t:
	}
	return nil
}
func (s *Strike) Load() {
	go s.load()
}
func (s *Strike) load() {
	for {
		if s.process() {
			waitKeepProcessing.Add(1)
			close(s.Trigger)
			if !isChanClosed(s.KillSignal) {
				close(s.KillSignal)
			}
			s.Online = false
			waitKeepProcessing.Done()
			break
		}
	}
}
func (s *Strike) process() (stop bool) {
	var a TrapTrigger
	select {
	case a = <-s.Trigger:
	case <-s.KillSignal:
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
func (s *Strike) handlerAppear(action trapMovement) (stop bool) {
	select {
	case <-time.After(time.Millisecond * 10 * time.Duration(action.Time)):
		s.Online = true
		if action.Func != nil {
			action.Func()
		}
		return false
	case <-s.KillSignal:
		return true
	}
}
func (s *Strike) handlerDisAppear(action trapMovement) (stop bool) {
	select {
	case <-time.After(time.Millisecond * 10 * time.Duration(action.Time)):
		s.Online = false
		if action.Func != nil {
			action.Func()
		}
		return false
	case <-s.KillSignal:
		return true
	}
}
func (s *Strike) handlerMovement(action trapMovement) (stop bool) {
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
			dx := action.Movement.X - PreferData.PosX
			dy := action.Movement.Y - PreferData.PosY
			s.Pos.x += dx / action.Time
			s.Pos.y += dy / action.Time
		} else {
			s.Pos.x += action.Movement.X / action.Time
			s.Pos.y += action.Movement.Y / action.Time
		}

		s.SizeX += (ifPositiveNum(s.SizeX, action.Movement.SizeX) - PreferData.sizeX) / action.Time
		s.SizeY += (ifPositiveNum(s.SizeY, action.Movement.SizeY) - PreferData.sizeY) / action.Time
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
			case <-s.KillSignal:
				return true
			}
		}
		if action.Movement.SetPos {
			s.Angle = action.Movement.Angle
			s.Pos.x = action.Movement.X
			s.Pos.y = action.Movement.Y
		} else {
			s.Angle = PreferData.Angle + action.Movement.Angle
			s.Pos.x = PreferData.PosX + action.Movement.X
			s.Pos.y = PreferData.PosY + action.Movement.Y
		}

		s.Shape.SetRotation(-getRadian(s.Angle))
		s.Shape.SetPosition(s.Pos.x, s.Pos.y)
		if action.Func != nil {
			action.Func()
		}
	}
	return false
}
func (s *Strike) handlerWaiting(action trapMovement) (stop bool) {
	select {
	case <-s.KillSignal:
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
