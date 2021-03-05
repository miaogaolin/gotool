package threading

import (
	"fmt"
	"runtime/debug"
	"syt-crawler/core/log"
)

func GoSafe(fn func()) {
	go runSafe(fn)
}

func runSafe(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Label("recover").Error(fmt.Sprintf("%v, %s", r, string(debug.Stack())))
		}
	}()
	fn()
}
