package iphotos

import (
	"os"
)

// 创建文件夹
func mkdirAll(path string) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) { //判断文件夹是否存在
		err = os.MkdirAll(path, os.ModePerm) //创建该文件夹
	}
	return err
}

// 对 os.ReadDir 处理，不进行排序
// filepath.WalkDir 需要把目录复制进内存
func readDir(name string) ([]os.DirEntry, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dirs, err := f.ReadDir(-1)
	return dirs, err
}
