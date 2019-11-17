package io

import (
	"github.com/liuguangw/m3u8_download/common"
	"io/ioutil"
)

func WriteTaskData(dataFilePath string, taskData *common.TaskData) error {
	headContent := taskData.M3u8Url + "\n"
	totalNodesCount := len(taskData.TaskStatus)
	resultArrTotalLength := len(headContent) + totalNodesCount
	if totalNodesCount%10 == 0 {
		resultArrTotalLength += totalNodesCount/10 - 1
	} else {
		resultArrTotalLength += totalNodesCount / 10
	}
	resultArr := make([]byte, 0, resultArrTotalLength)
	resultArr = append(resultArr, []byte(headContent)...)
	arrCurrentSize := len(resultArr)
	for nodeIndex, nodeStatus := range taskData.TaskStatus {
		resultArr = append(resultArr, nodeStatus)
		arrCurrentSize++
		if arrCurrentSize>= resultArrTotalLength {
			break
		}
		if nodeIndex%100 == 99 {
			resultArr = append(resultArr, 10) //\n
			arrCurrentSize++
		} else if nodeIndex%10 == 9 {
			resultArr = append(resultArr, 32) //space
			arrCurrentSize++
		}
	}
	return ioutil.WriteFile(dataFilePath, resultArr, 0644)
}
