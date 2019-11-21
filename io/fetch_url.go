package io

import (
	"errors"
	"github.com/liuguangw/m3u8_download/common"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

func doFetchUrl(targetUrl string, taskConfig *common.TaskConfig) (*http.Response, error) {
	// 配置proxy
	proxyFn := http.ProxyFromEnvironment
	if taskConfig.Proxy != "" {
		proxyUrl, err := url.Parse(taskConfig.Proxy)
		if err != nil {
			return nil, errors.New("Parse proxy Url failed: " + err.Error())
		}
		proxyFn = http.ProxyURL(proxyUrl)
	}
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(taskConfig.TimeOut) * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			Proxy:                 proxyFn,
			MaxIdleConns:          taskConfig.MaxTask,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}
	request, err := http.NewRequest("GET", targetUrl, nil)
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
