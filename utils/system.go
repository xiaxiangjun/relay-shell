package utils

import (
	"os"
	"path/filepath"
)

var appDirPath string
var configDirPath string

func init() {
	dir := filepath.Dir(os.Args[0])
	path, err := filepath.Abs(dir)
	if nil != err {
		appDirPath = dir
	} else {
		appDirPath = path
	}
}

// GetExePath 获取进程所在目录
func GetExePath(name string) string {
	return filepath.Join(appDirPath, name)
}
