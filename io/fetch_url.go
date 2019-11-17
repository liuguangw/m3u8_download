package io

import (
	"errors"
	"github.com/liuguangw/m3u8_download/common"
	"io/ioutil"
	"net/http"
	"time"
)

func doFetchUrl(url string, taskConfig *common.TaskConfig) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          taskConfig.MaxTask,
			IdleConnTimeout:       30 * time.Second,
			ResponseHeaderTimeout: time.Duration(taskConfig.TimeOut) * time.Second,
		},
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	//headers
	for key, value := range taskConfig.ExtraHeaders {
		request.Header.Add(key, value)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Http " + response.Status + " Error")
	}
	return response, nil
}

func FetchUrl(url string, taskConfig *common.TaskConfig) ([]byte, error) {
	response, err := doFetchUrl(url, taskConfig)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contentByteArr, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return contentByteArr, nil
}
