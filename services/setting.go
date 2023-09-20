package services

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type UserSettings struct {
	SelectedArea    string                `json:"selected_area"`
	SelectedStore   string                `json:"selected_store"`
	SelectedProduct string                `json:"selected_product"`
	BarkNotifyUrl   string                `json:"bark_notify_url"`
	ListenItems     map[string]ListenItem `json:"listen_items"`
}

// 保存配置到本地文件 SaveSettings saves settings to a file
func SaveSettings(settings UserSettings) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("user_settings.json", data, 0644)
}

// 加载缓存配置 LoadSettings loads settings from a file
func LoadSettings() (UserSettings, error) {
	var settings UserSettings
	data, err := ioutil.ReadFile("user_settings.json")
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(data, &settings)
	return settings, err
}

// 清空缓存配置 ClearSettings removes the settings file
func ClearSettings() error {
	return os.Remove("user_settings.json")
}
