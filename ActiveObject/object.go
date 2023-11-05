package ActiveObject

import (
	"fmt"

	"github.com/solarlune/resolv"
)

type Object struct {
	Object *resolv.Object
	// SignedImage []string
	WallSliding *resolv.Object
	OnGround    *resolv.Object
	SpeedX      float64
	SpeedY      float64
	Script      *ScriptRunner
}

var (
	NotJoint = fmt.Errorf("The object has NOT been registed")
)
var (
	friction = 0.5
	accel    = 0.5 + friction
	maxSpeed = 4.0
	jumpSpd  = 10.0
	gravity  = 0.75
)

func Register(object *resolv.Object, w *resolv.Space) *Object {
	w.Add(object)
	return &Object{
		Object: object,
	}
}
func (o *Object) Position(X, Y float64) error {
	o.Object.X += X
	o.Object.Y += Y
	o.Object.Update()
	return nil
}
func (o *Object) SetPosition(X, Y float64) error {
	o.Object.X, o.Object.Y = X, Y
	o.Object.Update()
	return nil
}
func (o *Object) Scale(xt, yt float64) error
func (o *Object) ReScale(xt, yt float64) error
func (o *Object) Size(x, y float64) error
func (o *Object) ReSize(x, y float64) error
func (o *Object) X() (X float64)
func (o *Object) Y() (Y float64)
func (o *Object) H() (H float64)
func (o *Object) W() (W float64)
func (o *Object) Sx() (Sx float64)
func (o *Object) Sy() (Sy float64)
func (o *Object) AddSpeedX() {

}

func (o *Object) AddSpeedY() {

}

// inspire from solarlune's `resolv` examples
func (o *Object) Resolv() error {
	if o.Object == nil {

	}
	if o.Object.Space == nil {

	}
	o.gravity(gravity)
	o.speedLimit()
	// dx是我们预计水平移动物体的距离
	dx := o.SpeedX
	if check := o.Object.Check(o.SpeedX, 0, "solid"); check != nil {
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
	return nil
}
func (o *Object) speedLimit() {
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
func (o *Object) gravity(float64) {
	o.SpeedY += gravity
	if o.WallSliding != nil && o.SpeedY > 1 {
		o.SpeedY = 1
	}
}
