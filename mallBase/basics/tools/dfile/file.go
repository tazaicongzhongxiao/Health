package dfile

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//
// GetCurrentDirectory
// @Description: 获取程序运行路径
// @return string
//
func GetCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}

//
// WriteFile
// @Description: 写文件
// @param filepath
// @param data
// @return error
//
func WriteFile(filepath string, data []byte) error {
	return ioutil.WriteFile(filepath, data, os.ModePerm)
}

//
// ReadFile
// @Description: 读文件
// @param filepath
// @return []byte
// @return error
//
func ReadFile(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

//
// Remove
// @Description: 删除文件
// @param path
// @return error
//
func Remove(path string) error {
	return os.Remove(path)
}

//
// RemoveAll
// @Description: 删除文件或文件夹
// @param path
// @return error
//
func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

//
// PathExists
// @Description: 检查目录或文件是否存在
// @param path
// @return isExist
//
func PathExists(path string) (isExist bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

//
// Mkdir
// @Description: 创建目录
// @param dir
// @return err
//
func Mkdir(dir string) (err error) {
	return os.MkdirAll(path.Dir(dir), os.ModePerm)
}

//
// EnsureDir
// @Description: 确保目录存在，如果没有，则创建它
// @param dir
// @return err
//
func EnsureDir(dir string) (err error) {
	parent := path.Dir(dir)
	if _, err = os.Stat(parent); os.IsNotExist(err) {
		if err = EnsureDir(parent); err != nil {
			return
		}
	}
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
	}
	return
}

//
// EnsureFile
// @Description:  确保文件存在，如果不存在，则创建它
// @param filepath
// @return err
//
func EnsureFile(filepath string) (err error) {
	var (
		file *os.File
	)
	if err = EnsureDir(path.Dir(filepath)); err != nil {
		return err
	}
	if _, err = os.Stat(filepath); os.IsNotExist(err) {
		file, err = os.Create(filepath)
		defer func() {
			file.Close()
		}()
	}
	return
}

//
// OuputFile
// @Description: 几乎与fs.WriteFile相同，不同之处在于如果目录不存在，则会创建该目录。
// @param filepath
// @param data
// @return error
//
func OuputFile(filepath string, data []byte) error {
	if err := EnsureDir(path.Dir(filepath)); err != nil {
		return err
	}
	return WriteFile(filepath, data)
}
