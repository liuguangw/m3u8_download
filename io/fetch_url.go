package io

import (
	"errors"
	"github.com/liuguangw/m3u8_download/common"
	"io/ioutil"
	"net/http"
)

func doFetchUrl(targetUrl string, client *http.Client, taskConfig *common.TaskConfig) (*http.Response, error) {
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
		return nil, errors.New("Http " + response.Status + " Error, URL: " + targetUrl)
	}
	return response, nil
}

func FetchUrl(url string, client *http.Client, taskConfig *common.TaskConfig) ([]byte, error) {
	response, err := doFetchUrl(url, client, taskConfig)
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
