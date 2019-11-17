package task

import (
	"errors"
	"github.com/liuguangw/m3u8_download/common"
	"github.com/liuguangw/m3u8_download/io"
	"io/ioutil"
)

func loadTaskM3u8Info(downloadTask *DownloadTask) (m3u8CacheExists bool,
	m3u8Info *common.M3u8Data,
	taskStatusArr []byte, err error) {
	m3u8CacheExists = false
	taskStatusArr = []byte{}
	taskConfig := downloadTask.TaskConfig
	//读取本地保存的任务状态
	if io.FileExists(downloadTask.TaskDataFilePath) {
		localTaskData, err := io.ReadTaskData(downloadTask.TaskDataFilePath)
		if err != nil {
			return m3u8CacheExists, m3u8Info, taskStatusArr, errors.New("Read Task Data File Error: " + err.Error())
		}
		taskStatusArr = localTaskData.TaskStatus
		if localTaskData.M3u8Url == taskConfig.M3u8Url && io.FileExists(downloadTask.ServerM3u8Path) {
			m3u8CacheExists = true
		}
	}
	var m3u8ContentBits []byte
	if m3u8CacheExists {
		//读取本地缓存的m3u8文件
		m3u8ContentBits, err = ioutil.ReadFile(downloadTask.ServerM3u8Path)
		if err != nil {
			return m3u8CacheExists, m3u8Info, taskStatusArr, errors.New("Read m3u8 File Error: " + err.Error())
		}
	} else {
		//请求m3u8文件
		m3u8ContentBits, err = io.FetchUrl(taskConfig.M3u8Url, taskConfig)
		if err != nil {
			return m3u8CacheExists, m3u8Info, taskStatusArr, errors.New("Fetch m3u8 url Error: " + err.Error())
		}
	}
	m3u8Info, err = io.ReadM3u8Content(string(m3u8ContentBits))
	if err != nil {
		return m3u8CacheExists, m3u8Info, taskStatusArr, errors.New("Parse m3u8 Error: " + err.Error())
	}
	return m3u8CacheExists, m3u8Info, taskStatusArr, nil
}
