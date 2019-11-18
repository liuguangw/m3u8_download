package main

import (
	"fmt"
	"github.com/liuguangw/m3u8_download/task"
	"github.com/liuguangw/m3u8_download/tools"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	tools.ShowCommonMessage("Powered by liuguang@github https://github.com/liuguangw")
	configPath := ""
	if len(os.Args) < 2 {
		noteStr := "Usage: m3u8_download [configFile]"
		tools.ShowErrorMessage(noteStr)
		return
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
		taskStartTime := time.Now()
		tools.ShowSuccessMessage(tools.FormatDownloadProgress(successCount, totalCount, &taskStartTime))
		for successCount < totalCount {
			successCount += <-downloadTask.DownloadSuccessCount
			//保存任务数据文件
			err = task.CacheTaskData(downloadTask)
			if err != nil {
				tools.ShowErrorMessage("cache task Data Error: " + err.Error())
			}
			tools.ShowSuccessMessage(tools.FormatDownloadProgress(successCount, totalCount, &taskStartTime))
		}
	}
	tools.ShowSuccessMessage("all downloaded")
	//开始转码
	localM3u8Name := filepath.Base(downloadTask.LocalM3u8Path)
	outputPath := "../../" + downloadTask.TaskConfig.SaveFileName
	cmd := exec.Command("ffmpeg", "-i", localM3u8Name, "-c", "copy", outputPath)
	//工作目录
	cmd.Dir = downloadTask.CacheDir
	//接管Stdout
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		tools.ShowErrorMessage("Get StdoutPipe error: " + err.Error())
		return
	}
	//接管Stderr
	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		tools.ShowErrorMessage("Get StderrPipe error: " + err.Error())
		return
	}
	//接管Stdin
	cmdIn, err := cmd.StdinPipe()
	if err != nil {
		tools.ShowErrorMessage("Get StdinPipe error: " + err.Error())
		return
	}
	//用于覆盖文件时确认
	_, err = cmdIn.Write([]byte("y\n"))
	if err != nil {
		tools.ShowErrorMessage("Stdin write error: " + err.Error())
		return
	}
	//启动
	err = cmd.Start()
	if err != nil {
		tools.ShowErrorMessage("Run ffmpeg error: " + err.Error())
		return
	}
	//展示stdout
	go func() {
		buf := make([]byte, 1024)
		for {
			n, cErr := cmdOut.Read(buf)
			if n > 0 {
				fmt.Print(string(buf[:n]))
			} else {
				break
			}
			if cErr == io.EOF {
				break
			} else if cErr != nil {
				tools.ShowErrorMessage("Read cmd output failed: " + cErr.Error())
				break
			}
		}
	}()
	//展示stderr
	go func() {
		buf := make([]byte, 1024)
		for {
			n, cErr := cmdErr.Read(buf)
			if n > 0 {
				fmt.Print(string(buf[:n]))
			} else {
				break
			}
			if cErr == io.EOF {
				break
			} else if cErr != nil {
				tools.ShowErrorMessage("Read cmd error failed: " + cErr.Error())
				break
			}
		}
	}()
	err = cmd.Wait()
	if err != nil {
		tools.ShowErrorMessage("Run ffmpeg error: " + err.Error())
		return
	}
	//删除缓存文件夹
	if downloadTask.TaskConfig.CleanCacheAfterSuccess {
		tools.ShowCommonMessage("clean cache....")
		err = os.RemoveAll(downloadTask.CacheDir)
		if err != nil {
			tools.ShowErrorMessage("Delete Cache Error: " + err.Error())
			return
		}
	}
	fileSavePath := filepath.Join(downloadTask.TaskConfig.SaveDir, downloadTask.TaskConfig.SaveFileName)
	tools.ShowSuccessMessage("all complete -> " + fileSavePath)
}
