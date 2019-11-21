package io

import (
	"errors"
	"github.com/liuguangw/m3u8_download/common"
	"net"
	"net/http"
	"net/url"
	"time"
)

func CreateHttpClient(taskConfig *common.TaskConfig) (*http.Client, error) {
	// 配置proxy
	proxyFn := http.ProxyFromEnvironment
	if taskConfig.Proxy != "" {
		proxyUrl, err := url.Parse(taskConfig.Proxy)
		if err != nil {
			return nil, errors.New("Parse proxy Url failed: " + err.Error())
		}
		proxyFn = http.ProxyURL(proxyUrl)
	}
	return &http.Client{
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
	}, nil
}
