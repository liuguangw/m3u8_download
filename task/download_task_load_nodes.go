package task

import (
	"github.com/liuguangw/m3u8_download/common"
	"github.com/liuguangw/m3u8_download/tools"
)

//加载任务节点，并返回已成功缓存过的文件数量
func (downloadTask *DownloadTask) loadTaskNodes(m3u8Info *common.M3u8Data,
	m3u8CacheExists bool, cachedTaskStatusArr []byte) int {
	successCachedCount := 0
	downloadTask.TaskNodes = make([]*common.DownloadTaskNode, len(m3u8Info.TsUrls))
	for tsIndex, tsUrl := range m3u8Info.TsUrls {
		tmpTaskNode := &common.DownloadTaskNode{
			//计算绝对链接
			TsUrl:  tools.GetItemUrl(downloadTask.TaskConfig.M3u8Url, tsUrl),
			Status: common.STATUS_NOT_RUNNING,
		}
		// 根据任务数据文件，标记已完成的任务
		if m3u8CacheExists && tsIndex < len(cachedTaskStatusArr) {
			if cachedTaskStatusArr[tsIndex] == common.STATUS_SUCCESS {
				tmpTaskNode.Status = common.STATUS_SUCCESS
				successCachedCount++
			}
		}
		downloadTask.TaskNodes[tsIndex] = tmpTaskNode
	}
	return successCachedCount
}
