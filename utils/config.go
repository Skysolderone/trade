package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"trade/model"
)

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*model.Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析 JSON
	var config model.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}
