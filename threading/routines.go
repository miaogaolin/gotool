package threading

import (
	"fmt"
	"runtime/debug"
	"github.com/miaogaolin/gotool/logx"
)

func GoSafe(fn func()) {
	go runSafe(fn)
}

func runSafe(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			logx.Label("recover").Error(fmt.Sprintf("%v, %s", r, string(debug.Stack())))
		}
	}()
	fn()
}
