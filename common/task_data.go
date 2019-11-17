package common

//表示保存到本地的任务数据文件的结构
type TaskData struct {
	M3u8Url    string
	TaskStatus []byte
}
