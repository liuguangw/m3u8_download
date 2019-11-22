package task

import (
	"encoding/base64"
	"errors"
	"github.com/liuguangw/m3u8_download/common"
	"github.com/liuguangw/m3u8_download/io"
	"github.com/liuguangw/m3u8_download/tools"
	"io/ioutil"
	"net/http"
)

func loadTaskM3u8Info(m3u8CacheExists bool, m3u8CachePath string,
	taskConfig *common.TaskConfig, client *http.Client) (*common.M3u8Data, error) {
	var m3u8Content string
	if m3u8CacheExists {
		//读取本地缓存的m3u8文件
		m3u8ContentBits, err := ioutil.ReadFile(m3u8CachePath)
		if err != nil {
			return nil, errors.New("Read m3u8 File Error: " + err.Error())
		}
		m3u8Content = string(m3u8ContentBits)
	} else {
		tools.ShowCommonMessage("downloading m3u8 file")
		//请求m3u8文件
		m3u8ContentBits, err := io.FetchUrl(taskConfig.M3u8Url, client, taskConfig)
		if err != nil {
			return nil, errors.New("Fetch m3u8 url Error: " + err.Error())
		}
		tools.ShowSuccessMessage("download m3u8 file success")
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
	return m3u8Info, nil
}
