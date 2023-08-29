package app

import (
	"fmt"
	"image/color"
	"strconv"
	"sync"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
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

var World = resolv.NewSpace(640, 480, 16, 16)

func (g *Game) Update() error {
	switch Status {
	case 0:
	default:
		if g.mainWorld.MainCharacter.Obj != nil {
			g.syncCharacterMovement()
			g.syncCharacter()
		}
		bulletList.updateBullets()
		g.checkCurrentPortals()
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
				if deadSoundPlayer.Rewind() != nil {
					return
				}
				deadSoundPlayer.Play()
			})
		} else {

			drawMidTextLineByLine(0, 10, color.Black, screen, "Welcome to", "I JUST WANT TO ESCAPE", "F1: VSync模式切换 F11: 全屏")
			g.HomePage.Update()
			g.HomePage.Draw(screen)
		}

	default:
		g.DrawGaming()
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

func (g *Game) Layout(int, int) (screenWidth int, screenHeight int) {
	return 640, 480
}

var Status = 0
var Dead = false
var waitKeepProcessing = new(sync.WaitGroup)
var WaitTime time.Duration
var KeyPressed = make(map[ebiten.Key]bool)

func newFill(w, h int, r, g, b, a uint8) (i *ebiten.Image) {
	i = ebiten.NewImage(w, h)
	i.Fill(color.RGBA{R: r, G: g, B: b, A: a})
	return
}
