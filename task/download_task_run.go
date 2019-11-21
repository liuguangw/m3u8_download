package task

import (
	"github.com/liuguangw/m3u8_download/common"
	"github.com/liuguangw/m3u8_download/io"
	"github.com/liuguangw/m3u8_download/tools"
	"log"
	"path/filepath"
	"strconv"
)

func (t *DownloadTask) RunDownload() {
	//创建httpClient
	httpClient, err := io.CreateHttpClient(t.TaskConfig)
	if err != nil {
		log.Fatalln("Create http client Error: " + err.Error())
	}
	for {
		taskIndex, ok := <-t.NextTaskIndex
		if !ok {
			break
		}
		taskInfo := t.TaskNodes[taskIndex]
		tsSavePath := filepath.Join(t.CacheDir, strconv.Itoa(taskIndex)+".ts")
		//download
		err := io.DownloadFile(taskInfo.TsUrl, httpClient, t.TaskConfig, tsSavePath)
		if err != nil {
			taskInfo.Status = common.STATUS_ERROR
			tools.ShowErrorMessage("save ts error: " + err.Error())
			continue
		}
		taskInfo.Status = common.STATUS_SUCCESS
		t.DownloadSuccessCount <- 1
	}
}
