package app

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"image/color"
	"log"
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/solarlune/resolv"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Character struct {
	Status   int // 1,2,3 for three pictures 4 for dead but not completed
	Obj      *resolv.Object
	Top      *resolv.Object
	Button   *resolv.Object
	OnGround bool
	SpeedX   float64
	SpeedY   float64
	Jump     Jump
	FaceAt   string //It can be l(eft) or r(ight)

}
type Jump struct {
	Jump   int
	Chance int
	Lock   sync.Mutex
}

//go:embed resource/*
var emFs embed.FS
var basicFont *opentype.Font
var startA float32

const sampleRate = 44100

func init() {
	go func() {
		for {
			if startA <= 0 {
				startA = 1
			}
			startA -= 0.01
			time.Sleep(time.Millisecond * 10)
		}
	}()
	fontData, err := emFs.ReadFile("resource/msyh.ttc")
	if err != nil {
		log.Fatal(err)
	}
	// 解析字体文件
	t, err := opentype.ParseCollection(fontData)
	if err != nil {
		log.Fatal(err)
	}
	basicFont, err = t.Font(0)
	// basicFont, err = opentype.Parse(fontData)

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
	{
		rf, err := emFs.ReadFile("resource/Portal.png")
		if err != nil {
			log.Fatal(err)
		}
		portalImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(rf))
		if err != nil {
			log.Fatal(err)
		}
	}
	{
		rf, err := emFs.ReadFile("resource/save.png")
		if err != nil {
			log.Fatal(err)
		}
		saveImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(rf))
		if err != nil {
			log.Fatal(err)
		}
	}
	Blood = ebiten.NewImage(1, 1)
	Blood.Fill(color.RGBA{255, 0, 0, 1})

	diedSound, err := emFs.ReadFile("resource/Duang.wav")
	if err != nil {
		log.Fatal(err)
	}
	deadSoundPlayer = audioContext.NewPlayerFromBytes(diedSound)

	strikeSound, err = emFs.ReadFile("resource/StrikeSound.wav")
	if err != nil {
		log.Fatal(err)
	}

}
func newSoundPlayer(b []byte) *audio.Player {
	return audioContext.NewPlayerFromBytes(b)
}

// Picture Image Put here
var (
	Stopper          *ebiten.Image
	Kid              map[int]*ebiten.Image
	GameOver         *ebiten.Image
	NormalBackground *ebiten.Image
	StrikePhoto      *ebiten.Image
	Blood            *ebiten.Image
	portalImage      *ebiten.Image
	saveImage        *ebiten.Image
	audioContext     *audio.Context = audio.NewContext(sampleRate)
	deadSoundPlayer  *audio.Player
	strikeSound      []byte
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

func (g *Game) box() {
	g.putBlocksLine(movement{0, 0}, 0.5, 1, 40)
	g.putBlocksLine(movement{0, 0}, 0.5, 2, 30)
	g.putBlocksLine(movement{0, 480 - 16}, 0.5, 1, 40)
	g.putBlocksLine(movement{640 - 16, 0}, 0.5, 2, 30)
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
	}(mid), Y+float64(y)*Ys, Xs/2, Ys/2, turn, c)
	f := GetFont(size * 2)
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

var errSendOnClosedChannel = errors.New("send on closed chan")

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
func (g *Game) resetMap() {
	for _, o := range g.mainWorld.Blocks {
		World.Remove(o.Obj)
	}
	g.mainWorld.Blocks = make([]*block, 0)
	for _, s := range strikeList {
		if !isChanClosed(s.KillSingal) {
			s.KillSingal <- 1
		}
	}
	strikeList = make(map[any]*strike)
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
	face, err := opentype.NewFace(basicFont, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     300,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

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

	obj := resolv.NewObject(0, 0, 17, 21, "character")
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
	g.mainWorld.MainCharacter.Obj.Y = y
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
	if collision := g.mainWorld.MainCharacter.Obj.Check(x, y, "deadly"); collision != nil {
		if contactSet := g.mainWorld.MainCharacter.Obj.Shape.Intersection(x, y, collision.Objects[0].Shape); contactSet != nil {
			Dead = true
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
		// if g.mainWorld.MainCharacter.SpeedY <= 4 && g.mainWorld.MainCharacter.SpeedY >= 0 {
		// 	g.mainWorld.MainCharacter.SpeedY += 0.15

		// } else if g.mainWorld.MainCharacter.SpeedY <= 0 {
		// 	g.mainWorld.MainCharacter.SpeedY += 0.2
		// }

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
	for _, sl := range strikeList {
		if sl.Online {
			if contactSet := g.mainWorld.MainCharacter.Obj.Shape.Intersection(0, 0, sl.Shape); contactSet != nil {
				Dead = true
			}
		}
	}
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

var strikeList = make(map[any]*strike)
var TriggerList = make(map[any]*strikeTrigger)
var rePressR bool

func (g *Game) DrawAll(screen *ebiten.Image) {
	renderStrikes(strikeList, screen)
	renderTrigger(TriggerList, screen)
	g.drawBox(screen)
	if !Dead {
		f1 := func(s string) float64 {
			if g.mainWorld.MainCharacter.FaceAt == "l" {
				return g.mainWorld.MainCharacter.Obj.X + 21
			}
			return g.mainWorld.MainCharacter.Obj.X - 4
		}
		f2 := func(s string) float64 {
			if g.mainWorld.MainCharacter.FaceAt == "l" {
				return -1
			}
			return 1
		}
		fa := g.mainWorld.MainCharacter.FaceAt
		screen.DrawImage(Kid[g.mainWorld.MainCharacter.Status],
			makeGeo(f1(fa), g.mainWorld.MainCharacter.Obj.Y, f2(fa), 1, 0, nil))
	}
}
