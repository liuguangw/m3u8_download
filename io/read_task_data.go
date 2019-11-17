package io

import (
	"errors"
	"github.com/liuguangw/m3u8_download/common"
	"io/ioutil"
	"strings"
)

func ReadTaskData(dataFilePath string) (*common.TaskData, error) {
	taskDataBits, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		return nil, err
	}
	fileContent := string(taskDataBits)
	pos := strings.Index(fileContent, "\n")
	if pos < 0 {
		return nil, errors.New("Borken Task Data file")
	}
	m3u8Url := fileContent[0:pos]
	taskStatus := make([]byte, 0, len(taskDataBits)-pos-1)
	for i := pos + 1; i < len(taskDataBits); i++ {
		nodeStatus := taskDataBits[i]
		if nodeStatus != common.STATUS_SUCCESS &&
			nodeStatus != common.STATUS_NOT_RUNNING &&
			nodeStatus != common.STATUS_RUNNING &&
			nodeStatus != common.STATUS_ERROR {
			continue
		}
		taskStatus = append(taskStatus, nodeStatus)
	}
	return &common.TaskData{
		M3u8Url:    m3u8Url,
		TaskStatus: taskStatus,
	}, nil
}
