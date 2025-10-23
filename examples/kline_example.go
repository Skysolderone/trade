package main

import (
	"log"
	"trade/db"
)

func main() {
	// 1. 初始化数据库连接
	log.Println("初始化数据库连接...")
	db.InitPostgreSql()

	// 2. 获取并保存BTC永续合约的1分钟K线数据（最近24小时，1440条）
	log.Println("开始获取并保存BTC K线数据...")
	err := db.FetchAndSaveKlines("BTCUSDT", "1m", 1440)
	if err != nil {
		log.Fatalf("获取并保存K线失败: %v", err)
	}

	// 3. 获取并保存ETH永续合约的1分钟K线数据
	log.Println("开始获取并保存ETH K线数据...")
	err = db.FetchAndSaveKlines("ETHUSDT", "1m", 1440)
	if err != nil {
		log.Fatalf("获取并保存K线失败: %v", err)
	}

	// 4. 从数据库查询最近100条BTC K线数据
	log.Println("从数据库查询BTC K线数据...")
	klines, err := db.GetKlinesFromDB("BTCUSDT", "1m", 100)
	if err != nil {
		log.Fatalf("查询K线失败: %v", err)
	}

	// 5. 打印前5条数据
	log.Printf("查询到 %d 条K线数据", len(klines))
	for i, kline := range klines {
		if i >= 5 {
			break
		}
		log.Printf("时间: %s, 开盘: %.2f, 最高: %.2f, 最低: %.2f, 收盘: %.2f, 成交量: %.4f",
			kline.OpenTimeDt.Format("2006-01-02 15:04:05"),
			kline.OpenPrice,
			kline.HighPrice,
			kline.LowPrice,
			kline.ClosePrice,
			kline.Volume)
	}

	log.Println("示例执行完成！")
}
