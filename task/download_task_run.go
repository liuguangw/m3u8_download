package task

import (
	"github.com/liuguangw/m3u8_download/common"
	"github.com/liuguangw/m3u8_download/io"
	"github.com/liuguangw/m3u8_download/tools"
	"log"
	"path/filepath"
	"strconv"
)

func (downloadTask *DownloadTask) runDownload() {
	//创建httpClient
	httpClient, err := io.CreateHttpClient(downloadTask.TaskConfig)
	if err != nil {
		log.Fatalln("Create http client Error: " + err.Error())
	}
	for {
		taskIndex, ok := <-downloadTask.NextTaskIndex
		if !ok {
			break
		}
		taskInfo := downloadTask.TaskNodes[taskIndex]
		tsSavePath := filepath.Join(downloadTask.CacheDir, "out"+strconv.Itoa(taskIndex)+".ts")
		//download
		err := io.DownloadFile(taskInfo.TsUrl, httpClient, downloadTask.TaskConfig, tsSavePath)
		if err != nil {
			taskInfo.Status = common.STATUS_ERROR
			tools.ShowErrorMessage("save ts error: " + err.Error())
			continue
		}
		taskInfo.Status = common.STATUS_SUCCESS
		downloadTask.DownloadSuccessCount <- 1
	}
}
