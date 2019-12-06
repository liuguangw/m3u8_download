package task

import (
	"github.com/liuguangw/m3u8_download/common"
	"github.com/liuguangw/m3u8_download/io"
	"os"
)

func cacheTaskStatus(downloadTask *DownloadTask) error {
	taskData := &common.TaskData{
		M3u8Url:    downloadTask.TaskConfig.M3u8Url,
		TaskStatus: make([]byte, len(downloadTask.TaskNodes)),
	}
	for taskIndex, taskInfo := range downloadTask.TaskNodes {
		taskData.TaskStatus[taskIndex] = taskInfo.Status
	}
	dataFilePath := downloadTask.TaskDataFilePath
	//先写到临时文件
	dataFileSwapPath := dataFilePath + ".swp"
	if err := io.WriteTaskData(dataFileSwapPath, taskData); err != nil {
		return err
	}
	//再rename
	return os.Rename(dataFileSwapPath, dataFilePath)
}
