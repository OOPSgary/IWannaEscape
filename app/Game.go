package app

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

func (g *Game) DrawGaming() {
	switch Status {
	case 1:
		g.game1()
	}
}

func (g *Game) game1() {

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
		10, 10, color.Black, g.screen,
		"康复训练:",
		"AD/<- ->左右移动",
		"空格键跳跃按R重开",
	)
	g.DrawAll(g.screen)
	if Dead {
		xg, yg := midImage(GameOver)
		geo := makeGeo(float64(g.screen.Bounds().Max.X/2)-float64(xg)*0.6, float64(g.screen.Bounds().Max.Y/2)-float64(yg)*0.6, 0.6, 0.6, 0, nil)
		geo.ColorScale.SetA(startA)
		g.screen.DrawImage(GameOver, geo)
		CallOnce(1, func() {
			for _, s := range strikeList {
				if !isChanClosed(s.KillSignal) {
					s.KillSignal <- 1
				}
			}
			cleanStrikes(strikeList)
			deadSoundPlayer.Rewind()
			deadSoundPlayer.Play()
		})
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

func game2() {}
