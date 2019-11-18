package tools

import (
	"strconv"
	"time"
)

func formatTimeDuration(hour, minute, second int) string {
	result := strconv.Itoa(second) + "s"
	if hour == 0 && minute == 0 {
		return result
	}
	result = strconv.Itoa(minute) + "m:" + result
	if hour == 0 {
		return result
	}
	result = strconv.Itoa(hour) + "h:" + result
	return result
}

func calcDuration(duration time.Duration) string {
	if duration < time.Second {
		return formatTimeDuration(0, 0, 0)
	}
	// < 1分钟
	if duration < time.Minute {
		timeS := int(duration / time.Second)
		return formatTimeDuration(0, 0, timeS)
	}
	//秒
	timeS := int((duration % time.Minute) / time.Second)
	if duration < time.Hour {
		timeM := int(duration / time.Minute)
		return formatTimeDuration(0, timeM, timeS)
	}
	//分
	timeM := int((duration % time.Hour) / time.Minute)
	//小时数
	timeH := int(duration / time.Hour)
	return formatTimeDuration(timeH, timeM, timeS)
}
