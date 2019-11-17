package task

import (
	"errors"
	"github.com/liuguangw/m3u8_download/common"
	"github.com/liuguangw/m3u8_download/io"
	"github.com/liuguangw/m3u8_download/tools"
)

type DownloadTask struct {
	TaskConfig           *common.TaskConfig
	TaskNodes            []*common.DownloadTaskNode
	CacheDir             string     //缓存目录
	ServerM3u8Path       string     //下载的m3u8文件
	LocalM3u8Path        string     //本地生成的m3u8文件
	TaskDataFilePath     string     //任务数据文件路径
	EncryptKeyPath       string     //加密的key保存路径
	CacahedSuccessCount  int        //之前已缓存成功的文件数
	DownloadSuccessCount chan int   //成功下载的文件数
	NextTaskIndex        chan int   //获取下个下载任务的索引
}

func NewDownloadTask(configFilePath string) (*DownloadTask, error) {
	taskConfig, err := ReadTaskConfig(configFilePath)
	if err != nil {
		return nil, errors.New("Read Config Error: " + err.Error())
	}
	downloadTask, err := loadBaseDownloadTask(taskConfig)
	if err != nil {
		return nil, err
	}
	//本地的m3u8文件是否存在
	m3u8CacheExists, m3u8Info, taskStatusArr, err := loadTaskM3u8Info(downloadTask)
	if err != nil {
		return nil, err
	}
	taskStatusArrLength := len(taskStatusArr)
	if !m3u8CacheExists {
		//缓存m3u8文件
		err = io.WriteM3u8Content(downloadTask.ServerM3u8Path, m3u8Info)
		if err != nil {
			return nil, errors.New("Cache m3u8 Error: " + err.Error())
		}
		if m3u8Info.EncryptKeyUri != "" {
			//缓存key
			keyFileUrl := tools.GetItemUrl(taskConfig.M3u8Url, m3u8Info.EncryptKeyUri)
			err = io.DownloadFile(keyFileUrl, taskConfig, downloadTask.EncryptKeyPath)
			if err != nil {
				return nil, errors.New("Download Key Error: " + err.Error())
			}
		}
	}
	err = SaveLocalM3u8(m3u8Info, downloadTask)
	if err != nil {
		return nil, errors.New("Create local m3u8 Error: " + err.Error())
	}
	downloadTask.TaskNodes = make([]*common.DownloadTaskNode, len(m3u8Info.TsUrls))
	for tsIndex, tsUrl := range m3u8Info.TsUrls {
		tmpTaskNode := &common.DownloadTaskNode{
			TsUrl:  tools.GetItemUrl(taskConfig.M3u8Url, tsUrl),
			Status: common.STATUS_NOT_RUNNING,
		}
		// 根据任务数据文件，标记已完成的任务
		if tsIndex < taskStatusArrLength {
			if taskStatusArr[tsIndex] == common.STATUS_SUCCESS {
				tmpTaskNode.Status = common.STATUS_SUCCESS
				downloadTask.CacahedSuccessCount++
			}
		}
		downloadTask.TaskNodes[tsIndex] = tmpTaskNode
	}

	downloadTask.DownloadSuccessCount = make(chan int)
	downloadTask.NextTaskIndex = make(chan int)
	return downloadTask, nil
}
