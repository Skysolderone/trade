package main

import (
	"fmt"

	"trade/db"
	"trade/model"
)

func main() {
	// 初始化数据库
	db.InitPostgreSql()

	intervals := []string{"1d", "8h", "4h", "2h", "1h"}
	symbols := []string{"BTCUSDT", "ETHUSDT"}

	fmt.Println("========== K线数据统计 ==========\n")
	for _, symbol := range symbols {
		fmt.Printf("交易对: %s\n", symbol)
		for _, interval := range intervals {
			var count int64
			var firstKline, lastKline model.Kline

			// 统计总数
			db.Pog.Model(&model.Kline{}).
				Where("symbol = ? AND interval = ?", symbol, interval).
				Count(&count)

			if count > 0 {
				// 查询最早的记录
				db.Pog.Where("symbol = ? AND interval = ?", symbol, interval).
					Order("open_time ASC").
					Limit(1).
					Find(&firstKline)

				// 查询最新的记录
				db.Pog.Where("symbol = ? AND interval = ?", symbol, interval).
					Order("open_time DESC").
					Limit(1).
					Find(&lastKline)

				fmt.Printf("  %s: %d 条数据\n", interval, count)
				fmt.Printf("      最早: %s\n", firstKline.OpenTime.Format("2006-01-02 15:04:05"))
				fmt.Printf("      最新: %s\n", lastKline.OpenTime.Format("2006-01-02 15:04:05"))
			} else {
				fmt.Printf("  %s: 0 条数据 ❌\n", interval)
			}
		}
		fmt.Println()
	}
}
