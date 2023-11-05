package utils

import (
	"fmt"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Text并不适用于格式化代码(使用自己写的DrawString进行渲染)
// 此为渲染一句话，而不是渲染多句
type Text []*struct {
	Size          int
	Text          string
	StrikeThrough bool
	Underline     bool
	//斜体
	Italic bool
	//下一行换行
	Endl bool

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

func WriteString(x, y float64, t Text) (i *ebiten.Image, err error) {

	getWH := func(str string, f font.Face) (x, y int) {
		b, _ := font.BoundString(f, str)
		return (b.Max.X - b.Min.X).Ceil(), (b.Max.Y - b.Min.Y).Ceil()
	}
	getFont := func(Size int) (f font.Face) {
		switch Size {
		case Small:
			return mplusSmallFont
		case Big:
			return mplusBigFont
		default:
			return mplusNormalFont
		}
	}

	//排版
	var lineHeight int
	var lineWidth int
	var textHeight int
	var textWidth int
	var lineHeightList []int = make([]int, 0, 6)
	for _, s := range t {
		if s.Italic {
			log.Fatal("not supported")
		}
		s.width, s.height = getWH(s.Text, getFont(s.Size))
		if s.height > lineHeight {
			lineHeight = s.height
		}
		lineWidth += s.width
		if s.Endl {
			textHeight += lineHeight
			if textWidth < lineWidth {
				textWidth = lineWidth
			}
			lineHeightList = append(lineHeightList, lineHeight)
			lineHeight = 0
			lineWidth = 0

		}
	}
	lineHeightList = append(lineHeightList, lineHeight)
	textHeight += lineHeight
	if textWidth < lineWidth {
		textWidth = lineWidth
	}
	var pointX int
	var pointY int
	var pointLine int
	i = ebiten.NewImage(textWidth, textHeight)
	for _, s := range t {

		// if dst.Bounds().Dx() < lineWidth || dst.Bounds().Dy() < lineHeight {
		// 	return 0,0,SizeNotFit
		// }
		f := getFont(s.Size)

		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(float64(pointX), float64(pointY+lineHeightList[pointLine]-s.height))
		text.DrawWithOptions(i, s.Text, f, opt)
		pointX += s.width
		if s.Endl {
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

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusSmallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}

var (
	mplusSmallFont  font.Face
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

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
