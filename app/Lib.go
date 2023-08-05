package app

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"image/color"
	"log"
	"sync"
	"time"

	"github.com/goki/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/solarlune/resolv"
	"golang.org/x/image/font"
)

var trapTrigger1 = TrapTrigger{
	[]trapmovement{

		{
			Mode: 1,
		},
		{
			Mode: 4,
			Time: 12,
		},
		{
			Mode: 3,
			Time: 25,
			Movement: movementPlus{
				SizeY: 40,
				Y:     -500,
			},
		}, {
			Mode: 3,
			Time: 52,
			Movement: movementPlus{
				SizeY: 4,
				Y:     500,
			},
		},
		// {
		// 	Mode: 2,
		// 	Time: 40,
		// },
		// {
		// 	Mode: 1,
		// 	Time: 40,
		// },
		{
			Mode: 3,
			Time: 50,
			Movement: movementPlus{
				X: 640,
			},
		},
		{
			Mode: 3,
			Time: 200,
			Movement: movementPlus{
				X: -680,
			},
		}, {
			Mode: 3,
			Time: 25,
			Movement: movementPlus{
				SizeY: 40,
				Y:     -500,
			},
		}, {
			Mode: 3,
			Time: 52,
			Movement: movementPlus{
				SizeY: 4,
				Y:     500,
			},
		}, {
			Mode: 3,
			Time: 10,
			Movement: movementPlus{
				X: 680,
			},
		},
		{
			Mode: 3,
			Time: 50,
			Movement: movementPlus{
				X:     -580,
				Y:     -250 + 16,
				SizeX: 8,
				SizeY: 8,
			},
		},
		{
			Mode: 3,
			Time: 90,
			Movement: movementPlus{
				Y: -400,
			},
		}, {
			Mode: 3,
			Time: 37,
			Movement: movementPlus{
				Y: 480,
			},
		}, {
			Mode: 3,
			Time: 100,
			Movement: movementPlus{
				SizeX: 17,
			},
		}, {
			Mode: 2,
			Time: 200,
		},
	},
}

//go:embed resource/*
var emFs embed.FS

var basicFont *truetype.Font
var startA float32 = 1

func init() {
	go func() {
		if startA <= 0 {
			startA = 1
		}
		startA -= 0.01
	}()
	fontData, err := emFs.ReadFile("resource/DingTalk.ttf")
	if err != nil {
		log.Fatal(err)
	}

	// 解析字体文件
	basicFont, err = truetype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	st, err := emFs.ReadFile("resource/Block.png")
	if err != nil {
		log.Fatal(err)
	}
	Stopper, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(st))
	if err != nil {
		log.Fatal(err)
	}

	sp, err := emFs.ReadFile("resource/Strike.png")
	if err != nil {
		log.Fatal(err)
	}
	StrikePhoto, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(sp))
	if err != nil {
		log.Fatal(err)
	}

	//Load map
	Kid = make(map[int]*ebiten.Image)
	for i := 1; i <= 3; i++ {
		kid, err := emFs.ReadFile(fmt.Sprintf("resource/Step%d.png", i))
		if err != nil {
			log.Fatal(err)
		}
		Kid[i], _, err = ebitenutil.NewImageFromReader(bytes.NewReader(kid))
		if err != nil {
			log.Fatal(err)
		}
	}

	GO, err := emFs.ReadFile("resource/GameOver.png")
	if err != nil {
		log.Fatal(err)
	}
	GameOver, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(GO))
	if err != nil {
		log.Fatal(err)
	}

	nb, err := emFs.ReadFile("resource/Background.png")
	if err != nil {
		log.Fatal(err)
	}
	NormalBackground, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(nb))
	if err != nil {
		log.Fatal(err)
	}
	Blood = ebiten.NewImage(1, 1)
	Blood.Fill(color.RGBA{255, 0, 0, 1})

	diedSound, err := emFs.ReadFile("resource/Duang.wav")
	if err != nil {
		log.Fatal(err)
	}
	const sampleRate = 44100

	Sound := audio.NewContext(sampleRate)
	{
		reader := bytes.NewReader(diedSound)

		deadSoundPlayer, err = Sound.NewPlayer(reader)
		if err != nil {
			log.Fatal(err)
		}
	}

}

// Picture Image Put here
var (
	Stopper          *ebiten.Image
	Kid              map[int]*ebiten.Image
	GameOver         *ebiten.Image
	NormalBackground *ebiten.Image
	StrikePhoto      *ebiten.Image
	Blood            *ebiten.Image
	deadSoundPlayer  *audio.Player
)

type movement struct{ x, y float64 }

var Co = struct {
	CallOnceMap map[any]bool
	Mutex       *sync.Mutex
}{
	CallOnceMap: make(map[any]bool),
	Mutex:       &sync.Mutex{},
}

func CallOnce(ID any, f func()) {
	Co.Mutex.Lock()
	if !Co.CallOnceMap[ID] {
		f()
		log.Println("Called for ", ID)
	}
	Co.CallOnceMap[ID] = true
	Co.Mutex.Unlock()
}
func CallExtra(ID any) {
	Co.Mutex.Lock()
	Co.CallOnceMap[ID] = false
	Co.Mutex.Unlock()
}
func (j *Jump) Update() bool {
	j.Lock.Lock()
	defer j.Lock.Unlock()

	if j.Jump > 0 && j.Chance == 0 {
		j.Jump--
		j.Chance = 1
		return true
	}
	return false
}
func (j *Jump) Add(i int) {
	j.Lock.Lock()
	defer j.Lock.Unlock()
	j.Jump += i
}
func (j *Jump) Reset() {
	j.Lock.Lock()
	defer j.Lock.Unlock()
	j.Jump = 2
}
func (j *Jump) ResetChance() {
	j.Lock.Lock()
	defer j.Lock.Unlock()
	j.Chance = 0
}
func (j *Jump) AddChance() bool {
	j.Lock.Lock()
	defer j.Lock.Unlock()
	if j.Chance <= 30 {
		j.Chance++
		return true
	} else {
		return false
	}

}

type safeMap struct {
	m    map[any]any
	lock *sync.RWMutex
}

func CreateSafeMap(m map[any]any) *safeMap {
	return &safeMap{
		m:    m,
		lock: &sync.RWMutex{},
	}
}
func (m safeMap) Swap(id any, data any) any {
	m.lock.Lock()
	d := m.m[id]
	m.m[id] = data
	m.lock.Unlock()
	return d
}
func (m safeMap) Store(id any, data any) {
	m.lock.Lock()
	m.m[id] = data
	m.lock.Unlock()
}
func (m safeMap) Get(id any) any {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.m[id]
}
func (m safeMap) SwapAndDelete(id any) any {
	m.lock.Lock()
	defer func() {
		m.m[id] = nil
		m.lock.Unlock()
	}()
	return m.m[id]
}
func (m safeMap) DeteleAll() []any {
	m.lock.Lock()
	defer m.lock.Unlock()
	var dump []any = make([]any, len(m.m))
	for _, a := range m.m {
		dump = append(dump, a)
	}
	return dump
}
func (m safeMap) SwapAllDelete() (any, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for s, a := range m.m {
		delete(m.m, s)
		return a, true
	}
	return nil, false

}

func (g *Game) box() {
	g.putBlocksLine(movement{0, 0}, 0.5, 1, 40)
	g.putBlocksLine(movement{0, 0}, 0.5, 2, 30)
	g.putBlocksLine(movement{0, 480 - 16}, 0.5, 1, 40)
	g.putBlocksLine(movement{640 - 16, 0}, 0.5, 2, 30)
}

type strike struct {
	Pos          movement
	SizeX, SizeY float64
	Angle        float64
	Obj          *resolv.Object
	Trigger      chan TrapTrigger
	KillSingal   chan interface{}
	Online       bool //Not using this again

	// Close  chan interface{} //Not using this again
	// Closed bool             //Not using this again
	// Mutex  *sync.Mutex      //Not using this again
}

// If you need to destory this strike just Close the channel
// Remind that delayed changes is recently not support Size
/*
func (s *strike) OldLoad() {
	s.Closed = false
	strikeWaitGroup.Add(1)
	go func() {
		defer strikeWaitGroup.Done()
		for S := range s.Trigger {
			for _, a := range S.Movement {
				if s.Closed {
					s.Mutex.Lock()

					if !s.Closed && World != nil && s.Obj != nil {
						World.Remove(s.Obj)
					} else {
						s.Mutex.Unlock()
						return
					}
					s.Mutex.Unlock()
					return
				}
				switch a.Mode {
				case 1:

					time.Sleep(time.Millisecond * time.Duration(a.Time) * 10)
					s.Mutex.Lock()
					if !s.Closed && World != nil && s.Obj != nil {
						World.Add(s.Obj)

					} else {
						s.Mutex.Unlock()
						return
					}
					s.Online = true
					s.Mutex.Unlock()
				case 2:
					time.Sleep(time.Millisecond * time.Duration(a.Time) * 10)
					s.Mutex.Lock()
					if !s.Closed && World != nil && s.Obj != nil {
						World.Remove(s.Obj)
					} else {
						s.Mutex.Unlock()
						return
					}
					s.Online = false
					s.Mutex.Unlock()
				case 3:
					if a.Time == 0 {

						s.Angle = a.Movement.Angle
						if a.Movement.X > 0 {
							s.Pos.x = a.Movement.X
						}
						if a.Movement.Y > 0 {
							s.Pos.y = a.Movement.Y
						}
						s.Mutex.Lock()

						if !s.Closed && World != nil && s.Obj != nil {
							World.Remove(s.Obj)
						} else {
							s.Mutex.Unlock()
							return
						}
						s.Mutex.Unlock()

						if a.Movement.SizeX > 0 {
							s.SizeX = a.Movement.SizeX
						}
						if a.Movement.SizeY > 0 {
							s.SizeY = a.Movement.SizeY
						}
						if a.Movement.SizeX > 0 && a.Movement.SizeY > 0 || (a.Movement.SizeX > 0 || a.Movement.SizeY > 0) {
							s.Obj = resolv.NewObject(s.Pos.x, s.Pos.y, 32*s.SizeX, 32*s.SizeY, "deadly")
							s.Obj.SetShape(resolv.NewConvexPolygon(
								0, 0,

								s.Obj.W/2, 0,
								s.Obj.W/2+1, 0,
								0, s.Obj.H,
								s.Obj.W, s.Obj.H,
							))
							World.Add(s.Obj)
						}
						s.Obj.X = s.Pos.x
						s.Obj.Y = s.Pos.y
						s.Mutex.Lock()
						if s.Obj != nil {
							s.Obj.Update()
						}
						s.Mutex.Unlock()
					} else {
						sizeX := s.SizeX
						sizeY := s.SizeY
						for i := float64(1); i <= float64(a.Time); i++ {
							t := time.After(time.Millisecond * 10)
							s.Angle += a.Movement.Angle / float64(a.Time)
							if a.Movement.X != 0 {
								s.Pos.x += a.Movement.X / float64(a.Time)
							}
							if a.Movement.Y != 0 {
								s.Pos.y += a.Movement.Y / float64(a.Time)
							}

							if a.Movement.SizeX > 0 && a.Movement.SizeY > 0 || (a.Movement.SizeX > 0 || a.Movement.SizeY > 0) {
								newSizeX := func() float64 {
									if a.Movement.SizeX > 0 {
										return (a.Movement.SizeX-sizeX)/float64(a.Time)*i + sizeX
									}
									return sizeX
								}()
								newSizeY := func() float64 {
									if a.Movement.SizeY > 0 {
										return (a.Movement.SizeY-sizeY)/float64(a.Time)*i + sizeY
									}
									return sizeY
								}()
								s.Mutex.Lock()

								if !s.Closed && World != nil && s.Obj != nil {
									World.Remove(s.Obj)

								} else {
									s.Mutex.Unlock()
									return
								}
								s.Mutex.Unlock()

								s.Obj = resolv.NewObject(s.Pos.x, s.Pos.y, 32*newSizeX, 32*newSizeY, "deadly")
								s.Obj.SetShape(resolv.NewConvexPolygon(
									0, 0,

									s.Obj.W/2, 0,
									s.Obj.W/2+1, 0,
									0, s.Obj.H,
									s.Obj.W, s.Obj.H,
								))

								s.Mutex.Lock()
								if !s.Closed && World != nil && s.Obj != nil {
									World.Add(s.Obj)
								} else {
									s.Mutex.Unlock()
									return
								}
								s.Mutex.Unlock()

								s.SizeX = newSizeX
								s.SizeY = newSizeY
							} else {
								s.Obj.X = s.Pos.x
								s.Obj.Y = s.Pos.y

							}
							s.Mutex.Lock()
							if s.Obj != nil {
								s.Obj.Update()
							}
							s.Mutex.Unlock()
							<-t
						}
					}
				case 4:
					time.Sleep(time.Millisecond * time.Duration(a.Time) * 10)
				}
			}
		}
	}()
	go func() {
		<-s.Close
		if s.Obj != nil && World != nil {
			s.Mutex.Lock()
			World.Remove(s.Obj)
			s.Mutex.Unlock()
		}
		s.Online = false
		s.Closed = true
		close(s.Trigger)
		close(s.Close)
	}()
}
*/
func (s *strike) Render(screen *ebiten.Image) {
	if s.Online {
		geo := &ebiten.DrawImageOptions{}
		// {
		// 	s2 := StrikePhoto.Bounds().Size()
		// 	geo.GeoM.Translate(-float64(s2.X)/2, -float64(s2.Y)/2)
		// }

		geo.GeoM.Rotate(getRadian(s.Angle))
		geo.GeoM.Scale(s.SizeX, s.SizeY)
		geo.GeoM.Translate(s.Obj.X, s.Obj.Y)
		screen.DrawImage(StrikePhoto, geo)
	}

}
func NewStrike(Pos movement, SizeX, SizeY float64) *strike {
	obj := resolv.NewObject(Pos.x, Pos.y, 4*8*SizeX, 4*8*SizeY, "deadly")
	obj.Update()
	obj.SetShape(resolv.NewConvexPolygon(
		0, 0,

		obj.W/2, 0,
		obj.W/2+1, 0,
		0, obj.H,
		obj.W, obj.H,
	))
	// obj := resolv.NewObject(Pos.x, Pos.y, 64, 32, "deadly")
	// obj.SetShape(resolv.NewConvexPolygon(
	// 	0, 0,
	// 	0, 0,
	// 	0, 64,
	// 	64, 64,
	// 	64, 64,
	// ))
	obj.Update()

	return &strike{
		Pos:     Pos,
		SizeX:   SizeX,
		SizeY:   SizeY,
		Angle:   0,
		Obj:     obj,
		Trigger: make(chan TrapTrigger),
		// Close:      make(chan interface{}),
		// Closed:     false,
		KillSingal: make(chan interface{}),
		Online:     false,
		// Mutex:      &sync.Mutex{},
	}
}

// DrawString函数作为绘画文字
func DrawString(s string, size int, X, Y, Xs, Ys, turn float64, c color.Color, Image *ebiten.Image, mid bool) int {
	var change int
	var y int
	change, y = midText(s, size)
	op := makeGeo(func(middle bool) float64 {
		if middle {
			return X - float64(change)*Xs
		}
		return X
	}(mid), Y+float64(y)*Ys, Xs, Ys, turn, c)
	f := GetFont(size)
	text.DrawWithOptions(Image, s, f, op)
	return y
}
func midText(s string, size int) (x, y int) {
	face := GetFont(size)
	b, _ := font.BoundString(face, s)
	return (b.Max.X - b.Min.X).Ceil() / 2, (b.Max.Y - b.Min.Y).Ceil()

}
func midImage(Image *ebiten.Image) (x, y int) {
	return (Image.Bounds().Max.X - Image.Bounds().Min.X) / 2, (Image.Bounds().Max.Y - Image.Bounds().Min.Y) / 2
}
func drawMidTextLineByLine(startHeight, size int, c color.Color, Image *ebiten.Image, s ...string) (Height int) {
	for _, str := range s {
		y := DrawString(str, size, float64(Image.Bounds().Dx())/2, float64(startHeight), 1, 1, 0, c, Image, true)
		startHeight += y + 10
	}
	return startHeight
}

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
func (g *Game) putBlocksLine(StartPos movement, Size float64, Line int, amount int) {
	for i := 0; i < amount; i++ {
		if Line == 1 {
			g.newBlock(movement{StartPos.x + float64(i)*Size*32, StartPos.y}, Size, 0)
		} else {
			g.newBlock(movement{StartPos.x, StartPos.y + float64(i)*Size*32}, Size, 0)
		}
	}

}
func (b block) Draw(screen *ebiten.Image) {
	geo := &ebiten.DrawImageOptions{}
	geo.GeoM.Scale(b.Size, b.Size)
	geo.GeoM.Translate(b.Obj.X, b.Obj.Y)
	geo.GeoM.Rotate(getRadian(float64(b.Rorate)))
	screen.DrawImage(Stopper, geo)
}

var errSendOnClosedChannel = errors.New("send on closed chan")

func (s *strike) Send(t TrapTrigger) error {
	select {
	case <-s.KillSingal:
		return errSendOnClosedChannel
	default:
		s.Trigger <- t
	}
	return nil
}
func (s *strike) Load() {
	go s.load()
}
func (s *strike) load() {
	for {
		if s.process() {
			waitKeepProcessing.Add(1)
			close(s.Trigger)
			if !isChanClosed(s.KillSingal) {
				close(s.KillSingal)
			}
			if s.Online {
				World.Remove(s.Obj)
			}
			waitKeepProcessing.Done()
			break
		}
	}
}
func (s *strike) process() (stop bool) {
	var a TrapTrigger
	select {
	case a = <-s.Trigger:
	case <-s.KillSingal:
		return true
	}
	for _, action := range a.Movement {
		log.Println(action)
		switch action.Mode {
		case 1:
			if s.handlerAppear(action) {
				return true
			}
		case 2:
			if s.handlerDisAppear(action) {
				return true
			}
		case 3:
			if s.handlerMovement(action) {
				return true
			}
		case 4:
			if s.handlerWaiting(action) {
				return true
			}
		}
	}
	return false
}
func (s *strike) handlerAppear(action trapmovement) (stop bool) {
	select {
	case <-time.After(time.Millisecond * 10 * time.Duration(action.Time)):
		World.Add(s.Obj)
		s.Online = true
		return false
	case <-s.KillSingal:
		return true
	}
}
func (s *strike) handlerDisAppear(action trapmovement) (stop bool) {
	select {
	case <-time.After(time.Millisecond * 10 * time.Duration(action.Time)):
		World.Remove(s.Obj)
		s.Online = false
		return false
	case <-s.KillSingal:
		return true
	}
}
func (s *strike) handlerMovement(action trapmovement) (stop bool) {
	PreferData := struct {
		sizeX, sizeY float64
	}{
		sizeX: s.SizeX,
		sizeY: s.SizeY,
	}
	delayProcess := func() {
		s.Pos.x += action.Movement.X / float64(action.Time)
		s.Pos.y += action.Movement.Y / float64(action.Time)
		s.SizeX += (ifPositiveNum(s.SizeX, action.Movement.SizeX) - PreferData.sizeX) / float64(action.Time)
		s.SizeY += (ifPositiveNum(s.SizeY, action.Movement.SizeY) - PreferData.sizeY) / float64(action.Time)
		s.Angle += action.Movement.Angle / action.Time
		s.Obj.X = s.Pos.x
		s.Obj.Y = s.Pos.y
		s.Obj.W = 32 * s.SizeX
		s.Obj.H = 32 * s.SizeY
		s.Obj.SetShape(resolv.NewConvexPolygon(
			0, 0,

			s.Obj.W/2, 0,
			s.Obj.W/2+1, 0,
			0, s.Obj.H,
			s.Obj.W, s.Obj.H,
		))
		// s.Obj.Shape.SetScale(s.SizeX, s.SizeY)
		s.Obj.Shape.SetRotation(s.Angle)
		s.Obj.Update()
	}
	if action.Time <= 0 {
		s.Pos.x += action.Movement.X
		s.Pos.y += action.Movement.Y
		s.SizeX = ifPositiveNum(s.SizeX, action.Movement.SizeX)
		s.SizeY = ifPositiveNum(s.SizeY, action.Movement.SizeY)
		s.Angle += action.Movement.Angle
		s.Obj.X = s.Pos.x
		s.Obj.Y = s.Pos.y
		s.Obj.W = 32 * s.SizeX
		s.Obj.H = 32 * s.SizeY
		s.Obj.SetShape(resolv.NewConvexPolygon(
			0, 0,

			s.Obj.W/2, 0,
			s.Obj.W/2+1, 0,
			0, s.Obj.H,
			s.Obj.W, s.Obj.H,
		))
		// s.Obj.Shape.SetScale(s.SizeX, s.SizeY)
		s.Obj.Shape.SetRotation(s.Angle)
		s.Obj.Update()
	} else {
		for i := 1; i < int(action.Time); i++ {
			select {
			case <-time.After(time.Millisecond * 10):
				delayProcess()
			case <-s.KillSingal:
				return true
			}
		}
	}
	return false
}
func (s *strike) handlerWaiting(action trapmovement) (stop bool) {
	select {
	case <-s.KillSingal:
		return true
	case <-time.After(time.Millisecond * 10 * time.Duration(action.Time)):
		return false
	}
}
func ifPositiveNum(value, replaceValue float64) float64 {
	if replaceValue > 0 {
		return replaceValue
	}
	return value
}

type TrapTrigger struct {
	Movement []trapmovement
}
type trapmovement struct {
	Mode int
	// 1 for appear (time will be used to sleep before it run)
	// 2 for dispear (time will be used to sleep before it run)
	// 3 for movement+rorate(Angle)
	// 4 for just sleep
	Time     float64 // Counted in 10Mill Must be integer
	Movement movementPlus
}
type movementPlus struct {
	Angle        float64
	X, Y         float64
	SizeX, SizeY float64
}
