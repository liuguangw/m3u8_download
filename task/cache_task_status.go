package task

import (
	"github.com/liuguangw/m3u8_download/common"
	"github.com/liuguangw/m3u8_download/io"
)

func cacheTaskStatus(downloadTask *DownloadTask) error {
	dataFilePath := downloadTask.TaskDataFilePath
	taskData := &common.TaskData{
		M3u8Url:    downloadTask.TaskConfig.M3u8Url,
		TaskStatus: make([]byte, len(downloadTask.TaskNodes)),
	}
	for taskIndex, taskInfo := range downloadTask.TaskNodes {
		taskData.TaskStatus[taskIndex] = taskInfo.Status
	}
	return io.WriteTaskData(dataFilePath, taskData)
}
