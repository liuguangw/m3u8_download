package task

import (
	"errors"
	"github.com/liuguangw/m3u8_download/common"
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
	//创建httpClient
	httpClient, err := io.CreateHttpClient(taskConfig)
	if err != nil {
		return errors.New("Create http client Error: " + err.Error())
	}
	m3u8Info, err := loadTaskM3u8Info(m3u8CacheExists,
		downloadTask.ServerM3u8Path, taskConfig, httpClient)
	if !m3u8CacheExists {
		//缓存m3u8文件
		err = io.WriteM3u8Content(downloadTask.ServerM3u8Path, m3u8Info)
		if err != nil {
			return errors.New("Cache m3u8 Error: " + err.Error())
		}
		if m3u8Info.EncryptKeyUri != "" {
			tools.ShowCommonMessage("downloading key file")
			//下载并缓存key
			keyFileUrl := tools.GetItemUrl(taskConfig.M3u8Url, m3u8Info.EncryptKeyUri)
			err = io.DownloadFile(keyFileUrl, httpClient, taskConfig, downloadTask.EncryptKeyPath)
			if err != nil {
				return errors.New("Download Key Error: " + err.Error())
			}
			tools.ShowSuccessMessage("download key file success")
		}
	}
	//生成并保存本地m3u8
	err = SaveLocalM3u8(m3u8Info, downloadTask)
	if err != nil {
		return errors.New("Create local m3u8 Error: " + err.Error())
	}
	//初始化：成功缓存的文件数、总文件数
	successCachedCount := 0
	totalCount := len(m3u8Info.TsUrls)
	downloadTask.TaskNodes = make([]*common.DownloadTaskNode, totalCount)
	for tsIndex, tsUrl := range m3u8Info.TsUrls {
		tmpTaskNode := &common.DownloadTaskNode{
			TsUrl:  tools.GetItemUrl(taskConfig.M3u8Url, tsUrl),
			Status: common.STATUS_NOT_RUNNING,
		}
		// 根据任务数据文件，标记已完成的任务
		if tsIndex < len(cachedTaskStatusArr) {
			if cachedTaskStatusArr[tsIndex] == common.STATUS_SUCCESS {
				tmpTaskNode.Status = common.STATUS_SUCCESS
				successCachedCount++
			}
		}
		downloadTask.TaskNodes[tsIndex] = tmpTaskNode
	}
	//初始化channel
	downloadTask.DownloadSuccessCount = make(chan int)
	downloadTask.NextTaskIndex = make(chan int)
	//
	if successCachedCount < totalCount {
		go downloadTask.RunBackend()
		for i := 0; i < downloadTask.TaskConfig.MaxTask; i++ {
			go downloadTask.RunDownload()
		}
		//保存任务数据文件
		err = CacheTaskData(downloadTask)
		if err != nil {
			tools.ShowErrorMessage("cache task Data Error: " + err.Error())
		}
		taskStartTime := time.Now()
		tools.ShowSuccessMessage(tools.FormatDownloadProgress(successCachedCount, totalCount, &taskStartTime))
		for successCachedCount < totalCount {
			successCachedCount += <-downloadTask.DownloadSuccessCount
			//保存任务数据文件
			err = CacheTaskData(downloadTask)
			if err != nil {
				tools.ShowErrorMessage("cache task Data Error: " + err.Error())
			}
			tools.ShowSuccessMessage(tools.FormatDownloadProgress(successCachedCount, totalCount, &taskStartTime))
		}
	}
	return nil
}
