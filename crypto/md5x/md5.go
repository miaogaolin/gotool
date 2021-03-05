package md5x

import (
	"crypto/md5"
	"encoding/hex"
)

// 生成md5
func New(str string) string {
	h := md5.New()
	_, err := h.Write([]byte(str))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}
