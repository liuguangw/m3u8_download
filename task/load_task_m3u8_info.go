package task

import (
	"encoding/base64"
	"errors"
	"github.com/liuguangw/m3u8_download/common"
	"github.com/liuguangw/m3u8_download/io"
	"github.com/liuguangw/m3u8_download/tools"
	"io/ioutil"
)

func loadTaskM3u8Info(m3u8CacheExists bool, downloadTask *DownloadTask) (*common.M3u8Data, error) {
	taskConfig := downloadTask.TaskConfig
	//创建httpClient
	httpClient, err := io.CreateHttpClient(taskConfig)
	if err != nil {
		return nil, errors.New("Create http client Error: " + err.Error())
	}
	var m3u8Content string
	if m3u8CacheExists {
		//读取本地缓存的m3u8文件
		m3u8ContentBits, err := ioutil.ReadFile(downloadTask.ServerM3u8Path)
		if err != nil {
			return nil, errors.New("Read m3u8 File Error: " + err.Error())
		}
		m3u8Content = string(m3u8ContentBits)
	} else {
		tools.ShowCommonMessage("Downloading m3u8 file …")
		//请求m3u8文件
		m3u8ContentBits, err := io.FetchUrl(taskConfig.M3u8Url, httpClient, taskConfig)
		if err != nil {
			return nil, errors.New("Download m3u8 file Error: " + err.Error())
		}
		tools.ShowSuccessMessage("Download m3u8 file success")
		if taskConfig.EncodeType == "base64" {
			m3u8ContentBits, err = base64.StdEncoding.DecodeString(string(m3u8ContentBits))
			if err != nil {
				return nil, errors.New("Base64 decode m3u8 content Error: " + err.Error())
			}
		}
		m3u8Content = string(m3u8ContentBits)
	}
	m3u8Info, err := io.ReadM3u8Content(m3u8Content)
	if err != nil {
		return nil, errors.New("Parse m3u8 Error: " + err.Error())
	}
	if !m3u8CacheExists {
		//缓存m3u8文件
		err = io.WriteM3u8Content(downloadTask.ServerM3u8Path, m3u8Info)
		if err != nil {
			return nil,errors.New("Cache m3u8 Error: " + err.Error())
		}
		if m3u8Info.EncryptKeyUri != "" {
			tools.ShowCommonMessage("Downloading Key File …")
			//下载并缓存key
			keyFileUrl := tools.GetItemUrl(taskConfig.M3u8Url, m3u8Info.EncryptKeyUri)
			err = io.DownloadFile(keyFileUrl, httpClient, taskConfig, downloadTask.EncryptKeyPath)
			if err != nil {
				return nil,errors.New("Download Key File Error: " + err.Error())
			}
			tools.ShowSuccessMessage("Download Key File Success")
		}
	}
	return m3u8Info, nil
}
