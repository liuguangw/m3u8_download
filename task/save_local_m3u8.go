package task

import (
	"github.com/liuguangw/m3u8_download/common"
	"github.com/liuguangw/m3u8_download/io"
	"path/filepath"
	"strconv"
	"strings"
)

func SaveLocalM3u8(m3u8Info *common.M3u8Data, downloadTask *DownloadTask) error {
	localM3u8Info := &common.M3u8Data{
		ExtinfArr:     m3u8Info.ExtinfArr,
		TsUrls:        make([]string, len(m3u8Info.TsUrls)),
		OtherHeaders:  m3u8Info.OtherHeaders,
		EncryptMethod: m3u8Info.EncryptMethod,
		EncryptIv:     m3u8Info.EncryptIv,
	}
	if m3u8Info.EncryptMethod != "" {
		encryptKeyPath := downloadTask.EncryptKeyPath
		localM3u8Info.EncryptKeyUri = strings.Replace(encryptKeyPath, "\\", "/", -1)
	}
	for i := 0; i < len(m3u8Info.TsUrls); i++ {
		tsFilePath := filepath.Join(downloadTask.CacheDir, strconv.Itoa(i)+".ts")
		localM3u8Info.TsUrls[i] = strings.Replace(tsFilePath, "\\", "/", -1)
	}
	return io.WriteM3u8Content(downloadTask.LocalM3u8Path, localM3u8Info)
}
