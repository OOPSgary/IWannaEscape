package app

import (
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

type bullet struct {
	X, Y   float64
	SpeedX float64
	Render bool
}
type bullets struct {
	b []bullet
	m *sync.Mutex
}

var blist *bullets = &bullets{
	m: &sync.Mutex{},
}

func (b *bullets) removeALL() {
	b.m.Lock()
	defer b.m.Unlock()
	clear(b.b)
}
func (b *bullets) addBullet(X, Y float64, left bool) {
	b.b = append(b.b, bullet{X: X, Y: Y, SpeedX: func(l bool) float64 {
		if l {
			return -3
		}
		return 3
	}(left), Render: true})
}
func (b *bullets) bulletRender(screen *ebiten.Image) {
	for _, bb := range b.b {
		bulletImage.DrawImage(screen, makeGeo(bb.X, bb.Y, 1, 1, 0, nil))
	}
}
func (b *bullets) updateBullets() {
	b.m.Lock()
	defer b.m.Unlock()
	for _, a := range b.b {
		a.X += a.SpeedX
	}

}

var bulletImage *ebiten.Image = ebiten.NewImage(1, 1)

func init() {
	vector.StrokeLine(bulletImage, 0, 0, 1, 1, 1, color.RGBA{0, 0, 0, 66}, true)
}

func (b *bullets) bulletCheckObjects(obj *resolv.Object) int {
	var newb []bullet
	var hit int
	for _, a := range b.b {
		if obj.X <= a.X && obj.X+obj.W >= a.X && obj.Y <= a.Y && obj.Y+obj.H >= a.Y {
			hit++
			continue
		}
		newb = append(newb, a)
	}
	b.m.Lock()
	defer b.m.Unlock()
	b.b = newb
	return hit
}
