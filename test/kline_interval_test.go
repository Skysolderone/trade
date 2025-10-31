package test

import (
	"fmt"
	"testing"

	"trade/db"
	"trade/model"
)

// TestKlineIntervals 测试不同时间周期的K线数据
func TestKlineIntervals(t *testing.T) {
	// 初始化数据库
	db.InitPostgreSql()

	intervals := []string{"1d", "8h", "4h", "2h", "1h"}
	symbols := []string{"BTCUSDT", "ETHUSDT"}

	fmt.Println("========== K线数据统计 ==========")
	for _, symbol := range symbols {
		fmt.Printf("\n交易对: %s\n", symbol)
		for _, interval := range intervals {
			var count int64
			var klines []model.Kline

			// 统计总数
			db.Pog.Model(&model.Kline{}).
				Where("symbol = ? AND interval = ?", symbol, interval).
				Count(&count)

			// 查询最新的5条记录
			db.Pog.Where("symbol = ? AND interval = ?", symbol, interval).
				Order("open_time DESC").
				Limit(5).
				Find(&klines)

			fmt.Printf("  %s: 共 %d 条数据\n", interval, count)
			if len(klines) > 0 {
				fmt.Printf("    最新记录: %s\n", klines[0].OpenTime.Format("2006-01-02 15:04:05"))
				fmt.Printf("    最早记录: %s\n", klines[len(klines)-1].OpenTime.Format("2006-01-02 15:04:05"))
			}
		}
	}
}
