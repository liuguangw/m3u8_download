package io

import (
	"github.com/liuguangw/m3u8_download/common"
	sysio "io"
	"os"
)

func DownloadFile(url string, taskConfig *common.TaskConfig, savePath string) error {
	response, err := doFetchUrl(url, taskConfig)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	fp, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = sysio.Copy(fp, response.Body)
	if err != nil {
		return err
	}
	return nil
}
