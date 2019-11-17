package task

import (
	"errors"
	"github.com/liuguangw/m3u8_download/common"
	"github.com/liuguangw/m3u8_download/io"
	"path/filepath"
	"strings"
)

func loadBaseDownloadTask(taskConfig *common.TaskConfig) (*DownloadTask, error) {
	if taskConfig.SaveDir == "" {
		return nil, errors.New("Config Error: save_dir Can't Be Empty")
	}
	if taskConfig.SaveFileName == "" {
		return nil, errors.New("Config Error: save_file_name Can't Be Empty")
	}
	//保存文件夹
	err := io.EnsureDir(taskConfig.SaveDir)
	if err != nil {
		return nil, errors.New("Create Save Dir(" + taskConfig.SaveDir + ") Error: " + err.Error())
	}
	//缓存文件夹名称
	cacheDirName := ""
	pos := strings.Index(taskConfig.SaveFileName, ".")
	if pos < 0 {
		cacheDirName = taskConfig.SaveFileName
	} else {
		cacheDirName = taskConfig.SaveFileName[0:pos]
	}
	//缓存文件夹路径
	cacheDir := filepath.Join(taskConfig.SaveDir, "cache", cacheDirName)
	err = io.EnsureDir(cacheDir)
	if err != nil {
		return nil, errors.New("Create Cache Dir(" + cacheDir + ") Error: " + err.Error())
	}
	return &DownloadTask{
		TaskConfig:       taskConfig,
		TaskNodes:        nil,
		CacheDir:         cacheDir,
		ServerM3u8Path:   filepath.Join(cacheDir, "000server.m3u8"),
		LocalM3u8Path:    filepath.Join(cacheDir, "000local.m3u8"),
		TaskDataFilePath: filepath.Join(cacheDir, "000task.txt"),
		EncryptKeyPath:   filepath.Join(cacheDir, "000key.ts"),
	}, nil
}
