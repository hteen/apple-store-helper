package config

import (
	"embed"
	"io/fs"
)

//go:embed files/*.json
var Files embed.FS

// ReadConfigFile 读取配置文件
func ReadConfigFile(filename string) ([]byte, error) {
	return Files.ReadFile("files/" + filename)
}

// MustReadConfigFile 读取配置文件
func MustReadConfigFile(filename string) []byte {
	data, err := Files.ReadFile("files/" + filename)
	if err != nil {
		panic(err)
	}
	return data
}

// ReadConfigDir 读取配置目录
func ReadConfigDir() (fs.ReadDirFS, error) {
	subFS, err := fs.Sub(Files, "files")
	if err != nil {
		return nil, err
	}
	return subFS.(fs.ReadDirFS), nil
}

// GetConfigFiles 获取所有配置文件
func GetConfigFiles() embed.FS {
	return Files
}
