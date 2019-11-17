package io

import (
	"errors"
	"github.com/liuguangw/m3u8_download/common"
	"strings"
)

func ReadM3u8Content(content string) (*common.M3u8Data, error) {
	m3u8Data := &common.M3u8Data{
		ExtinfArr:    []string{},
		TsUrls:       []string{},
		OtherHeaders: []string{},
	}
	lines := strings.Split(content, "\n")
	for _, lineContent := range lines {
		if lineContent == "" {
			continue
		}
		if strings.HasPrefix(lineContent, "#EXTINF:") {
			m3u8Data.ExtinfArr = append(m3u8Data.ExtinfArr, lineContent)
		} else if strings.HasPrefix(lineContent, "#EXT-X-KEY:") {
			pos := strings.Index(lineContent, ":")
			keyFields := strings.Split(lineContent[pos+1:], ",") //[METHOD=AES-128,URI="f90b70050200392f.ts",IV=0x1fd6e4d84de6fcbeb718173108446356]
			keyFieldsMap := map[string]string{}
			for _, keyFieldStr := range keyFields {
				keyFieldKvArr := strings.Split(keyFieldStr, "=")     // [METHOD,AES-128]
				keyFieldK := keyFieldKvArr[0]                        //METHOD
				keyFieldV := strings.Trim(keyFieldKvArr[1], "\" \r") //AES-128
				keyFieldsMap[keyFieldK] = keyFieldV
			}
			m3u8Data.EncryptMethod = keyFieldsMap["METHOD"]
			m3u8Data.EncryptKeyUri = keyFieldsMap["URI"]
			if mapV, ok := keyFieldsMap["IV"]; ok {
				m3u8Data.EncryptIv = mapV
			}
		} else if strings.HasPrefix(lineContent, "#") {
			if lineContent != "#EXT-X-ENDLIST" {
				m3u8Data.OtherHeaders = append(m3u8Data.OtherHeaders, lineContent)
			}
		} else {
			m3u8Data.TsUrls = append(m3u8Data.TsUrls, lineContent)
		}
	}
	if len(m3u8Data.TsUrls) == 0{
		return m3u8Data,errors.New("no ts url found")
	}
	return m3u8Data, nil
}
