package app

import (
	"image/color"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/goki/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
	"golang.org/x/image/font"
)

type Game struct {
	HomePage  *ebitenui.UI
	mainWorld struct {
		MainCharacter Character
		Blocks        []*block
	}
	screen *ebiten.Image
	Wait   *sync.WaitGroup
}

var mainProcess *sync.WaitGroup = &sync.WaitGroup{}

func (g *Game) Update() error {
	mainProcess.Wait()
	switch Status {
	case 1:
		g.syncCharacterMovement()
		g.syncCharacter()
	}

	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	// var waitChannel <-chan time.Time
	// if !ebiten.IsVsyncEnabled() && ChangeFPSfps > 0 {
	// 	waitChannel = time.After(time.Second / time.Duration(ChangeFPSfps))
	// }
	g.screen = screen
	screen.DrawImage(NormalBackground, nil)
	g.syncExtraKeys()
	switch Status {
	case 0:
		//I JUST WANT TO ESCAPE TITLE

		if Dead {
			xg, yg := midImage(GameOver)
			opG := &ebiten.DrawImageOptions{}
			opG.GeoM.Scale(0.6, 0.6)
			opG.GeoM.Translate(float64(screen.Bounds().Max.X/2)-float64(xg)*0.6, float64(screen.Bounds().Max.Y/2)-float64(yg)*0.6)
			opG.ColorScale.SetA(startA)
			screen.DrawImage(GameOver, opG)
			if Keys.R && !rePressR {
				rePressR = true
				Status = 1
				mainProcess.Add(1)
			}
			if !Keys.R {
				rePressR = false
			}
			CallOnce("Trap-1", func() {
				deadSoundPlayer.Rewind()
				deadSoundPlayer.Play()
			})
		} else {
			drawMidTextLineByLine(0, 10, color.Black, screen, "Welcome to", "I JUST WANT TO ESCAPE", "按下F1可调整最大FPS")
			g.HomePage.Update()
			g.HomePage.Draw(screen)
		}
		if ebiten.IsKeyPressed(ebiten.KeyF1) && !ChaneFPSPressed {
			ChaneFPSPressed = true
			if !ChangeFPSMode {
				ChangeFPSMode = true
			} else {
				if ChangeFPSfps > 0 {
					ebiten.SetVsyncEnabled(false)
				}
				ChangeFPSMode = false
			}
		} else if !ebiten.IsKeyPressed(ebiten.KeyF1) {
			ChaneFPSPressed = false
		}
		if ChangeFPSMode {
			drawMidTextLineByLine(300, 5, color.Black, screen, "请输入目标FPS 再次按F1确认 按回车键重置", strconv.FormatInt(ChangeFPSfps, 10))
			ChangeFps()
		}
	case 1:
		if Keys.R && !rePressR {
			mainProcess.Add(1)
			rePressR = true
			CallExtra(0)
			CallExtra(1)
		}
		if !Keys.R {
			rePressR = false
		}
		CallOnce(0, func() {
			Dead = false
			g.resetMap()
			waitKeepProcessing.Wait()
			g.resetCharacter()
			g.moveCharacter(60, 40)
			g.box()
			g.putBlocksLine(movement{200, 480 - 3*16}, 0.5, 2, 2)
			strikeList[0] = NewStrike(movement{-10, 400}, 4, 4)
			strikeList[0].Load()
			to := resolv.NewObject(
				400, 400, 100, 100,
			)
			ti := ebiten.NewImage(100, 100)
			ti.Fill(color.Black)
			TriggerList[0] = NewTrigger(
				strikeList[0],
				nil,
				to,
				ti,
				trapTrigger1,
				0.2,
			)
			go TriggerList[0].Process()
			mainProcess.Done()
		})

		//Setting up main Wide-World

		drawMidTextLineByLine(
			10, 10, color.Black, screen,
			"康复训练:",
			"AD/<- ->左右移动",
			"空格键跳跃按R重开",
		)

		g.drawBox(screen)
		if !Dead {
			op := &ebiten.DrawImageOptions{}
			if g.mainWorld.MainCharacter.FaceAt == "l" {
				op.GeoM.Scale(-1, 1)
				op.GeoM.Translate(+25, 0)
			}
			op.GeoM.Translate(g.mainWorld.MainCharacter.Obj.X-4, g.mainWorld.MainCharacter.Obj.Y)

			screen.DrawImage(Kid[g.mainWorld.MainCharacter.Status], op)
			TriggerList[0].Render(screen)
			strikeList[0].Render(screen)
		} else {
			xg, yg := midImage(GameOver)
			opG := &ebiten.DrawImageOptions{}
			opG.GeoM.Scale(0.6, 0.6)
			opG.GeoM.Translate(float64(screen.Bounds().Max.X/2)-float64(xg)*0.6, float64(screen.Bounds().Max.Y/2)-float64(yg)*0.6)
			opG.ColorScale.SetA(startA)
			screen.DrawImage(GameOver, opG)
			CallOnce(1, func() {
				deadSoundPlayer.Rewind()
				deadSoundPlayer.Play()
				if !isChanClosed(strikeList[0].KillSingal) {
					strikeList[0].KillSingal <- 1
				}

			})
		}
	}

	ebitenutil.DebugPrint(screen, strconv.Itoa(int(ebiten.ActualFPS())))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

type Character struct {
	Status   int // 1,2,3 for three pictures 4 for dead but not completed
	Obj      *resolv.Object
	Top      *resolv.Object
	Button   *resolv.Object
	OnGround bool
	SpeedX   float64
	SpeedY   float64
	Jump     Jump
	FaceAt   string //It can be l(eft) or r(ight)

}
type Jump struct {
	Jump   int
	Chance int
	Lock   sync.Mutex
}

var Status int = 0

var Dead bool = false

var waitKeepProcessing = new(sync.WaitGroup)

var ChangeFPSMode bool
var ChangeFPSfps int64
var ChaneFPSPressed bool
var WaitTime time.Duration
var KeyPressed = make(map[ebiten.Key]bool)

func ChangeFps() {
	for _, k := range NumKeys {
		ChangeFPSfps = TypeNumer(k, ChangeFPSfps)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		ChangeFPSfps = 0
	}
}
func (g *Game) drawBox(screen *ebiten.Image) {
	for _, a := range g.mainWorld.Blocks {
		a.Draw(screen)
	}
}

// func newFill(w, h int, r, g, b, a uint8) (i *ebiten.Image) {
// 	i = ebiten.NewImage(w, h)
// 	i.Fill(color.RGBA{r, g, b, a})
// 	return
// }

func TypeNumer(key ebiten.Key, num int64) int64 {
	if key-ebiten.Key0 >= 0 || key-ebiten.Key0 <= 9 {
		if ebiten.IsKeyPressed(key) && !KeyPressed[key] {
			num = num*10 + int64(key-ebiten.Key0)
			KeyPressed[key] = true
		}
		if !ebiten.IsKeyPressed(key) {
			KeyPressed[key] = false
		}

	}
	return num
}

var NumKeys = []ebiten.Key{
	ebiten.Key0,
	ebiten.Key1,
	ebiten.Key2,
	ebiten.Key3,
	ebiten.Key4,
	ebiten.Key5,
	ebiten.Key6,
	ebiten.Key7,
	ebiten.Key8,
	ebiten.Key9,
}

type fontdata struct {
	lock  *sync.RWMutex
	fonts map[int]font.Face
}

var fontData fontdata = fontdata{
	lock:  new(sync.RWMutex),
	fonts: make(map[int]font.Face),
}

func (g *Game) resetCharacter() {
	obj := resolv.NewObject(0, 0, 17, 21-2*2, "character")
	obj.SetShape(resolv.NewConvexPolygon(
		0, 0,
		0, 0,
		17, 0,
		17, 21,
		0, 21,
	))
	top := resolv.NewObject(0, 0, 15, 1, "character")
	top.SetShape(resolv.NewConvexPolygon(
		0, 0,
		0, 0,
		15, 0,
		15, 1,
		0, 1,
	))
	button := resolv.NewObject(0, 0, 15, 1, "character")
	button.SetShape(resolv.NewConvexPolygon(
		0, 0,
		0, 0,
		15, 0,
		15, 1,
		0, 1,
	))
	World.Add(obj)
	World.Add(top)
	World.Add(button)
	g.mainWorld.MainCharacter = Character{
		Status: 1,
		Obj:    obj,
		Top:    top,
		Button: button,
	}
}
func (g *Game) moveCharacter(x, y float64) {
	// g.mainWorld.MainCharacter.Obj.X = x + 4
	// g.mainWorld.MainCharacter.Obj.Y = y + 19
	// g.mainWorld.MainCharacter.Top.X = x + 4
	// g.mainWorld.MainCharacter.Top.Y = y
	// g.mainWorld.MainCharacter.Button.X = x + 4
	// g.mainWorld.MainCharacter.Button.Y = y + 21
	g.mainWorld.MainCharacter.Obj.X = x + 4
	g.mainWorld.MainCharacter.Obj.Y = y + 2
	g.mainWorld.MainCharacter.Top.X = x + 5
	g.mainWorld.MainCharacter.Top.Y = y
	g.mainWorld.MainCharacter.Button.X = x + 5
	g.mainWorld.MainCharacter.Button.Y = y + 21
	g.mainWorld.MainCharacter.Obj.Update()
	g.mainWorld.MainCharacter.Top.Update()
	g.mainWorld.MainCharacter.Button.Update()
}

type strikeTrigger struct {
	s         *strike
	Press     *ebiten.Key
	obj       *resolv.Object
	action    TrapTrigger
	triggered bool
	Image     *ebiten.Image
	ImageA    float32
	close     chan any
}

func NewTrigger(s *strike, key *ebiten.Key, obj *resolv.Object, Image *ebiten.Image, action TrapTrigger, RGBA float32) *strikeTrigger {
	return &strikeTrigger{
		s:         s,
		Press:     key,
		obj:       obj,
		action:    action,
		triggered: false,
		close:     make(chan any),
		Image:     Image,
		ImageA:    RGBA,
	}
}
func (st *strikeTrigger) Process() {
	var Closed bool
	go func() {
		<-st.close
		Closed = true
	}()
	if st.obj != nil {
		World.Add(st.obj)
		for {
			tc := time.After(time.Millisecond * 10)
			if Closed {
				waitKeepProcessing.Add(1)
				close(st.close)
				World.Remove(st.obj)
				st.obj = nil
				st.Image = nil
				waitKeepProcessing.Done()
				break
			}
			if !st.triggered {
				if co := st.obj.Check(0, 0, "character"); co != nil {
					playerobj := co.Objects[0]
					if st.obj.Shape != nil && playerobj.Shape != nil {
						if cos := st.obj.Shape.Intersection(0, 0, co.Objects[0].Shape); cos != nil {
							st.triggered = true
							st.s.Send(st.action)
							st.close <- 1
						}
					} else {
						st.triggered = true
						st.s.Send(st.action)
						st.close <- 1
					}

				}
			}

			<-tc
		}

	} else {
		for {
			if Closed {
				close(st.close)
				break
			}
			tc := time.After(time.Millisecond * 10)
			if !st.triggered {
				if ebiten.IsKeyPressed(*st.Press) {
					st.triggered = true
					st.s.Trigger <- st.action
					st.close <- 1
				}
			}

			<-tc
		}

	}

}

// it must be a Object trigger if it Renders
func (st *strikeTrigger) Render(screen *ebiten.Image) {
	// make sure that im not stupid
	if st.obj != nil && st.Image != nil {
		geo := &ebiten.DrawImageOptions{}
		geo.GeoM.Translate(st.obj.X, st.obj.Y)
		geo.ColorScale.SetA(st.ImageA)
		screen.DrawImage(st.Image, geo)
	}
}

func (g *Game) syncCharacter() (RenderX, RenderY float64) {

	x, y := g.mainWorld.MainCharacter.SpeedX, g.mainWorld.MainCharacter.SpeedY
	if collision := g.mainWorld.MainCharacter.Obj.Check(x, y, "deadly"); collision != nil {
		if contactSet := g.mainWorld.MainCharacter.Obj.Shape.Intersection(x, y, collision.Objects[0].Shape); contactSet != nil {
			Dead = true
		}
	}

	if collision := g.mainWorld.MainCharacter.Obj.Check(x, y, "Stopper"); collision != nil {
		s := collision.SlideAgainstCell(collision.Cells[0], "Stopper")
		if s != nil {
			x = s.X()
		} else {
			x = 0
		}
		// x = collision.ContactWithObject(collision.Objects[0]).X()
		g.mainWorld.MainCharacter.SpeedX = 0
	}
	if collisionTop := g.mainWorld.MainCharacter.Top.Check(x, y, "Stopper"); collisionTop != nil {
		if dy := collisionTop.ContactWithObject(collisionTop.Objects[0]).Y(); dy != 0 {
			y = dy
			g.mainWorld.MainCharacter.SpeedY = 0
		}
	}
	if collisionButton := g.mainWorld.MainCharacter.Button.Check(x, y, "Stopper"); collisionButton != nil && g.mainWorld.MainCharacter.SpeedY >= 0 && g.mainWorld.MainCharacter.Button.Bottom() < collisionButton.Objects[0].Y+2 {

		g.mainWorld.MainCharacter.Jump.Reset()
		if dy := collisionButton.ContactWithObject(collisionButton.Objects[0]).Y(); dy != 0 {
			y = dy
			g.mainWorld.MainCharacter.SpeedY = 0
		}

	} else {

		//if collision := g.mainWorld.MainCharacter.Button.Check(x, y+1, "Stopper"); collision == nil
		if g.mainWorld.MainCharacter.SpeedY <= 4 && g.mainWorld.MainCharacter.SpeedY >= 0 {
			g.mainWorld.MainCharacter.SpeedY += 0.15

		} else if g.mainWorld.MainCharacter.SpeedY <= 0 {
			g.mainWorld.MainCharacter.SpeedY += 0.2
		}
		// y = g.mainWorld.MainCharacter.SpeedY

	}

	g.mainWorld.MainCharacter.Obj.X += x
	g.mainWorld.MainCharacter.Obj.Y += y

	g.mainWorld.MainCharacter.Top.X += x
	g.mainWorld.MainCharacter.Top.Y += y
	g.mainWorld.MainCharacter.Button.X += x
	g.mainWorld.MainCharacter.Button.Y += y
	g.mainWorld.MainCharacter.Obj.Update()
	g.mainWorld.MainCharacter.Top.Update()
	g.mainWorld.MainCharacter.Button.Update()
	return g.mainWorld.MainCharacter.Obj.X - 4, g.mainWorld.MainCharacter.Obj.Y
}

func (g *Game) syncCharacterMovement() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || Keys.A {
		g.mainWorld.MainCharacter.SpeedX = -2
		if g.mainWorld.MainCharacter.Status < 2 {
			g.mainWorld.MainCharacter.Status++
		} else {
			g.mainWorld.MainCharacter.Status = 1
		}
		g.mainWorld.MainCharacter.FaceAt = "l"
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || Keys.D {
		g.mainWorld.MainCharacter.SpeedX = 2
		if g.mainWorld.MainCharacter.Status < 2 {
			g.mainWorld.MainCharacter.Status++
		} else {
			g.mainWorld.MainCharacter.Status = 1
		}
		g.mainWorld.MainCharacter.FaceAt = "r"
	} else {
		if g.mainWorld.MainCharacter.SpeedX >= 1 {
			g.mainWorld.MainCharacter.SpeedX -= 0.5
			if g.mainWorld.MainCharacter.Status < 2 {
				g.mainWorld.MainCharacter.Status++
			} else {
				g.mainWorld.MainCharacter.Status = 1
			}
		} else if g.mainWorld.MainCharacter.SpeedX <= -1 {
			g.mainWorld.MainCharacter.SpeedX += 0.5
			if g.mainWorld.MainCharacter.Status < 2 {
				g.mainWorld.MainCharacter.Status++
			} else {
				g.mainWorld.MainCharacter.Status = 1
			}
		} else {
			g.mainWorld.MainCharacter.Status = 1
			g.mainWorld.MainCharacter.SpeedX = 0
		}
	}
	if Keys.Space {
		if g.mainWorld.MainCharacter.Jump.Update() {
			g.mainWorld.MainCharacter.SpeedY = -4
		} else if g.mainWorld.MainCharacter.Jump.AddChance() {
			g.mainWorld.MainCharacter.SpeedY -= 0.05
		}
	}
}

var strikeList map[int]*strike = make(map[int]*strike)
var TriggerList = make(map[any]*strikeTrigger)
var rePressR bool

func (g *Game) resetMap() {
	for _, o := range g.mainWorld.Blocks {
		World.Remove(o.Obj)
	}
	g.mainWorld.Blocks = nil
	for _, s := range strikeList {
		if !isChanClosed(s.KillSingal) {
			s.KillSingal <- 1
		}
	}
	strikeList = make(map[int]*strike)
	for _, s := range TriggerList {
		if !isChanClosed(s.close) {
			s.close <- 1
		}

	}
	TriggerList = make(map[any]*strikeTrigger)

}
func isChanClosed(c chan any) bool {
	select {
	case _, ok := <-c:
		return !ok
	default:
		return false
	}
}
func GetFont(size int) font.Face {
	fontData.lock.RLock()
	if fontData.fonts[size] != nil {
		defer fontData.lock.RUnlock()
		return fontData.fonts[size]
	}
	fontData.lock.RUnlock()

	face := truetype.NewFace(basicFont, &truetype.Options{
		Size:    float64(size),
		DPI:     300,
		Hinting: font.HintingFull,
	})
	fontData.lock.Lock()
	fontData.fonts[size] = face
	fontData.lock.Unlock()
	return face
}

// 角度->弧度 转换
func getRadian(turn float64) float64 {
	return math.Pi / 180 * turn
}

// 制作针对文字的位移
// Xs,Ys为横向纵向拉伸程序
func makeGeo(X, Y, Xs, Ys, turn float64, c color.Color) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(getRadian(turn))
	op.GeoM.Scale(Xs, Ys)
	op.GeoM.Translate(X, Y)
	if c != nil {
		op.ColorScale.ScaleWithColor(c)
	}

	return op
}

var World = resolv.NewSpace(640, 480, 16, 16)
var Keys struct {
	R     bool
	Space bool
	A, D  bool
}

func (g *Game) syncExtraKeys() {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		Keys.R = true
	} else {
		Keys.R = false
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		Keys.Space = true
	} else {
		Keys.Space = false
		g.mainWorld.MainCharacter.Jump.ResetChance()
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		Keys.A = true
	} else {
		Keys.A = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		Keys.D = true
	} else {
		Keys.D = false
	}
}

// but it must use as TriggerList
