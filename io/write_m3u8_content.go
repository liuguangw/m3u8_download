package io

import (
	"github.com/liuguangw/m3u8_download/common"
	"io/ioutil"
)

func WriteM3u8Content(savePath string, m3u8Info *common.M3u8Data) error {
	fileContent := ""
	for _, extHeader := range m3u8Info.OtherHeaders {
		fileContent += extHeader + "\n"
	}
	if m3u8Info.EncryptMethod != "" {
		fileContent += "#EXT-X-KEY:METHOD=" + m3u8Info.EncryptMethod
		fileContent += ",URI=\"" + m3u8Info.EncryptKeyUri + "\""
		if m3u8Info.EncryptIv != "" {
			fileContent += ",IV=\"" + m3u8Info.EncryptIv + "\""
		}
		fileContent += "\n"
	}
	for infIndex, infLine := range m3u8Info.ExtinfArr {
		fileContent += infLine + "\n"
		fileContent += m3u8Info.TsUrls[infIndex] + "\n"
	}
	fileContent += "#EXT-X-ENDLIST\n"
	return ioutil.WriteFile(savePath, []byte(fileContent), 0644)
}
