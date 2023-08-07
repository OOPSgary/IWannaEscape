package app

import (
	"fmt"
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

func (g *Game) Update() error {
	switch Status {
	case 1:
		if g.mainWorld.MainCharacter.Obj != nil {
			g.syncCharacterMovement()
			g.syncCharacter()
		}

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
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
			}
			if !Keys.R {
				rePressR = false
			}
			CallOnce("Trap-1", func() {
				deadSoundPlayer.Rewind()
				deadSoundPlayer.Play()
			})
		} else {

			drawMidTextLineByLine(0, 10, color.Black, screen, "Welcome to", "I JUST WANT TO ESCAPE", "F1: VSync模式切换 F11: 全屏")

			g.HomePage.Update()
			g.HomePage.Draw(screen)
		}

	case 1:
		if Keys.R && !rePressR {
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
			g.characterJoin()
			go TriggerList[0].Process()
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
	Vsync := ebiten.IsVsyncEnabled()
	if Keys.F1.Press && !Keys.F1.Pressed {
		Keys.F1.Pressed = true
		if Vsync {
			ebiten.SetVsyncEnabled(false)
		} else {
			ebiten.SetVsyncEnabled(true)
		}
	}
	fps := strconv.Itoa(int(ebiten.ActualFPS()))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS:%v VSync:%t", fps, Vsync))
	if Keys.F11.Press && !Keys.F11.Pressed {
		Keys.F11.Pressed = true
		if ebiten.IsFullscreen() {
			ebiten.SetFullscreen(false)
		} else {
			ebiten.SetFullscreen(true)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

var Status int = 0
var Dead bool = false
var waitKeepProcessing = new(sync.WaitGroup)
var WaitTime time.Duration
var KeyPressed = make(map[ebiten.Key]bool)

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
	if g.mainWorld.MainCharacter.Obj != nil {
		g.characterLeave()
	}

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
	g.mainWorld.MainCharacter = Character{
		Status: 1,
		Obj:    obj,
		Top:    top,
		Button: button,
	}
}
func (g *Game) moveCharacter(x, y float64) {
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
func (g *Game) characterJoin() {
	World.Add(g.mainWorld.MainCharacter.Obj)
	World.Add(g.mainWorld.MainCharacter.Top)
	World.Add(g.mainWorld.MainCharacter.Button)
}
func (g *Game) characterLeave() {
	World.Remove(g.mainWorld.MainCharacter.Obj)
	World.Remove(g.mainWorld.MainCharacter.Top)
	World.Remove(g.mainWorld.MainCharacter.Button)
}

func (g *Game) syncCharacter() (RenderX, RenderY float64) {

	x, y := g.mainWorld.MainCharacter.SpeedX, g.mainWorld.MainCharacter.SpeedY
	// if collision := g.mainWorld.MainCharacter.Obj.Check(x, y, "deadly"); collision != nil {
	// 	if contactSet := g.mainWorld.MainCharacter.Obj.Shape.Intersection(x, y, collision.Objects[0].Shape); contactSet != nil {
	// 		Dead = true
	// 	}
	// }
	for _, sl := range strikeList {
		if sl.Online {
			if contactSet := g.mainWorld.MainCharacter.Obj.Shape.Intersection(x, y, sl.Shape); contactSet != nil {
				Dead = true
			}
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
		if g.mainWorld.MainCharacter.SpeedY <= 5 && g.mainWorld.MainCharacter.SpeedY >= 0 {
			g.mainWorld.MainCharacter.SpeedY += 0.25

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
			g.mainWorld.MainCharacter.SpeedX -= 1
			if g.mainWorld.MainCharacter.Status < 2 {
				g.mainWorld.MainCharacter.Status++
			} else {
				g.mainWorld.MainCharacter.Status = 1
			}
		} else if g.mainWorld.MainCharacter.SpeedX <= -1 {
			g.mainWorld.MainCharacter.SpeedX += 1
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

var strikeList = make(map[int]*strike)

//	var strikeList = struct {
//		M       *map[int]*strike
//		RWMutex *sync.RWMutex
//	}{
//
//		M:       new(map[int]*strike),
//		RWMutex: &sync.RWMutex{},
//	}
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
		if !isChanClosed(s.KillSingal) {
			s.KillSingal <- 1
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
	F1    struct {
		Press   bool
		Pressed bool
	}
	F11 struct {
		Press   bool
		Pressed bool
	}
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
	if ebiten.IsKeyPressed(ebiten.KeyF1) && !Keys.F1.Pressed {
		Keys.F1.Press = true
	} else if !ebiten.IsKeyPressed(ebiten.KeyF1) {
		Keys.F1.Pressed = false
	} else {
		Keys.F1.Press = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyF11) && !Keys.F11.Pressed {
		Keys.F11.Press = true
	} else if !ebiten.IsKeyPressed(ebiten.KeyF11) {
		Keys.F11.Pressed = false
	} else {
		Keys.F11.Press = false
	}
}

// but it must use as TriggerList
