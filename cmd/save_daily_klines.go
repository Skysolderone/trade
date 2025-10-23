package main

import (
	"log"
	"trade/db"
)

func main() {
	// 1. 初始化数据库连接
	log.Println("正在初始化数据库...")
	db.InitPostgreSql()

	// 2. 获取并保存BTC的所有1d K线数据（最多1500天）
	log.Println("开始获取BTC 1d K线数据...")
	err := db.FetchAndSaveKlines("BTCUSDT", "1d", 1500)
	if err != nil {
		log.Fatalf("获取并保存K线失败: %v", err)
	}

	log.Println("✅ BTC 1d K线数据保存成功！")

	// 3. 查询并显示统计信息
	klines, err := db.GetKlinesFromDB("BTCUSDT", "1d", 10)
	if err != nil {
		log.Fatalf("查询K线失败: %v", err)
	}

	log.Printf("数据库中共有 BTC 1d K线数据，最新10条如下：")
	for i, kline := range klines {
		log.Printf("[%d] 时间: %s, 开盘: %.2f, 最高: %.2f, 最低: %.2f, 收盘: %.2f, 成交量: %.2f",
			i+1,
			kline.OpenTimeDt.Format("2006-01-02"),
			kline.OpenPrice,
			kline.HighPrice,
			kline.LowPrice,
			kline.ClosePrice,
			kline.Volume)
	}

	log.Println("✅ 完成！")
}
