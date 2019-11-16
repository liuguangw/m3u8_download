package download

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type TaskConfig struct {
	//mu38地址
	M3u8Url string `json:"m3u8_url"`

	ExtraHeaders map[string]string `json:"extra_headers"`
	//加密类型
	EncodeType string `json:"encode_type"`
	//连接超时时间
	TimeOut int `json:"time_out"`
	//最大同时下载数
	MaxTask int `json:"max_task"`
	//保存文件名
	FileName string `json:"file_name"`
	//保存文件夹路径
	SaveDir string `json:"save_dir"`
}

func ReadTaskConfig(coonfigPath string) (*TaskConfig, error) {
	config := &TaskConfig{
		M3u8Url:      "",
		ExtraHeaders: map[string]string{},
		EncodeType:   "",
		TimeOut:      10,
		MaxTask:      8,
		FileName:     "movie",
		SaveDir:      "F:\\movie",
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

func FetchUrl(config *TaskConfig, url string) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    config.MaxTask,
			IdleConnTimeout: 30 * time.Second,
		},
		//Timeout: time.Duration(config.TimeOut) * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range config.ExtraHeaders {
		req.Header.Add(key, value)
	}
	//fmt.Println(req,client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Http " + resp.Status + " Error")
	}
	return resp, nil
}
