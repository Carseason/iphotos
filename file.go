package iphotos

import "os"

// 创建文件夹
func mkdirAll(path string) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) { //判断文件夹是否存在
		err = os.MkdirAll(path, os.ModePerm) //创建该文件夹
	}
	return err
}
