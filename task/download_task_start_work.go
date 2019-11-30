package task

import (
	"errors"
	"github.com/liuguangw/m3u8_download/io"
	"github.com/liuguangw/m3u8_download/tools"
	"time"
)

func (downloadTask *DownloadTask) StartWork() error {
	taskConfig := downloadTask.TaskConfig
	//创建需要的文件夹
	err := io.EnsureDir(taskConfig.SaveDir)
	if err != nil {
		return errors.New("Create Save Dir(" + taskConfig.SaveDir + ") Error: " + err.Error())
	}
	err = io.EnsureDir(downloadTask.CacheDir)
	if err != nil {
		return errors.New("Create Cache Dir(" + downloadTask.CacheDir + ") Error: " + err.Error())
	}
	//读取本地保存的任务状态
	var cachedTaskStatusArr []byte
	cachedTaskUrl := ""
	if io.FileExists(downloadTask.TaskDataFilePath) {
		cachedInfo, err := io.ReadTaskData(downloadTask.TaskDataFilePath)
		if err != nil {
			return errors.New("Read Task Data File Error: " + err.Error())
		}
		cachedTaskStatusArr = cachedInfo.TaskStatus
		cachedTaskUrl = cachedInfo.M3u8Url
	}
	m3u8CacheExists := (cachedTaskUrl == taskConfig.M3u8Url) && io.FileExists(downloadTask.ServerM3u8Path)
	//加载m3u8信息
	m3u8Info, err := loadTaskM3u8Info(m3u8CacheExists, downloadTask)
	if err != nil {
		return err
	}
	//生成并保存本地m3u8
	err = SaveLocalM3u8(m3u8Info, downloadTask)
	if err != nil {
		return errors.New("Create local m3u8 Error: " + err.Error())
	}
	//总文件数
	totalCount := len(m3u8Info.TsUrls)
	//已成功缓存的文件数
	successCachedCount := downloadTask.loadTaskNodes(m3u8Info, cachedTaskStatusArr)
	//初始化channel
	downloadTask.DownloadSuccessCount = make(chan int)
	downloadTask.NextTaskIndex = make(chan int)
	//
	if successCachedCount < totalCount {
		go downloadTask.runBackend()
		for i := 0; i < downloadTask.TaskConfig.MaxTask; i++ {
			go downloadTask.runDownload()
		}
		//保存任务数据文件
		err = cacheTaskStatus(downloadTask)
		if err != nil {
			tools.ShowErrorMessage("cache task Data Error: " + err.Error())
		}
		//开始计时
		taskStartTime := time.Now()
		tools.ShowSuccessMessage(tools.FormatDownloadProgress(successCachedCount, totalCount, &taskStartTime))
		for successCachedCount < totalCount {
			successCachedCount += <-downloadTask.DownloadSuccessCount
			//保存任务数据文件
			err = cacheTaskStatus(downloadTask)
			if err != nil {
				tools.ShowErrorMessage("cache task Data Error: " + err.Error())
			}
			tools.ShowSuccessMessage(tools.FormatDownloadProgress(successCachedCount, totalCount, &taskStartTime))
		}
	}
	return nil
}
