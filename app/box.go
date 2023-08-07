package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type block struct {
	Obj    *resolv.Object
	Size   float64
	Rorate int
}

func (g *Game) newBlock(pos movement, size float64, rorate int) *block {
	obj := resolv.NewObject(float64(pos.x), float64(pos.y), 32*size, 32*size, "Stopper")
	World.Add(obj)

	b := block{
		Obj:    obj,
		Size:   size,
		Rorate: rorate,
	}
	g.mainWorld.Blocks = append(g.mainWorld.Blocks, &b)
	return &b
}

// 1 stands for X-Line()--- 2 stands for Y-Line |

func (b block) Draw(screen *ebiten.Image) {
	geo := &ebiten.DrawImageOptions{}
	geo.GeoM.Scale(b.Size, b.Size)
	geo.GeoM.Translate(b.Obj.X, b.Obj.Y)
	geo.GeoM.Rotate(getRadian(float64(b.Rorate)))
	screen.DrawImage(Stopper, geo)
}
func (g *Game) putBlocksLine(StartPos movement, Size float64, Line int, amount int) {
	for i := 0; i < amount; i++ {
		if Line == 1 {
			g.newBlock(movement{StartPos.x + float64(i)*Size*32, StartPos.y}, Size, 0)
		} else {
			g.newBlock(movement{StartPos.x, StartPos.y + float64(i)*Size*32}, Size, 0)
		}
	}

}
func (g *Game) drawBox(screen *ebiten.Image) {
	for _, a := range g.mainWorld.Blocks {
		if a != nil {
			a.Draw(screen)
		}
	}
}
