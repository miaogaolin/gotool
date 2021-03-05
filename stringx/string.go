package stringx

import (
	"fmt"
)

func Interface(data interface{}) string {
	return fmt.Sprintf("%v", data)
}
