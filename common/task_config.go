package common

//任务配置信息
type TaskConfig struct {
	//mu3u8地址
	M3u8Url string `json:"m3u8_url"`

	//下载时需要设置的额外http头
	ExtraHeaders map[string]string `json:"extra_headers"`

	//加密类型 有的m3u8文件内容是base64加密的
	EncodeType string `json:"encode_type"`

	//连接超时时间
	TimeOut int `json:"time_out"`

	//最大同时下载数
	MaxTask int `json:"max_task"`

	//保存文件名
	SaveFileName string `json:"save_file_name"`

	//保存文件夹路径
	SaveDir string `json:"save_dir"`
}