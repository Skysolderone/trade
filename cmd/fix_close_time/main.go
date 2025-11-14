package main

import (
	"fmt"
	"log"
	"time"

	"trade/db"
	"trade/model"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          修复详细记录的收盘时间                                ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// 初始化数据库
	db.InitPostgreSql()

	// 修复策略一的详细记录
	fixStrategy1DetailRecords()

	// 修复策略二的详细记录
	fixStrategy2DetailRecords()

	fmt.Println()
	fmt.Println("✅ 所有记录修复完成！")
}

// fixStrategy1DetailRecords 修复策略一的详细记录
func fixStrategy1DetailRecords() {
	fmt.Println("开始修复策略一详细记录...")

	// 查询所有需要修复的记录（close_time 早于 1900年）
	var records []model.Strategy1DetailRecord
	earlyTime := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	result := db.Pog.Where("close_time < ?", earlyTime).Find(&records)
	if result.Error != nil {
		log.Printf("查询策略一记录失败: %v", result.Error)
		return
	}

	fmt.Printf("找到 %d 条需要修复的记录\n", len(records))

	successCount := 0
	failCount := 0

	for i, record := range records {
		// 查询对应的结果记录，获取 symbol 和 interval
		var strategyResult model.Strategy1Result
		if err := db.Pog.Where("id = ?", record.ResultID).First(&strategyResult).Error; err != nil {
			log.Printf("查询策略结果失败 (result_id=%d): %v", record.ResultID, err)
			failCount++
			continue
		}

		// 从 year 字段解析日期（year 字段实际存储的是完整日期 YYYY-MM-DD）
		// 从 K 线表查询对应的收盘时间
		var kline model.Kline
		err := db.Pog.Where("symbol = ? AND interval = ? AND date = ?",
			strategyResult.Symbol,
			strategyResult.Interval,
			record.Year).
			First(&kline).Error

		if err != nil {
			log.Printf("未找到对应K线 (symbol=%s, interval=%s, date=%s): %v",
				strategyResult.Symbol, strategyResult.Interval, record.Year, err)
			failCount++
			continue
		}

		// 更新 close_time
		if err := db.Pog.Model(&record).Update("close_time", kline.CloseTime).Error; err != nil {
			log.Printf("更新记录失败 (id=%d): %v", record.ID, err)
			failCount++
			continue
		}

		successCount++

		// 每100条显示一次进度
		if (i+1)%100 == 0 {
			fmt.Printf("  进度: %d/%d\n", i+1, len(records))
		}
	}

	fmt.Printf("策略一修复完成: 成功 %d 条, 失败 %d 条\n\n", successCount, failCount)
}

// fixStrategy2DetailRecords 修复策略二的详细记录
func fixStrategy2DetailRecords() {
	fmt.Println("开始修复策略二详细记录...")

	// 查询所有需要修复的记录（close_time 早于 1900年）
	var records []model.Strategy2DetailRecord
	earlyTime := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	result := db.Pog.Where("close_time < ?", earlyTime).Find(&records)
	if result.Error != nil {
		log.Printf("查询策略二记录失败: %v", result.Error)
		return
	}

	fmt.Printf("找到 %d 条需要修复的记录\n", len(records))

	successCount := 0
	failCount := 0

	for i, record := range records {
		// 查询对应的结果记录，获取 symbol, interval 和 hour
		var strategyResult model.Strategy2Result
		if err := db.Pog.Where("id = ?", record.ResultID).First(&strategyResult).Error; err != nil {
			log.Printf("查询策略结果失败 (result_id=%d): %v", record.ResultID, err)
			failCount++
			continue
		}

		// 从 date 字段和 hour 查询对应的 K 线
		// date 字段存储的是完整日期 YYYY-MM-DD
		hourStr := fmt.Sprintf("%d", strategyResult.Hour)

		var kline model.Kline
		err := db.Pog.Where("symbol = ? AND interval = ? AND date = ? AND hour = ?",
			strategyResult.Symbol,
			strategyResult.Interval,
			record.Date,
			hourStr).
			First(&kline).Error

		if err != nil {
			log.Printf("未找到对应K线 (symbol=%s, interval=%s, date=%s, hour=%s): %v",
				strategyResult.Symbol, strategyResult.Interval, record.Date, hourStr, err)
			failCount++
			continue
		}

		// 更新 close_time
		if err := db.Pog.Model(&record).Update("close_time", kline.CloseTime).Error; err != nil {
			log.Printf("更新记录失败 (id=%d): %v", record.ID, err)
			failCount++
			continue
		}

		successCount++

		// 每100条显示一次进度
		if (i+1)%100 == 0 {
			fmt.Printf("  进度: %d/%d\n", i+1, len(records))
		}
	}

	fmt.Printf("策略二修复完成: 成功 %d 条, 失败 %d 条\n\n", successCount, failCount)
}
