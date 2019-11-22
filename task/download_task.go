package task

import (
	"github.com/liuguangw/m3u8_download/common"
	"path/filepath"
	"strings"
)

type DownloadTask struct {
	TaskConfig           *common.TaskConfig
	TaskNodes            []*common.DownloadTaskNode
	CacheDir             string   //缓存目录
	ServerM3u8Path       string   //下载的m3u8文件
	LocalM3u8Path        string   //本地生成的m3u8文件
	TaskDataFilePath     string   //任务数据文件路径
	EncryptKeyPath       string   //加密的key保存路径
	DownloadSuccessCount chan int //成功下载的文件数
	NextTaskIndex        chan int //获取下个下载任务的索引
}

func NewDownloadTask(taskConfig *common.TaskConfig) (*DownloadTask, error) {
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
	return &DownloadTask{
		TaskConfig:       taskConfig,
		TaskNodes:        nil,
		CacheDir:         cacheDir,
		ServerM3u8Path:   filepath.Join(cacheDir, "000server_cache.data"),
		LocalM3u8Path:    filepath.Join(cacheDir, "000local.m3u8"),
		TaskDataFilePath: filepath.Join(cacheDir, "000task.txt"),
		EncryptKeyPath:   filepath.Join(cacheDir, "000key.ts"),
	}, nil
}
