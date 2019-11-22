package task

import (
	"encoding/json"
	"errors"
	"github.com/liuguangw/m3u8_download/common"
	"io/ioutil"
	"strings"
)

func LoadTaskConfig(configPath string) (*common.TaskConfig, error) {
	//配置默认值
	config := &common.TaskConfig{
		ExtraHeaders: map[string]string{},
		TimeOut:      10,
		MaxTask:      8,
	}
	if configPath != "" {
		// 读取文件
		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		// json解析
		err = json.Unmarshal(data, config)
		if err != nil {
			return nil, err
		}
	}
	//检测各字段
	if err := checkTaskConfig(config); err != nil {
		return nil, err
	}
	return config, nil
}

func checkTaskConfig(taskConfig *common.TaskConfig) error {
	if taskConfig.M3u8Url == "" {
		return errors.New("m3u8_url can't be empty")
	}
	if taskConfig.EncodeType != "" && taskConfig.EncodeType != "base64" {
		return errors.New("un support encode_type: " + taskConfig.EncodeType)
	}
	if taskConfig.TimeOut < 0 {
		return errors.New("time_out must >= 0")
	}
	if taskConfig.MaxTask <= 0 {
		return errors.New("max_task must > 0")
	}
	if taskConfig.SaveFileName == "" {
		return errors.New("save_file_name can't be empty")
	} else if strings.Index(taskConfig.SaveFileName, ".") <= 0 {
		return errors.New("invalid save_file_name")
	}
	if taskConfig.SaveDir == "" {
		return errors.New("save_dir can't be empty")
	}
	if taskConfig.Proxy != "" {
		proxyPrefixValid := false
		supportPrefixArr := []string{
			"http://",
			"https://",
			"socks5://",
		}
		for _, proxyPrefix := range supportPrefixArr {
			if strings.HasPrefix(taskConfig.Proxy, proxyPrefix) {
				proxyPrefixValid = true
				break
			}
		}
		if !proxyPrefixValid {
			return errors.New("invalid proxy url")
		}
	}
	return nil
}
