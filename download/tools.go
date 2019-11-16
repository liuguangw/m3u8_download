package download

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

func getItemUrl(currentUrl, filename string) string {
	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		return filename
	}
	urlInfo, _ := url.Parse(currentUrl)
	prefixUrl := urlInfo.Scheme + "://" + urlInfo.Host
	if !strings.HasPrefix(filename, "/") {
		sPos := strings.LastIndex(urlInfo.Path, "/")
		prefixUrl += urlInfo.Path[:sPos] + "/"
	}
	return prefixUrl + filename
}

func readKeyFields(config *TaskConfig, lineContent string) (*EncryptInfo, error) {
	resultMap := map[string]string{}
	content := lineContent[strings.Index(lineContent, ":")+1:]
	contentArr := strings.Split(content, ",")
	for _, sContent := range contentArr {
		sContentArr := strings.Split(sContent, "=")
		k := sContentArr[0]
		v := sContentArr[1]
		if v[0:1] == "\"" {
			v = v[1 : len(v)-1]
		}
		resultMap[k] = v
	}
	keyFileUrl := getItemUrl(config.M3u8Url, resultMap["URI"])
	resp, err := FetchUrl(config, keyFileUrl)
	if err != nil {
		return nil, errors.New("fetch key error: " + err.Error())
	}
	contentByteArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read key error: " + err.Error())
	}
	resp.Body.Close()
	iv, err := hex.DecodeString(resultMap["IV"][2:])
	if err != nil {
		return nil, errors.New("decode iv(" + resultMap["IV"] + ") error: " + err.Error())
	}
	return &EncryptInfo{
		Key:    contentByteArr,
		Iv:     iv,
		Method: resultMap["METHOD"],
	}, nil
}

func WriteTsList(filePath string, nodeCount int) error {
	fileContent := ""
	for i := 0; i < nodeCount; i++ {
		fileContent += " file '" + strconv.Itoa(i) + ".ts'"
		if i < nodeCount-1 {
			fileContent += "\n"
		}
	}
	return ioutil.WriteFile(filePath, []byte(fileContent), 0644)
}
