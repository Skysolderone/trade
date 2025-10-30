package main

import (
	"fmt"
	"log"

	"trade/db"
	"trade/kline"
	"trade/strategy"
	"trade/utils"
)

func main() {
	// 初始化币安客户端(如果需要API密钥,请在这里填写)
	db.InitBinance("", "")

	// 初始化数据库连接
	db.InitPostgreSql()

	// 加载配置文件
	config, err := utils.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	fmt.Printf("开始更新 %d 个交易对的K线数据...\n", len(config.Symbols))

	// 遍历配置文件中的所有交易对和时间区间
	for _, symbolConfig := range config.Symbols {
		fmt.Printf("\n========== 处理交易对: %s ==========\n", symbolConfig.Symbol)

		for _, interval := range symbolConfig.Intervals {
			fmt.Printf("\n--- 时间区间: %s ---\n", interval)
			// 更新K线数据(从数据库最新记录开始更新到昨天)
			kline.UpdateKline(symbolConfig.Symbol, interval)
		}
	}

	fmt.Println("\n========== 所有数据更新完成 ==========")

	// 运行策略一: 分析历史同期涨跌概率（传入配置文件）
	fmt.Println("\n========== 开始运行策略一 ==========")
	strategy.Strategy1(config)
}
