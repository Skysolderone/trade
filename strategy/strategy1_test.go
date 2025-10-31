package strategy

import (
	"testing"

	"trade/db"
	"trade/utils"
)

func TestStrategy1(t *testing.T) {
	db.InitPostgreSql()
	config, err := utils.LoadConfig("config.json")
	if err != nil {
		t.Fatalf("加载配置文件失败: %v", err)
	}
	if config == nil || len(config.Symbols) == 0 {
		t.Fatalf("配置文件为空")
	}
	Strategy1(config)
}
