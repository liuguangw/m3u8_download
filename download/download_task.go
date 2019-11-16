package download

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync"
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
	t.TaskStatus.Nodes[taskIndex].Status = STATUS_SUCCESS
	resultArr := make([]byte, len(t.TaskStatus.Nodes))
	for nodeIndex, nodeInfo := range t.TaskStatus.Nodes {
		resultArr[nodeIndex] = nodeInfo.Status
	}
	dataFile := filepath.Join(t.Config.SaveDir, "tmp", t.Config.FileName, "000task.data")
	_ = ioutil.WriteFile(dataFile, resultArr, 0644)
	_ = t.Bar.Incr()
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
		//fmt.Println("[index]downloading: ", taskIndex)
		taskInfo := t.TaskStatus.Nodes[taskIndex]
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
		if t.TaskStatus.FileEncryptInfo != nil {
			contentByteArr, err = DecryptData(contentByteArr,
				t.TaskStatus.FileEncryptInfo.Key, t.TaskStatus.FileEncryptInfo.Iv)
			if err != nil {
				taskInfo.Status = STATUS_NOT_RUNNING
				fmt.Println("decode ts error: " + err.Error())
				continue
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
	}
	t.DownloadComplete <- true
}
