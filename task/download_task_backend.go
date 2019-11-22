package task

import (
	"github.com/liuguangw/m3u8_download/common"
	"time"
)

func (downloadTask *DownloadTask) RunBackend() {
	for {
		hasWork := false
		for i, taskNode := range downloadTask.TaskNodes {
			if taskNode.Status != common.STATUS_SUCCESS {
				hasWork = true
			}
			if taskNode.Status == common.STATUS_NOT_RUNNING || taskNode.Status == common.STATUS_ERROR {
				taskNode.Status = common.STATUS_RUNNING
				downloadTask.NextTaskIndex <- i
			}
		}
		if !hasWork {
			close(downloadTask.NextTaskIndex)
			break
		}
		time.Sleep(time.Duration(1500) * time.Millisecond)
	}
}
