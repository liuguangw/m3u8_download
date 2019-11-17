package main

import (
	"github.com/liuguangw/m3u8_download/task"
	"github.com/liuguangw/m3u8_download/tools"
	"os"
	"strconv"
)

func main() {
	configPath := "F:\\movie\\config.json"
	if len(os.Args) < 2 {
		noteStr := "Usage: m3u8_download [configFile]"
		tools.ShowErrorMessage(noteStr)
		//return
	} else {
		configPath = os.Args[1]
	}
	downloadTask, err := task.NewDownloadTask(configPath)
	if err != nil {
		tools.ShowError(err)
		return
	}
	//fmt.Println(downloadTask)
	successCount := downloadTask.CacahedSuccessCount
	totalCount := len(downloadTask.TaskNodes)
	if successCount < totalCount {
		go downloadTask.RunBackend()
		for i := 0; i < downloadTask.TaskConfig.MaxTask; i++ {
			go downloadTask.RunDownload()
		}
		//保存任务数据文件
		err = task.CacheTaskData(downloadTask)
		if err != nil {
			tools.ShowErrorMessage("cache task Data Error: " + err.Error())
		}
		tools.ShowSuccessMessage(strconv.Itoa(successCount) + "/" + strconv.Itoa(totalCount) + " files downloaded")
		for successCount < totalCount {
			successCount += <-downloadTask.DownloadSuccessCount
			//保存任务数据文件
			err = task.CacheTaskData(downloadTask)
			if err != nil {
				tools.ShowErrorMessage("cache task Data Error: " + err.Error())
			}
			tools.ShowSuccessMessage(strconv.Itoa(successCount) + "/" + strconv.Itoa(totalCount) + " files downloaded")
		}
	}
	tools.ShowSuccessMessage("all downloaded")
}
