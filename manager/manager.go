package manager

import (
	"embed"
	"fmt"
	"image"
	"io"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	FormatNotExist = fmt.Errorf("format is not support in this platform")
	SignNotExist   = fmt.Errorf("sign is not exist")
	NotLoad        = fmt.Errorf("file not load because format/file not found")
	audioContext   = audio.NewContext(sampleRate)
	sampleRate     = 44400
)

type files struct {
	store map[string]*info
	// data  map[string]*[]byte
	m *sync.RWMutex
}
type info struct {
	form          int
	filename      string
	fileInterface interface{}
}
type audiofile struct {
	data []byte
}
type imagefile struct {
	image       *ebiten.Image
	actualImage *image.Image
}

func (a audiofile) Player() *audio.Player { return audioContext.NewPlayerFromBytes(a.data) }

const (
	ImageFile = iota
	AudioFile
	TextFile
	DataFile
)

func (f *files) Sign(FileName, SignName string, Format int) {
	f.m.Lock()
	defer f.m.Unlock()
	f.store[SignName] = &info{
		form:     Format,
		filename: FileName,
	}
}
func (f *files) SignBatchFiles(list []struct {
	FileName string
	SignName string
	Format   int
}) {
	f.m.Lock()
	defer f.m.Unlock()
	for _, s := range list {
		f.store[s.SignName] = &info{
			form:     s.Format,
			filename: s.FileName,
		}
	}
}
func (f *files) PreLoadFile() error {
	f.m.Lock()
	defer f.m.Unlock()
	for _, s := range f.store {
		if s.fileInterface != nil {
			continue
		}
		inter, err := loadfile(s.filename, s.form)
		if err != nil {
			return err
		}
		s.fileInterface = inter
	}
	return nil
}
func loadfile(filename string, format int) (interface{}, error) {
	b, err := assetsFile(filename)
	if err != nil {
		return nil, err
	}
	inter, err := formatfile(b, format)
	b.Close()
	return inter, err
}

func (f *files) GetBytes(SignName string) ([]byte, error) {
	f.m.RLocker()
	if f.store[SignName] == nil {
		f.m.RUnlock()
		return nil, SignNotExist
	}
	if DataFile != f.store[SignName].form {
		f.m.RUnlock()
		return nil, NotLoad
	}
	b, ok := f.store[SignName].fileInterface.(*[]byte)
	if !ok {
		inter, err := loadfile(f.store[SignName].filename, f.store[SignName].form)
		if err != nil {
			f.m.RUnlock()
			return nil, err
		}
		if b, ok := inter.(*[]byte); !ok {
			f.m.RUnlock()
			return nil, NotLoad
		} else {
			f.m.RUnlock()
			f.m.Lock()
			defer f.m.Unlock()
			f.store[SignName].fileInterface = inter
			return *b, nil
		}

	}
	f.m.RUnlock()
	return *b, nil
}

func (f *files) GetText(SignName string) (string, error) {
	f.m.RLocker()
	if f.store[SignName] == nil {
		f.m.RUnlock()
		return "", SignNotExist
	}
	if TextFile != f.store[SignName].form {
		f.m.RUnlock()
		return "", NotLoad
	}
	if s, ok := f.store[SignName].fileInterface.(*string); !ok {
		inter, err := loadfile(f.store[SignName].filename, f.store[SignName].form)
		if err != nil {
			f.m.RUnlock()
			return "", err
		}
		if s, ok := inter.(*string); !ok {
			f.m.RUnlock()
			return "", NotLoad
		} else {
			f.m.RUnlock()
			f.m.Lock()
			defer f.m.Unlock()
			f.store[SignName].fileInterface = inter
			return *s, nil
		}

	} else {
		f.m.RUnlock()
		return *s, nil
	}
}

func (f *files) GetImage(SignName string) (*ebiten.Image, *image.Image, error) {
	f.m.RLocker()

	if f.store[SignName] == nil {
		f.m.RUnlock()
		return nil, nil, SignNotExist
	}
	if ImageFile != f.store[SignName].form {
		f.m.RUnlock()
		return nil, nil, NotLoad
	}
	if i, ok := f.store[SignName].fileInterface.(*imagefile); !ok {
		inter, err := loadfile(f.store[SignName].filename, f.store[SignName].form)
		if err != nil {
			f.m.RUnlock()
			return nil, nil, err
		}
		if image, ok := inter.(*imagefile); !ok {
			f.m.RUnlock()
			return nil, nil, NotLoad
		} else {
			f.m.RUnlock()
			f.m.Lock()
			defer f.m.Unlock()
			f.store[SignName].fileInterface = inter
			return image.image, image.actualImage, nil
		}

	} else {
		f.m.RUnlock()
		return i.image, i.actualImage, nil
	}

}

func (f *files) FileSign(FileName, SignName string) (interface{}, error)
func Init(e *embed.FS) {

}
func assetsFile(filename string) (b *os.File, err error) {
	return os.Open(fmt.Sprintf("assets\\%v", filename))
}

func formatfile(data io.Reader, format int) (inter interface{}, err error) {
	switch format {
	case ImageFile:
		image, file, err := ebitenutil.NewImageFromReader(data)
		return &imagefile{
			image:       image,
			actualImage: &file,
		}, err
	case AudioFile:
		b, err := io.ReadAll(data)
		return &audiofile{data: b}, err
	case TextFile:
		b, err := io.ReadAll(data)
		s := string(b)
		return &s, err
	case DataFile:
		b, err := io.ReadAll(data)
		return &b, err
	default:
		return nil, FormatNotExist
	}
}
