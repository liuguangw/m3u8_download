package common

type M3u8Data struct {
	ExtinfArr    []string
	TsUrls       []string
	OtherHeaders []string

	EncryptMethod string //加密方式
	EncryptKeyUri string //加密的key的uri
	EncryptIv     string //iv
}
