package common

//任务状态
const (
	STATUS_NOT_RUNNING byte = 48  // ASCII 0
	STATUS_RUNNING     byte = 37  // ASCII %
	STATUS_SUCCESS     byte = 49  // ASCII 1
	STATUS_ERROR       byte = 120 // ASCII x
)

//单个任务信息
type DownloadTaskNode struct {
	TsUrl  string //ts文件url(绝对链接)
	Status byte
}
