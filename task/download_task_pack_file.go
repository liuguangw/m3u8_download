package task

import (
	"errors"
	"fmt"
	"github.com/liuguangw/m3u8_download/tools"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func (downloadTask *DownloadTask) PackFile() error {
	localM3u8Name := filepath.Base(downloadTask.LocalM3u8Path)
	outputPath := "../../" + downloadTask.TaskConfig.SaveFileName
	cmd := exec.Command("ffmpeg", "-i", localM3u8Name, "-c", "copy",
		"-metadata", "description=Packed by liuguangw/m3u8_download", outputPath)
	//工作目录
	cmd.Dir = downloadTask.CacheDir
	//接管Stdout
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return errors.New("Get StdoutPipe error: " + err.Error())
	}
	//接管Stderr
	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		return errors.New("Get StderrPipe error: " + err.Error())
	}
	//接管Stdin
	cmdIn, err := cmd.StdinPipe()
	if err != nil {
		return errors.New("Get StdinPipe error: " + err.Error())
	}
	//用于覆盖文件时确认
	_, err = cmdIn.Write([]byte("y\n"))
	if err != nil {
		return errors.New("Stdin write error: " + err.Error())
	}
	//启动
	err = cmd.Start()
	if err != nil {
		return errors.New("Run ffmpeg error: " + err.Error())
	}
	//展示stdout
	go showCommandOutput("cmd output", cmdOut)
	//展示stderr
	go showCommandOutput("cmd error", cmdErr)
	err = cmd.Wait()
	if err != nil {
		return errors.New("Exec ffmpeg error: " + err.Error())
	}
	//删除缓存文件夹
	if downloadTask.TaskConfig.CleanCacheAfterSuccess {
		tools.ShowCommonMessage("Clean cache....")
		err = os.RemoveAll(downloadTask.CacheDir)
		if err != nil {
			return errors.New("Clean Cache Error: " + err.Error())
		}
	}
	return nil
}

func showCommandOutput(commandType string, closer io.ReadCloser) {
	buf := make([]byte, 1024)
	for {
		n, err := closer.Read(buf)
		if n > 0 {
			fmt.Print(string(buf[:n]))
		}
		if err == io.EOF {
			break
		} else if err != nil {
			tools.ShowErrorMessage("Read " + commandType + " failed: " + err.Error())
			break
		}
	}
}
