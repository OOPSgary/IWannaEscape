package script

import "IJustWantToEscape/utils"

type Runner interface {
	Update() error
}

type Script []struct {
	Type string
	Act  interface{}
	//Tick定义运行此动作所执行时间(按照60tick/s计算)
	Tick     uint
	WithWait bool
}

type Movement struct {
	// X,Y最佳为Tick的倍数
	X float64
	Y float64
	//MoveTo是将XY设置为目标坐标，而不是更改的坐标
	MoveTo bool
	//AutoCorrect 自动校正移动后的坐标 较为适合可能会导致偏差的坐标
	//方法简单粗暴，记录了原始坐标，在最后一个tick设置其为原始坐标+移动坐标/目标坐标
	//当X,Y为tick的倍数时，不推荐设置为True
	AutoCorrect bool
}

//Text直接使用util的text 方便渲染器使用
type Text utils.Text
