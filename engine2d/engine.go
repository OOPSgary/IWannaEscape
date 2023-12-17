package engine2d

import (
	"IJustWantToEscape/method"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type drawLevel []*struct {
	level   int
	Objects []method.Object
}

// 元素个数
func (t drawLevel) Len() int {
	return len(t)
}

// 比较结果
func (t drawLevel) Less(i, j int) bool {
	return t[i].level < t[j].level
}

// 交换方式
func (t drawLevel) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type Game struct {
	Layers drawLevel
}

const (
	ScreenHeight = 960
	ScreenWidth  = 1280
)

func (g *Game) Update() error {
	return nil
}
func (g *Game) Draw(dst *ebiten.Image) {
	sort.Sort(g.Layers)
	for _, layer := range g.Layers {
		for _, o := range layer.Objects {
			o.Draw(dst)
		}
	}
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
