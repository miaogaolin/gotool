package filex

import (
	"fmt"
	"github.com/miaogaolin/gotool/randx"
	"os"
	"time"
)

//调用os.MkdirAll递归创建文件夹
func CreateDir(dir string) error {
	if !isExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func CreateFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
}

// RandomName 随机文件名称
func RandomName() string {
	return fmt.Sprintf("%d%s", time.Now().Unix(), randx.RandomString(10))
}
