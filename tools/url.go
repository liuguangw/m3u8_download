package tools

import (
	"net/url"
	"strings"
)

//根据相对url计算绝对url
func GetItemUrl(currentUrl, filename string) string {
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