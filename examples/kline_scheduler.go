package main

import (
	"log"
	"time"
	"trade/db"
)

// SchedulerConfig 定时任务配置
type SchedulerConfig struct {
	Symbol   string        // 交易对
	Interval string        // K线周期
	Limit    int           // 获取数量
	Duration time.Duration // 执行间隔
}

func main() {
	// 初始化数据库
	log.Println("初始化数据库连接...")
	db.InitPostgreSql()

	// 定义需要定时更新的交易对配置
	configs := []SchedulerConfig{
		{
			Symbol:   "BTCUSDT",
			Interval: "1m",
			Limit:    1440, // 24小时的1分钟数据
			Duration: 1 * time.Minute,
		},
		{
			Symbol:   "ETHUSDT",
			Interval: "1m",
			Limit:    1440,
			Duration: 1 * time.Minute,
		},
		{
			Symbol:   "BNBUSDT",
			Interval: "1m",
			Limit:    1440,
			Duration: 1 * time.Minute,
		},
	}

	// 启动定时任务
	for _, config := range configs {
		go startScheduler(config)
	}

	// 保持主程序运行
	log.Println("定时任务已启动，按 Ctrl+C 停止...")
	select {}
}

// startScheduler 启动单个交易对的定时任务
func startScheduler(config SchedulerConfig) {
	// 立即执行一次
	log.Printf("[%s] 立即执行初始数据同步...", config.Symbol)
	err := db.FetchAndSaveKlines(config.Symbol, config.Interval, config.Limit)
	if err != nil {
		log.Printf("[%s] 初始同步失败: %v", config.Symbol, err)
	}

	// 创建定时器
	ticker := time.NewTicker(config.Duration)
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("[%s] 开始定时同步K线数据...", config.Symbol)
		err := db.FetchAndSaveKlines(config.Symbol, config.Interval, config.Limit)
		if err != nil {
			log.Printf("[%s] 同步失败: %v", config.Symbol, err)
		} else {
			log.Printf("[%s] 同步成功", config.Symbol)
		}
	}
}
