package tools

import (
	"strconv"
	"time"
)

// 格式化下载进度
func FormatDownloadProgress(successCount, totalCount int, taskStartTime *time.Time) string {
	result := strconv.Itoa(successCount) + "/" + strconv.Itoa(totalCount) + " files downloaded"
	percent := calcPercent(successCount, totalCount)
	result += "(" + percent + ")"
	duration := time.Now().Sub(*taskStartTime)
	result += " [" + calcDuration(duration) + "]"
	return result
}
