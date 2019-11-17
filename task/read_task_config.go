package task

import (
	"encoding/json"
	"github.com/liuguangw/m3u8_download/common"
	"io/ioutil"
)

func ReadTaskConfig(coonfigPath string) (*common.TaskConfig, error) {
	//配置默认值
	config := &common.TaskConfig{
		M3u8Url:                "",
		ExtraHeaders:           map[string]string{},
		EncodeType:             "",
		TimeOut:                10,
		MaxTask:                8,
		SaveFileName:           "output.mp4",
		SaveDir:                "F:\\movie",
		CleanCacheAfterSuccess: false,
	}
	if coonfigPath != "" {
		// 读取文件
		data, err := ioutil.ReadFile(coonfigPath)
		if err != nil {
			return nil, err
		}
		// json解析
		err = json.Unmarshal(data, config)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}
