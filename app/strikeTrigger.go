package app

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type strikeTrigger struct {
	s          *strike
	Press      *ebiten.Key
	obj        *resolv.Object
	action     TrapTrigger
	Image      *ebiten.Image
	ImageA     float32
	KillSingal chan interface{}
}

func NewTrigger(s *strike, key *ebiten.Key, obj *resolv.Object, Image *ebiten.Image, action TrapTrigger, RGBA float32) *strikeTrigger {
	return &strikeTrigger{
		s:          s,
		Press:      key,
		obj:        obj,
		action:     action,
		KillSingal: make(chan interface{}, 1),
		Image:      Image,
		ImageA:     RGBA,
	}
}
func (st *strikeTrigger) Process() {
	if st.obj != nil {
		go st.processObj()
		return
	}
	go st.processKeyListening()
}
func (st *strikeTrigger) processObj() {
	World.Add(st.obj)
	for {
		select {
		case <-time.After(time.Millisecond * 10):
		case <-st.KillSingal:
			waitKeepProcessing.Add(1)
			close(st.KillSingal)
			World.Remove(st.obj)
			st.obj = nil
			st.Image = nil
			waitKeepProcessing.Done()
			return
		}
		if co := st.obj.Check(0, 0, "character"); co != nil && st.collision(co) {
			st.s.Send(st.action)
			st.KillSingal <- 1
		}
	}
}
func (st *strikeTrigger) processKeyListening() {
	for {
		select {
		case <-st.KillSingal:
			waitKeepProcessing.Add(1)
			close(st.KillSingal)
			st.obj = nil
			st.Image = nil
			waitKeepProcessing.Done()
			return
		case <-time.After(time.Millisecond * 10):
		}
		if ebiten.IsKeyPressed(*st.Press) {
			st.s.Trigger <- st.action
			st.KillSingal <- 1
		}
	}
}

func (st *strikeTrigger) collision(co *resolv.Collision) bool {
	if st.obj.Shape != nil {
		if cos := st.obj.Shape.Intersection(0, 0, co.Objects[0].Shape); cos != nil {
			return true
		}
	} else if co.Objects[0].Shape != nil {
		return true
	}
	return false
}

func renderTrigger(tri map[any]*strikeTrigger, screen *ebiten.Image) {
	for _, t := range tri {
		t.Render(screen)
	}
}

// it must be a Object trigger if it Renders
func (st *strikeTrigger) Render(screen *ebiten.Image) {
	if st.obj != nil && st.Image != nil {
		geo := &ebiten.DrawImageOptions{}
		geo.GeoM.Translate(st.obj.X, st.obj.Y)
		geo.ColorScale.SetA(st.ImageA)
		screen.DrawImage(st.Image, geo)
	}
}

var trapTrigger1 = TrapTrigger{
	[]trapmovement{

		{
			Mode: 1,
		},
		{
			Mode: 4,
			Time: 12,
		},
		{
			Mode: 3,
			Time: 22,
			Movement: movementPlus{
				SizeY: 40,
				Y:     -500,
			},
		}, {
			Mode: 3,
			Time: 52,
			Movement: movementPlus{
				SizeY: 4,
				Y:     500,
			},
		},
		// {
		// 	Mode: 2,
		// 	Time: 40,
		// },
		// {
		// 	Mode: 1,
		// 	Time: 40,
		// },
		{
			Mode: 3,
			Time: 50,
			Movement: movementPlus{
				X: 640,
			},
		},
		{
			Mode: 3,
			Time: 200,
			Movement: movementPlus{
				X: -680,
			},
		}, {
			Mode: 3,
			Time: 20,
			Movement: movementPlus{
				SizeY: 40,
				Y:     -500,
			},
		}, {
			Mode: 3,
			Time: 200,
			Movement: movementPlus{
				SizeY: 4,
				Y:     500,
			},
		}, {
			Mode: 3,
			Time: 10,
			Movement: movementPlus{
				X: 680,
			},
		},
		{
			Mode: 3,
			Time: 25,
			Movement: movementPlus{
				X:     -670,
				Y:     -250 + 16,
				SizeX: 8,
				SizeY: 8,
			},
		},
		{
			Mode: 3,
			Time: 90,
			Movement: movementPlus{
				Y: -400,
			},
		}, {
			Mode: 3,
			Time: 30,
			Movement: movementPlus{
				Y: 440,
			},
		}, {
			Mode: 3,
			Time: 70,
			Movement: movementPlus{
				SizeX: 17,
			},
		}, {
			Mode: 4,
			Time: 50,
		}, {
			Mode: 3,
			Time: 100,
			Movement: movementPlus{
				SizeX: 20,
				SizeY: 2,
			},
		}, {
			Mode: 3,
			Time: 200,
			Movement: movementPlus{
				X:     300,
				Y:     150,
				SizeX: 6,
				SizeY: 4,
			},
		}, {
			Mode: 3,
			Time: 200,
			Movement: movementPlus{
				Angle: -90,
				Y:     100,
				SizeX: 2,
				SizeY: 2,
			},
		}, {
			Mode: 3,
			Time: 80,
			Movement: movementPlus{
				X: 300,
			},
		}, {
			Mode: 3,
			Time: 80,
			Movement: movementPlus{
				X: -800,
			},
		}, {
			Mode: 3,
			Time: 20,
			Movement: movementPlus{
				Angle:  90,
				SizeX:  2,
				SizeY:  2,
				SetPos: true,
				X:      200,
				Y:      200,
			},
		},
	},
}
