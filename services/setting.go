package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type UserSettings struct {
	SelectedArea    string                `json:"selected_area"`
	SelectedStore   string                `json:"selected_store"`
	SelectedProduct string                `json:"selected_product"`
	BarkNotifyUrl   string                `json:"bark_notify_url"`
	ListenItems     map[string]ListenItem `json:"listen_items"`
}

// SaveSettings 保存配置到本地文件
func SaveSettings(settings UserSettings) error {
	// 获取可执行文件所在目录
	execDir, err := os.Executable()
	if err != nil {
		// 如果获取可执行文件路径失败，使用当前工作目录
		execDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %v", err)
		}
	} else {
		execDir = filepath.Dir(execDir)
	}

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	filePath := filepath.Join(execDir, "user_settings.json")
	return os.WriteFile(filePath, data, 0644)
}

// LoadSettings 加载缓存配置
func LoadSettings() (UserSettings, error) {
	var settings UserSettings

	// 获取可执行文件所在目录
	execDir, err := os.Executable()
	if err != nil {
		// 如果获取可执行文件路径失败，使用当前工作目录
		execDir, err = os.Getwd()
		if err != nil {
			return settings, fmt.Errorf("failed to get current directory: %v", err)
		}
	} else {
		execDir = filepath.Dir(execDir)
	}

	filePath := filepath.Join(execDir, "user_settings.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(data, &settings)
	return settings, err
}

// ClearSettings 清空缓存配置
func ClearSettings() error {
	// 获取可执行文件所在目录
	execDir, err := os.Executable()
	if err != nil {
		// 如果获取可执行文件路径失败，使用当前工作目录
		execDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %v", err)
		}
	} else {
		execDir = filepath.Dir(execDir)
	}

	filePath := filepath.Join(execDir, "user_settings.json")
	return os.Remove(filePath)
}
