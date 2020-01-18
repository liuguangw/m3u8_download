package task

import (
	"github.com/liuguangw/m3u8_download/common"
	"github.com/liuguangw/m3u8_download/io"
	"path/filepath"
	"strconv"
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
		localM3u8Info.EncryptKeyUri = filepath.Base(downloadTask.EncryptKeyPath)
	}
	for i := 0; i < len(m3u8Info.TsUrls); i++ {
		localM3u8Info.TsUrls[i] = "out" + strconv.Itoa(i) + ".ts"
	}
	return io.WriteM3u8Content(downloadTask.LocalM3u8Path, localM3u8Info)
}
