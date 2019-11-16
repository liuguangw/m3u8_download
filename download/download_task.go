package download

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MainDownloadTask struct {
	TaskStatus       *TaskStatus
	Config           *TaskConfig
	NewTaskIndex     chan int
	DownloadComplete chan bool
	mutex            sync.Mutex
	Bar              *uiprogress.Bar
}

func (t *MainDownloadTask) reportSuccess(taskIndex int) {
	t.mutex.Lock()
	totalNodes := t.TotalTaskCount()
	resultArrLength := totalNodes
	if totalNodes%10 == 0 {
		resultArrLength += totalNodes/10 - 1
	} else {
		resultArrLength += totalNodes / 10
	}
	resultArr := make([]byte, resultArrLength)
	resultArrIndex := 0
	for nodeIndex, nodeInfo := range t.TaskStatus.Nodes {
		if nodeIndex == taskIndex {
			resultArr[resultArrIndex] = STATUS_SUCCESS
		} else {
			resultArr[resultArrIndex] = nodeInfo.Status
		}
		if nodeIndex%100 == 99 {
			resultArrIndex++
			if resultArrIndex< resultArrLength{
				resultArr[resultArrIndex] = 10 //\n
			}

		} else if nodeIndex%10 == 9 {
			resultArrIndex++
			if resultArrIndex< resultArrLength{
				resultArr[resultArrIndex] = 32 //space
			}

		}
		resultArrIndex++
	}
	dataFile := filepath.Join(t.Config.SaveDir, "tmp", t.Config.FileName, "000task.data")
	_ = ioutil.WriteFile(dataFile, resultArr, 0644)
	_ = t.Bar.Incr()
	t.TaskStatus.Nodes[taskIndex].Status = STATUS_SUCCESS //
	t.mutex.Unlock()
}

func (t *MainDownloadTask) TotalTaskCount() int {
	return len(t.TaskStatus.Nodes)
}

func (t *MainDownloadTask) RunDownload() {
	for {
		taskIndex, ok := <-t.NewTaskIndex
		if !ok {
			break
		}
		taskInfo := t.TaskStatus.Nodes[taskIndex]
		//fmt.Println("[index]downloading: ", taskIndex,"  ",taskInfo.DownloadUrl)
		tsSavePath := filepath.Join(t.Config.SaveDir, "tmp", t.Config.FileName, strconv.Itoa(taskIndex)+".ts")
		//download
		resp, err := FetchUrl(t.Config, taskInfo.DownloadUrl)
		if err != nil {
			taskInfo.Status = STATUS_NOT_RUNNING
			fmt.Println("fetch ts error: " + err.Error())
			continue
		}
		contentByteArr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			taskInfo.Status = STATUS_NOT_RUNNING
			fmt.Println("read ts error: " + err.Error())
			continue
		}
		resp.Body.Close()
		encryptInfo := t.TaskStatus.FileEncryptInfo
		if encryptInfo != nil {
			if strings.HasPrefix(encryptInfo.Method, "AES") {
				//AES解密
				contentByteArr, err = AesDecryptData(contentByteArr, encryptInfo.Key, encryptInfo.Iv)
				if err != nil {
					taskInfo.Status = STATUS_NOT_RUNNING
					fmt.Println("decode ts error: " + err.Error())
					continue
				}
			}
		}
		err = ioutil.WriteFile(tsSavePath, contentByteArr, 0644)
		if err != nil {
			taskInfo.Status = STATUS_NOT_RUNNING
			fmt.Println("write ts error: " + err.Error())
			continue
		}
		t.reportSuccess(taskIndex)
	}
}

func (t *MainDownloadTask) RunBackend() {
	for {
		hasWork := false
		for i, taskNode := range t.TaskStatus.Nodes {
			if taskNode.Status != STATUS_SUCCESS {
				hasWork = true
			}
			if taskNode.Status == STATUS_NOT_RUNNING {
				taskNode.Status = STATUS_RUNNING
				t.NewTaskIndex <- i
			}
		}
		if !hasWork {
			close(t.NewTaskIndex)
			break
		}
		time.Sleep(time.Duration(1500) * time.Millisecond)
	}
	t.DownloadComplete <- true
}
