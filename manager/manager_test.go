package manager

import (
	"sync"
	"testing"
)

func TestFiles(t *testing.T) {
	//设置manager工作目录
	Init("./testassets")
	// 创建一个新的 files 实例
	f := &files{
		store: make(map[string]*info),
		m:     new(sync.RWMutex),
	}

	// 在测试中使用 Init 方法，模拟加载资源
	// 你可以根据需要传入一个 embed.FS 参数

	// 在测试中添加一些文件标志
	f.Sign("example.txt", "textfile", TextFile)
	f.Sign("example.txt", "textByte", DataFile)
	f.Sign("example.png", "imagefile", ImageFile)

	// 测试 PreLoadFile 方法
	if err := f.PreLoadFile(); err != nil {
		t.Errorf("PreLoadFile failed: %v", err)
	}

	// 测试 GetText 方法
	text, err := f.GetText("textfile")
	if err != nil {
		t.Errorf("GetText failed: %v", err)
	} else {
		// 检查返回的文本是否符合预期
		expectedText := "This is an example text file."
		if text != expectedText {
			t.Errorf("GetText returned unexpected text: got %q, expected %q", text, expectedText)
		}
	}

	// 测试 GetImage 方法
	image, _, err := f.GetImage("imagefile")
	if err != nil {
		t.Errorf("GetImage failed: %v", err)
	} else {
		// 检查返回的图像是否非空
		if image == nil {
			t.Error("GetImage returned a nil image")
		}
	}

	// 测试 GetBytes 方法
	bytes, err := f.GetBytes("textByte")
	if err != nil {
		t.Errorf("GetBytes failed: %v", err)
	} else {
		// 检查返回的字节是否非空
		if len(bytes) == 0 {
			t.Error("GetBytes returned empty bytes")
		}
	}

	// 测试一个不存在的标志
	_, err = f.GetText("nonexistentfile")
	if err != SignNotExist {
		t.Errorf("GetText for nonexistent file did not return SignNotExist error")
	}
}
