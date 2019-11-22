package main

import (
	"github.com/liuguangw/m3u8_download/task"
	"github.com/liuguangw/m3u8_download/tools"
	"os"
	"path/filepath"
)

func main() {
	tools.ShowCommonMessage("Powered by liuguang@github https://github.com/liuguangw")
	//获取配置文件路径
	configPath := ""
	if len(os.Args) < 2 {
		noteStr := "Usage: m3u8_download [configFile]"
		tools.ShowErrorMessage(noteStr)
		return
	} else {
		configPath = os.Args[1]
	}
	//读取配置
	taskConfig, err := task.LoadTaskConfig(configPath)
	if err != nil {
		tools.ShowErrorMessage("Load Config Error: " + err.Error())
		return
	}
	//创建任务
	downloadTask, err := task.NewDownloadTask(taskConfig)
	if err != nil {
		tools.ShowError(err)
		return
	}
	//启动任务
	err = downloadTask.StartWork()
	if err != nil {
		tools.ShowError(err)
		return
	}
	tools.ShowSuccessMessage("all downloaded")
	//开始合并转码
	err = downloadTask.PackFile()
	if err != nil {
		tools.ShowError(err)
		return
	}
	fileSavePath := filepath.Join(downloadTask.TaskConfig.SaveDir, downloadTask.TaskConfig.SaveFileName)
	tools.ShowSuccessMessage("all complete -> " + fileSavePath)
}
