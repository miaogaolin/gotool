package timex

import "time"

var (
	dayLayout     = "2006-01-02"
	secondsLayout = "2006-01-02 15:04:05"
)

func DayString() string {
	return time.Now().Format(dayLayout)
}

func MustStringToTime(timeStr string) time.Time {
	t, _ := time.ParseInLocation(secondsLayout, timeStr, time.Local)
	return t
}
