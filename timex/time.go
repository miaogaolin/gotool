package timex

import "time"

func DayString() string {
	return time.Now().Format("2006-01-02")
}
