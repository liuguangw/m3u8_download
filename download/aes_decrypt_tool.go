package download

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func DecryptData(binData ,key,iv []byte)([]byte,error)  {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil,errors.New("key error: "+ err.Error())
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(binData, binData)
	return binData,nil
}
