package imagex

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMime(t *testing.T) {
	file, err := os.Open("./image_test.gif")
	if err != nil {
		assert.Fail(t, "图片地址不存在")
	}
	mime, err := Mime(file)
	assert.Equal(t, "image/gif", mime)

	file, err = os.Open("./image_test.jpg")
	if err != nil {
		assert.Fail(t, "图片地址不存在")
	}
	mime, err = Mime(file)
	assert.Equal(t, "image/jpeg", mime)
}

func TestTransferJPEG(t *testing.T) {
	file, err := os.Open("./image_test.gif")
	if err != nil {
		assert.Fail(t, "图片地址不存在")
	}

	file, err = TransferJPEG(file, 200, 0)
	if err != nil {
		assert.Fail(t, "图片转化失败")
	}
}
