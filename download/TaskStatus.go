package download

import (
	"errors"
	"io/ioutil"
	"strings"
)

const (
	STATUS_NOT_RUNNING byte = 48
	STATUS_RUNNING     byte = 37
	STATUS_SUCCESS     byte = 49
)

type TaskNode struct {
	DownloadUrl string
	Status      byte
}
type TaskStatus struct {
	FileEncryptInfo *EncryptInfo
	Nodes           []*TaskNode
}

func GetTaskStatus(config *TaskConfig) (*TaskStatus, error) {
	resp, err := FetchUrl(config, config.M3u8Url)
	if err != nil {
		return nil, errors.New("fetch m3u8 error: " + err.Error())
	}
	contentByteArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read m3u8 error: " + err.Error())
	}
	resp.Body.Close()
	m3u8Content := string(contentByteArr)
	//fmt.Println("[" + m3u8Content + "]")
	lines := strings.Split(m3u8Content, "\n")
	var encryptInfo *EncryptInfo
	taskStatusResult := &TaskStatus{
		Nodes: []*TaskNode{},
	}
	for _, lineContent := range lines {
		if strings.HasPrefix(lineContent, "#EXT-X-KEY:") {
			encryptInfo, err = readKeyFields(config, lineContent)
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(lineContent, "#") {
			continue
		} else {
			fileName := strings.Trim(lineContent, "\r")
			if fileName == "" {
				continue
			}
			taskStatusResult.Nodes = append(taskStatusResult.Nodes, &TaskNode{
				DownloadUrl: getItemUrl(config.M3u8Url, lineContent),
				Status:      STATUS_NOT_RUNNING,
			})
		}
	}
	taskStatusResult.FileEncryptInfo = encryptInfo
	/*fmt.Println(taskStatusResult.FileEncryptInfo)
	for _, debugInfo := range taskStatusResult.Nodes {
		fmt.Println("<" + debugInfo.DownloadUrl + ">")
	}*/
	return taskStatusResult, nil
}
