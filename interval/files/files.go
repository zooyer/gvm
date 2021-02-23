package files

import (
	"github.com/zooyer/gvm/interval/debug"
	"os"
)

func Exists(filename string) bool {
	if _, err := os.Lstat(filename); err != nil && os.IsNotExist(err) {
		debug.Println("files: exists error:", err.Error())
		return false
	}
	return true
}

func IsDir(filename string) bool {
	stat, err := os.Lstat(filename)
	if err != nil {
		debug.Println("files: is dir error:", err.Error())
		return false
	}
	return stat.IsDir()
}

func IsFile(filename string) bool {
	stat, err := os.Lstat(filename)
	if err != nil {
		debug.Println("files: is file error:", err.Error())
		return false
	}
	return !stat.IsDir()
}

func AppendFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		debug.Println("files: append file error:", err.Error())
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
