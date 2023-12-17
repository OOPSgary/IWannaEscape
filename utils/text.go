package utils

import (
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Text并不适用于格式化代码(使用自己写的DrawString进行渲染)
// 此为渲染一句话，而不是渲染多句
type Speech []*Text
type Text struct {
	Size          int
	Text          string
	StrikeThrough bool
	Underline     bool

	//居中时默认换行
	Middle bool
	//下一行换行
	Endl  bool
	Color color.Color
	//Be set After WriteString
	width  int
	height int
}

const (
	Normal = iota
	Small
	Big
)

var basicFont *opentype.Font
var DefaultTextSize = 5
var SizeNotFit = fmt.Errorf("the size of the image is smaller than required")

var OnlyOne = new(Once)

func WriteString(t Speech) (i *ebiten.Image, err error) {

	getWH := func(str string, f font.Face) (x, y int) {
		b, _ := font.BoundString(f, str)
		return (b.Max.X - b.Min.X).Ceil(), (b.Max.Y - b.Min.Y).Ceil()
	}
	getFont := func(Size int) (f font.Face) {
		// switch Size {
		// case Small:
		// 	return mplusSmallFont
		// case Big:
		// 	return mplusBigFont
		// default:
		// 	return mplusNormalFont
		// }
		return storedFonts[Size]
	}

	//排版
	var lineHeight int
	var lineWidth int
	var textHeight int
	var textWidth int
	var lineHeightList []int = make([]int, 0, 6)
	for _, s := range t {
		s.width, s.height = getWH(s.Text, getFont(s.Size))
		if s.Underline {
			s.height += 2
		}
		if s.height > lineHeight {
			lineHeight = s.height
		}
		s.width += 5
		lineWidth += s.width

		if s.Endl || s.Middle {
			textHeight += lineHeight + 10
			if textWidth < lineWidth {
				textWidth = lineWidth
			}
			lineHeightList = append(lineHeightList, lineHeight)
			lineHeight = 0
			lineWidth = 0

		}
		if s.Color == nil {
			s.Color = color.White
		}
	}

	lineHeightList = append(lineHeightList, lineHeight)
	textHeight += lineHeight + 10
	if textWidth < lineWidth {
		textWidth = lineWidth
	}
	var pointX int
	var pointY int
	var pointLine int
	i = ebiten.NewImage(textWidth, textHeight)
	for _, s := range t {
		OnlyOne.Do(s.Text, func() {
			println(pointY+lineHeightList[pointLine]-s.height, "LineHeight", lineHeightList[pointLine])
		})
		// if dst.Bounds().Dx() < lineWidth || dst.Bounds().Dy() < lineHeight {
		// 	return 0,0,SizeNotFit
		// }
		f := getFont(s.Size)
		opt := &ebiten.DrawImageOptions{}
		if s.Middle {
			pointX = (textWidth - s.width) / 2
		}
		opt.GeoM.Translate(float64(pointX), float64(pointY+lineHeightList[pointLine]))
		opt.ColorScale.ScaleWithColor(s.Color)
		text.DrawWithOptions(i, s.Text, f, opt)
		if s.StrikeThrough {
			vector.StrokeLine(i, float32(pointX), float32(pointY+lineHeightList[pointLine]-s.height/2+6),
				float32(pointX+s.width), float32(pointY+lineHeightList[pointLine]-s.height/2+6), 4, s.Color, true)
		}
		if s.Underline {
			vector.StrokeLine(i, float32(pointX), float32(pointY+lineHeightList[pointLine]+2),
				float32(pointX+s.width), float32(pointY+lineHeightList[pointLine]+2), 4, s.Color, true)
		}
		pointX += s.width
		if s.Endl || s.Middle {
			pointX = 0
			pointY += lineHeightList[pointLine]
			pointLine++
		}
	}
	return i, nil
}

type fontdata struct {
	lock  *sync.RWMutex
	fonts map[int]font.Face
}

var fontData = fontdata{
	lock:  new(sync.RWMutex),
	fonts: make(map[int]font.Face),
}

var storedFonts []font.Face = make([]font.Face, 3)

func init() {
	yahei, err := os.ReadFile("C:\\Windows\\Fonts\\msyh.ttc")
	t, err := opentype.ParseCollection(yahei)
	if err != nil {
		log.Fatal(err)
	}
	tt, err := t.Font(0)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	storedFonts[Small], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	storedFonts[Normal], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	storedFonts[Big], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// var (
// 	mplusSmallFont  font.Face
// 	mplusNormalFont font.Face
// 	mplusBigFont    font.Face
// )

/*
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
*/
