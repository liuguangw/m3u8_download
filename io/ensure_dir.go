package io

import (
	"os"
)

//保证文件夹存在
func EnsureDir(dirName string) error {
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
