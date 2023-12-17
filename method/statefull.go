package method

import (
	"IJustWantToEscape/manager"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type StatefullObject struct {
	Object *resolv.Object
	// SignedImage []string
	WallSliding                    *resolv.Object
	OnGround                       *resolv.Object
	ScaleX, ScaleY, SpeedX, SpeedY float64

	OutOfCheck  bool
	Script      *ScriptRunner
	FacingRight bool
	TargetImage string
}

var (
	NotJoint               = fmt.Errorf("The object has NOT been registed")
	StatefullWithOutObject = fmt.Errorf("Why StateFull Object dont have an ResolvObject?")
)
var (
	friction = 0.5
	accel    = 0.5 + friction
	maxSpeed = 4.0
	jumpSpd  = 10.0
	gravity  = 0.75
)

func (o *StatefullObject) Run(world *resolv.Space) error {
	//if Scipt is here, we let the scipt to Init the Object
	if o.Script != nil {
		return nil
	}
	if o.Object == nil {
		return StatefullWithOutObject
	}
	world.Add(o.Object)
	return nil
}

func (o *StatefullObject) Position() (X, Y float64) {
	return o.Object.X, o.Object.Y
}

// func (o *StatefullObject) SetPosition(X, Y float64) error {
// 	o.Object.X, o.Object.Y = X, Y
// 	o.Object.Update()
// 	return nil
// }

// inspire from solarlune's `resolv` examples
func (o *StatefullObject) Resolv() error {
	if o.Object == nil || o.OutOfCheck {
		return nil
	}
	if o.Object.Space == nil {
		return NotJoint
	}

	o.gravity(gravity)
	o.speedLimit()
	// dx是我们预计水平移动物体的距离
	dx := o.SpeedX
	if check := o.Object.Check(o.SpeedX, 0, "stateless"); check != nil {
		// 水平移动相对简单；我们只需检查前方是否有固态物体。如果有，则移动到接触并停止水平移动速度。如果没有，则可以向前移动。

		dx = check.ContactWithCell(check.Cells[0]).X()
		o.SpeedX = 0

		// 滑墙
		if o.OnGround == nil {
			o.WallSliding = check.Objects[0]
		}
	}
	//将dx应用到对象上
	o.Object.X += dx

	//对于水平移动 ,我们将有不同的处理
	// 首先，我们将 OnGround 设置为 nil，以防我们最终没有站在任何东西上。
	o.OnGround = nil

	//设置我们的纵向移动速度 并限制到最大水平
	dy := math.Max(math.Min(o.SpeedY, 16), -16)
	// 我们将使用 dy（垂直移动速度）检查碰撞，但在向下移动时添加一个单位以更深入地查找地面上的固体物体，具体来说。
	checkDistance := dy
	if dy >= 0 {
		checkDistance++
	}

	if check := o.Object.Check(0, checkDistance); check != nil {
		slide := check.SlideAgainstCell(check.Cells[0], "stateless")
		if dy < 0 && check.Cells[0].ContainsTags("stateless") && slide != nil && math.Abs(slide.X()) <= 8 {
			// 如果我们可以在这里滑动，那就这样做。没有接触，并且垂直速度（dy）保持向上。
			o.Object.X += slide.X()
		} else {
			//检查斜坡
			if ramps := check.ObjectsByTags("ramp"); len(ramps) > 0 {
				for _, ramp := range ramps {
					if contactSet := o.Object.Shape.Intersection(dx, 8, ramp.Shape); dy >= 0 && contactSet != nil {
						o.SpeedY = 0
						if n := contactSet.TopmostPoint()[1] - o.Object.Bottom() + 0.1; n > dy {
							o.OnGround = ramp
							dy = n
						}
					}
				}
			}

			if solids := check.ObjectsByTags("stateless"); len(solids) > 0 && (o.OnGround == nil || o.OnGround.Y >= solids[0].Y) {
				dy = check.ContactWithObject(solids[0]).Y()
				o.SpeedY = 0
				// 只有当我们着陆时才算在地面上（如果物体的 Y 大于玩家的 Y）。
				if solids[0].Y > o.Object.Y {
					o.OnGround = solids[0]
				}
			}
			if o.OnGround != nil {
				o.WallSliding = nil // 玩家在地面上，所以不再进行壁滑。
			}
		}
	}
	o.Object.Y = dy

	wallNext := 1.0
	if !o.FacingRight {
		wallNext = -1
	}

	// 如果玩家旁边的墙消失，停止壁滑。
	if c := o.Object.Check(wallNext, 0, "solid"); o.WallSliding != nil && c == nil {
		o.WallSliding = nil
	}

	o.Object.Update() // 更新玩家在空间中的位置。
	return nil
}
func (o *StatefullObject) speedLimit() {
	if o.SpeedX > friction {
		o.SpeedX -= friction
	} else if o.SpeedX < -friction {
		o.SpeedX += friction
	} else {
		o.SpeedX = 0
	}

	if o.SpeedX > maxSpeed {
		o.SpeedX = maxSpeed
	} else if o.SpeedX < -maxSpeed {
		o.SpeedX = -maxSpeed
	}
}
func (o *StatefullObject) gravity(float64) {
	o.SpeedY += gravity
	if o.WallSliding != nil && o.SpeedY > 1 {
		o.SpeedY = 1
	}
}
func (o *StatefullObject) Update() error {
	return o.Resolv()
}
func (o *StatefullObject) Draw(dst *ebiten.Image) error {
	if o.TargetImage == "" {
		return nil
	}
	option := &ebiten.DrawImageOptions{}
	option.GeoM.Translate(o.Position())
	option.GeoM.Scale(o.ScaleX, o.ScaleY)
	i, _, err := manager.Manager.GetImage(o.TargetImage)
	if err != nil {
		return err
	}
	i.DrawImage(dst, option)
	return nil
}
