package strategy

import (
	"testing"

	"trade/db"
	"trade/model"
	"trade/utils"
)

// TestStrategy2 测试小时级别策略
func TestStrategy2(t *testing.T) {
	// 初始化数据库
	db.InitPostgreSql()

	// 加载配置文件
	config, err := utils.LoadConfig("../config.json")
	if err != nil {
		t.Fatalf("加载配置文件失败: %v", err)
	}

	// 运行小时级别策略
	Strategy2(config)
}

// TestStrategy2SingleSymbol 测试单个交易对的小时级别分析
func TestStrategy2SingleSymbol(t *testing.T) {
	// 初始化数据库
	db.InitPostgreSql()

	// 创建测试配置
	config := &model.Config{
		Symbols: []model.SymbolConfig{
			{
				Symbol:    "BTCUSDT",
				Intervals: []string{"1h"},
			},
		},
	}

	// 运行小时级别策略
	Strategy2(config)
}
